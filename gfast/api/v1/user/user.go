package user

import (
	"gfast/internal/app/user/model/entity"

	commonApi "gfast/api/v1/common"

	"github.com/gogf/gf/v2/frame/g"
)

type UserListReq struct {
	g.Meta `path:"/list" tags:"用户管理" method:"get" summary:"用户列表"`
	Name   string `p:"name"` //用户名
	commonApi.PageReq
}

type UserListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	List []*entity.Users `json:"list"`
}

type UserUpdatePointReq struct {
	g.Meta `path:"/updatePoint" tags:"用户管理" method:"post" summary:"更新用户积分"`
	Id     int `p:"id"`    //id
	Point  int `p:"point"` //积分
	commonApi.PageReq
}

type UserUpdatePointRes struct {
	g.Meta `mime:"application/json"`
}

type UserUpdatePasswordReq struct {
	g.Meta   `path:"/updatePassword" tags:"用户管理" method:"post" summary:"更新用户密码"`
	Id       int    `p:"id"`       //id
	Password string `p:"password"` //密码
	commonApi.PageReq
}
type UserUpdatePasswordRes struct {
	g.Meta `mime:"application/json"`
}
