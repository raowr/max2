// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Users is the golang structure for table users.
type Users struct {
	Id       uint   `json:"id"       orm:"id"       description:"自增ID"`
	Name     string `json:"name"     orm:"name"     description:"用户名称"`
	Password string `json:"password" orm:"password" description:"用户密码"`
	Point    int64  `json:"point"    orm:"point"    description:"用户分数"`
	Token    string `json:"token"    orm:"token"    description:"登录token"`
}
