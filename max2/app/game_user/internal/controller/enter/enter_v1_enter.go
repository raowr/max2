package enter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/util/grand"

	v1 "game_user/api/enter/v1"
	log_gamev1 "game_user/api/log_game/v1"
	"game_user/internal/consts"
	"game_user/internal/controller/log_game"
	"game_user/internal/controller/room"
	"game_user/internal/message"
	"game_user/internal/service"

	"github.com/go-faker/faker/v4"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g" // 修正：GFv2 正确导入路径
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
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

	// 获取或生成用户ID
	userName := r.GetQuery("user_name").String()
	token := r.GetQuery("token").String()

	//通过username从缓存获取token
	// userToken := service.Cache().Get(ctx, userName).String()
	userInfo := service.Cache().Get(ctx, "user:"+userName).Bytes()
	entUser := room.Users{}
	if userInfo == nil {
		g.Log().Errorf(ctx, "用户不存在")
		return
	}
	//缓存有，验证密码正确
	if err = json.Unmarshal(userInfo, &entUser); err != nil {
		g.Log().Errorf(ctx, "用户信息json错误")
		return
	}
	if entUser.Token != token {
		g.Log().Errorf(ctx, "token不存在")
		ws.Close()
		return
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.JwtKey), nil
	})
	if err != nil || !jwtToken.Valid {
		g.Log().Errorf(ctx, "token验证失败: %v", err)
		ws.Close()
		return
	}

	// 升级HTTP连接为WebSocket
	ws, err = wsUpGrader.Upgrade(r.Response.Writer, r.Request, nil)
	if err != nil {
		g.Log().Errorf(ctx, "WebSocket升级失败: %v", err)
		r.Response.WriteHeader(http.StatusInternalServerError)
		r.Response.Write([]byte("连接建立失败"))
		return
	}

	// 创建连接专用上下文
	connCtx, cancel := context.WithCancel(context.Background())

	// 创建客户端实例（增大消息缓冲）
	client := &Client{
		conn:      ws,
		userName:  userName, // 初始化关闭标记为0（未关闭）
		cancel:    cancel,
		subClient: g.Redis(),
		pubClient: g.Redis(),
	}

	// 加锁更新客户端连接（重连逻辑）
	clientsMu.Lock()
	if oldClient, err := getClient(userName); err == nil && oldClient != nil {
		if oldClient.cancel != nil {
			oldClient.cancel() // 触发旧连接的 <-ctx.Done()
		}
		if oldClient.conn != nil {
			oldClient.conn.Close() // 关闭旧连接
		}

		g.Log().Infof(ctx, "用户 %s 重连，关闭旧连接", userName)
	}

	if err = addClient(client); err != nil {
		g.Log().Errorf(ctx, "添加客户端到缓存失败: %v", err)
		ws.Close()
		return
	}
	clientsMu.Unlock()

	g.Log().Infof(ctx, "用户 %s 连接成功", userName)

	// 启动协程（增加退出日志）
	go client.readLoop(connCtx)
	go client.writeLoop(connCtx)
	go client.heartbeatCheck(connCtx)

	// 返回空响应
	return &v1.EnterRes{}, nil
}

// 读消息循环：处理客户端消息
func (c *Client) readLoop(ctx context.Context) {
	defer c.closeConnection("读循环退出")

	for {
		// 监听 ctx 是否已关闭 (调用 Close 方法时触发)
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "Casbin 停止监听")
			return
		default:
			// 继续执行后续功能
		}
		// 读取客户端消息
		mt, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				g.Log().Errorf(ctx, "用户 %s 消息读取错误: %v", c.userName, err)
			}
			g.Log().Infof(ctx, "用户 %s 连接已关闭，读循环退出", c.userName)
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
				From: c.userName,
			}
			c.safeSendMessage(ctx, errMsg)
			continue
		}
		g.Log().Infof(ctx, "用户 %s 接收消息: %s（类型: %d）", c.userName, data, mt)
		// 根据消息类型处理业务（所有操作加锁保护全局rm）
		switch msg.Type {
		case consts.InitRoom:
			c.handleInitRoom(ctx)
		case consts.CreateRoom:
			c.handleCreateRoom(ctx)
		case consts.JoinRoom:
			c.handleJoinRoom(ctx, msg.Data)
		case consts.LeaveRoom:
			c.handleLeaveRoom(ctx)
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
				From: c.userName,
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
	humanPlayer := roomInfo.CreatePlayer(c.userName, room.Human)
	humanPlayer.UserName = c.userName
	c.pid = humanPlayer.ID
	rm.PlayerList[humanPlayer.UserName] = humanPlayer // 关联用户与玩家
	//缓存当前玩家信息
	playerJsonStr, err := json.Marshal(humanPlayer)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.PlayerInfoPrefix+humanPlayer.UserName, playerJsonStr, 0)

	//发送日志 创建房间
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     roomInfo.ID,
		Type:       int32(roomInfo.Type),
		Status:     int32(roomInfo.Status),
		UserID:     humanPlayer.UserName,
		Point:      humanPlayer.Point, //积分
		Action:     consts.InitRoom,   //行为
		Remain:     "",                //剩余牌
		OutCardIds: "",                //玩家单次打出的牌id
		Text:       "",                //完整信息
	})

	// 创建AI玩家，随机ai人数
	aiCount := grand.N(0, 3)
	for i := 0; i < aiCount; i++ {
		// aiName := fake.Name()
		aiName := faker.ChineseName()
		roomInfo.CreatePlayer(aiName, room.AI)
	}

	//缓存当前房间信息
	roomJsonStr, err := json.Marshal(roomInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

	playerDTOs := getPlayers(roomInfo)

	players, err := json.Marshal(playerDTOs)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userName, err)
		return
	}
	msgData := message.ChatMsg{
		Type: consts.InitRoom,
		Data: gconv.String(players),
		From: c.userName,
	}
	c.safeSendMessage(ctx, msgData)
	if aiCount >= 3 {
		return
	}
	// 延迟添加额外AI（带上下文超时，避免协程泄漏）
	go func(roomInfo *room.Room, userName string, clientPid int) {
		// 使用Background上下文确保协程能执行完
		localCtx := context.Background()
		// 添加日志记录协程启动
		g.Log().Infof(localCtx, "用户 %s 开始延迟添加AI到房间 %s", userName, roomInfo.ID)

		<-time.After(2 * time.Second)
		rmMu.Lock()
		defer rmMu.Unlock()
		aiNum := len(roomInfo.Players)
		for i := 0; i < 4-aiNum; i++ {
			// aiName := fake.Name()
			aiName := faker.ChineseName()
			roomInfo.CreatePlayer(aiName, room.AI)
		}

		//缓存当前房间信息
		roomJsonStr, err := json.Marshal(roomInfo)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		service.Cache().Set(ctx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

		// 创建一个临时结构体，只包含可序列化的字段
		playerDTOs := getPlayers(roomInfo)

		players, err := json.Marshal(playerDTOs)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userName, err)
			return
		}
		msgData := message.ChatMsg{
			Type: consts.InitRoom,
			Data: gconv.String(players),
			From: c.userName,
		}
		c.safeSendMessage(ctx, msgData)
	}(roomInfo, c.userName, c.pid)

}

// 处理创建房间
func (c *Client) handleCreateRoom(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()
	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）
	c.clearRoom(ctx)

	roomInfo := rm.CreateRoom(2) //创建好友房
	// c.safeSendToRoomChan(ctx, roomInfo)

	playerName := c.userName
	humanPlayer := roomInfo.CreatePlayer(playerName, room.Human)
	humanPlayer.UserName = c.userName
	c.pid = humanPlayer.ID
	rm.PlayerList[humanPlayer.UserName] = humanPlayer // 关联用户与玩家

	//缓存当前玩家信息
	jsonStr, err := json.Marshal(humanPlayer)
	if err != nil {
		g.Log().Error(ctx, err)
	}

	service.Cache().Set(ctx, consts.PlayerInfoPrefix+humanPlayer.UserName, jsonStr, 0)

	//缓存当前房间信息
	roomJsonStr, err := json.Marshal(roomInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

	// 推送玩家列表（处理JSON错误）
	playerDTOs := getPlayers(roomInfo)

	players, err := json.Marshal(playerDTOs)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", c.userName, err)
		return
	}
	msgData := message.ChatMsg{
		Type: consts.CreateRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(players),
		From: c.userName,
	}
	c.safeSendMessage(ctx, msgData)

	//发送日志 创建房间
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     roomInfo.ID,
		Type:       int32(roomInfo.Type),
		Status:     int32(roomInfo.Status),
		UserID:     humanPlayer.UserName,
		Point:      humanPlayer.Point, //积分
		Action:     consts.InitRoom,   //行为
		Remain:     "",                //剩余牌
		OutCardIds: "",                //玩家单次打出的牌id
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
		g.Log().Errorf(ctx, "用户 %s 解析加入房间数据失败: %v", c.userName, err)
		return
	}
	roomID := reqData.RoomID
	if roomID == "" {
		g.Log().Errorf(ctx, "用户 %s 加入房间ID为空", c.userName)
		return
	}
	// roomInfo, ok := rm.Rooms[roomID]
	roomInfoBytes := service.Cache().Get(ctx, consts.RoomInfoPrefix+roomID).Bytes()
	if roomInfoBytes == nil {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 不存在", c.userName, roomID)
		return
	}
	var roomInfo room.Room
	if err = json.Unmarshal(roomInfoBytes, &roomInfo); err != nil {
		g.Log().Errorf(ctx, "房间信息json错误: %v", err)
		return
	}
	// 检查房间是否已满
	if len(roomInfo.Players) >= 4 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已满", c.userName, roomID)
		return
	}
	// 检查房间是否正在游戏中
	if roomInfo.IsPlaying {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 正在游戏中", c.userName, roomID)
		return
	}
	// 检查房间是否已开始游戏
	if roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已开始游戏", c.userName, roomID)
		return
	}
	//判断玩家已在房间
	inRoom := false
	var existingPlayer *room.Player // 用于存储已存在的玩家
	for _, player := range roomInfo.Players {
		if c.userName == player.UserName {
			g.Log().Errorf(ctx, "用户 %s 已在房间 %s", c.userName, roomID)
			inRoom = true
			existingPlayer = player // 记录已存在的玩家
			break
		}
	}
	//不在房间再加入房间，在房间直接返回房间数据
	var humanPlayer *room.Player
	if !inRoom {
		playerName := c.userName
		humanPlayer = roomInfo.CreatePlayer(playerName, room.Human)
		humanPlayer.UserName = c.userName
		c.pid = humanPlayer.ID
		humanPlayer.Name = c.userName
		rm.PlayerList[humanPlayer.UserName] = humanPlayer // 关联用户与玩家
		// rm.JoinRoom(humanPlayer, roomID)
	} else {
		humanPlayer = existingPlayer
	}

	//记录玩家所在的房间
	jsonStr, err := json.Marshal(humanPlayer)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.PlayerInfoPrefix+humanPlayer.UserName, jsonStr, 0)

	//缓存当前房间信息
	roomJsonStr, err := json.Marshal(&roomInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

	playerDTOs := getPlayers(&roomInfo)

	//加入成功通知所有人
	msgData := message.ChatMsg{
		Type: consts.JoinRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(playerDTOs),
		From: c.userName,
	}
	// c.safeSendMessage(ctx, msgData)
	for _, player := range roomInfo.Players {
		msgData.From = player.UserName
		_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", c.userName, err)
		}
	}
}

// 离开房间
func (c *Client) handleLeaveRoom(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx)
	if playerInfo == nil || roomInfo == nil {
		return
	}
	players := make([]*room.Player, 0) //临时玩家数组
	var isHomeowner bool
	for _, player := range roomInfo.Players {
		if player.UserName != playerInfo.UserName {
			players = append(players, player)
		}
		//如果是房主离开,重新牌ID
		if player.UserName == playerInfo.UserName {
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
		//移除房间信息
		delete(rm.Rooms, playerInfo.RoomID)
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

	//移除rm中用户列表
	delete(rm.PlayerList, playerInfo.UserName)
	// service.Cache().Remove(ctx, consts.PlayerInfoPrefix+playerInfo.UserName)

	//缓存当前房间信息
	roomJsonStr, err := json.Marshal(roomInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

	// 创建一个临时结构体，只包含可序列化的字段
	playerDTOs := getPlayers(roomInfo)
	//如果时房主离开，还需要发房间信息到新房主

	//应该通知到房间内每个人
	// roomInfo.SendRoomMessage(consts.JoinRoom, playerDTOs)
	// c.safeSendMessage(ctx, msgData)
	msgData := message.ChatMsg{
		Type: consts.LeaveRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(playerDTOs),
		From: c.userName,
	}
	for _, player := range roomInfo.Players {
		msgData.From = player.UserName
		_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", c.userName, err)
		}
	}
	//房间没人清理房间,或者是段位房
	c.clearRoom(ctx)

	//发送日志 离开房间
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     roomInfo.ID,
		Type:       int32(roomInfo.Type),
		Status:     int32(roomInfo.Status),
		UserID:     playerInfo.UserName,
		Point:      playerInfo.Point, //积分
		Action:     consts.LeaveRoom, //行为
		Remain:     "",               //剩余牌
		OutCardIds: "",               //玩家单次打出的牌id
		Text:       "",               //完整信息
	})

}

// 处理开始游戏
func (c *Client) handlePlay(ctx context.Context) {
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx)
	if playerInfo == nil || roomInfo == nil {
		return
	}
	//判断是否为房主
	if playerInfo.ID != 0 {
		g.Log().Errorf(ctx, "用户 %s 不是房主", c.userName)
		return
	}

	// 启动游戏逻辑（传入上下文，支持取消）
	//判断人数是否足够
	if len(roomInfo.Players) < 4 {
		g.Log().Errorf(ctx, "用户 %s 房间人数不足", c.userName)
		return
	}

	//判断游戏进行中不能再次play,避免同一房间多次play
	if roomInfo.IsPlaying || roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 进行中不能再次开始比赛", c.userName, roomInfo.ID)
		return
	}

	rmMu.Lock()
	roomInfo.IsPlaying = true
	roomInfo.Status = 1 //游戏中
	rmMu.Unlock()
	//如果时好友房,通知房间内每个人开始游戏
	if roomInfo.Type == 2 {
		// roomInfo.SendRoomMessage(consts.Play, "")
		msgData := message.ChatMsg{
			Type: consts.Play, //此处有修改注意，客户端跳转
			Data: "",
			From: c.userName,
		}
		for _, player := range roomInfo.Players {
			msgData.From = player.UserName
			_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
			if err != nil {
				g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", c.userName, err)
			}
		}
	}

	//初始化房间定时器
	roomInfo.Rgtimer = gtimer.New()
	roomInfo.Rgtimer.Add(context.Background(), 1*time.Second, roomInfo.GameLoop)
	roomInfo.Rgtimer.Stop()
	roomInfo.MsgQueue = make(chan *message.ChatMsg, 10)

	room.PlayOneGame(roomInfo)

	//发送日志 开始游戏
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     roomInfo.ID,
		Type:       int32(roomInfo.Type),
		Status:     int32(roomInfo.Status),
		UserID:     c.userName,
		Point:      playerInfo.Point, //积分
		Action:     consts.Play,      //行为
		Remain:     "",               //剩余牌
		OutCardIds: "",               //玩家单次打出的牌id
		Text:       "",               //完整信息
	})
}

// 处理出牌
func (c *Client) handlePlayCard(ctx context.Context, data string) {
	rmMu.RLock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx)
	if playerInfo == nil || roomInfo == nil {
		return
	}
	rmMu.RUnlock()
	// 解析出牌数据（严格错误处理）
	var reqData struct {
		Pid     int   `json:"pid"`
		CardIds []int `json:"cardIds"`
		Pass    int   `json:"pass"`
	}
	if err := json.Unmarshal([]byte(data), &reqData); err != nil {
		g.Log().Errorf(ctx, "用户 %s 解析出牌数据失败: %v", c.userName, err)
		return
	}
	if reqData.Pid != c.pid {
		g.Log().Errorf(ctx, "用户 %s PID不匹配（%d vs %d）", c.userName, reqData.Pid, c.pid)
		return
	}

	//如果不是人类出牌时间，返回
	if playerInfo.ID != roomInfo.Current {
		g.Log().Errorf(ctx, "用户 %s 不是该玩家出牌阶段", c.userName)
		return
	}

	// 更新玩家出牌信息
	// rmMu.Lock()
	// player.OutCardIds = reqData.CardIds
	// player.Pass = reqData.Pass
	// rmMu.Unlock()

	//发布出牌
	msgData := message.ChatMsg{
		Type: consts.PlayCard, //此处有修改注意，客户端跳转
		Data: gconv.String(reqData),
		From: c.userName,
	}
	_, err := c.pubClient.Publish(ctx, consts.PlayerPlayCardPrefix+playerInfo.UserName, c.encodeMessage(ctx, msgData))
	g.Log().Infof(ctx, "用户 %s 发送出牌消息: %s", c.userName, gconv.String(reqData))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送出牌消息失败: %v", c.userName, err)
	}
}

// 处理获取房间信息
func (c *Client) handleGetInfo(ctx context.Context) {
	rmMu.Lock()
	defer rmMu.Unlock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx)
	if playerInfo == nil || roomInfo == nil {
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
	// if isPlaying {
	// 	c.safeSendToRoomChan(ctx, roomInfo)
	// }

	// 4. 玩家卡片数据（通常比较稳定，可以不加锁）
	cards := make([]int, 0)
	cardsNum := make([]*message.PlayData, 0)
	var playerPoint int64
	var mustPid int

	for _, p := range roomInfo.Players {
		if p.ID == playerInfo.ID {
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
		g.Log().Errorf(ctx, "用户 %s 序列化房间信息失败: %v", c.userName, err)
		return
	}

	// 7. 最后在锁外发送消息
	msgData := message.ChatMsg{
		Type: consts.GetInfo,
		Data: gconv.String(resData),
		From: c.userName,
	}
	_, err = c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+playerInfo.UserName, c.encodeMessage(ctx, msgData))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", playerInfo.UserName, err)
	}
}

// 处理心跳
func (c *Client) handleHeartbeat(ctx context.Context, data string) {
	if data == "ping" {
		msgData := message.ChatMsg{
			Type: consts.Heartbeat,
			Data: "pong",
			From: c.userName,
		}
		c.safeSendMessage(ctx, msgData)
	}
}

// 写消息循环：发送消息到客户端
func (c *Client) writeLoop(ctx context.Context) {
	defer g.Log().Infof(ctx, "用户 %s writeLoop退出", c.userName)
	sub, _, err := c.subClient.Subscribe(ctx, consts.PlayerMsgPrefix+c.userName)
	if err != nil {
		g.Log().Error(ctx, "writeLoop 订阅失败:", err)
		return
	}
	// 确保在函数退出时关闭订阅
	defer func() {
		_ = sub.Close(ctx)
	}()

	// 循环接收消息
	for {
		// 监听 ctx 是否已关闭 (调用 Close 方法时触发)
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "writeLoop 停止监听")
			return
		default:
			// 继续执行后续功能
		}

		// Receive 接收消息
		msg, err := sub.Receive(ctx)
		if err != nil {
			// 如果是 Context 被取消导致的错误，直接退出
			if ctx.Err() != nil {
				return
			}
			// 其它错误（如网络断开），可以考虑简单的重试或记录日志
			g.Log().Error(ctx, "Casbin Watcher 接收消息错误:", err)
			// 如果出错直接退出，等待下一次重启或手动干预
			return
		}
		msgData, success := room.ParseRedisSubscribeMessage(msg.String(), c.userName, ctx)
		if !success {
			continue
		}

		// 检查连接是否已关闭
		c.mutex.RLock()
		conn := c.conn
		c.mutex.RUnlock()
		if conn == nil {
			g.Log().Infof(ctx, "用户 %s 连接已关闭，跳过消息发送", c.userName)
			return
		}
		msgDataBytes, err := gjson.Encode(msgData)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", c.userName, err)
		}
		g.Log().Infof(ctx, "用户 %s 发送消息: %s", c.userName, string(msgDataBytes))
		if err := c.conn.WriteMessage(ghttp.WsMsgText, msgDataBytes); err != nil {
			g.Log().Errorf(ctx, " writeLoop 用户 %s 消息发送失败: %v,消息内容: %s", c.userName, err, string(msgDataBytes))
			return
		}
	}
}

// 心跳检测：超时断开连接
func (c *Client) heartbeatCheck(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	defer g.Log().Infof(ctx, "用户 %s 心跳检测退出", c.userName)
	g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
	for {
		select {
		case <-ticker.C:
			g.Log().Infof(ctx, "heartbeatCheck: %v", time.Since(c.heartbeat))
			if time.Since(c.heartbeat) > 60*time.Second {
				g.Log().Infof(ctx, "用户 %s 心跳超时（60秒），断开连接", c.userName)
				//c.conn.Close()
				c.closeConnection("心跳超时")
				return
			}
		case <-ctx.Done():
			g.Log().Infof(ctx, "用户 %s 心跳检测被取消: %v", c.userName, ctx.Err())
			return
		}
	}
}

// 安全发送消息的函数
func (c *Client) safeSendMessage(ctx context.Context, msg message.ChatMsg) {
	_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+c.userName, c.encodeMessage(ctx, msg))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", c.userName, err)
	}
}

// 消息编码（处理错误）
func (c *Client) encodeMessage(ctx context.Context, msg message.ChatMsg) []byte {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", c.userName, err)
		return []byte(`{"type":"error","data":"消息编码失败"}`)
	}
	return msgBytes
}

// 关闭客户端连接
func (c *Client) closeConnection(reason string) {

	// 原子操作检查是否已关闭
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		g.Log().Infof(context.Background(), "用户 %s 连接已关闭，跳过重复关闭", c.userName)
		return
	}

	rmMu.Lock()
	defer rmMu.Unlock()
	logCtx := context.Background()
	g.Log().Infof(logCtx, "用户 %s 关闭连接，原因: %s", c.userName, reason)

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

	// 最后清理资源
	clientsMu.Lock()
	err := removeClient(c.userName)
	if err != nil {
		g.Log().Infof(logCtx, "用户 %s 删除客户端失败", c.userName)
	}
	clientsMu.Unlock()
}

// 清理旧房间
func (c *Client) clearRoom(ctx context.Context) {
	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）

	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx)
	if playerInfo == nil || roomInfo == nil {
		return
	}

	if roomInfo.Type == 1 || len(roomInfo.Players) <= 0 {
		service.Cache().Remove(ctx, consts.RoomInfoPrefix+playerInfo.RoomID)
	}
	service.Cache().Remove(ctx, consts.PlayerInfoPrefix+c.userName)

}

func (c *Client) getPlayerAndRoomInfo(ctx context.Context) (*room.Player, *room.Room) {
	playerInfo := &room.Player{} // 或 new(room.Player)
	jsonStr := service.Cache().Get(ctx, consts.PlayerInfoPrefix+c.userName).Bytes()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := json.Unmarshal(jsonStr, playerInfo); err != nil {
		g.Log().Error(ctx, err)
	}

	roomInfoBytes := service.Cache().Get(ctx, consts.RoomInfoPrefix+playerInfo.RoomID).Bytes()
	if roomInfoBytes == nil {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 不存在", c.userName, playerInfo.RoomID)
		return nil, nil
	}
	var roomInfo room.Room
	if err = json.Unmarshal(roomInfoBytes, &roomInfo); err != nil {
		g.Log().Errorf(ctx, "房间信息json错误: %v", err)
		return nil, nil
	}
	return playerInfo, &roomInfo
}

/*
{"type":"initRoom","data":"","name":""}
*/
