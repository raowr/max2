// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// LogGame is the golang structure for table log_game.
type LogGame struct {
	Id        int         `json:"id"        orm:"id"         description:"自增id"`
	RoomId    string      `json:"roomId"    orm:"room_id"    description:"房间id"`
	Type      int         `json:"type"      orm:"type"       description:"房间类型：1段位房，2好友房"`
	Status    int         `json:"status"    orm:"status"     description:"房间状态"`
	UserId    string      `json:"userId"    orm:"user_id"    description:"用户id"`
	Point     int64       `json:"point"     orm:"point"      description:"用户分数"`
	Action    string      `json:"action"    orm:"action"     description:"行为"`
	Remain    string      `json:"remain"    orm:"remain"     description:"剩余的牌"`
	OutCards  string      `json:"outCards"  orm:"out_cards"  description:"打出的牌"`
	Text      string      `json:"text"      orm:"text"       description:"完整信息"`
	CreatTime *gtime.Time `json:"creatTime" orm:"creat_time" description:"创建时间"`
}
