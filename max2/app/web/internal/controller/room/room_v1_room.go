package room

import (
	"context"

	v1 "web/api/room/v1"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Room(ctx context.Context, req *v1.RoomReq) (res *v1.RoomRes, err error) {
	// return nil, gerror.NewCode(gcode.CodeNotImplemented)
	g.RequestFromCtx(ctx).Response.WriteTpl("room.html", g.Map{
		"header":    "This is header",
		"container": "This is container",
		"footer":    "This is footer",
	})
	return
}
