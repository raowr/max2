package v1

import "github.com/gogf/gf/v2/frame/g"

type RegisterReq struct {
	g.Meta    `path:"/register" tags:"Register" method:"post" summary:"Register"`
	Username  string `json:"username" v:"required"`
	Password  string `json:"password" v:"required"`
	Password2 string `json:"password2" v:"required"`
}
type RegisterRes struct {
	g.Meta `mime:"application/json"`
}
