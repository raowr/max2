package consts

import "time"

const (
	InitRoom             = "initRoom"
	Play                 = "play"
	PlayCard             = "playCard"
	GetInfo              = "getInfo"
	ChatSessionName      = "ChatName"
	ChatSessionNameTemp  = "ChatNameTemp"
	ChatSessionNameError = "ChatNameError"
	ChatTypeSend         = "send"
	ChatTypeList         = "list"
	ChatTypeError        = "error"
	ChatIntervalLimit    = time.Second
)
