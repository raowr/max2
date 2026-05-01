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
	ChatSessionName      = "ChatName"
	ChatSessionNameTemp  = "ChatNameTemp"
	ChatSessionNameError = "ChatNameError"
	ChatTypeSend         = "send"
	ChatTypeList         = "list"
	ChatTypeError        = "error"
	ChatIntervalLimit    = time.Second
	PlayerMsgPrefix      = "player_msg_"
	PlayerInfoPrefix     = "player_info_"
	RoomInfoPrefix       = "room_info_"
	RoomFriendMsgPrefix  = "room_friend_msg_"
	PlayerRoom           = "player_room"       //redis中玩家所在的房间
	PlayerPlayCardPrefix = "player_play_card_" //redis中玩家出牌的牌
)

var JwtKey = "db03d23b03ec405793b38f10592a2f34" // jwt密钥
