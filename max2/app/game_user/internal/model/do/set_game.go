// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SetGame is the golang structure of table set_game for DAO operations like Where/Data.
type SetGame struct {
	g.Meta `orm:"table:set_game, do:true"`
	Id     any // 自增id
	Name   any // 设置名称
	Key    any // 设置的key
	Value  any // 设置值
}
