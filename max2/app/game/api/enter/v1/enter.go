package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type EnterReq struct {
	g.Meta `path:"/enter" tags:"Enter" method:"get" summary:"You  enter max2"`
}
type EnterRes struct {
	g.Meta `mime:"json/xml" example:"string"`
}
