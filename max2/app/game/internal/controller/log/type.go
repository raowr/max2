package log

import "context"

// 行为：创建，加入，开始，发牌，出牌，不出，结算，离开
// 房间号，	房间类型，房间状态，用户ID，积分，	行为，	剩余牌, 完整信息
type LogInfo struct {
	RoomID     string `json:"roomID"`     //房间号
	Type       int    `json:"type"`       //房间类型 1比赛房，2好友房
	Status     int    `json:"status"`     //房间状态 0 未开始 1 游戏中 2 结算中
	UserID     string `json:"userID"`     //用户ID
	Point      int64  `json:"point"`      //积分
	Action     string `json:"action"`     //行为
	Remain     string `json:"remain"`     //剩余牌
	OutCardIds string `json:"outCardIds"` //玩家单次打出的牌id
	Text       string `json:"text"`       //完整信息
	// 其他字段...
}

// LogChan 日志通道，用于异步处理日志
var (
	LogChan chan LogInfo
	ctx     context.Context
	cancel  context.CancelFunc
)
