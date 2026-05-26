package consts

import "time"

const (
	InitRoom             = "initRoom"
	CreateRoom           = "createRoom"
	JoinRoom             = "joinRoom"
	LeaveRoom            = "leaveRoom"
	Play                 = "play"
	PlayCard             = "playCard"
	GetInfo              = "getInfo"
	Heartbeat            = "heartbeat"
	HealthTip            = "healthTip"
	ChatSessionName      = "ChatName"
	ChatSessionNameTemp  = "ChatNameTemp"
	ChatSessionNameError = "ChatNameError"
	ChatTypeSend         = "send"
	ChatTypeList         = "list"
	ChatTypeError        = "error"
	ChatIntervalLimit    = time.Second
	PlayerMsgPrefix      = "player:msg:"
	PlayerInfoPrefix     = "player:info:"
	RoomInfoPrefix       = "room:info:"
	RoomFriendMsgPrefix  = "room:friend:msg:"
	PlayerPlayCardPrefix = "player:play:card:" //redis中玩家出牌的牌
	PlayTimePrefix       = "playtime:"
	LastRemindPrefix     = "last_remind:"
)

var JwtKey = "db03d23b03ec405793b38f10592a2f34" // jwt密钥
