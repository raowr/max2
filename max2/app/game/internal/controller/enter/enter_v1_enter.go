package enter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v1 "game/api/enter/v1"
	"game/internal/consts"
	"game/internal/controller/room"
	"game/internal/message"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
)

func (c *ControllerV1) Enter(ctx context.Context, req *v1.EnterReq) (res *v1.EnterRes, err error) {
	var (
		r          = g.RequestFromCtx(ctx)
		ws         *websocket.Conn
		msg        message.ChatMsg
		name       string
		msgByte    []byte
		rm         *room.RoomManager
		pid        int
		wsUpGrader = websocket.Upgrader{
			// CheckOrigin allows any origin in development
			// In production, implement proper origin checking for security
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			// Error handler for upgrade failures
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				// Implement error handling logic here
			},
		}
	)
	ws, err = wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		r.Response.Write(err.Error())
		return
	}
	for {
		// Blocking reading message from current websocket.
		_, msgByte, err = ws.ReadMessage()
		if err != nil {
			return nil, nil
		}
		// Message decoding.
		if err = gjson.DecodeTo(msgByte, &msg); err != nil {
			msgData := message.ChatMsg{
				Type: consts.ChatTypeError,
				Data: fmt.Sprintf(`invalid message: %s`, err.Error()),
				From: "",
			}
			_ = c.write(ws, msgData)
			continue
		}
		msg.From = name

		g.Log().Print(ctx, msg)

		if rm == nil {
			rm = room.NewRoomManager()
		}

		switch msg.Type {
		case consts.InitRoom:
			// Checks sending interval limit.
			// rm := room.NewRoomManager()
			playerName := "美女"
			humanPlayer := rm.CreatePlayer(playerName, room.Human)
			pid = humanPlayer.ID
			aiCount := 2
			roomInfo := rm.CreateRoom(humanPlayer, aiCount)
			players, _ := json.Marshal(roomInfo.Players)

			msgData := message.ChatMsg{
				Type: consts.InitRoom,
				Data: gconv.String(players),
				From: gconv.String(humanPlayer.ID),
			}
			_ = c.write(ws, msgData)
			select {
			case <-time.After(2 * time.Second):
				fmt.Println("\n时间到，已超时")
				//3秒后地3个人加入
				aiName := fmt.Sprintf("机器人%d号", len(roomInfo.Players))
				aiPlayer := rm.CreatePlayer(aiName, room.AI)
				rm.JoinRoom(aiPlayer, roomInfo.ID)
			}
			// room.PlayOneGame(roomInfo)
			//返回玩家列表
			players, _ = json.Marshal(roomInfo.Players)
			msgData = message.ChatMsg{
				Type: consts.InitRoom,
				Data: gconv.String(players),
				From: gconv.String(humanPlayer.ID),
			}
			_ = c.write(ws, msgData)
			// room.PlayOneGame(roomInfo)
		case consts.Play:
			var roomId string
			for _, v := range rm.PlayerList {
				if v.ID == pid {
					roomId = v.RoomID
					break
				}
			}
			roomInfo := rm.Rooms[roomId]
			go func() {
				for {
					roomMsg := <-roomInfo.MsgChan
					msgData := message.ChatMsg{
						Type: roomMsg.Type,
						Data: roomMsg.Data,
						From: gconv.String(pid),
					}
					_ = c.write(ws, msgData)
				}
			}()
			roomInfo.IsPlaying = true
			room.PlayOneGame(roomInfo)
		case consts.PlayCard:
			var roomId string
			player := &room.Player{}
			for _, v := range rm.PlayerList {
				if v.ID == pid {
					player = v
					roomId = v.RoomID
					break
				}
			}
			err := json.Unmarshal([]byte(msg.Data), &player.OutCards)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			roomInfo := rm.Rooms[roomId]
			go func() {
				roomMsg := <-roomInfo.MsgChan
				msgData := message.ChatMsg{
					Type: roomMsg.Type,
					Data: roomMsg.Data,
					From: gconv.String(pid),
				}
				_ = c.write(ws, msgData)
			}()
		}
	}

}

// write sends message to current client.
func (c *ControllerV1) write(ws *websocket.Conn, msg message.ChatMsg) error {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		return err
	}
	return ws.WriteMessage(ghttp.WsMsgText, msgBytes)
}

/*
{"type":"initRoom","data":"","name":""}
*/
