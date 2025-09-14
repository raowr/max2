package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"web/internal/controller/game"
	"web/internal/controller/hello"
	"web/internal/controller/index"
	"web/internal/controller/room"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.AddSearchPath("C:/Users/rao/go/src/game1/max2/app/web/template")
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
					index.NewV1(),
					room.NewV1(),
					game.NewV1(),
				)
			})
			s.SetPort(8000)
			s.Run()
			return nil
		},
	}
)
