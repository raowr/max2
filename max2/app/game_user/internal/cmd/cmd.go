package cmd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gcmd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"game_user/internal/controller/enter"
	"game_user/internal/controller/log_game"
	"game_user/internal/controller/set_game"
	"game_user/internal/controller/settle"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http and grpc servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			var wg sync.WaitGroup

			// 初始化游戏配置缓存
			set_game.InitCache()

			// 读取配置，决定注册模式
			manualAddress := g.Cfg().MustGet(ctx, "game_user_svc.address").String()
			autoRegister := (manualAddress == "") // 为空时自动注册，非空时手动注册

			// 1. 启动 HTTP 服务
			wg.Add(1)
			go func() {
				defer wg.Done()

				// 如果是自动注册模式，提前设置 registry
				if autoRegister {
					gsvc.SetRegistry(etcd.New("127.0.0.1:2379"))
					g.Log().Info(ctx, "自动注册模式：已设置 etcd registry")
				} else {
					// 手动注册模式：确保 HTTP 服务启动时 registry 为 nil，避免自动注册
					// 注意：此时全局 registry 应该是 nil（默认），如果有其他组件已经设置，可以在这里临时清空？
					// 由于 gsvc.SetRegistry(nil) 会 panic，我们只能通过顺序保证。
					// 因为 gRPC 服务尚未启动，所以 registry 一定为 nil，无需额外操作。
					g.Log().Info(ctx, "手动注册模式：HTTP 服务不会自动注册")
				}

				s := g.Server()
				s.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(ghttp.MiddlewareHandlerResponse)
					group.Bind(enter.NewV1())
				})
				s.Run()
			}()

			// 等待 HTTP 服务完全启动（防止手动注册时 etcd 写入过快）
			time.Sleep(2 * time.Second)

			// 2. 启动 gRPC 服务器（它会设置 registry，但对已启动的 HTTP 无影响）
			wg.Add(1)
			go func() {
				defer wg.Done()

				c := grpcx.Server.NewConfig()
				if c.Address == "" {
					c.Address = ":9000"
				}

				grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
				c.Options = append(c.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(grpcx.Server.UnaryValidate),
				}...)
				gs := grpcx.Server.New(c)
				set_game.Register(gs)
				gs.Run()
			}()

			// 3. 手动注册模式：写入公网 IP 到 etcd
			if !autoRegister {
				go func() {
					address := manualAddress
					cli, err := clientv3.New(clientv3.Config{
						Endpoints:   []string{"127.0.0.1:2379"},
						DialTimeout: 5 * time.Second,
					})
					if err != nil {
						g.Log().Errorf(ctx, "创建 etcd client 失败: %v", err)
						return
					}
					// 不要 defer cli.Close()，保持连接用于续约

					leaseResp, err := cli.Grant(ctx, 15)
					if err != nil {
						g.Log().Errorf(ctx, "创建租约失败: %v", err)
						return
					}

					keepCtx := context.Background()
					keepAliveChan, err := cli.KeepAlive(keepCtx, leaseResp.ID)
					if err != nil {
						g.Log().Errorf(ctx, "启动续约失败: %v", err)
						return
					}
					go func() {
						for range keepAliveChan {
							// 续约中
						}
						g.Log().Warning(ctx, "租约续约通道已关闭")
					}()

					key := fmt.Sprintf("/service/default/default/%s/latest/%s", "game_user.svc", address)
					value := `{"insecure":true,"protocol":"http"}`
					_, err = cli.Put(ctx, key, value, clientv3.WithLease(leaseResp.ID))
					if err != nil {
						g.Log().Errorf(ctx, "注册失败: %v", err)
						return
					}
					g.Log().Infof(ctx, "手动注册公网 IP 成功: %s", key)

					// ---------- 新增：删除自动注册的内网 IP ----------
					prefix := "/service/default/default/game_user.svc/latest/"
					resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
					if err == nil {
						for _, kv := range resp.Kvs {
							if !strings.Contains(string(kv.Key), address) {
								// 删除非公网 IP 的记录
								if _, err := cli.Delete(ctx, string(kv.Key)); err != nil {
									g.Log().Errorf(ctx, "删除残留记录失败: %v", err)
								} else {
									g.Log().Infof(ctx, "已删除内网 IP 记录: %s", string(kv.Key))
								}
							}
						}
					}
					// -------------------------------------

					select {} // 阻塞，保持续约
				}()
			}
			// 等待所有服务结束
			wg.Wait()

			log_game.ShutdownLog()
			settle.ShutdownSettle()
			return nil
		},
	}
)
