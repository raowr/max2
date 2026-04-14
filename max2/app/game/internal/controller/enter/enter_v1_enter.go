package enter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/util/grand"

	v1 "game/api/enter/v1"
	"game/internal/consts"
	"game/internal/controller/log"
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

	// 创建连接专用上下文
	connCtx, cancel := context.WithCancel(context.Background())

	// 创建客户端实例（增大消息缓冲）
	client := &Client{
		conn:      ws,
		userID:    userID,
		heartbeat: time.Now(),
		sendChan:  make(chan []byte, 10000),   // 缓冲增大至10000，减少阻塞
		roomChan:  make(chan *room.Room, 100), //新建房间通道
		cancel:    cancel,
		mutex:     sync.RWMutex{}, // 显式初始化互斥锁
		closed:    0,              // 初始化关闭标记为0（未关闭）
	}

	// 加锁更新客户端连接（重连逻辑）
	clientsMu.Lock()
	if oldClient, err := getClient(userID); err == nil && oldClient != nil {
		if oldClient.cancel != nil {
			oldClient.cancel() // 触发旧连接的 <-ctx.Done()
		}
		if oldClient.conn != nil {
			oldClient.conn.Close() // 关闭旧连接
		}
		// oldClient.conn.Close() // 关闭旧连接

		oldClient.mutex.Lock()
		// 关闭sendChan
		if oldClient.sendChan != nil {
			close(oldClient.sendChan)
			oldClient.sendChan = nil
		}
		// 关闭roomChan
		if oldClient.roomChan != nil {
			close(oldClient.roomChan)
			oldClient.roomChan = nil
		}
		oldClient.mutex.Unlock()

		g.Log().Infof(ctx, "用户 %s 重连，关闭旧连接", userID)
	}

	if err = addClient(client); err != nil {
		clientsMu.Unlock()
		g.Log().Errorf(ctx, "添加客户端到缓存失败: %v", err)
		ws.Close()
		return
	}
	clientsMu.Unlock()

	g.Log().Infof(ctx, "用户 %s 连接成功", userID)

	// 启动协程（增加退出日志）
	go client.readLoop(connCtx)
	go client.readServe(connCtx)
	go client.writeLoop(connCtx)
	go client.heartbeatCheck(connCtx)

	// 返回空响应
	return &v1.EnterRes{}, nil
}

// 读消息循环：处理客户端消息
func (c *Client) readLoop(ctx context.Context) {
	defer c.closeConnection("读循环退出")

	for {
		// 检查连接是否已关闭
		if atomic.LoadInt32(&c.closed) == 1 {
			g.Log().Infof(ctx, "用户 %s 连接已关闭，读循环退出", c.userID)
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
				From: c.userID,
			}
			c.safeSendMessage(ctx, errMsg)
			continue
		}
		g.Log().Infof(ctx, "用户 %s 接收消息: %s（类型: %d）", c.userID, data, mt)
		// 根据消息类型处理业务（所有操作加锁保护全局rm）
		switch msg.Type {
		case consts.InitRoom:
			c.handleInitRoom(ctx)
		case consts.CreateRoom:
			c.handleCreateRoom(ctx)
		case consts.JoinRoom:
			c.handleJoinRoom(ctx, msg.Data)
		case consts.LeaveRoom:
			c.handleLeaveRoom(ctx, msg.Data)
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
				From: c.userID,
			}
			c.safeSendMessage(ctx, errMsg)
		}
	}
}

// 处理房间初始化
func (c *Client) handleInitRoom(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()
	c.clearRoom(ctx)
	// 创建新房间和玩家
	roomInfo := rm.CreateRoom(1) //创建比赛房
	c.safeSendToRoomChan(ctx, roomInfo)
	playerName := "美女"
	humanPlayer := roomInfo.CreatePlayer(playerName, room.Human)
	humanPlayer.UserId = c.userID
	c.pid = humanPlayer.ID
	rm.PlayerList[humanPlayer.UserId] = humanPlayer // 关联用户与玩家

	//发送日志 创建房间
	log.SendLog(log.LogInfo{
		RoomID:     roomInfo.ID,
		Type:       roomInfo.Type,
		Status:     roomInfo.Status,
		UserID:     humanPlayer.UserId,
		Point:      humanPlayer.Point, //积分
		Action:     consts.InitRoom,   //行为
		Remain:     make([]string, 0), //剩余牌
		OutCardIds: make([]int, 0),    //玩家单次打出的牌id
		Text:       "",                //完整信息
	})

	// 创建AI玩家，随机ai人数
	aiCount := grand.N(0, 3)
	for i := 0; i < aiCount; i++ {
		aiName := fmt.Sprintf("帅锅%d号", i+1)
		roomInfo.CreatePlayer(aiName, room.AI)
	}

	playerDTOs := getPlayers(roomInfo)

	players, err := json.Marshal(playerDTOs)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
		return
	}
	msgData := message.ChatMsg{
		Type: consts.InitRoom,
		Data: gconv.String(players),
		From: c.userID,
	}
	c.safeSendMessage(ctx, msgData)
	if aiCount >= 3 {
		return
	}
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
		aiNum := len(roomInfo.Players)
		for i := 0; i < 4-aiNum; i++ {
			aiName := fmt.Sprintf("帅锅%d号", aiNum+i+1)
			roomInfo.CreatePlayer(aiName, room.AI)
		}

		// 创建一个临时结构体，只包含可序列化的字段
		playerDTOs := getPlayers(roomInfo)

		players, err := json.Marshal(playerDTOs)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
			return
		}
		msgData := message.ChatMsg{
			Type: consts.InitRoom,
			Data: gconv.String(players),
			From: c.userID,
		}
		c.safeSendMessage(ctx, msgData)
	}(roomInfo.ID, c.userID, c.pid)

}

// 处理创建房间
func (c *Client) handleCreateRoom(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()
	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）
	c.clearRoom(ctx)

	roomInfo := rm.CreateRoom(2) //创建比赛房
	c.safeSendToRoomChan(ctx, roomInfo)

	playerName := "美女"
	humanPlayer := roomInfo.CreatePlayer(playerName, room.Human)
	humanPlayer.UserId = c.userID
	c.pid = humanPlayer.ID
	rm.PlayerList[humanPlayer.UserId] = humanPlayer // 关联用户与玩家

	// 推送玩家列表（处理JSON错误）
	playerDTOs := getPlayers(roomInfo)

	players, err := json.Marshal(playerDTOs)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
		return
	}
	msgData := message.ChatMsg{
		Type: consts.InitRoom,
		Data: gconv.String(players),
		From: c.userID,
	}
	c.safeSendMessage(ctx, msgData)

	//发送日志 创建房间
	log.SendLog(log.LogInfo{
		RoomID:     roomInfo.ID,
		Type:       roomInfo.Type,
		Status:     roomInfo.Status,
		UserID:     humanPlayer.UserId,
		Point:      humanPlayer.Point, //积分
		Action:     consts.InitRoom,   //行为
		Remain:     make([]string, 0), //剩余牌
		OutCardIds: make([]int, 0),    //玩家单次打出的牌id
		Text:       "",                //完整信息
	})

}

// 处理加入房间
func (c *Client) handleJoinRoom(ctx context.Context, data string) {
	rmMu.Lock()
	defer rmMu.Unlock()
	// 解析加入房间数据（严格错误处理）
	var reqData struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal([]byte(data), &reqData); err != nil {
		g.Log().Errorf(ctx, "用户 %s 解析加入房间数据失败: %v", c.userID, err)
		return
	}
	roomID := reqData.RoomID
	if roomID == "" {
		g.Log().Errorf(ctx, "用户 %s 加入房间ID为空", c.userID)
		return
	}
	roomInfo, ok := rm.Rooms[roomID]
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 不存在", c.userID, roomID)
		return
	}
	// 检查房间是否已满
	if len(roomInfo.Players) >= 4 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已满", c.userID, roomID)
		return
	}
	// 检查房间是否正在游戏中
	if roomInfo.IsPlaying {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 正在游戏中", c.userID, roomID)
		return
	}
	// 检查房间是否已开始游戏
	if roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已开始游戏", c.userID, roomID)
		return
	}
	//判断玩家已在房间
	inRoom := false
	for _, player := range roomInfo.Players {
		if player.UserId == c.userID {
			g.Log().Errorf(ctx, "用户 %s 已在房间 %s", c.userID, roomID)
			inRoom = true
			break
		}
	}
	//不在房间再加入房间，在房间直接返回房间数据
	var humanPlayer *room.Player
	if !inRoom {
		playerName := "好友"
		humanPlayer = roomInfo.CreatePlayer(playerName, room.Human)
		humanPlayer.UserId = c.userID
		c.pid = humanPlayer.ID
		humanPlayer.Name = playerName + gconv.String(humanPlayer.ID)
		rm.PlayerList[humanPlayer.UserId] = humanPlayer // 关联用户与玩家
		// rm.JoinRoom(humanPlayer, roomID)
	}

	c.safeSendToRoomChan(ctx, roomInfo)

	// 推送玩家列表（处理JSON错误）
	// 创建一个临时结构体，只包含可序列化的字段
	playerDTOs := getPlayers(roomInfo)

	//players, err := json.Marshal(playerDTOs)
	//if err != nil {
	//	g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userID, err)
	//	return
	//}
	//msgData := message.ChatMsg{
	//	Type: consts.JoinRoom,
	//	Data: gconv.String(players),
	//	From: gconv.String(humanPlayer.ID),
	//}
	//应该通知到房间内每个人
	roomInfo.SendRoomMessage(consts.JoinRoom, playerDTOs)
	//c.safeSendMessage(ctx, msgData)

	//发送日志 加入房间
	log.SendLog(log.LogInfo{
		RoomID:     roomInfo.ID,
		Type:       roomInfo.Type,
		Status:     roomInfo.Status,
		UserID:     c.userID,
		Point:      humanPlayer.Point, //积分
		Action:     consts.JoinRoom,   //行为
		Remain:     make([]string, 0), //剩余牌
		OutCardIds: make([]int, 0),    //玩家单次打出的牌id
		Text:       "",                //完整信息
	})
}

// 离开房间
func (c *Client) handleLeaveRoom(ctx context.Context, data string) {
	rmMu.Lock()
	defer rmMu.Unlock()
	//移除房间内该玩家
	player, ok := rm.PlayerList[c.userID]
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
		return
	}
	isHomeowner := false
	roomInfo, ok := rm.Rooms[player.RoomID]
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 未找到房间 %s 信息", c.userID, player.RoomID)
		return
	}
	players := make([]*room.Player, 0) //临时玩家数组
	for _, player := range roomInfo.Players {
		if player.UserId != c.userID {
			players = append(players, player)
		}
		//如果是房主离开,重新牌ID
		if player.UserId == c.userID {
			if player.ID == 0 {
				isHomeowner = true
			}
		}

	}
	roomInfo.Players = players
	roomInfo.NextPlayerID = 0
	if isHomeowner {
		if len(roomInfo.Players) > 0 {
			for _, player := range roomInfo.Players {
				player.ID = roomInfo.NextPlayerID
				roomInfo.NextPlayerID++
			}
		}
	} else {
		if len(roomInfo.Players) > 0 {
			roomInfo.NextPlayerID = 1
			for _, player := range roomInfo.Players {
				if player.ID == 0 {
					continue
				} else {
					player.ID = roomInfo.NextPlayerID
					roomInfo.NextPlayerID++
				}
			}
		}
	}

	if player.MsgChan != nil {
		c.mutex.Lock()
		ch := player.MsgChan
		player.MsgChan = nil
		c.mutex.Unlock()
		close(ch)
		g.Log().Infof(context.Background(), "用户 %s 离开房间 %s 后 MsgChan 关闭", c.userID, roomInfo.ID)
	}

	//移除rm中用户列表
	delete(rm.PlayerList, c.userID)

	// 创建一个临时结构体，只包含可序列化的字段
	playerDTOs := getPlayers(roomInfo)

	//应该通知到房间内每个人
	roomInfo.SendRoomMessage(consts.JoinRoom, playerDTOs)

	//发送日志 离开房间
	log.SendLog(log.LogInfo{
		RoomID:     roomInfo.ID,
		Type:       roomInfo.Type,
		Status:     roomInfo.Status,
		UserID:     c.userID,
		Point:      player.Point,      //积分
		Action:     consts.LeaveRoom,  //行为
		Remain:     make([]string, 0), //剩余牌
		OutCardIds: make([]int, 0),    //玩家单次打出的牌id
		Text:       "",                //完整信息
	})

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

	// 启动游戏逻辑（传入上下文，支持取消）
	//判断人数是否足够
	if len(roomInfo.Players) < 4 {
		g.Log().Errorf(ctx, "用户 %s 房间人数不足", c.userID)
		return
	}

	//判断游戏进行中不能再次play,避免同一房间多次play
	if roomInfo.IsPlaying || roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 进行中不能再次开始比赛", c.userID, roomInfo.ID)
		return
	}

	rmMu.Lock()
	roomInfo.IsPlaying = true
	roomInfo.Status = 1 //游戏中
	rmMu.Unlock()
	//如果时好友房,通知房间内每个人开始游戏
	if roomInfo.Type == 2 {
		roomInfo.SendRoomMessage(consts.Play, "")
	}
	room.PlayOneGame(roomInfo)

	//发送日志 开始游戏
	log.SendLog(log.LogInfo{
		RoomID:     roomInfo.ID,
		Type:       roomInfo.Type,
		Status:     roomInfo.Status,
		UserID:     c.userID,
		Point:      player.Point,      //积分
		Action:     consts.Play,       //行为
		Remain:     make([]string, 0), //剩余牌
		OutCardIds: make([]int, 0),    //玩家单次打出的牌id
		Text:       "",                //完整信息
	})
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
	rmMu.RLock()
	roomInfo := rm.Rooms[player.RoomID]
	rmMu.RUnlock()
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

	//如果不是人类出牌时间，返回
	if player.ID != roomInfo.Current {
		g.Log().Errorf(ctx, "用户 %s 不是该玩家出牌阶段", c.userID)
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
	rmMu.Lock()
	defer rmMu.Unlock()
	// 1. 基础数据访问（假设是安全的）
	player, ok := rm.PlayerList[c.userID]
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
		return
	}

	roomInfo, ok := rm.Rooms[player.RoomID]
	if !ok {
		g.Log().Errorf(ctx, "用户 %s 房间不存在", c.userID)
		return
	}

	// 2. 如果涉及到可能被其他goroutine修改的数据，才加锁
	var current, outStarTime int
	var isPlaying bool
	var status int
	var roomType int

	// 快速获取易变数据
	current = roomInfo.Current
	outStarTime = roomInfo.OutStarTime
	isPlaying = roomInfo.IsPlaying
	status = roomInfo.Status
	roomType = roomInfo.Type

	// 3. 在锁外进行其他操作
	if isPlaying {
		c.safeSendToRoomChan(ctx, roomInfo)
	}

	// 4. 玩家卡片数据（通常比较稳定，可以不加锁）
	cards := make([]int, 0)
	cardsNum := make([]*message.PlayData, 0)
	var playerPoint int64
	var mustPid int

	for _, p := range roomInfo.Players {
		if p.ID == player.ID {
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
		if p.Must {
			mustPid = p.ID
		}
	}

	// 5. 计算剩余信息，在锁外进行
	outCardTimeout := room.GetOutCardTimeout()
	remainOutCardTimeout := outCardTimeout
	if outStarTime > 0 {
		now := int(time.Now().Unix())
		remainOutCardTimeout = outCardTimeout - (now - outStarTime)
		if remainOutCardTimeout < 0 {
			remainOutCardTimeout = 0
		}
	}

	// 上一次出牌
	outCards := make([]int, 0)
	for _, card := range roomInfo.LastCards {
		outCards = append(outCards, card.Id)
	}

	// 上一位出牌玩家
	lastPid := (current - 1 + 4) % 4

	// 创建一个临时结构体，只包含可序列化的字段
	playerDTOs := getPlayers(roomInfo)

	// 6. 序列化并推送（处理JSON错误）
	resData, err := json.Marshal(struct {
		RoomId               string              `json:"roomId"`
		Players              []*PlayerDTO        `json:"players"`
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
		Status               int                 `json:"status"`
		Type                 int                 `json:"type"` //房间类型 1比赛房，2好友房
	}{
		RoomId:               roomInfo.ID,
		Players:              playerDTOs,
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
		Status:               status,
		Type:                 roomType,
	})
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化房间信息失败: %v", c.userID, err)
		return
	}

	// 7. 最后在锁外发送消息
	msgData := message.ChatMsg{
		Type: consts.GetInfo,
		Data: gconv.String(resData),
		From: c.userID,
	}
	c.safeSendMessage(ctx, msgData)
}

// 处理心跳
func (c *Client) handleHeartbeat(ctx context.Context, data string) {
	if data == "ping" {
		// 加读锁保护对sendChan的访问
		c.mutex.RLock()
		sendChan := c.sendChan
		c.mutex.RUnlock()
		c.heartbeat = time.Now()
		// 非阻塞发送pong，避免通道满导致阻塞
		if sendChan != nil {
			select {
			case sendChan <- []byte(`{"type":"heartbeat","data":"pong"}`):
			default:
				g.Log().Warningf(ctx, "用户 %s 心跳响应发送阻塞（通道满）", c.userID)
			}
		}

	}
}

// 读取服务端消息并推送给客户端（修复空指针和资源泄漏）
func (c *Client) readServe(ctx context.Context) {
	// 循环读取玩家消息（带退出机制）
	var roomInfo *room.Room
	// 使用局部变量存储当前roomChan的引用
	var currentRoomChan chan *room.Room
	c.mutex.RLock()
	currentRoomChan = c.roomChan
	c.mutex.RUnlock()
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s readServe退出（上下文关闭）房间", c.userID)
			return

		case newRoom := <-currentRoomChan:
			// 切换到新房间前先检查newRoom是否为空
			if newRoom == nil {
				g.Log().Warningf(ctx, "用户 %s 接收到空的房间信息", c.userID)
				continue
			}
			// 切换到新房间
			c.mutex.Lock()
			roomInfo = newRoom
			c.mutex.Unlock()

			// 找到当前客户端对应的玩家
			var currentPlayer *room.Player
			for _, p := range roomInfo.Players {
				if p.UserId == c.userID {
					currentPlayer = p
					break
				}
			}

			// 如果找到了玩家且玩家有消息通道，创建一个goroutine来监听该玩家的消息
			if currentPlayer != nil && currentPlayer.MsgChan != nil {
				go func(player *room.Player) {
					for {
						select {
						case playerMsg, ok := <-player.MsgChan:
							if !ok {
								return
							}

							// 推送玩家消息给客户端
							msgData := message.ChatMsg{
								Type: playerMsg.Type,
								Data: playerMsg.Data,
								From: c.userID,
							}
							g.Log().Infof(ctx, "玩家 %d 向用户 %s 发送消息: %+v", player.ID, c.userID, msgData)

							c.safeSendMessage(ctx, msgData)
						case <-ctx.Done():
							return
						}
					}
				}(currentPlayer)
			}

		default:
			// 短暂休眠避免忙等待
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// 写消息循环：发送消息到客户端
func (c *Client) writeLoop(ctx context.Context) {
	defer g.Log().Infof(ctx, "用户 %s writeLoop退出", c.userID)
	g.Log().Infof(ctx, "用户 %s 收到用户消息: ", c.userID)
	for {
		// 检查连接是否已关闭
		if atomic.LoadInt32(&c.closed) == 1 {
			g.Log().Infof(ctx, "用户 %s 连接已关闭，写循环退出", c.userID)
			return
		}
		var roomID string
		player, ok := rm.PlayerList[c.userID]
		if ok && player.RoomID != "" {
			roomID = player.RoomID
		}

		// 在 writeLoop 函数中修改为
		c.mutex.RLock()
		sendChan := c.sendChan
		c.mutex.RUnlock()
		if sendChan != nil {
			select {
			case data, ok := <-sendChan:
				if !ok {
					g.Log().Infof(ctx, "用户 %s sendChan已关闭，writeLoop退出", c.userID)
					return
				}

				g.Log().Infof(ctx, " writeLoop 房间 %s 向用户 %s 发送消息: %s", roomID, c.userID, string(data))

				// 使用互斥锁保护并发访问
				c.mutex.Lock()
				conn := c.conn
				c.mutex.Unlock()

				if conn == nil {
					g.Log().Errorf(ctx, " writeLoop 房间 %s 用户 %s 连接已关闭，无法发送消息", roomID, c.userID)
					return
				}

				if err := conn.WriteMessage(ghttp.WsMsgText, data); err != nil {
					g.Log().Errorf(ctx, " writeLoop 房间 %s 用户 %s 消息发送失败: %v,消息内容: %s", roomID, c.userID, err, string(data))
					return
				}
			case <-ctx.Done():
				g.Log().Infof(ctx, "用户 %s writeLoop退出2 房间 %s", c.userID, roomID)
				return
			}
		}
	}
}

// 心跳检测：超时断开连接
func (c *Client) heartbeatCheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer g.Log().Infof(ctx, "用户 %s 心跳检测退出", c.userID)
	g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
	for {
		select {
		case <-ticker.C:
			g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
			if time.Since(c.heartbeat) > 60*time.Second {
				g.Log().Infof(ctx, "用户 %s 心跳超时（60秒），断开连接", c.userID)
				//c.conn.Close()
				c.closeConnection("心跳超时")
				return
			}
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s 心跳检测被取消: %v", c.userID, ctx.Err())
			return
		}
	}
}

// 安全发送消息的函数
func (c *Client) safeSendMessage(ctx context.Context, msg message.ChatMsg) {
	c.mutex.RLock()
	sendChan := c.sendChan
	c.mutex.RUnlock()

	if sendChan != nil {
		// 使用匿名函数包装发送操作并添加 recover 机制
		func() {
			defer func() {
				if r := recover(); r != nil {
					g.Log().Errorf(ctx, "用户 %s 消息发送发生panic: %v", c.userID, r)
				}
			}()

			select {
			case sendChan <- c.encodeMessage(ctx, msg):
				// 发送成功
			default:
				g.Log().Warningf(ctx, "用户 %s 消息发送阻塞（通道满）", c.userID)
			}
		}()
	}
}

// 2. 安全地向 roomChan 发送消息的方法
func (c *Client) safeSendToRoomChan(ctx context.Context, roomInfo *room.Room) bool {
	// 加读锁保护对 roomChan 的访问
	c.mutex.RLock()
	roomChan := c.roomChan
	c.mutex.RUnlock()

	// 检查通道是否存在
	if roomChan == nil {
		g.Log().Warningf(ctx, "用户 %s 尝试发送消息到未初始化的 roomChan", c.userID)
		return false
	}

	// 使用匿名函数包装发送操作并添加 recover 机制
	result := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				g.Log().Errorf(ctx, "用户 %s 向 roomChan 发送消息发生panic: %v", c.userID, r)
				result = false
			}
		}()

		// 非阻塞发送或阻塞发送取决于业务需求
		select {
		case roomChan <- roomInfo:
			result = true
		case <-ctx.Done():
			g.Log().Warningf(ctx, "用户 %s 向 roomChan 发送消息超时", c.userID)
			result = false
		default:
			g.Log().Warningf(ctx, "用户 %s 向 roomChan 发送消息阻塞", c.userID)
			result = false
		}
	}()

	return result
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

// 关闭客户端连接
func (c *Client) closeConnection(reason string) {

	// 原子操作检查是否已关闭
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		g.Log().Infof(context.Background(), "用户 %s 连接已关闭，跳过重复关闭", c.userID)
		return
	}

	rmMu.Lock()
	defer rmMu.Unlock()
	logCtx := context.Background()
	g.Log().Infof(logCtx, "用户 %s 关闭连接，原因: %s", c.userID, reason)

	// 先取消上下文
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil // 设置为 nil，避免重复调用
	}

	// 然后关闭连接（防止重复关闭）
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil // 关键：设置为 nil，防止重复关闭
	}

	// 在 closeConnection 函数中修改为（确保互斥锁保护）
	if c.sendChan != nil {
		c.mutex.Lock()
		ch := c.sendChan
		c.sendChan = nil // 先设置为 nil，防止其他 goroutine 再次使用
		c.mutex.Unlock()

		close(ch)
		g.Log().Infof(context.Background(), "用户 %s sendChan 已关闭", c.userID)
	}

	// 关闭 roomChan
	if c.roomChan != nil {
		c.mutex.Lock()
		ch := c.roomChan
		c.roomChan = nil
		c.mutex.Unlock()
		close(ch)
		g.Log().Infof(context.Background(), "用户 %s roomChan 已关闭", c.userID)
	}

	// 最后清理资源
	clientsMu.Lock()
	err := removeClient(c.userID)
	if err != nil {
		g.Log().Infof(logCtx, "用户 %s 删除客户端失败", c.userID)
	}
	clientsMu.Unlock()
}

// 清理旧房间
func (c *Client) clearRoom(ctx context.Context) {
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
		//关闭用户消息通道
		player, ok := rm.PlayerList[c.userID]
		if !ok {
			g.Log().Errorf(ctx, "用户 %s 未找到玩家信息", c.userID)
			return
		}
		if player.MsgChan != nil {
			c.mutex.Lock()
			ch := player.MsgChan
			player.MsgChan = nil
			c.mutex.Unlock()
			close(ch)
			g.Log().Infof(context.Background(), "用户 %s MsgChan 已关闭", c.userID)
		}
		// 从玩家列表移除旧玩家
		delete(rm.PlayerList, c.userID)
	}
}

/*
{"type":"initRoom","data":"","name":""}
*/
