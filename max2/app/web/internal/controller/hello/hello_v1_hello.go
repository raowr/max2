package hello

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "web/api/hello/v1"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	// g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	g.RequestFromCtx(ctx).Response.WriteTpl("index.html", g.Map{
		"header":    "This is header",
		"container": "This is container",
		"footer":    "This is footer",
	})
	return
}
