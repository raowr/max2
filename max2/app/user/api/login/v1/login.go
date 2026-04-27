package v1

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/login" tags:"Login" method:"post" summary:"Login"`
	Username string `json:"username" v:"required"`
	Password string `json:"password"`
}
type LoginRes struct {
	g.Meta   `mime:"application/json"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Node     string `json:"node"`
	Point    int64  `json:"point"`
}

type LogoutReq struct {
	g.Meta   `path:"/logout" tags:"Logout" method:"post" summary:"Logout"`
	Username string `json:"username" v:"required"`
	Token    string `json:"token"`
}
type LogoutRes struct {
	g.Meta `mime:"application/json"`
}
