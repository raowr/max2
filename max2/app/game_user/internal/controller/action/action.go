package action

import (
	"context"
	"encoding/json"
	v1 "game_user/api/action/v1"
	"game_user/internal/consts"
	"game_user/internal/controller/log_game"
	"game_user/internal/controller/room"
	"game_user/internal/message"
	"game_user/internal/service"
	"time"

	log_gamev1 "game_user/api/log_game/v1"

	"github.com/go-faker/faker/v4"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
)

type Controller struct {
	v1.UnimplementedActionServer
	subClient *gredis.Redis // 订阅客户端
	pubClient *gredis.Redis // 发布客户端
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterActionServer(s.Server, &Controller{
		subClient: g.Redis(),
		pubClient: g.Redis(),
	})
}

func (c *Controller) SendAction(ctx context.Context, req *v1.SendActionReq) (res *v1.SendActionRes, err error) {
	// 处理请求
	switch req.Type {
	case consts.InitRoom:
		c.handleInitRoom(ctx, req)
	case consts.CreateRoom:
		c.handleCreateRoom(ctx, req)
	case consts.JoinRoom:
		c.handleJoinRoom(ctx, req)
	case consts.LeaveRoom:
		c.handleLeaveRoom(ctx, req)
	case consts.Play:
		c.handlePlay(ctx, req)
	case consts.PlayCard:
		c.handlePlayCard(ctx, req)
	case consts.GetInfo:
		c.handleGetInfo(ctx, req)
	default:
		return nil, gerror.NewCode(gcode.CodeNotImplemented) // ← 只在未匹配时返回错误
	}
	return nil, nil // ← 匹配到 case 时正常返回
}

// 处理房间初始化
func (c *Controller) handleInitRoom(ctx context.Context, req *v1.SendActionReq) {
	rmMu.Lock()
	defer rmMu.Unlock()

	//判断是否可以创建房间
	_, oldRoomInfo := c.getPlayerAndRoomInfo(ctx, req)
	if oldRoomInfo != nil && (oldRoomInfo.Status == 1 || oldRoomInfo.IsPlaying) {
		g.Log().Infof(ctx, "用户 %s 已在房间 %s 进行中，不能创建新房间", req.From, oldRoomInfo.ID)
		return
	}

	c.clearRoom(ctx, req)

	// 创建新房间和玩家
	roomInfo := rm.CreateRoom(1) //创建段位房
	humanPlayer := roomInfo.CreatePlayer(req.From, room.Human)
	humanPlayer.UserName = req.From
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
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", req.From, err)
		return
	}

	msgData := &v1.SendActionReq{
		Type: consts.InitRoom,
		Data: gconv.String(players),
		From: req.From,
	}
	c.safeSendMessage(ctx, msgData)

	if aiCount >= 3 {
		return
	}

	// 延迟添加额外AI（使用独立上下文，不阻塞主请求）
	go func(roomInfo *room.Room, userName string) {
		// 创建独立的上下文，不受原始请求影响
		localCtx, localCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer localCancel()

		g.Log().Infof(localCtx, "用户 %s 开始延迟添加AI到房间 %s", userName, roomInfo.ID)

		// 使用 select 避免阻塞
		select {
		case <-time.After(2 * time.Second):
			// 继续执行
		case <-localCtx.Done():
			g.Log().Warningf(localCtx, "延迟操作超时，房间: %s", roomInfo.ID)
			return
		}

		rmMu.Lock()
		defer rmMu.Unlock()

		aiNum := len(roomInfo.Players)
		for i := 0; i < 4-aiNum; i++ {
			aiName := faker.ChineseName()
			roomInfo.CreatePlayer(aiName, room.AI)
		}

		//缓存当前房间信息（使用 localCtx）
		roomJsonStr, err := json.Marshal(roomInfo)
		if err != nil {
			g.Log().Error(localCtx, err)
			return
		}
		service.Cache().Set(localCtx, consts.RoomInfoPrefix+roomInfo.ID, roomJsonStr, 0)

		// 创建一个临时结构体，只包含可序列化的字段
		playerDTOs := getPlayers(roomInfo)
		players, err := json.Marshal(playerDTOs)
		if err != nil {
			g.Log().Errorf(localCtx, "用户 %s 序列化玩家列表失败: %v", userName, err)
			return
		}

		msgData := &v1.SendActionReq{
			Type: consts.InitRoom,
			Data: gconv.String(players),
			From: userName,
		}
		c.safeSendMessage(localCtx, msgData)

	}(roomInfo, req.From)

}

// 处理创建房间
func (c *Controller) handleCreateRoom(ctx context.Context, req *v1.SendActionReq) {
	rmMu.Lock()
	defer rmMu.Unlock()

	//判断是否可以创建房间
	_, oldRoomInfo := c.getPlayerAndRoomInfo(ctx, req)
	if oldRoomInfo != nil && (oldRoomInfo.Status == 1 || oldRoomInfo.IsPlaying) {
		g.Log().Infof(ctx, "用户 %s 已在房间 %s 进行中，不能创建新房间", req.From, oldRoomInfo.ID)
		return
	}

	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）
	c.clearRoom(ctx, req)

	roomInfo := rm.CreateRoom(2) //创建好友房
	// c.safeSendToRoomChan(ctx, roomInfo)

	playerName := req.From
	humanPlayer := roomInfo.CreatePlayer(playerName, room.Human)
	humanPlayer.UserName = playerName
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
		g.Log().Errorf(ctx, "用户 %s 序列化玩家列表失败: %v", req.From, err)
		return
	}
	msgData := &v1.SendActionReq{
		Type: consts.CreateRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(players),
		From: req.From,
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
func (c *Controller) handleJoinRoom(ctx context.Context, req *v1.SendActionReq) {
	rmMu.Lock()
	defer rmMu.Unlock()
	// 解析加入房间数据（严格错误处理）
	var reqData struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		g.Log().Errorf(ctx, "用户 %s 解析加入房间数据失败: %v", req.From, err)
		return
	}
	roomID := reqData.RoomID
	if roomID == "" {
		g.Log().Errorf(ctx, "用户 %s 加入房间ID为空", req.From)
		return
	}
	// roomInfo, ok := rm.Rooms[roomID]
	roomInfoBytes := service.Cache().Get(ctx, consts.RoomInfoPrefix+roomID).Bytes()
	if roomInfoBytes == nil {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 不存在", req.From, roomID)
		return
	}
	var roomInfo room.Room
	if err = json.Unmarshal(roomInfoBytes, &roomInfo); err != nil {
		g.Log().Errorf(ctx, "房间信息json错误: %v", err)
		return
	}
	// 检查房间是否已满
	if len(roomInfo.Players) >= 4 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已满", req.From, roomID)
		return
	}
	// 检查房间是否正在游戏中
	if roomInfo.IsPlaying {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 正在游戏中", req.From, roomID)
		return
	}
	// 检查房间是否已开始游戏
	if roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 已开始游戏", req.From, roomID)
		return
	}
	//判断玩家已在房间
	inRoom := false
	var existingPlayer *room.Player // 用于存储已存在的玩家
	for _, player := range roomInfo.Players {
		if req.From == player.UserName {
			g.Log().Errorf(ctx, "用户 %s 已在房间 %s", req.From, roomID)
			inRoom = true
			existingPlayer = player // 记录已存在的玩家
			break
		}
	}
	//不在房间再加入房间，在房间直接返回房间数据
	var humanPlayer *room.Player
	if !inRoom {
		playerName := req.From
		humanPlayer = roomInfo.CreatePlayer(playerName, room.Human)
		humanPlayer.UserName = req.From
		humanPlayer.Name = req.From
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
	msgData := &v1.SendActionReq{
		Type: consts.JoinRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(playerDTOs),
		From: req.From,
	}
	// c.safeSendMessage(ctx, msgData)
	for _, player := range roomInfo.Players {
		msgData.From = player.UserName
		_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", req.From, err)
		}
	}
}

// 离开房间
func (c *Controller) handleLeaveRoom(ctx context.Context, req *v1.SendActionReq) {
	rmMu.Lock()
	defer rmMu.Unlock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx, req)
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
	msgData := &v1.SendActionReq{
		Type: consts.LeaveRoom, //此处有修改注意，客户端跳转
		Data: gconv.String(playerDTOs),
		From: req.From,
	}
	for _, player := range roomInfo.Players {
		msgData.From = player.UserName
		_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", req.From, err)
		}
	}
	//房间没人清理房间,或者是段位房
	c.clearRoom(ctx, req)

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
func (c *Controller) handlePlay(ctx context.Context, req *v1.SendActionReq) {
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx, req)
	if playerInfo == nil || roomInfo == nil {
		return
	}
	//判断是否为房主
	if playerInfo.ID != 0 {
		g.Log().Errorf(ctx, "用户 %s 不是房主", req.From)
		return
	}

	// 启动游戏逻辑（传入上下文，支持取消）
	//判断人数是否足够
	if len(roomInfo.Players) < 4 {
		g.Log().Errorf(ctx, "用户 %s 房间人数不足", req.From)
		return
	}

	//判断游戏进行中不能再次play,避免同一房间多次play
	if roomInfo.IsPlaying || roomInfo.Status == 1 {
		g.Log().Errorf(ctx, "用户 %s 房间 %s 进行中不能再次开始比赛", req.From, roomInfo.ID)
		return
	}

	rmMu.Lock()
	roomInfo.IsPlaying = true
	roomInfo.Status = 1 //游戏中
	rmMu.Unlock()
	//如果时好友房,通知房间内每个人开始游戏
	if roomInfo.Type == 2 {
		// roomInfo.SendRoomMessage(consts.Play, "")
		msgData := &v1.SendActionReq{
			Type: consts.Play, //此处有修改注意，客户端跳转
			Data: "",
			From: req.From,
		}
		for _, player := range roomInfo.Players {
			msgData.From = player.UserName
			_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, c.encodeMessage(ctx, msgData))
			if err != nil {
				g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", req.From, err)
			}
		}
	}

	//初始化房间定时器
	roomInfo.Rgtimer = gtimer.New()
	roomInfo.Rgtimer.Add(context.Background(), 1*time.Second, roomInfo.GameLoop)
	roomInfo.Rgtimer.Stop()
	roomInfo.MsgQueue = make(chan *v1.SendActionReq, 10)

	room.PlayOneGame(roomInfo)

	//发送日志 开始游戏
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     roomInfo.ID,
		Type:       int32(roomInfo.Type),
		Status:     int32(roomInfo.Status),
		UserID:     req.From,
		Point:      playerInfo.Point, //积分
		Action:     consts.Play,      //行为
		Remain:     "",               //剩余牌
		OutCardIds: "",               //玩家单次打出的牌id
		Text:       "",               //完整信息
	})
}

// 处理出牌
func (c *Controller) handlePlayCard(ctx context.Context, req *v1.SendActionReq) {
	rmMu.RLock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx, req)
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
	if err := json.Unmarshal([]byte(req.Data), &reqData); err != nil {
		g.Log().Errorf(ctx, "用户 %s 解析出牌数据失败: %v", req.From, err)
		return
	}
	if reqData.Pid != playerInfo.ID {
		g.Log().Errorf(ctx, "用户 %s PID不匹配（%d vs %d）", req.From, reqData.Pid, playerInfo.ID)
		return
	}

	//如果不是人类出牌时间，返回
	if playerInfo.ID != roomInfo.Current {
		g.Log().Errorf(ctx, "用户 %s 不是该玩家出牌阶段", req.From)
		return
	}

	// 更新玩家出牌信息
	// rmMu.Lock()
	// player.OutCardIds = reqData.CardIds
	// player.Pass = reqData.Pass
	// rmMu.Unlock()

	//发布出牌
	msgData := &v1.SendActionReq{
		Type: consts.PlayCard, //此处有修改注意，客户端跳转
		Data: gconv.String(reqData),
		From: req.From,
	}
	_, err := c.pubClient.Publish(ctx, consts.PlayerPlayCardPrefix+playerInfo.UserName, c.encodeMessage(ctx, msgData))
	g.Log().Infof(ctx, "用户 %s 发送出牌消息: %s", req.From, gconv.String(reqData))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送出牌消息失败: %v", req.From, err)
	}
}

// 处理获取房间信息
func (c *Controller) handleGetInfo(ctx context.Context, req *v1.SendActionReq) {
	rmMu.Lock()
	defer rmMu.Unlock()
	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx, req)
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
		g.Log().Errorf(ctx, "用户 %s 序列化房间信息失败: %v", req.From, err)
		return
	}

	// 7. 最后在锁外发送消息
	msgData := &v1.SendActionReq{
		Type: consts.GetInfo,
		Data: gconv.String(resData),
		From: req.From,
	}
	_, err = c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+playerInfo.UserName, c.encodeMessage(ctx, msgData))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", playerInfo.UserName, err)
	}
}

// 清理旧房间
func (c *Controller) clearRoom(ctx context.Context, req *v1.SendActionReq) {
	// 清理当前用户关联的旧房间（优化：定向清理，避免全量遍历）

	playerInfo, roomInfo := c.getPlayerAndRoomInfo(ctx, req)
	if playerInfo == nil || roomInfo == nil {
		return
	}

	if roomInfo.Type == 1 || len(roomInfo.Players) <= 0 {
		service.Cache().Remove(ctx, consts.RoomInfoPrefix+playerInfo.RoomID)
	}
	service.Cache().Remove(ctx, consts.PlayerInfoPrefix+req.From)

}

func (c *Controller) getPlayerAndRoomInfo(ctx context.Context, req *v1.SendActionReq) (*room.Player, *room.Room) {
	playerInfo := &room.Player{} // 或 new(room.Player)
	jsonStr := service.Cache().Get(ctx, consts.PlayerInfoPrefix+req.From).Bytes()
	if len(jsonStr) == 0 || jsonStr == nil {
		g.Log().Infof(ctx, "用户 %s 玩家不在房间内", req.From)
		return nil, nil
	}
	if err := json.Unmarshal(jsonStr, playerInfo); err != nil {
		g.Log().Errorf(ctx, "用户 %s 玩家信息json错误: %v", req.From, err)
		return nil, nil
	}

	roomInfoBytes := service.Cache().Get(ctx, consts.RoomInfoPrefix+playerInfo.RoomID).Bytes()
	if roomInfoBytes == nil {
		g.Log().Infof(ctx, "用户 %s 房间 %s 不存在", req.From, playerInfo.RoomID)
		return nil, nil
	}
	var roomInfo room.Room
	if err = json.Unmarshal(roomInfoBytes, &roomInfo); err != nil {
		g.Log().Errorf(ctx, "房间信息json错误: %v", err)
		return nil, nil
	}
	return playerInfo, &roomInfo
}

// 安全发送消息的函数
func (c *Controller) safeSendMessage(ctx context.Context, msg *v1.SendActionReq) {
	_, err := c.pubClient.Publish(ctx, consts.PlayerMsgPrefix+msg.From, c.encodeMessage(ctx, msg))
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", msg.From, err)
	}
}

// 消息编码（处理错误）
func (c *Controller) encodeMessage(ctx context.Context, msg *v1.SendActionReq) []byte {
	msgBytes, err := gjson.Encode(msg)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", msg.From, err)
		return []byte(`{"type":"error","data":"消息编码失败"}`)
	}
	return msgBytes
}
