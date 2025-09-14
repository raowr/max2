package message

type ChatMsg struct {
	Type string      `json:"type" v:"required"`
	Data interface{} `json:"data" v:"required"`
	From string      `json:"name" v:""`
}

type PlayData struct {
	Id int `json:"id"` //用户id
}
