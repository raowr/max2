package cmd

import (
	"context"
	"sync"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"user/internal/controller/hello"
	"user/internal/controller/settle"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				s := g.Server()
				s.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(ghttp.MiddlewareHandlerResponse)
					group.Bind(
						hello.NewV1(),
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
				settle.Register(gs)
				gs.Run()
			}()
			wg.Wait()
			return nil
		},
	}
)
