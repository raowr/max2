package room

import (
	"context"
	"encoding/json"
	"fmt"
	actionv1 "game_user/api/action/v1"
	log_gamev1 "game_user/api/log_game/v1"
	v1 "game_user/api/settle/v1"
	"game_user/internal/consts"
	"game_user/internal/controller/log_game"
	"game_user/internal/controller/settle"
	"game_user/internal/service"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

// 创建新的房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms:      make(map[string]*Room),
		PlayerList: make(map[string]*Player),
	}
}

// 创建新玩家
func (room *Room) CreatePlayer(name string, playerType PlayerType) *Player {
	player := &Player{
		ID:          room.NextPlayerID,
		Name:        name,
		Type:        playerType,
		handPattern: make(map[int][][]Card),
		CardNum:     13,
		Point:       0, //GetGameInitPoint(), //初始积分
		RoomID:      room.ID,
		//UserId: 	 generateUserID(),//初始化用户ID
	}
	//如果是好友房初始化默认积分，如果是段位房使用玩家积累的积分
	if room.Type == 1 && playerType == Human {
		userInfo := service.Cache().Get(ctx, "user:"+name)
		entUser := Users{}
		if userInfo == nil {
			g.Log().Errorf(ctx, "用户不存在")

		}
		//缓存有，验证密码正确
		if err := json.Unmarshal(userInfo.Bytes(), &entUser); err != nil {
			g.Log().Errorf(ctx, "用户信息解析错误")

		}
		player.Point = entUser.Point
	}
	if room.Type == 2 {
		player.Point = GetGameInitPoint()
	}

	room.Players = append(room.Players, player) //加入房间
	room.NextPlayerID++
	return player
}

// 创建新房间（支持添加AI）
func (rm *RoomManager) CreateRoom(roomType int) *Room {
	// 生成唯一房间ID
	roomID := generateRoomID()

	// 创建新房间
	room := &Room{
		ID:        roomID,
		Players:   []*Player{},
		IsPlaying: false,
		// Rgtimer:   gtimer.New(),
		Status:   0, //未开始
		Type:     roomType,
		MsgQueue: make(chan *actionv1.SendActionReq, 10), // 缓冲大小为 10
	}

	// 添加AI机器人
	//for i := 0; i < aiCount; i++ {
	//	aiName := fmt.Sprintf("帅锅%d号", i+1)
	//	aiPlayer := rm.CreatePlayer(aiName, AI)
	//	room.Players = append(room.Players, aiPlayer)
	//	aiPlayer.RoomID = roomID
	//}

	// 将玩家加入房间
	//player.RoomID = roomID
	// room.Rgtimer.Add(context.Background(), 1*time.Second, room.GameLoop)
	// room.Rgtimer.Stop() //先停止

	// 保存房间
	rm.Rooms[roomID] = room

	return room
}

// 加入房间
func (rm *RoomManager) JoinRoom(player *Player, roomID string) bool {
	room, exists := rm.Rooms[roomID]
	if !exists {
		return false
	}

	// 检查房间是否已满
	if len(room.Players) >= 4 {
		return false
	}

	// 检查房间是否正在游戏中
	if room.IsPlaying {
		return false
	}

	// 将玩家加入房间
	room.Players = append(room.Players, player)
	player.RoomID = roomID

	return true
}

// 离开房间
func (rm *RoomManager) LeaveRoom(player *Player) {
	if player.RoomID == "" {
		return
	}

	room, exists := rm.Rooms[player.RoomID]
	if !exists {
		return
	}
	//停止房间定时器
	room.Rgtimer.Close()

	// 从房间中移除玩家
	for i, p := range room.Players {
		if p.ID == player.ID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}

	// 清空玩家的房间ID
	player.RoomID = ""

	delete(rm.Rooms, room.ID)
	room = nil
}

// 生成唯一房间ID
func generateRoomID() string {
	rand.Seed(time.Now().UnixNano())
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// 生成唯一玩家ID
func GenerateUserID() string {
	result := generateRoomID() + gconv.String(gtime.Now().TimestampMilli())
	return result
}

// 初始化一副牌
func initDeck() []Card {
	var deck []Card
	suits := []string{"♦", "♣", "♥", "♠"}
	values := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	cardId := 0
	for i, valueName := range values {
		for suit, suitName := range suits {
			cardId++
			card := Card{
				Value:    valueName,
				Suit:     suit,
				Name:     suitName + valueName,
				Id:       cardId,
				Rank:     3 + i,
				SuitName: suitName,
			}
			deck = append(deck, card)
		}
	}

	return deck
}

// 洗牌
func shuffleDeck(deck []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

// 清空玩家手中的牌
func clearPlayerCards(players []*Player) {
	for _, player := range players {
		player.Cards = []Card{}
		player.CardNum = 13
		player.handPattern = make(map[int][][]Card)
		player.Must = false
		player.OutCardIds = []int{}
	}
}

// 发牌
func dealCards(deck []Card, players []*Player) {
	//每位玩家13张牌
	for i := 0; i < len(deck); i++ {
		players[i%4].Cards = append(players[i%4].Cards, deck[i])
	}
}

// 显示玩家的牌
func (room *Room) showPlayerCards(player *Player) {
	var logCardsStr string
	if player.Type == Human {
		g.Log().Infof(ctx, "%s的牌: ", player.Name)
		// 创建牌名切片
		cardNames := make([]string, len(player.Cards))
		for i, card := range player.Cards {
			cardNames[i] = card.Name
		}

		// 使用 strings.Join 连接，自动处理分隔符
		cardsList := strings.Join(cardNames, ",")
		g.Log().Infof(ctx, "%s", cardsList)
		logCardsStr = cardsList
	} else {
		g.Log().Infof(ctx, "%s的牌: ", player.Name)
		logCardsSlice := []string{}
		for i, v := range player.handPattern {
			g.Log().Infof(ctx, "牌型：%d,具体的牌：%v ", i, v)
			var cardNameGroupSlice []string
			for _, group := range v {
				var cardNameSlice []string
				var cardNameStr string
				for _, card := range group {
					cardNameSlice = append(cardNameSlice, card.Name)
				}
				cardNameStr = "[" + strings.Join(cardNameSlice, ",") + "]"
				cardNameGroupSlice = append(cardNameGroupSlice, cardNameStr)
			}
			logCardsSlice = append(logCardsSlice, fmt.Sprintf("牌型：%d,具体的牌：%v", i, cardNameGroupSlice))
		}
		logCardsStr = strings.Join(logCardsSlice, ";")
	}
	//发送日志 发牌 显示玩家的牌
	log_game.SendLog(&log_gamev1.SendLogReq{
		RoomID:     room.ID,
		Type:       int32(room.Type),
		Status:     int32(room.Status),
		UserID:     player.UserName,
		Point:      player.Point, //积分
		Action:     "showCard",   //行为
		Remain:     logCardsStr,  //剩余牌
		OutCardIds: "",           //玩家单次打出的牌id
		Text:       "",           //完整信息
	})
}

// 确定谁先出牌
func bidLandlord(room *Room) {
	g.Log().Infof(ctx, "\n确定谁先出牌...")
	currentPlayer := 0
	for _, v := range room.Players {
		for _, card := range v.Cards {
			if card.Id == 1 {
				currentPlayer = v.ID
				//必须出牌
				v.Must = true
			}
		}
		v.CardNum = len(v.Cards)
	}
	// ♦3先出牌
	room.Current = currentPlayer
}

// 验证玩家输入的牌
func parseCardIndices(player *Player) bool {
	if player.Type == AI {
		return true
	}
	if len(player.OutCardIds) == 0 {
		return false // 输入无效
	}
	for _, cardId := range player.OutCardIds {
		isHas := false
		for _, v := range player.Cards {
			if cardId == v.Id {
				isHas = true
				break
			}
		}
		if !isHas {
			return false
		}
	}
	return true
}

// 从玩家手中获取选中的牌
func getSelectedCards(player *Player, indices []int) []Card {
	var selected []Card
	for _, idx := range indices {
		for _, card := range player.Cards {
			if idx == card.Id {
				selected = append(selected, card)
			}
		}
	}
	return selected
}

// 从玩家手中移除牌
func removeCards(player *Player, cardType int, indices []int) {
	// 创建 ID 集合用于快速查找
	idSet := make(map[int]bool)
	for _, id := range indices {
		idSet[id] = true
	}

	if player.Type == Human {
		newCards := make([]Card, 0, len(player.Cards))
		for _, card := range player.Cards {
			if !idSet[card.Id] {
				newCards = append(newCards, card)
			}
		}
		player.Cards = newCards
	} else {
		for patternType, patternGroups := range player.handPattern {
			newPatternGroups := make([][]Card, 0)
			for _, group := range patternGroups {
				newGroup := make([]Card, 0)
				for _, card := range group {
					if !idSet[card.Id] {
						newGroup = append(newGroup, card)
					}
				}
				// 只保留非空的牌组
				if len(newGroup) > 0 {
					newPatternGroups = append(newPatternGroups, newGroup)
				}
			}
			player.handPattern[patternType] = newPatternGroups
		}
		//如果需要重新整理牌
		if player.ReCard {
			player.Cards = make([]Card, 0)
			player.ReCard = false
			if len(player.handPattern) > 0 {
				for _, patternGroups := range player.handPattern {
					for _, group := range patternGroups {
						for _, card := range group {
							player.Cards = append(player.Cards, card)
						}
					}
				}
			}
			//重新整理牌
			if len(player.Cards) > 0 {
				player.handPattern = make(map[int][][]Card)
				player.Solve()
			}
		}
	}

}

// 判断牌型是否有效
func isValidCardType(cards []Card) (int, bool) {
	if len(cards) == 0 {
		return 0, true // 不出牌
	}

	// 单牌
	if len(cards) == 1 {
		return SINGLE, true
	}

	// 对子
	if len(cards) == 2 {
		if judgeTwoCards(cards) == PAIR {
			return PAIR, true
		}
	}

	//5张
	//花色(5张相同的花色)
	//福禄(3带2)
	//4条(4带1)
	//同花顺(5张点数连续的牌,花色相同)
	if len(cards) == 5 {
		if pokerHand := judgeFiveCards(cards); pokerHand > 0 {
			return pokerHand, true
		}
	}
	return 0, false
}

// AI出牌决策
// player Ai玩家
// Landlord 人类玩家
func aiDecideCards(player, landlord *Player, lastPH int, lastCards []Card) (pokerHandType int, playCards []Card) {

	//新做法
	//情况1：如果AI是大就是必出的，选择出最小牌所在的牌组
	//情况2：当玩家牌数大于5时，和情况1相同做法，如果玩家等于5张，选择小于5张出牌，如果玩家牌数小于5张，优先出比玩家剩余牌数多的牌组
	//情况3：玩家剩余一张牌时，优先出牌组，AI也剩余全是单牌，优先重大到小出牌，即顶牌
	//情况4：AI出牌选择相同牌型最小的牌组出

	// playCards := make([]Card, 0)
	//正常情况,正常出牌，但要比上一家大，而且牌型相同的最小牌组出

	return DecideCards(player, landlord, lastPH, lastCards)
}

// 检查游戏是否结束,算奖，输家剩几张牌，就输多少积分
// 赢玩家剩余多少张牌，就赢多少积分
func isGameOver(room *Room) (isOver bool, overMsgData []*overMsg, winName string) {
	var win int64 = 0
	var winer *Player
	for _, player := range room.Players {
		win += int64(player.CardNum)
		if player.CardNum == 0 {
			isOver = true
			winer = player
			winName = player.Name
		}
	}
	if winer != nil {
		winer.Win = win
		room.LastPH = 0
		for _, player := range room.Players {
			if player.Type == Human {
				var playerWin int64
				if player.ID == winer.ID {
					player.Point += winer.Win //赢的玩家
					playerWin = winer.Win     //输的是正数
				} else {
					player.Point -= int64(player.CardNum) //输的玩家
					playerWin = -int64(player.CardNum)    //输的是负数
				}
				//计算抽水,段位房才抽水
				if room.Type == 1 {
					commission := GetGameCommission()
					if player.Point > commission {
						player.Point -= commission
					} else {
						player.Point = 0
					}
				}
				overMsgData = append(overMsgData, &overMsg{
					WinName: player.Name,
					Winner:  player.ID,
					Win:     playerWin,
					Point:   player.Point,
				})
			}
		}
	}
	return isOver, overMsgData, winName
}

// 玩一局游戏
func PlayOneGame(room *Room) {
	g.Log().Infof(ctx, "\n===== 房间 %s 游戏开始 =====", room.ID)

	// 初始化游戏
	room.Deck = initDeck()

	// 洗牌
	shuffleDeck(room.Deck)

	// 清空玩家手牌
	clearPlayerCards(room.Players)

	// 发牌，每位玩家13张牌
	dealCards(room.Deck, room.Players)

	// 确定谁先出牌，既方块3在谁那
	bidLandlord(room)

	//发完牌开始整理机器人的牌
	for _, player := range room.Players {
		if int(player.Type) == int(AI) {
			player.Solve()
		}
	}
	//开始监听玩家出牌
	// 修复资源泄漏：先关闭旧的订阅
	// 检查是否已经启动了消息接收 goroutine
	go room.startMessageReceiver()

	room.pubClient = g.Redis()

	room.subClient = g.Redis()

	// ... existing code ...
	// 显示玩家的牌（除了底牌）
	for _, player := range room.Players {
		//只能显示自己的牌
		if player.Type == Human {
			room.Landlord = player
			cards := make([]int, 0)
			for _, card := range room.Landlord.Cards {
				cards = append(cards, card.Id)
			}
			go func(humanPlayer *Player) {
				data, _ := json.Marshal(struct {
					Cards          []int `json:"cards"`
					Current        int   `json:"current"`
					PlayerPoint    int64 `json:"playerPoint"`    //玩家总瓜子数
					OutCardTimeout int   `json:"outCardTimeout"` //出牌最大时间(单位秒) /s
				}{
					Cards:          cards,
					Current:        room.Current,
					PlayerPoint:    room.Landlord.Point,
					OutCardTimeout: GetOutCardTimeout(), //出牌最大时间(单位秒) /s
				})
				msgType := "showCard"

				//发送玩家牌信息
				msg := RoomMsg{
					Type: msgType,
					Data: gconv.String(data),
				}
				msgBytes, err := gjson.Encode(msg)
				if err != nil {
					g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", player.UserName, err)
				}
				_, err = room.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, msgBytes)
				if err != nil {
					g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", player.UserName, err)
				}
			}(player) // 传入当前玩家
		}
		room.showPlayerCards(player)
	}
	room.Rgtimer.Start()

}

// 在 startMessageReceiver 中（只启动一次）
func (room *Room) startMessageReceiver() {
	defer func() {
		g.Log().Error(ctx, "room startMessageReceiver 退出:")
	}()
	for {
		// 监听 ctx 是否已关闭 (调用 Close 方法时触发)
		select {
		case <-ctx.Done():
			g.Log().Info(ctx, "startMessageReceiver 停止监听")
			return
		default:
			// 继续执行后续功能
		}
		// 使用 Pattern 订阅当前出牌玩家出牌消息
		var currentName string
		if room.subConn == nil {
			currentName = room.Players[room.Current].Name
			sub, _, err := room.subClient.Subscribe(ctx, consts.PlayerPlayCardPrefix+currentName)
			if err != nil {
				g.Log().Error(ctx, "room startMessageReceiver Subscribe 失败:", err)
				return
			}
			room.subConn = sub
		}

		room.mutex.RLock()
		subConn := room.subConn
		room.mutex.RUnlock()

		msg, err := subConn.Receive(ctx)
		if err != nil {
			// 如果是 Context 被取消导致的错误，直接退出
			if ctx.Err() != nil {
				g.Log().Infof(ctx, "room ctx Watcher 订阅关闭:", ctx.Err())
			}
			// 其它错误（如网络断开），可以考虑简单的重试或记录日志
			g.Log().Infof(ctx, "room Casbin Watcher 订阅关闭:", err)
			// 如果出错直接退出，等待下一次重启或手动干预
			if room.Status == 1 || room.IsPlaying {
				continue
			} else {
				return
			}

		}

		// 解析消息
		msgData, success := ParseRedisSubscribeMessage(msg.String(), currentName, ctx)
		if success && msgData != nil {
			// 放入消息队列
			select {
			case room.MsgQueue <- msgData:
			default:
			}
		}
	}
}

// 房间循环定时器
func (room *Room) GameLoop(ctx context.Context) {
	// 检查游戏是否结束
	over, overMsgData, winName := isGameOver(room)
	if over {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done() // 使用defer确保即使发生panic也会调用Done()
			data, _ := json.Marshal(overMsgData)
			// 使用select和超时避免永久阻塞
			msgType := "over"
			room.safeSendRoomMessage(msgType, data)
			//发送日志 结算游戏
			for _, v := range overMsgData {
				g.Log().Infof(ctx, "游戏每个用户结算信息：%v, %v, %v, %v\n", v.Winner, v.Win, v.Point, v.WinName)
				for _, player := range room.Players {
					if player.ID == v.Winner {
						log_game.SendLog(&log_gamev1.SendLogReq{
							RoomID: room.ID,
							Type:   int32(room.Type),
							Status: int32(room.Status),
							UserID: player.UserName,
							Point:  player.Point,             //积分
							Action: "over",                   //行为
							Remain: getPlayerCardStr(player), //剩余牌
						})
					}
				}
			}
			//更新玩家积分
			if room.Type == 1 {
				for _, v := range overMsgData {
					userInfo := service.Cache().Get(ctx, "user:"+v.WinName)
					entUser := Users{}
					if userInfo == nil {
						g.Log().Errorf(ctx, "用户不存在")

					}
					//缓存有，验证密码正确
					if err := json.Unmarshal(userInfo.Bytes(), &entUser); err != nil {
						g.Log().Errorf(ctx, "用户名或密码错误")

					}
					entUser.Point = v.Point
					service.Cache().Set(ctx, "user:"+v.WinName, entUser, 7*24*time.Hour)

					//发送到user服保存数据库
					//发送日志 发牌 显示玩家的牌
					settle.SendSettle(&v1.SendSettleReq{
						Id:       int64(entUser.Id),
						UserName: v.WinName,
						Point:    v.Point,
					})
				}
			}
		}()

		room.Rgtimer.Stop()
		room.Rgtimer.Close()
		room.IsPlaying = false
		room.Status = 2 //结算中
		room.LastCards = make([]Card, 0)
		room.OutStarTime = 0
		room.passCount = 0
		//修复资源泄漏：先关闭旧的订阅
		room.subClientClose()
		g.Log().Infof(ctx, "\n游戏结束！恭喜%s！获胜\n", winName)
		//缓存当前房间信息
		wg.Add(1)
		go func() {
			defer wg.Done()
			roomJsonStr, err := json.Marshal(room)
			if err != nil {
				g.Log().Error(ctx, err)
			}
			service.Cache().Set(ctx, consts.RoomInfoPrefix+room.ID, roomJsonStr, 24*time.Hour)
		}()
		wg.Wait()
		return
	}

	currentPlayer := room.Players[room.Current]
	g.Log().Infof(ctx, "\n%s的回合 (当前手牌数: %d)\n", currentPlayer.Name, currentPlayer.CardNum)
	// room.showPlayerCards(currentPlayer)
	g.Log().Infof(ctx, "上一手牌，牌型: %v, 牌是：%s\n", room.LastPH, showCards(room.LastCards))
	g.Log().Infof(ctx, "请选择要出的牌 (输入牌的序号，用逗号分隔，0表示不出): ")

	var indices []int
	var selectedCards []Card

	// 记录开始出牌时间
	now := int(time.Now().Unix())
	if room.OutStarTime == 0 {
		room.OutStarTime = now
	}

	if currentPlayer.Type == AI {
		// AI决策，随机秒数
		//thingTime := grand.N(1,5)
		//time.Sleep(time.Duration(thingTime) * time.Second) // 模拟思考时间
		//记录开始出牌时间

		thinkTime := grand.N(2, 5) //模拟思考秒数
		if now-room.OutStarTime < thinkTime {
			return
		}
		// room.OutStarTime = 0
		room.LastPH, selectedCards = aiDecideCards(currentPlayer, room.Landlord, room.LastPH, room.LastCards)
		// selectedCards = getSelectedCards(currentPlayer, indices)
		if len(selectedCards) > 0 {
			for _, v := range selectedCards {
				indices = append(indices, v.Id)
			}
		}

	} else {
		// 非阻塞读取：读取一次消息，如果没有消息则跳过
		var msgData *actionv1.SendActionReq
		var hasMessage bool

		// 非阻塞读取消息队列（不会创建新 goroutine）
		select {
		case msgData = <-room.MsgQueue:
			hasMessage = true
			g.Log().Infof(ctx, "从队列收到消息: %v", msgData.Type)
		default:
			// 没有消息，直接返回（下次循环再检查）
			// return
		}

		var reqData struct {
			Pid     int   `json:"pid"`
			CardIds []int `json:"cardIds"`
			Pass    int   `json:"pass"`
		}
		g.Log().Infof(ctx, "用户 %s 接收出牌消息1: %v", currentPlayer.Name, msgData)
		if hasMessage && msgData != nil {
			err := gconv.Struct(msgData.Data, &reqData)
			if err != nil {
				// 处理错误
				g.Log().Errorf(ctx, "用户 %s 消息解码失败: %v", currentPlayer.Name, err)
			}
		}
		g.Log().Infof(ctx, "用户 %s 接收出牌消息2: %v", currentPlayer.Name, reqData)
		currentPlayer.OutCardIds = reqData.CardIds
		currentPlayer.Pass = reqData.Pass

		outCardTimeout := GetOutCardTimeout()
		//玩家不操作逻辑
		if len(currentPlayer.OutCardIds) <= 0 && //玩家还没出牌
			now-room.OutStarTime < outCardTimeout && //还在倒计时中
			currentPlayer.Pass == 0 { //玩家还没不出
			return
		} else {
			//有出牌
			if len(currentPlayer.OutCardIds) > 0 && now-room.OutStarTime < outCardTimeout {
				g.Log().Debug(ctx, currentPlayer.OutCardIds)
				//判断出牌数量是否符合规则
				if len(currentPlayer.OutCardIds) != 1 &&
					len(currentPlayer.OutCardIds) != 2 &&
					len(currentPlayer.OutCardIds) != 5 {
					return
				}
				//验证牌是否存在
				if !parseCardIndices(currentPlayer) {
					currentPlayer.OutCardIds = make([]int, 0)
					g.Log().Infof(ctx, "输入的牌不存在%v \n", currentPlayer.OutCardIds)
					msg := fmt.Sprintf("输入的牌不存在:%v", currentPlayer.OutCardIds)
					//通知出牌失败
					go func() {
						data, _ := json.Marshal(struct {
							Pid  int    `json:"pid"`
							Code int    `json:"code"` //0成功,非零失败
							Msg  string `json:"msg"`  //提示消息
						}{
							Pid:  currentPlayer.ID,
							Code: 1,
							Msg:  msg,
						})
						msgType := "outCard"
						room.safeSendRoomMessage(msgType, data)
					}()
					currentPlayer.OutCardIds = make([]int, 0)
					return
				}
				//获取出牌索引
				indices = currentPlayer.OutCardIds
				selectedCards = getSelectedCards(currentPlayer, indices)
				//判断牌型是否正确
				//判断牌型LastPH
				LastPH, isHas := isValidCardType(selectedCards)
				if !isHas {
					currentPlayer.OutCardIds = make([]int, 0)
					g.Log().Debug(ctx, "牌型无效,selectedCards:%v", selectedCards)
					msg := fmt.Sprintf("牌型无效,selectedCards:%v", selectedCards)
					//通知出牌失败
					go func() {
						data, _ := json.Marshal(struct {
							Pid  int    `json:"pid"`
							Code int    `json:"code"` //0成功,非零失败
							Msg  string `json:"msg"`  //提示消息
						}{
							Pid:  currentPlayer.ID,
							Code: 1,
							Msg:  msg,
						})
						msgType := "outCard"
						room.safeSendRoomMessage(msgType, data)
					}()
					return
				}
				//判断和上一手牌型是否相同
				if room.LastPH > 0 {
					if room.LastPH == SINGLE || room.LastPH == PAIR {
						if LastPH != room.LastPH {
							currentPlayer.OutCardIds = make([]int, 0)
							g.Log().Debug(ctx, "牌型和上一手牌型不一致1,LastPH:%v,room.LastPH:%v", LastPH, room.LastPH)
							msg := fmt.Sprintf("牌型和上一手牌型不一致1,LastPH:%v,room.LastPH:%v", LastPH, room.LastPH)
							//通知出牌失败
							go func() {
								data, _ := json.Marshal(struct {
									Pid  int    `json:"pid"`
									Code int    `json:"code"` //0成功,非零失败
									Msg  string `json:"msg"`  //提示消息
								}{
									Pid:  currentPlayer.ID,
									Code: 1,
									Msg:  msg,
								})
								msgType := "outCard"
								room.safeSendRoomMessage(msgType, data)
							}()
							return
						}
					} else {
						if LastPH < room.LastPH {
							currentPlayer.OutCardIds = make([]int, 0)
							g.Log().Debug(ctx, "牌型和上一手牌型不一致2,LastPH:%v,room.LastPH:%v", LastPH, room.LastPH)
							msg := fmt.Sprintf("牌型和上一手牌型不一致2,LastPH:%v,room.LastPH:%v", LastPH, room.LastPH)
							//通知出牌失败
							go func() {
								data, _ := json.Marshal(struct {
									Pid  int    `json:"pid"`
									Code int    `json:"code"` //0成功,非零失败
									Msg  string `json:"msg"`  //提示消息
								}{
									Pid:  currentPlayer.ID,
									Code: 1,
									Msg:  msg,
								})
								msgType := "outCard"
								room.safeSendRoomMessage(msgType, data)
							}()
							return
						}
					}

				}

			}

			// selectedCards = getSelectedCards(currentPlayer, indices)
		}
		//倒计时结束，自动出牌逻辑
		if now-room.OutStarTime >= outCardTimeout {
			currentPlayer.OutCardIds = make([]int, 0)
			indices = make([]int, 0)
			selectedCards = make([]Card, 0)                               //超过出牌时间过
			if len(currentPlayer.OutCardIds) <= 0 && currentPlayer.Must { //如果是必出牌最小的一张
				minCardId := currentPlayer.Cards[0].Id  //最小牌ID
				for _, v := range currentPlayer.Cards { //找出最小牌Id
					if v.Id < minCardId {
						minCardId = v.Id
					}
				}
				indices = append(indices, minCardId)
				selectedCards = getSelectedCards(currentPlayer, indices)
				currentPlayer.OutCardIds = indices
			}
		}
		//玩家点不出逻辑
		if currentPlayer.Pass == 1 {
			currentPlayer.Pass = 0
			if currentPlayer.Must {
				if now-room.OutStarTime < outCardTimeout {
					//超时出牌
					return
				} //未超时不做处理,上面已经处理
			} else {
				//玩家不出
				currentPlayer.OutCardIds = make([]int, 0)
				indices = make([]int, 0)
				selectedCards = make([]Card, 0)
			}
		}

	}
	code := 0
	// 验证牌型
	cardType, valid := isValidCardType(selectedCards)
	if !valid {
		g.Log().Infof(ctx, "牌型无效，请重新选择 %v \n", selectedCards)
		code = 1

	}
	// 验证是否能压过上一手牌
	if !compareCards(room.LastPH, cardType, room.LastCards, selectedCards) {
		g.Log().Infof(ctx, "不能压过上5一手牌,请重新选择")
		code = 1

	}
	if code != 0 {
		//通知出牌失败
		go func() {
			data, _ := json.Marshal(struct {
				Pid  int `json:"pid"`
				Code int `json:"code"` //0成功,非零失败
			}{
				Pid:  currentPlayer.ID,
				Code: code,
			})
			msgType := "outCard"
			room.safeSendRoomMessage(msgType, data)
		}()
		currentPlayer.OutCardIds = make([]int, 0)
		return
	}
	var msgType = ""
	// 如果是不出牌
	if len(selectedCards) == 0 {
		g.Log().Infof(ctx, "%s不出\n", currentPlayer.Name)
		//通知用户,某位机器人不出牌，
		msgType = "pass"
		room.passCount++
		// 如果三家都不出，重置上一手牌
		if room.passCount >= 3 {
			room.LastPH = 0
			room.LastCards = []Card{}
			room.passCount = 0
		}
		//发送日志 玩家不出牌
		log_game.SendLog(&log_gamev1.SendLogReq{
			RoomID: room.ID,
			Type:   int32(room.Type),
			Status: int32(room.Status),
			UserID: currentPlayer.UserName,
			Point:  currentPlayer.Point,             //积分
			Action: "pass",                          //行为
			Remain: getPlayerCardStr(currentPlayer), //剩余牌
		})
	} else {
		// 出牌
		g.Log().Infof(ctx, "%s出了: %s (%v)\n", currentPlayer.Name, showCards(selectedCards), cardType)
		msgType = "outCard"
		//其他人改为非必出
		for _, v := range room.Players {
			if v.ID == currentPlayer.ID {
				v.Must = true //谁出牌，谁就暂定是必出
			} else {
				v.Must = false //其他人可以不出牌
			}
		}
		removeCards(currentPlayer, cardType, indices)
		currentPlayer.CardNum -= len(selectedCards)
		room.LastPH = cardType
		room.LastCards = selectedCards
		room.passCount = 0
		//发送日志 玩家出牌
		log_game.SendLog(&log_gamev1.SendLogReq{
			RoomID:     room.ID,
			Type:       int32(room.Type),
			Status:     int32(room.Status),
			UserID:     currentPlayer.UserName,
			Point:      currentPlayer.Point,             //积分
			Action:     "outCard",                       //行为
			Remain:     getPlayerCardStr(currentPlayer), //剩余牌
			OutCardIds: showCards(selectedCards),
		})
	}
	room.OutStarTime = now
	// room.OutStarTime = 0 //人类出牌时间恢复为0
	currentPlayer.OutCardIds = make([]int, 0)
	// 下一个玩家
	room.Current = (room.Current + 1) % 4
	room.Turn++

	var mustPid int
	for _, v := range room.Players {
		if v.Must {
			mustPid = v.ID //必须出牌的玩家
		}
	}

	//修复资源泄漏：先关闭旧的订阅
	room.subClientClose()

	//缓存当前房间信息
	roomJsonStr, err := json.Marshal(room)
	if err != nil {
		g.Log().Error(ctx, err)
	}
	service.Cache().Set(ctx, consts.RoomInfoPrefix+room.ID, roomJsonStr, 24*time.Hour)

	//通知用户,出牌，
	data, _ := json.Marshal(struct {
		Pid            int   `json:"pid"`
		CardIds        []int `json:"cardIds"`
		CardType       int   `json:"card_type"`
		CardsNum       int   `json:"cards_num"`      //剩余牌数
		Code           int   `json:"code"`           //0成功,非零失败
		Current        int   `json:"current"`        //当前出牌玩家
		MustPid        int   `json:"mustPid"`        //必须出牌的玩家ID
		OutCardTimeout int   `json:"outCardTimeout"` //出牌最大时间(单位秒) /s
	}{
		Pid:            currentPlayer.ID,
		CardIds:        indices,
		CardType:       cardType,
		CardsNum:       currentPlayer.CardNum,
		Code:           0,
		Current:        room.Current,
		MustPid:        mustPid,
		OutCardTimeout: GetOutCardTimeout(), //出牌最大时间(单位秒) /s
	})
	room.safeSendRoomMessage(msgType, data)

}

// 安全发送消息的函数
func (room *Room) safeSendRoomMessage(msgType string, data any) {
	// 先加锁保护Players切片的访问
	// room.mutex.RLock()
	players := room.Players
	// room.mutex.RUnlock()

	// 转换数据为字符串形式
	// dataStr := gconv.String(data)

	// 遍历每个玩家，分别发送消息
	for _, player := range players {
		if player.Type == AI {
			continue
		}
		//发送玩家牌信息
		msg := RoomMsg{
			Type: msgType,
			Data: gconv.String(data),
		}
		msgBytes, err := gjson.Encode(msg)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 消息编码失败: %v", player.UserName, err)
		}
		_, err = room.pubClient.Publish(ctx, consts.PlayerMsgPrefix+player.UserName, msgBytes)
		if err != nil {
			g.Log().Errorf(ctx, "用户 %s 发送消息失败: %v", player.UserName, err)
		}
	}
}

// 提供给外部使用
func (room *Room) SendRoomMessage(msgType string, data any) {
	room.safeSendRoomMessage(msgType, data)
}

func (room *Room) subClientClose() {
	room.mutex.Lock()
	old := room.subConn
	room.subConn = nil
	room.mutex.Unlock()
	if old != nil {
		_ = old.Close(ctx) // 在锁外关闭，避免死锁或长时间持锁
	}
}

// 显示房间列表
func showRooms(rm *RoomManager) {
	g.Log().Infof(ctx, "\n===== 房间列表 =====")
	if len(rm.Rooms) == 0 {
		g.Log().Infof(ctx, "暂无可用房间")
		return
	}

	for id, room := range rm.Rooms {
		status := "等待中"
		if room.IsPlaying {
			status = "游戏中"
		}
		g.Log().Infof(ctx, "房间ID: %s, 玩家数: %d/3, 状态: %s\n", id, len(room.Players), status)
	}
}

// 开始游戏主程序
/*func startGame() {
	g.Log().Infof(ctx,"欢迎来到斗地主游戏！")
	rm := NewRoomManager()

	// 创建第一个玩家（人类）
	g.Log().Infof(ctx,"请输入你的名字: ")
	var playerName string
	fmt.Scanln(&playerName)
	humanPlayer := rm.CreatePlayer(playerName, Human)

	for {
		g.Log().Infof(ctx,"\n===== 主菜单 =====")
		g.Log().Infof(ctx,"1. 创建房间（与机器人对战）")
		g.Log().Infof(ctx,"2. 加入房间（与其他玩家对战）")
		g.Log().Infof(ctx,"3. 查看房间列表")
		g.Log().Infof(ctx,"4. 退出游戏")

		if humanPlayer.RoomID != "" {
			room := rm.Rooms[humanPlayer.RoomID]
			g.Log().Infof(ctx,"\n你当前在房间 %s 中，已有 %d/3 名玩家\n", room.ID, len(room.Players))
			g.Log().Infof(ctx,"5. 离开房间")
			if len(room.Players) == 3 && !room.IsPlaying {
				g.Log().Infof(ctx,"6. 开始游戏")
			}
		}

		g.Log().Infof(ctx,"请选择操作: ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if humanPlayer.RoomID != "" {
				g.Log().Infof(ctx,"你已经在一个房间中，请先离开当前房间")
				break
			}

			// 选择机器人数量
			g.Log().Infof(ctx,"请选择机器人数量 (1-2): ")
			var aiCountStr string
			fmt.Scanln(&aiCountStr)
			aiCount, err := strconv.Atoi(aiCountStr)
			if err != nil || aiCount < 1 || aiCount > 2 {
				g.Log().Infof(ctx,"无效的数量，默认创建2个机器人")
				aiCount = 2
			}

			room := rm.CreateRoom(humanPlayer, aiCount)
			g.Log().Infof(ctx,"房间创建成功！房间ID: %s\n", room.ID)
			g.Log().Infof(ctx,"房间内有 %d 个机器人，共 %d 名玩家\n", aiCount, len(room.Players))
			if len(room.Players) == 3 {
				g.Log().Infof(ctx,"可以选择开始游戏了")
			}

		case "2":
			if humanPlayer.RoomID != "" {
				g.Log().Infof(ctx,"你已经在一个房间中，请先离开当前房间")
				break
			}
			g.Log().Infof(ctx,"请输入房间ID: ")
			var roomID string
			fmt.Scanln(&roomID)
			success := rm.JoinRoom(humanPlayer, roomID)
			if success {
				room := rm.Rooms[roomID]
				g.Log().Infof(ctx,"成功加入房间 %s！当前房间有 %d/3 名玩家\n", roomID, len(room.Players))
			} else {
				g.Log().Infof(ctx,"加入房间失败，房间不存在或已满或正在游戏中")
			}

		case "3":
			showRooms(rm)

		case "4":
			// 退出前离开房间
			if humanPlayer.RoomID != "" {
				rm.LeaveRoom(humanPlayer)
			}
			g.Log().Infof(ctx,"谢谢游玩，再见！")
			return

		case "5":
			if humanPlayer.RoomID == "" {
				g.Log().Infof(ctx,"你不在任何房间中")
				break
			}
			rm.LeaveRoom(humanPlayer)
			g.Log().Infof(ctx,"已离开房间")

		case "6":
			if humanPlayer.RoomID == "" {
				g.Log().Infof(ctx,"你不在任何房间中")
				break
			}
			room := rm.Rooms[humanPlayer.RoomID]
			if len(room.Players) != 3 {
				g.Log().Infof(ctx,"玩家不足3人，无法开始游戏")
				break
			}
			if room.IsPlaying {
				g.Log().Infof(ctx,"游戏已经开始")
				break
			}
			room.IsPlaying = true
			PlayOneGame(room)

		default:
			g.Log().Infof(ctx,"无效的操作，请重新选择")
		}
	}
}*/
