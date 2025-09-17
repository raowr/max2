package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"web/internal/controller/game"
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
			s.AddSearchPath("D:/gowork/max2/max2/app/web/template")
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					index.NewV1(),
					room.NewV1(),
					game.NewV1(),
				)
			})
			s.SetPort(8001)
			s.Run()
			return nil
		},
	}
)
