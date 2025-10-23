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
			//删除所有房间
			for roomID, room := range rm.Rooms {
				room.Rgtimer.Close()
				delete(rm.Rooms, roomID)
			}

			rm.NextPlayerID = 0
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
				aiName := fmt.Sprintf("帅锅%d号", len(roomInfo.Players))
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
					//玩家输完钱,需要离开房间
					for _, player := range roomInfo.Players {
						if player.Type == room.Human && player.Point <= 0 {
							rm.LeaveRoom(player)
						}
					}
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
			data := struct {
				Pid     int   `json:"pid"`
				CardIds []int `json:"cardIds"`
				Pass    int   `json:"pass"`
			}{}
			err := json.Unmarshal([]byte(msg.Data), &data)
			if err != nil {
				g.Log().Error(ctx, err)
				return nil, err
			}
			if data.Pid != pid {
				g.Log().Error(ctx, fmt.Errorf("pid not equal"))
				return nil, fmt.Errorf("pid not equal")
			}
			player.OutCardIds = data.CardIds
			player.Pass = data.Pass
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
		case consts.GetInfo: //获取房间信息,用于再game页刷新
			var roomId string
			for _, v := range rm.PlayerList {
				if v.ID == pid {
					roomId = v.RoomID
					break
				}
			}
			roomInfo := rm.Rooms[roomId]
			// 显示玩家的牌（除了底牌）
			cards := make([]int, 0)
			cardsNum := make([]*message.PlayData, 0)
			playerPoint := int64(0)
			for _, player := range roomInfo.Players {
				//只能显示自己的牌
				if player.ID == pid {
					playerPoint = player.Point
					for _, card := range player.Cards {
						cards = append(cards, card.Id)
					}
				} else {
					cardsNum = append(cardsNum, &message.PlayData{
						Id:      player.ID,
						CardNum: player.CardNum,
					})
				}
			}
			//计算剩余出牌时间
			var remainOutCardTimeout, outCardTimeout int
			outCardTimeout = room.GetOutCardTimeout()
			if roomInfo.OutStarTime > 0 {
				now := int(time.Now().Unix())
				remainOutCardTimeout = outCardTimeout - (now - roomInfo.OutStarTime)
			} else {
				remainOutCardTimeout = outCardTimeout
			}

			go func() {
				data, _ := json.Marshal(struct {
					Cards                []int               `json:"cards"`
					Current              int                 `json:"current"`
					PlayerPoint          int64               `json:"playerPoint"`          //玩家总瓜子数
					OutCardTimeout       int                 `json:"outCardTimeout"`       //出牌最大时间(单位秒) /s
					RemainOutCardTimeout int                 `json:"remainOutCardTimeout"` //剩余出牌时间(单位秒) /s
					CardsNum             []*message.PlayData `json:"cardsNum"`             //机器人剩余牌数
					IsPlaying            bool                `json:"isPlaying"`            //是否进行中游戏
				}{
					Cards:                cards,
					Current:              roomInfo.Current,
					PlayerPoint:          playerPoint,
					OutCardTimeout:       outCardTimeout,       //出牌最大时间(单位秒) /s
					RemainOutCardTimeout: remainOutCardTimeout, //剩余出牌时间(单位秒) /s
					CardsNum:             cardsNum,
					IsPlaying:            roomInfo.IsPlaying,
				})
				msgData := message.ChatMsg{
					Type: consts.GetInfo,
					Data: gconv.String(data),
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
