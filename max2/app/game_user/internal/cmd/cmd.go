package cmd

import (
	"context"
	"sync"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gcmd"
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

			// 在 goroutine 中启动 HTTP 服务器
			wg.Add(1)
			go func() {
				defer wg.Done()
				gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))
				s := g.Server(`game_user.svc`)
				// 从配置文件读取服务地址
				if addr := g.Cfg().MustGet(ctx, `server:game_user.svc:address`).String(); addr != "" {
					s.SetAddr(addr)
				}
				// s := g.Server()
				s.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(ghttp.MiddlewareHandlerResponse)
					group.Bind(
						enter.NewV1(),
					)
				})
				s.Run()
			}()

			// 在 goroutine 中启动 gRPC 服务器
			wg.Add(1)
			go func() {
				defer wg.Done()

				// 使用配置文件中的设置或默认值
				c := grpcx.Server.NewConfig()
				if c.Address == "" {
					c.Address = ":9000" // 默认 gRPC 端口
				}

				grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
				c.Options = append(c.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(
						grpcx.Server.UnaryValidate,
					)}...,
				)
				gs := grpcx.Server.New(c)
				set_game.Register(gs)
				gs.Run()
			}()

			// 等待所有服务器完成（通常是无限期等待，直到收到中断信号）
			wg.Wait()

			log_game.ShutdownLog()  // 服务器停止后关闭日志系统
			settle.ShutdownSettle() // 服务器停止后关闭结算系统

			return nil
		},
	}
)
