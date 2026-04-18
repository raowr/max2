package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type SetReq struct {
	g.Meta `path:"/set" tags:"Set" method:"post" summary:"You first set api"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
type SetRes struct {
	g.Meta `mime:"application/json" example:"string"`
}
