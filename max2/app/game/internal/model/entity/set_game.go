// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// SetGame is the golang structure for table set_game.
type SetGame struct {
	Id    int    `json:"id"    orm:"id"    description:"自增id"`
	Name  string `json:"name"  orm:"name"  description:"设置名称"`
	Key   string `json:"key"   orm:"key"   description:"设置的key"`
	Value string `json:"value" orm:"value" description:"设置值"`
}
