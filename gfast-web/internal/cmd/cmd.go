package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"web/internal/controller/index"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.View().SetPath("./resource/dist")
			s := g.Server()
			s.AddSearchPath("./resource/dist")
			s.AddSearchPath("./resource/public")            //处理微信可以访问
			s.AddStaticPath("/public", "./resource/public") // 映射公共目录
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					index.NewV1(),
				)
			})
			s.SetPort(8005)
			s.Run()
			return nil
		},
	}
)
