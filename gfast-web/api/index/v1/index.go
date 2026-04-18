package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type IndexReq struct {
	g.Meta `path:"/index" tags:"Index" method:"get" summary:"You first index temp"`
}
type IndexRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
