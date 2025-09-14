package game

import (
	"context"

	v1 "web/api/game/v1"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Game(ctx context.Context, req *v1.GameReq) (res *v1.GameRes, err error) {
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
	g.RequestFromCtx(ctx).Response.WriteTpl("game.html", g.Map{
		"header":    "This is header",
		"container": "This is container",
		"footer":    "This is footer",
	})
	return
}
