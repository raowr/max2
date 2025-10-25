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
	"github.com/gogf/gf/v2/frame/g" // 修正：GFv2 正确导入路径
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
)

func (c *ControllerV1) Enter(ctx context.Context, req *v1.EnterReq) (res *v1.EnterRes, err error) {
	var (
		r          = g.RequestFromCtx(ctx)
		ws         *websocket.Conn
		wsUpGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 从配置读取 debug 模式（推荐：配置文件中配置 app.debug = true/false）
				// MustGet: 键不存在时 panic；若想避免 panic，用 Get().Bool() 并处理默认值
				//isDebug := g.Cfg().MustGet(ctx, "app.debug").Bool()
				//if isDebug {
				//	return true // 调试模式允许所有跨域
				//}
				//// 2. Android 原生客户端（Origin 为空字符串）
				//origin := r.Header.Get("Origin")
				//allowedOrigins := map[string]bool{
				//	allowedOrigin: true, // 浏览器前端域名
				//	"": true,            // 允许空 Origin（Android 原生）
				//}
				//return allowedOrigins[origin]
				return true
			},
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				g.Log().Errorf(ctx, "WebSocket升级失败: %v", reason)
				w.WriteHeader(status)
				w.Write([]byte(reason.Error()))
			},
		}
	)

	// 升级HTTP连接为WebSocket
	ws, err = wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Errorf(ctx, "WebSocket升级失败: %v", err)
		r.Response.WriteHeader(http.StatusInternalServerError)
		r.Response.Write([]byte("连接建立失败"))
		return
	}

	// 获取或生成用户ID
	userID := r.GetQuery("user_id").String()
	if userID == "" {
		userID = room.GenerateUserID()
		g.Log().Infof(ctx, "用户生成新ID: %s", userID)
	}

	// 创建客户端实例（增大消息缓冲）
	client := &Client{
		conn:       ws,
		userID:     userID,
		heartbeat:  time.Now(),
		sendChan:   make(chan []byte, 10000), // 缓冲增大至10000，减少阻塞
		roomIdChan: make(chan string, 1),     // 房间ID通道（缓冲1条）
	}

	// 加锁更新客户端连接（重连逻辑）
	clientsMu.Lock()
	if oldClient, err := getClient(userID); err == nil && oldClient != nil {
		oldClient.conn.Close() // 关闭旧连接
		oldClient.closed = true
		g.Log().Infof(ctx, "用户 %s 重连，关闭旧连接", userID)
	}
	if err = addClient(client); err != nil {
		g.Log().Errorf(ctx, "添加客户端到缓存失败: %v", err)
		ws.Close()
		return
	}
	clientsMu.Unlock()

	g.Log().Infof(ctx, "用户 %s 连接成功", userID)

	// 启动协程（增加退出日志）
	go client.readLoop(ctx)
	go client.readServe(ctx)
	go client.writeLoop(ctx)
	go client.heartbeatCheck(ctx)

	// 优化后
	<-ctx.Done()
	g.Log().Infof(ctx, "用户 退出（上下文关闭）")
	return
}

// 读消息循环：处理客户端消息
func (c *Client) readLoop(ctx context.Context) {
	defer func() {
		// 清理资源（带锁）
		clientsMu.Lock()
		delete(clients, c.userID)
		clientsMu.Unlock()

		c.conn.Close()
		g.Log().Infof(ctx, "用户 %s 断开连接（readLoop退出）", c.userID)
	}()

	for {
		if c.closed {
			return
		}
		// 读取客户端消息
		mt, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				g.Log().Errorf(ctx, "用户 %s 消息读取错误: %v", c.userID, err)
			}
			return // 连接断开，退出循环
		}

		// 只要成功接收到消息，无论后续处理是否成功，都更新心跳
		c.heartbeat = time.Now()

		// 解析消息（严格错误处理）
		var msg message.ChatMsg
		if err := gconv.Struct(data, &msg); err != nil {
			errMsg := message.ChatMsg{
				Type: consts.ChatTypeError,
				Data: fmt.Sprintf("无效消息格式: %v", err),
				From: "",
			}
			c.sendChan <- c.encodeMessage(ctx, errMsg)
			continue
		}
		// 根据消息类型处理业务（所有操作加锁保护全局rm）
		switch msg.Type {
		case consts.InitRoom:
			c.handleInitRoom(ctx)
		case consts.Play:
			c.handlePlay(ctx)
		case consts.PlayCard:
			c.handlePlayCard(ctx, msg.Data)
		case consts.GetInfo:
			c.handleGetInfo(ctx)
		case consts.Heartbeat:
			c.handleHeartbeat(ctx, msg.Data)
		default:
			errMsg := message.ChatMsg{
				Type: consts.ChatTypeError,
				Data: "未知消息类型",
				From: "",
			}
			c.sendChan <- c.encodeMessage(ctx, errMsg)
		}

		g.Log().Infof(ctx, "用户 %s 接收消息: %s（类型: %d）", c.userID, data, mt)
	}
}

// 处理房间初始化
func (c *Client) handleInitRoom(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()

	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）
	var oldRoomID string
	for _, player := range rm.PlayerList {
		if player.UserId == c.userID {
			oldRoomID = player.RoomID
			break
		}
	}
	if oldRoomID != "" {
		if roomInfo, ok := rm.Rooms[oldRoomID]; ok {
			if roomInfo.Rgtimer != nil {
				roomInfo.Rgtimer.Stop()
				roomInfo.Rgtimer.Close() // 停止定时器（确保资源释放）
			}
			delete(rm.Rooms, oldRoomID)
			g.Log().Infof(ctx, "用户 %s 清理旧房间: %s", c.userID, oldRoomID)
		}
		// 从玩家列表移除旧玩家
		delete(rm.PlayerList, c.userID)
	}

	// 创建新房间和玩家
	roomInfo := rm.CreateRoom()
	playerName := "美女"
	humanPlayer := roomInfo.CreatePlayer(playerName, room.Human)
	humanPlayer.UserId = c.userID
	c.pid = humanPlayer.ID
	rm.PlayerList[humanPlayer.UserId] = humanPlayer // 关联用户与玩家
	// 发送房间ID到客户端
	c.roomIdChan <- roomInfo.ID

	// 创建AI玩家
	aiCount := 2
	for i := 0; i < aiCount; i++ {
		aiName := fmt.Sprintf("帅锅%d号", i+1)
		roomInfo.CreatePlayer(aiName, room.AI)
	}

	// 推送玩家列表（处理JSON错误）
	players, err := json.Marshal(roomInfo.Players)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
		return
	}
	msgData := message.ChatMsg{
		Type: consts.InitRoom,
		Data: gconv.String(players),
		From: gconv.String(humanPlayer.ID),
	}
	c.sendChan <- c.encodeMessage(ctx, msgData)

	// 延迟添加额外AI（带上下文超时，避免协程泄漏）
	go func(roomID string, userID string, clientPid int) {
		// 使用Background上下文确保协程能执行完
		localCtx := context.Background()
		// 添加日志记录协程启动
		g.Log().Infof(localCtx, "用户 %s 开始延迟添加AI到房间 %s", userID, roomID)

		<-time.After(2 * time.Second)
		rmMu.Lock()
		defer rmMu.Unlock()

		roomInfo, ok := rm.Rooms[roomID]
		if !ok {
			g.Log().Warningf(ctx, "房间 %s 已不存在，停止添加AI", roomID)
			return
		}
		aiName := fmt.Sprintf("帅锅%d号", len(roomInfo.Players))
		roomInfo.CreatePlayer(aiName, room.AI)

		// 再次推送更新后的玩家列表
		players, err := json.Marshal(roomInfo.Players)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
			return
		}
		msgData := message.ChatMsg{
			Type: consts.InitRoom,
			Data: gconv.String(players),
			From: gconv.String(clientPid),
		}
		c.sendChan <- c.encodeMessage(localCtx, msgData)
	}(roomInfo.ID, c.userID, c.pid)

}

// 处理开始游戏
func (c *Client) handlePlay(ctx context.Context) {
	rmMu.RLock()
	player, ok := rm.PlayerList[c.userID]
	if !ok {
		rmMu.RUnlock()
		g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
		return
	}
	roomInfo, ok := rm.Rooms[player.RoomID]
	rmMu.RUnlock()

	if !ok {
		g.Log().Errorf(ctx, "用户 %s 房间不存在", c.userID)
		return
	}

	rmMu.Lock()
	roomInfo.IsPlaying = true
	rmMu.Unlock()

	// 启动游戏逻辑（传入上下文，支持取消）
	room.PlayOneGame(roomInfo)
}

// 处理出牌
func (c *Client) handlePlayCard(ctx context.Context, data string) {
	rmMu.RLock()
	player, ok := rm.PlayerList[c.userID]
	rmMu.RUnlock()
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
		return
	}

	// 解析出牌数据（严格错误处理）
	var reqData struct {
		Pid     int   `json:"pid"`
		CardIds []int `json:"cardIds"`
		Pass    int   `json:"pass"`
	}
	if err := json.Unmarshal([]byte(data), &reqData); err != nil {
		g.Log().Errorf(ctx, "用户 %s 解析出牌数据失败: %v", c.userID, err)
		return
	}
	if reqData.Pid != c.pid {
		g.Log().Errorf(ctx, "用户 %s PID不匹配（%d vs %d）", c.userID, reqData.Pid, c.pid)
		return
	}

	// 更新玩家出牌信息
	rmMu.Lock()
	player.OutCardIds = reqData.CardIds
	player.Pass = reqData.Pass
	rmMu.Unlock()
}

// 处理获取房间信息
func (c *Client) handleGetInfo(ctx context.Context) {
	rmMu.RLock()
	player, ok := rm.PlayerList[c.userID]
	if !ok {
		rmMu.RUnlock()
		g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
		return
	}
	roomInfo, ok := rm.Rooms[player.RoomID]
	rmMu.RUnlock()
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 房间不存在", c.userID)
		return
	}
	// 发送房间ID到客户端
	c.roomIdChan <- roomInfo.ID
	// 整理房间信息（仅暴露当前玩家可见数据）
	cards := make([]int, 0)
	cardsNum := make([]*message.PlayData, 0)
	playerPoint := int64(0)
	rmMu.RLock()
	for _, p := range roomInfo.Players {
		if p.ID == c.pid {
			playerPoint = p.Point
			for _, card := range p.Cards {
				cards = append(cards, card.Id)
			}
		} else {
			cardsNum = append(cardsNum, &message.PlayData{
				Id:      p.ID,
				CardNum: p.CardNum,
			})
		}
	}
	current := roomInfo.Current
	isPlaying := roomInfo.IsPlaying
	outStarTime := roomInfo.OutStarTime
	rmMu.RUnlock()

	// 计算剩余出牌时间
	outCardTimeout := room.GetOutCardTimeout()
	remainOutCardTimeout := outCardTimeout
	if outStarTime > 0 {
		now := int(time.Now().Unix())
		remainOutCardTimeout = outCardTimeout - (now - outStarTime)
		if remainOutCardTimeout < 0 {
			remainOutCardTimeout = 0
		}
	}

	//上一次出牌
	outCards := make([]int, 0)
	for _, card := range roomInfo.LastCards {
		outCards = append(outCards, card.Id)
	}
	var mustPid int
	for _, v := range roomInfo.Players {
		if v.Must {
			mustPid = v.ID //必须出牌的玩家
		}
	}
	//上一位出牌玩家
	lastPid := (roomInfo.Current - 1 + 4) % 4

	// 序列化并推送（处理JSON错误）
	resData, err := json.Marshal(struct {
		RoomId               string              `json:"roomId"`
		Players              []*room.Player      `json:"players"`
		Cards                []int               `json:"cards"`
		Current              int                 `json:"current"`
		PlayerPoint          int64               `json:"playerPoint"`
		OutCardTimeout       int                 `json:"outCardTimeout"`
		RemainOutCardTimeout int                 `json:"remainOutCardTimeout"`
		CardsNum             []*message.PlayData `json:"cardsNum"`
		IsPlaying            bool                `json:"isPlaying"`
		OutCards             []int               `json:"outCards"`
		MustPid              int                 `json:"mustPid"`
		LastPid              int                 `json:"lastPid"`
	}{
		RoomId:               roomInfo.ID,
		Players:              roomInfo.Players,
		Cards:                cards,
		Current:              current,
		PlayerPoint:          playerPoint,
		OutCardTimeout:       outCardTimeout,
		RemainOutCardTimeout: remainOutCardTimeout,
		CardsNum:             cardsNum,
		IsPlaying:            isPlaying,
		OutCards:             outCards,
		MustPid:              mustPid,
		LastPid:              lastPid,
	})
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化房间信息失败: %v", c.userID, err)
		return
	}

	msgData := message.ChatMsg{
		Type: consts.GetInfo,
		Data: gconv.String(resData),
		From: gconv.String(c.pid),
	}
	c.sendChan <- c.encodeMessage(ctx, msgData)
}

// 处理心跳
func (c *Client) handleHeartbeat(ctx context.Context, data string) {
	if data == "ping" {
		c.heartbeat = time.Now()
		// 非阻塞发送pong，避免通道满导致阻塞
		select {
		case c.sendChan <- []byte(`{"type":"heartbeat","data":"pong"}`):
		default:
			g.Log().Warningf(ctx, "用户 %s 心跳响应发送阻塞（通道满）", c.userID)
		}
	}
}

// 读取服务端消息并推送给客户端（修复空指针和资源泄漏）
func (c *Client) readServe(ctx context.Context) {

	var roomID string
	roomInfo := &room.Room{}
	// 循环读取房间消息（带退出机制）
	for {
		if c.closed {
			g.Log().Warningf(ctx, "用户 %s 房间1 已不存在，退出readServe", c.userID)
			return
		}

		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s readServe退出（上下文关闭）", c.userID)
			return
		case roomID = <-c.roomIdChan: //监听房间Id变化
			rmMu.RLock()
			roomInfo = rm.Rooms[roomID]
			rmMu.RUnlock()
		case roomMsg, ok := <-roomInfo.MsgChan:
			if !ok {
				g.Log().Infof(ctx, "用户 %s 房间消息通道已关闭", c.userID)
				return
			}
			// 推送消息给客户端
			msgData := message.ChatMsg{
				Type: roomMsg.Type,
				Data: roomMsg.Data,
				From: gconv.String(c.pid),
			}
			g.Log().Infof(ctx, "房间 %s 向用户 %s 发送消息: %s", roomID, c.userID, msgData)
			select {
			case c.sendChan <- c.encodeMessage(ctx, msgData):
			default:
				g.Log().Warningf(ctx, "用户 %s 消息发送阻塞（通道满）", c.userID)
			}

			// 处理玩家余额不足（带锁）
			rmMu.Lock()
			for _, p := range roomInfo.Players {
				if p.Type == room.Human && p.Point <= 0 {
					rm.LeaveRoom(p)
					g.Log().Infof(ctx, "玩家 %d 余额不足，已离开房间", p.ID)
				}
			}
			rmMu.Unlock()
		}
	}
}

// 写消息循环：发送消息到客户端
func (c *Client) writeLoop(ctx context.Context) {
	defer g.Log().Infof(ctx, "用户 %s writeLoop退出", c.userID)
	g.Log().Infof(ctx, "用户 %s 收到用户消息: ", c.userID)
	for {
		var roomID string
		player, ok := rm.PlayerList[c.userID]
		if ok && player.RoomID != "" {
			roomID = player.RoomID
		}
		if c.closed {
			g.Log().Infof(ctx, "用户 %s writeLoop退出1", c.userID)
			return
		}
		//先执行到这里等待sendChan有数据
		select {
		// 使用多返回值语法检测通道关闭
		case data, ok := <-c.sendChan:
			// 如果通道已关闭，退出循环
			if !ok {
				g.Log().Infof(ctx, "用户 %s sendChan已关闭，writeLoop退出", c.userID)
				return
			}

			g.Log().Infof(ctx, " writeLoop 房间 %s 向用户 %s 发送消息: %s", roomID, c.userID, string(data))
			if err := c.conn.WriteMessage(ghttp.WsMsgText, data); err != nil {
				g.Log().Errorf(ctx, " writeLoop 房间 %s 用户 %s 消息发送失败: %v,消息内容: %s", roomID, c.userID, err, string(data))
				return
			}
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s writeLoop退出2 房间 %s", c.userID, roomID)
			return
		}
	}
}

// 心跳检测：超时断开连接
func (c *Client) heartbeatCheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer g.Log().Infof(ctx, "用户 %s 心跳检测退出", c.userID)
	fmt.Println("heartbeatCheck:", c.heartbeat)
	for {
		if c.closed {
			return
		}
		select {
		case <-ticker.C:
			g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
			if time.Since(c.heartbeat) > 60*time.Second {
				g.Log().Infof(ctx, "用户 %s 心跳超时（60秒），断开连接", c.userID)
				c.conn.Close()
				return
			}
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s 心跳检测被取消: %v", c.userID, ctx.Err())
			return
		}
	}
}

// 消息编码（处理错误）
func (c *Client) encodeMessage(ctx context.Context, msg message.ChatMsg) []byte {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", c.userID, err)
		return []byte(`{"type":"error","data":"消息编码失败"}`)
	}
	return msgBytes
}

/*
{"type":"initRoom","data":"","name":""}
*/
