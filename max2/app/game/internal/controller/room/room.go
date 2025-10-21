package room

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/util/grand"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/util/gconv"
)

type MsgChan chan RoomMsg

type RoomMsg struct {
	Type string
	Data string
}

// 牌的结构体
type Card struct {
	Value    string // 牌值 3-15，分别对应3-A
	Suit     string // 花色 0-3 方块、梅花、红桃、黑桃
	Name     string // 牌的名称
	Id       int    //牌的id
	Rank     int    //牌的点数
	SuitName string //花色名称
}

// 玩家类型
type PlayerType int

const (
	Human PlayerType = iota
	AI
)

const (
	SINGLE   = iota + 1 // 单牌
	PAIR                // 对子
	STRAIGHT            // 顺子
	SUIT                // 花色
	THREE               // 3带2(葫芦)
	FOUR                // 4带1
	FLUSH               // 同花顺
)

// 玩家结构体
type Player struct {
	ID          int
	Cards       []Card
	Name        string
	RoomID      string           // 玩家所在房间ID
	Type        PlayerType       // 玩家类型（人类或AI）
	OutCardIds  []int            //玩家单次打出的牌id
	Must        bool             //是否必须要出牌
	Win         int64            //奖励,单局奖励
	Pass        int              //是否跳过
	handPattern map[int][][]Card //整理的牌型放数组中
	CardNum     int              //牌数
	ReCard      bool             //是否需要重新整理牌
	Point       int64            //积分，总积分
}

// 房间结构体
type Room struct {
	ID          string
	Players     []*Player
	Deck        []Card
	Landlord    *Player
	Farmers     []*Player
	Current     int    // 当前出牌玩家索引
	LastCards   []Card // 上一手牌
	LastPH      int    //上一手牌型
	Turn        int    // 轮次
	IsPlaying   bool   // 房间是否正在游戏中
	MsgChan     chan RoomMsg
	Rgtimer     *gtimer.Timer
	OutStarTime int //出牌开始时间
	passCount   int //不出次数
}

// 房间管理器
type RoomManager struct {
	Rooms        map[string]*Room
	PlayerList   map[int]*Player
	NextPlayerID int
}

// 创建新的房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms:        make(map[string]*Room),
		PlayerList:   make(map[int]*Player),
		NextPlayerID: 0,
	}
}

// 创建新玩家
func (rm *RoomManager) CreatePlayer(name string, playerType PlayerType) *Player {
	player := &Player{
		ID:          rm.NextPlayerID,
		Name:        name,
		Type:        playerType,
		handPattern: make(map[int][][]Card),
		CardNum:     13,
		Point:       100, //初始100积分
	}
	rm.PlayerList[player.ID] = player
	rm.NextPlayerID++
	return player
}

// 创建新房间（支持添加AI）
func (rm *RoomManager) CreateRoom(player *Player, aiCount int) *Room {
	// 生成唯一房间ID
	roomID := generateRoomID()

	// 创建新房间
	room := &Room{
		ID:        roomID,
		Players:   []*Player{player},
		IsPlaying: false,
		MsgChan:   make(chan RoomMsg),
		Rgtimer:   gtimer.New(),
	}
	room.Rgtimer.Add(context.Background(), 1*time.Second, room.GameLoop)
	room.Rgtimer.Stop() //先停止

	// 添加AI机器人
	for i := 0; i < aiCount; i++ {
		aiName := fmt.Sprintf("帅锅%d号", i+1)
		aiPlayer := rm.CreatePlayer(aiName, AI)
		room.Players = append(room.Players, aiPlayer)
		aiPlayer.RoomID = roomID
	}

	// 将玩家加入房间
	player.RoomID = roomID

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

	// 从房间中移除玩家
	for i, p := range room.Players {
		if p.ID == player.ID {
			room.Players = append(room.Players[:i], room.Players[i+1:]...)
			break
		}
	}

	// 清空玩家的房间ID
	player.RoomID = ""

	// 如果房间为空，删除房间
	if len(room.Players) == 0 {
		delete(rm.Rooms, room.ID)
	} else if room.IsPlaying {
		// 如果游戏正在进行中且玩家离开，结束当前游戏
		room.IsPlaying = false
		fmt.Printf("玩家 %s 已离开房间，游戏结束\n", player.Name)
	}
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

// 初始化一副牌
func initDeck() []Card {
	var deck []Card
	suits := []string{"♦", "♣", "♥", "♠"}
	values := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	cardId := 0
	for i, valueName := range values {
		for _, suitName := range suits {
			cardId++
			card := Card{
				Value:    valueName,
				Suit:     suitName,
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
func showPlayerCards(player *Player) {
	if player.Type == Human {
		fmt.Printf("%s的牌: ", player.Name)
		for i, card := range player.Cards {
			fmt.Printf("%d.%s ", i+1, card.Name)
		}
		fmt.Println()
	} else {
		fmt.Printf("%s的牌: ", player.Name)
		for i, v := range player.handPattern {
			fmt.Printf("牌型：%d,具体的牌：%v ", i, v)
		}
		fmt.Println()
	}

}

// 确定谁先出牌
func bidLandlord(room *Room) {
	fmt.Println("\n确定谁先出牌...")
	currentPlayer := 0
	for _, v := range room.Players {
		for _, card := range v.Cards {
			if card.Id == 1 {
				currentPlayer = v.ID
				//必须出牌
				v.Must = true
			}
		}
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

// 比较牌的大小
func compareCards(lastPH, pokerHand int, last []Card, current []Card) bool {
	if len(last) == 0 {
		return true // 上一手没牌，当前任何有效牌型都可以出
	}

	if len(current) == 0 {
		return true // 不出牌
	}
	//牌型相同比较
	if lastPH == pokerHand {
		return getMaxValue(pokerHand, current) > getMaxValue(lastPH, last)
	} else {
		//牌型不同，当前牌型必须大于上一手牌型
		return pokerHand > lastPH
	}
}

// 获取牌组中的最大值
func getMaxValue(pokerHand int, cards []Card) (maxVal int) {

	//分牌型获取最大值
	if pokerHand == SINGLE || //单牌最大值
		pokerHand == PAIR || //对子最大值
		pokerHand == STRAIGHT || //顺子最大值
		pokerHand == SUIT || //同花最大值
		pokerHand == FLUSH { //同花顺最大值
		for _, card := range cards {
			if card.Id > maxVal {
				maxVal = card.Id
			}
		}
	}
	if pokerHand == THREE { //三带二最大值
		rankCount := countRanks(cards)
		for rank, count := range rankCount {
			if count == 3 {
				maxVal = rank
			}
		}
	}
	if pokerHand == FOUR { //四带一最大值
		rankCount := countRanks(cards)
		for rank, count := range rankCount {
			if count == 4 {
				maxVal = rank
			}
		}
	}
	return maxVal
}

// 显示牌组
func showCards(cards []Card) string {
	if len(cards) == 0 {
		return "不出"
	}

	var names []string
	for _, card := range cards {
		names = append(names, card.Name)
	}
	return strings.Join(names, " ")
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
	if len(lastCards) > 0 {
		pokerHandCards := make([][]Card, 0)
		if lastPH == SINGLE || lastPH == PAIR { //单牌和对子这样比比较
			for pokerHand, v := range player.handPattern {
				if pokerHand == lastPH {
					for _, v1 := range v {
						if compareCards(lastPH, pokerHand, lastCards, v1) {
							pokerHandCards = append(pokerHandCards, v1)
						}
					}
				}
			}
			if len(pokerHandCards) > 0 {
				playCards = pokerHandCards[0] //最小牌组
				for _, v := range pokerHandCards {
					if !compareCards(lastPH, lastPH, playCards, v) {
						playCards = v
					}
				}
			}
		} else { //5张的这样比较
			pokerHandNew := make(map[int][][]Card, 0)
			for pokerHand, v := range player.handPattern {
				if pokerHand > lastPH { //牌型比上一手牌大
					pokerHandNew[pokerHand] = v
				} else if pokerHand == lastPH { //牌型形同
					for _, v1 := range v {
						if compareCards(lastPH, pokerHand, lastCards, v1) {
							pokerHandCards = append(pokerHandCards, v1)
						}
					}
				}
			}
			for pokerHand, v := range pokerHandNew {
				for _, v1 := range v {
					if compareCards(lastPH, pokerHand, lastCards, v1) {
						pokerHandCards = append(pokerHandCards, v1)
					}
				}
			}
			if len(pokerHandCards) > 0 {
				playCards = pokerHandCards[0] //最小牌组
				for _, v := range pokerHandCards {
					if !compareCards(lastPH, lastPH, playCards, v) {
						playCards = v
					}
				}
			}
		}

		//如果时3号玩家，上次出牌是单牌，玩家剩余1张牌，选择比上家大的最大的单牌出
		if player.ID == 3 && lastPH == SINGLE && len(landlord.Cards) == 1 {
			if len(pokerHandCards) > 0 {
				playCards = pokerHandCards[0] //最大牌组
				for _, v := range pokerHandCards {
					if compareCards(lastPH, lastPH, playCards, v) {
						playCards = v
					}
				}
			}
		}
		//如果时3号玩家，如果打的是单张,玩家3没有单张，需要拆出单张，如果打的是对，需要从3条和4条中拆
		if player.ID == 3 && lastPH == SINGLE && len(landlord.Cards) == 1 && len(playCards) <= 0 {
			//选择最大牌所在的牌组
			maxCardId := 0 //最大牌id
			maxCard := make([]Card, 0)
			if player.CardNum > 0 {
				for _, v := range player.handPattern {
					for _, v1 := range v {
						for _, v2 := range v1 {
							if v2.Id > maxCardId {
								maxCardId = v2.Id
							}
						}
					}
				}
				if maxCardId > 0 {
					findCard := false
					for _, v := range player.handPattern {
						for _, v1 := range v {
							for _, v2 := range v1 {
								if v2.Id == maxCardId {
									maxCard = append(maxCard, v2)
									findCard = true
									break
								}
							}
							if findCard {
								break
							}
						}
						if findCard {
							break
						}
					}
					//如果比上手大，就需要重新整理牌
					if compareCards(lastPH, lastPH, lastCards, maxCard) {
						playCards = maxCard
						player.ReCard = true
					}
				}
			}

		}
		pokerHandType = lastPH
	}
	//情况1
	if player.Must || len(lastCards) <= 0 {
		//先取最小牌在牌组中，
		minCardId := 52 //最小牌id
		for _, v := range player.handPattern {
			for _, v1 := range v {
				for _, v2 := range v1 {
					if v2.Id < minCardId {
						minCardId = v2.Id
					}
				}
			}
		}
		findCard := false
		for pokerHand, v := range player.handPattern {
			for _, v1 := range v {
				for _, v2 := range v1 {
					if v2.Id == minCardId {
						findCard = true
						break
					}
				}
				if findCard {
					playCards = v1
					break
				}
			}
			if findCard {
				pokerHandType = pokerHand
				break
			}
		}
		//这里判断玩家牌数量，如果有小于5张的牌组，优先出这些牌组
		//情况2,玩家牌数等于5张时
		if len(landlord.Cards) == 5 {
			pokerHandCards := make([]Card, 0)
			for pokerHand, v := range player.handPattern {
				if pokerHand <= PAIR {
					for _, v1 := range v {
						pokerHandCards = append(pokerHandCards, v1...)
					}
				}
			}
			var minCardId int
			if len(pokerHandCards) > 0 {
				minCardId = pokerHandCards[0].Id //最小牌id
				for _, v := range pokerHandCards {
					if v.Id < minCardId {
						minCardId = v.Id
					}
				}
			}
			for pokerHand, v := range player.handPattern {
				if pokerHand <= PAIR {
					for _, v1 := range v {
						for _, v2 := range v1 {
							if v2.Id == minCardId {
								playCards = v1
								pokerHandType = pokerHand
								break
							}
						}
					}
				}
			}
		}
		//情况2,玩家牌数小于5张时,优先出比玩家剩余牌数多的牌组,从多的往少的找
		if 1 < len(landlord.Cards) && len(landlord.Cards) < 5 {
			pokerHandCards := make([]Card, 0)
			for pokerHand, v := range player.handPattern {
				if pokerHand > PAIR {
					for _, v1 := range v {
						pokerHandCards = append(pokerHandCards, v1...)
					}
				}
			}
			var minCardId int
			if len(pokerHandCards) > 0 {
				minCardId = pokerHandCards[0].Id //最小牌id
				for _, v := range pokerHandCards {
					if v.Id < minCardId {
						minCardId = v.Id
					}
				}
			}
			for pokerHand, v := range player.handPattern {
				if pokerHand <= PAIR {
					for _, v1 := range v {
						for _, v2 := range v1 {
							if v2.Id == minCardId {
								playCards = v1
								pokerHandType = pokerHand
								break
							}
						}
					}
				}
			}
		}
		//如果玩家剩余1张牌
		if len(landlord.Cards) == 1 {
			//优先出大于单牌数的,全部Ai从大打到小
			pokerHandCards := make([]Card, 0)
			for pokerHand, v := range player.handPattern {
				if pokerHand > SINGLE {
					for _, v1 := range v {
						pokerHandCards = append(pokerHandCards, v1...)
					}
				}
			}
			var minCardId int
			if len(pokerHandCards) > 0 {
				minCardId = pokerHandCards[0].Id //最小牌id
				for _, v := range pokerHandCards {
					if v.Id < minCardId {
						minCardId = v.Id
					}
				}
				for pokerHand, v := range player.handPattern {
					if pokerHand <= PAIR {
						for _, v1 := range v {
							for _, v2 := range v1 {
								if v2.Id == minCardId {
									playCards = v1
									pokerHandType = pokerHand
									break
								}
							}
						}
					}
				}
			} else {
				pokerHandCards := make([][]Card, 0)
				for pokerHand, v := range player.handPattern {
					if pokerHand == SINGLE {
						pokerHandCards = v
					}
				}
				if len(pokerHandCards) > 0 {
					playCards = pokerHandCards[0] //最大牌组
					for _, v := range pokerHandCards {
						if compareCards(SINGLE, SINGLE, playCards, v) {
							playCards = v
						}
					}
				}
			}
		}

	}

	return
}

// 检查游戏是否结束,算奖，输家剩几张牌，就输多少积分
// 赢玩家剩余多少张牌，就赢多少积分
func isGameOver(room *Room) (isOver bool, winer *Player, playerPoint, playerWin int64) {
	var win int64 = 0
	for _, player := range room.Players {
		win += int64(player.CardNum)
		if player.CardNum == 0 {
			isOver = true
			winer = player
		}
	}
	if winer != nil {
		winer.Win = win
		room.LastPH = 0
		for _, player := range room.Players {
			if player.Type == Human {
				if player.ID == winer.ID {
					player.Point += winer.Win //人类赢
					playerWin = winer.Win
				} else {
					player.Point -= int64(player.CardNum) //人类输
					playerWin = int64(player.CardNum)
				}
				//计算抽水
				commission := GetGameCommission()
				if player.Point > commission {
					player.Point -= commission
				} else {
					player.Point = 0
				}
				playerPoint = player.Point
			}
		}
	}
	return isOver, winer, playerPoint, playerWin
}

// 玩一局游戏
func PlayOneGame(room *Room) {
	fmt.Printf("\n===== 房间 %s 游戏开始 =====", room.ID)

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

	// 显示玩家的牌（除了底牌）
	for _, player := range room.Players {
		//只能显示自己的牌
		if player.ID == int(Human) {
			room.Landlord = player
			cards := make([]int, 0)
			for _, card := range room.Landlord.Cards {
				cards = append(cards, card.Id)
			}
			go func() {
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
				room.MsgChan <- RoomMsg{
					Type: "showCard",
					Data: gconv.String(data),
				}
			}()
		}
		showPlayerCards(player)
	}

	// fmt.Println("\n底牌是:", showCards(room.Deck))
	// fmt.Println("地主是:", room.Landlord.Name)
	// fmt.Println("游戏开始!")

	// 显示地主的牌
	// showPlayerCards(room.Landlord)

	// 游戏主循环
	// passCount := 0
	room.Rgtimer.Start()

}

// 房间循环定时器
func (room *Room) GameLoop(ctx context.Context) {
	// 检查游戏是否结束
	over, winner, playerPoint, playerWin := isGameOver(room)
	if over {
		go func() {
			data, _ := json.Marshal(struct {
				WinName     string `json:"winName"`
				Winner      int    `json:"winner"`
				Point       int64  `json:"point"`
				Win         int64  `json:"win"`         //当次赢分
				PlayerPoint int64  `json:"playerPoint"` //玩家总瓜子数
				PlayerWin   int64  `json:"playerWin"`   //玩家本局结算分数(排除抽佣)
			}{
				WinName:     winner.Name,
				Winner:      winner.ID,
				Point:       winner.Point,
				Win:         winner.Win,
				PlayerPoint: playerPoint,
				PlayerWin:   playerWin,
			})
			room.MsgChan <- RoomMsg{
				Type: "over",
				Data: gconv.String(data),
			}
		}()
		room.Rgtimer.Stop()
		room.IsPlaying = false
		room.LastCards = make([]Card, 0)
		room.OutStarTime = 0
		room.passCount = 0
		fmt.Printf("\n游戏结束！恭喜%s！获胜\n", winner.Name)
	}

	currentPlayer := room.Players[room.Current]
	fmt.Printf("\n%s的回合 (当前手牌数: %d)\n", currentPlayer.Name, currentPlayer.CardNum)
	showPlayerCards(currentPlayer)
	fmt.Printf("上一手牌，牌型: %v, 牌是：%s\n", room.LastPH, showCards(room.LastCards))
	fmt.Println("请选择要出的牌 (输入牌的序号，用逗号分隔，0表示不出): ")

	var indices []int
	var selectedCards []Card
	if currentPlayer.Type == AI {
		// AI决策，随机秒数
		//thingTime := grand.N(1,5)
		//time.Sleep(time.Duration(thingTime) * time.Second) // 模拟思考时间
		//记录开始出牌时间
		now := int(time.Now().Unix())
		if room.OutStarTime == 0 {
			room.OutStarTime = now
		}
		thinkTime := grand.N(1, 2) //模拟思考秒数
		if now-room.OutStarTime < thinkTime {
			return
		}
		room.OutStarTime = 0
		room.LastPH, selectedCards = aiDecideCards(currentPlayer, room.Landlord, room.LastPH, room.LastCards)
		// selectedCards = getSelectedCards(currentPlayer, indices)
		if len(selectedCards) > 0 {
			for _, v := range selectedCards {
				indices = append(indices, v.Id)
			}
		}

	} else {
		//记录开始出牌时间
		now := int(time.Now().Unix())
		if room.OutStarTime == 0 {
			room.OutStarTime = now
		}
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
					fmt.Printf("输入的牌不存在%v \n", currentPlayer.OutCardIds)
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
						room.MsgChan <- RoomMsg{
							Type: "outCard",
							Data: gconv.String(data),
						}
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
						room.MsgChan <- RoomMsg{
							Type: "outCard",
							Data: gconv.String(data),
						}
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
								room.MsgChan <- RoomMsg{
									Type: "outCard",
									Data: gconv.String(data),
								}
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
								room.MsgChan <- RoomMsg{
									Type: "outCard",
									Data: gconv.String(data),
								}
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
		fmt.Printf("牌型无效，请重新选择 %v \n", selectedCards)
		code = 1

	}
	// 验证是否能压过上一手牌
	if !compareCards(room.LastPH, cardType, room.LastCards, selectedCards) {
		fmt.Println("不能压过上5一手牌,请重新选择")
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
			room.MsgChan <- RoomMsg{
				Type: "outCard",
				Data: gconv.String(data),
			}
		}()
		currentPlayer.OutCardIds = make([]int, 0)
		return
	}
	var msgType = ""
	// 如果是不出牌
	if len(selectedCards) == 0 {
		fmt.Printf("%s不出\n", currentPlayer.Name)
		//通知用户,某位机器人不出牌，
		msgType = "pass"
		room.passCount++
		// 如果三家都不出，重置上一手牌
		if room.passCount >= 3 {
			room.LastPH = 0
			room.LastCards = []Card{}
			room.passCount = 0
		}
	} else {
		// 出牌
		fmt.Printf("%s出了: %s (%v)\n", currentPlayer.Name, showCards(selectedCards), cardType)
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
	}
	room.OutStarTime = 0 //人类出牌时间恢复为0
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

	//通知用户,出牌，
	go func() {
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
		room.MsgChan <- RoomMsg{
			Type: msgType,
			Data: gconv.String(data),
		}
	}()

}

// 显示房间列表
func showRooms(rm *RoomManager) {
	fmt.Println("\n===== 房间列表 =====")
	if len(rm.Rooms) == 0 {
		fmt.Println("暂无可用房间")
		return
	}

	for id, room := range rm.Rooms {
		status := "等待中"
		if room.IsPlaying {
			status = "游戏中"
		}
		fmt.Printf("房间ID: %s, 玩家数: %d/3, 状态: %s\n", id, len(room.Players), status)
	}
}

// 开始游戏主程序
func startGame() {
	fmt.Println("欢迎来到斗地主游戏！")
	rm := NewRoomManager()

	// 创建第一个玩家（人类）
	fmt.Print("请输入你的名字: ")
	var playerName string
	fmt.Scanln(&playerName)
	humanPlayer := rm.CreatePlayer(playerName, Human)

	for {
		fmt.Println("\n===== 主菜单 =====")
		fmt.Println("1. 创建房间（与机器人对战）")
		fmt.Println("2. 加入房间（与其他玩家对战）")
		fmt.Println("3. 查看房间列表")
		fmt.Println("4. 退出游戏")

		if humanPlayer.RoomID != "" {
			room := rm.Rooms[humanPlayer.RoomID]
			fmt.Printf("\n你当前在房间 %s 中，已有 %d/3 名玩家\n", room.ID, len(room.Players))
			fmt.Println("5. 离开房间")
			if len(room.Players) == 3 && !room.IsPlaying {
				fmt.Println("6. 开始游戏")
			}
		}

		fmt.Print("请选择操作: ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if humanPlayer.RoomID != "" {
				fmt.Println("你已经在一个房间中，请先离开当前房间")
				break
			}

			// 选择机器人数量
			fmt.Print("请选择机器人数量 (1-2): ")
			var aiCountStr string
			fmt.Scanln(&aiCountStr)
			aiCount, err := strconv.Atoi(aiCountStr)
			if err != nil || aiCount < 1 || aiCount > 2 {
				fmt.Println("无效的数量，默认创建2个机器人")
				aiCount = 2
			}

			room := rm.CreateRoom(humanPlayer, aiCount)
			fmt.Printf("房间创建成功！房间ID: %s\n", room.ID)
			fmt.Printf("房间内有 %d 个机器人，共 %d 名玩家\n", aiCount, len(room.Players))
			if len(room.Players) == 3 {
				fmt.Println("可以选择开始游戏了")
			}

		case "2":
			if humanPlayer.RoomID != "" {
				fmt.Println("你已经在一个房间中，请先离开当前房间")
				break
			}
			fmt.Print("请输入房间ID: ")
			var roomID string
			fmt.Scanln(&roomID)
			success := rm.JoinRoom(humanPlayer, roomID)
			if success {
				room := rm.Rooms[roomID]
				fmt.Printf("成功加入房间 %s！当前房间有 %d/3 名玩家\n", roomID, len(room.Players))
			} else {
				fmt.Println("加入房间失败，房间不存在或已满或正在游戏中")
			}

		case "3":
			showRooms(rm)

		case "4":
			// 退出前离开房间
			if humanPlayer.RoomID != "" {
				rm.LeaveRoom(humanPlayer)
			}
			fmt.Println("谢谢游玩，再见！")
			return

		case "5":
			if humanPlayer.RoomID == "" {
				fmt.Println("你不在任何房间中")
				break
			}
			rm.LeaveRoom(humanPlayer)
			fmt.Println("已离开房间")

		case "6":
			if humanPlayer.RoomID == "" {
				fmt.Println("你不在任何房间中")
				break
			}
			room := rm.Rooms[humanPlayer.RoomID]
			if len(room.Players) != 3 {
				fmt.Println("玩家不足3人，无法开始游戏")
				break
			}
			if room.IsPlaying {
				fmt.Println("游戏已经开始")
				break
			}
			room.IsPlaying = true
			PlayOneGame(room)

		default:
			fmt.Println("无效的操作，请重新选择")
		}
	}
}
