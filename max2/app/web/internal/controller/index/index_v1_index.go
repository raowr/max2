package index

import (
	"context"

	v1 "web/api/index/v1"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	g.RequestFromCtx(ctx).Response.WriteTpl("index.html", g.Map{
		"header":    "This is header",
		"container": "This is container",
		"footer":    "This is footer",
	})
	return
}
