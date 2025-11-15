package consts

import "time"

const (
	InitRoom             = "initRoom"
	CreateRoom           = "createRoom"
	JoinRoom             = "joinRoom"
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
)
