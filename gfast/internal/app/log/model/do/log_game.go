// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogGame is the golang structure of table log_game for DAO operations like Where/Data.
type LogGame struct {
	g.Meta    `orm:"table:log_game, do:true"`
	Id        any         // 自增id
	RoomId    any         // 房间id
	Type      any         // 房间类型：1段位房，2好友房
	Status    any         // 房间状态
	UserId    any         // 用户id
	Point     any         // 用户分数
	Action    any         // 行为
	Remain    any         // 剩余的牌
	OutCards  any         // 打出的牌
	Text      any         // 完整信息
	CreatTime *gtime.Time // 创建时间
}
