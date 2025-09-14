package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type RoomReq struct {
	g.Meta `path:"/room" tags:"Room" method:"get" summary:"You first room temp"`
}
type RoomRes struct {
	g.Meta `mime:"text/html" example:"string"`
}
