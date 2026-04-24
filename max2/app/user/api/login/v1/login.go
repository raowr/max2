package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/login" tags:"Login" method:"post" summary:"Login"`
	Username string `json:"username" v:"required"`
	Password string `json:"password"`
}
type LoginRes struct {
	g.Meta `mime:"application/json"`
	Token  string `json:"token"`
	Node   string `json:"node"`
}
