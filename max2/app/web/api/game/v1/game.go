package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GameReq struct {
	g.Meta `path:"/game" tags:"Game" method:"get" summary:"You first game temp"`
}
type GameRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
