package cmd

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/zeromicro/go-zero/core/discov"

	"gate-service/internal/controller/action"
	"gate-service/internal/controller/enter"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			var wg sync.WaitGroup
			// 1. HTTP 服务器
			wg.Add(1)
			go func() {
				defer wg.Done()
				s := g.Server()
				s.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(ghttp.MiddlewareHandlerResponse)
					group.Bind(
						enter.NewV1(),
					)
				})
				s.Run()
			}()

			// 2. 使用 go-zero Publisher 注册服务到 etcd
			wg.Add(1)
			go func() {
				defer wg.Done()

				serverName := g.Cfg().MustGet(ctx, "server.name").String()
				address := g.Cfg().MustGet(ctx, "gate_service_svc.address").String()
				if address == "" {
					address = "127.0.0.1:8010"
				}

				// 正确的 key 格式: /discov/service_name/
				key := "/discov/" + serverName + "/"

				g.Log().Infof(ctx, "准备注册服务: serverName=%s, address=%s, key=%s", serverName, address, key)

				// 创建 Publisher
				publisher := discov.NewPublisher(
					[]string{"127.0.0.1:2379"},
					key,
					address,
				)

				// 启动自动续期（阻塞）
				g.Log().Info(ctx, "开始服务注册...")
				if err := publisher.KeepAlive(); err != nil {
					g.Log().Errorf(ctx, "服务注册失败: %v", err)
					return
				}

				// 这里应该不会执行到，除非 KeepAlive 意外返回
				g.Log().Warningf(ctx, "KeepAlive 意外返回，服务 [%s] 已开始", serverName)
			}()
			wg.Wait()
			action.ShutdownAction()
			return nil
		},
	}
)
