package room

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

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
	Value int    // 牌值 3-15，分别对应3-A
	Suit  int    // 花色 0-3 方块、梅花、红桃、黑桃
	Name  string // 牌的名称
}

// 玩家类型
type PlayerType int

const (
	Human PlayerType = iota
	AI
)

const (
	OutCard = 10 //出牌最大时间
)

// 玩家结构体
type Player struct {
	ID       int
	Cards    []Card
	Name     string
	RoomID   string     // 玩家所在房间ID
	Type     PlayerType // 玩家类型（人类或AI）
	OutCards []Card     //玩家单次打出的牌
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
	Turn        int    // 轮次
	IsPlaying   bool   // 房间是否正在游戏中
	MsgChan     chan RoomMsg
	Rgtimer     *gtimer.Timer
	OutStarTime int //出牌开始时间
}

// 房间管理器
type RoomManager struct {
	Rooms        map[string]*Room
	PlayerList   map[int]*Player
	nextPlayerID int
}

// 创建新的房间管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		Rooms:        make(map[string]*Room),
		PlayerList:   make(map[int]*Player),
		nextPlayerID: 0,
	}
}

// 创建新玩家
func (rm *RoomManager) CreatePlayer(name string, playerType PlayerType) *Player {
	player := &Player{
		ID:   rm.nextPlayerID,
		Name: name,
		Type: playerType,
	}
	rm.PlayerList[player.ID] = player
	rm.nextPlayerID++
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

	// 添加AI机器人
	for i := 0; i < aiCount; i++ {
		aiName := fmt.Sprintf("机器人%d号", i+1)
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
	suits := []string{"方块", "梅花", "红桃", "黑桃"}
	values := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}

	for suit, suitName := range suits {
		for i, valueName := range values {
			card := Card{
				Value: 3 + i,
				Suit:  suit,
				Name:  suitName + valueName,
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
	}
}

// 发牌
func dealCards(deck []Card, players []*Player) {
	//每位玩家13张牌
	for i := 0; i < len(deck); i++ {
		players[i%4].Cards = append(players[i%4].Cards, deck[i])
	}

	// 对每个玩家的牌进行排序
	for _, player := range players {
		sort.Slice(player.Cards, func(i, j int) bool {
			return player.Cards[i].Value < player.Cards[j].Value
		})
	}

}

// 显示玩家的牌
func showPlayerCards(player *Player) {
	fmt.Printf("%s的牌: ", player.Name)
	for i, card := range player.Cards {
		fmt.Printf("%d.%s ", i+1, card.Name)
	}
	fmt.Println()
}

// AI叫地主决策
func aiDecideBidLandlord(player *Player) bool {
	// 简单策略：根据手中大牌数量决定是否叫地主
	highCards := 0
	for _, card := range player.Cards {
		if card.Value >= 14 || card.Value == 16 || card.Value == 17 {
			highCards++
		}
	}

	// 有3张以上大牌则叫地主，否则有30%概率叫地主
	if highCards >= 3 {
		return true
	}
	return rand.Float64() < 0.3
}

// 确定谁先出牌
func bidLandlord(room *Room) {
	fmt.Println("\n确定谁先出牌...")
	currentPlayer := 0
	for _, v := range room.Players {
		for _, card := range v.Cards {
			if card.Suit == 0 && card.Value == 3 {
				currentPlayer = v.ID
			}
		}
	}
	// 地主先出牌
	room.Current = currentPlayer
}

// 解析玩家输入的牌索引
func parseCardIndices(input string, maxIndex int) []int {
	if input == "0" {
		return []int{} // 不出牌
	}

	parts := strings.Split(input, ",")
	var indices []int

	for _, part := range parts {
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 || idx > maxIndex {
			return nil // 无效输入
		}
		indices = append(indices, idx-1) // 转换为0-based索引
	}

	// 去重并排序
	uniqueIndices := make(map[int]bool)
	for _, idx := range indices {
		uniqueIndices[idx] = true
	}
	indices = []int{}
	for idx := range uniqueIndices {
		indices = append(indices, idx)
	}
	sort.Ints(indices)

	return indices
}

// 从玩家手中获取选中的牌
func getSelectedCards(player *Player, indices []int) []Card {
	var selected []Card
	for _, idx := range indices {
		selected = append(selected, player.Cards[idx])
	}
	return selected
}

// 从玩家手中移除牌
func removeCards(player *Player, indices []int) {
	// 按降序排序索引，以便从后往前删除
	sort.Sort(sort.Reverse(sort.IntSlice(indices)))

	for _, idx := range indices {
		// 移除指定索引的牌
		player.Cards = append(player.Cards[:idx], player.Cards[idx+1:]...)
	}
}

// 判断牌型是否有效
func isValidCardType(cards []Card) (string, bool) {
	if len(cards) == 0 {
		return "pass", true // 不出牌
	}

	// 单牌
	if len(cards) == 1 {
		return "single", true
	}

	// 对子
	if len(cards) == 2 {
		if cards[0].Value == cards[1].Value {
			return "pair", true
		}
		return "", false
	}

	//福禄(3带2)

	//4条(4带1)

	//顺子(5张点数连续的牌)

	//花色(5张相同的花色)

	//同花顺
	if len(cards) == 5 {

		// 检查是否连续
		values := make([]int, len(cards))
		for i, card := range cards {
			values[i] = card.Value
		}
		sort.Ints(values)

		for i := 1; i < len(values); i++ {
			if values[i] != values[i-1]+1 {
				return "", false
			}
		}
		return "straight", true
	}

	// 其他牌型可以在这里继续实现...
	return "", false
}

// 比较牌的大小
func compareCards(last []Card, current []Card) bool {
	if len(last) == 0 {
		return true // 上一手没牌，当前任何有效牌型都可以出
	}

	if len(current) == 0 {
		return true // 不出牌
	}

	lastType, _ := isValidCardType(last)
	currentType, _ := isValidCardType(current)

	// 相同牌型且长度相同才能比较
	if lastType == currentType && len(last) == len(current) {
		if getMaxValue(current) > getMaxValue(last) {
			return true
		} else if getMaxValue(current) == getMaxValue(last) {
			//比较花色
			return getMaxSuit(current) > getMaxSuit(last)
		}
	}

	return false
}

// 获取牌组中的最大值
func getMaxValue(cards []Card) int {
	maxVal := 0
	for _, card := range cards {
		if card.Value > maxVal {
			maxVal = card.Value
		}
	}
	//相同点数牌，继续判断花色
	return maxVal
}

// 获取牌组中的最大花色
func getMaxSuit(cards []Card) int {
	maxVal := 0
	for _, card := range cards {
		if card.Suit > maxVal {
			maxVal = card.Suit
		}
	}
	//相同点数牌，继续判断花色
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
func aiDecideCards(player *Player, lastCards []Card) []int {
	// 简单AI策略：尝试找出能压过上一手的最小牌组
	// 如果上一手没牌，尝试出最小的牌组

	// 生成所有可能的有效牌组合
	possiblePlays := generateAllPossiblePlays(player.Cards)

	// 过滤掉不能压过上一手的牌
	validPlays := []struct {
		indices []int
		cards   []Card
	}{}

	for _, play := range possiblePlays {
		if compareCards(lastCards, play.cards) {
			validPlays = append(validPlays, play)
		}
	}

	// 如果没有可以出的牌，返回空（不出）
	if len(validPlays) == 0 {
		return []int{}
	}

	// 简单策略：优先出最小的牌组
	bestPlay := validPlays[0]
	for _, play := range validPlays[1:] {
		if isBetterPlay(bestPlay.cards, play.cards) {
			bestPlay = play
		}
	}

	return bestPlay.indices
}

// 判断是否是更好的出牌（用于AI决策）
func isBetterPlay(currentBest, candidate []Card) bool {
	// 简单策略：牌数少的更好；牌数相同则最大值小的更好
	if len(candidate) < len(currentBest) {
		return true
	}
	if len(candidate) > len(currentBest) {
		return false
	}
	return getMaxValue(candidate) < getMaxValue(currentBest)
}

// 生成所有可能的有效牌组合（简化版）
func generateAllPossiblePlays(cards []Card) []struct {
	indices []int
	cards   []Card
} {
	var plays []struct {
		indices []int
		cards   []Card
	}

	// 添加不出牌选项
	plays = append(plays, struct {
		indices []int
		cards   []Card
	}{[]int{}, []Card{}})

	// 单牌
	for i, card := range cards {
		plays = append(plays, struct {
			indices []int
			cards   []Card
		}{[]int{i}, []Card{card}})
	}

	// 对子
	for i := 0; i < len(cards); i++ {
		for j := i + 1; j < len(cards); j++ {
			if cards[i].Value == cards[j].Value {
				plays = append(plays, struct {
					indices []int
					cards   []Card
				}{[]int{i, j}, []Card{cards[i], cards[j]}})
			}
		}
	}

	// 这里可以扩展更多牌型的生成逻辑

	//福禄(3带2)

	//4条(4带1)

	//顺子(5张点数连续的牌)

	//花色(5张相同的花色)

	//同花顺

	return plays
}

// 检查游戏是否结束
func isGameOver(room *Room) (bool, *Player) {
	for _, player := range room.Players {
		if len(player.Cards) == 0 {
			return true, player
		}
	}
	return false, nil
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

	// 显示玩家的牌（除了底牌）
	for _, player := range room.Players {
		//只能显示自己的牌
		if player.ID == int(Human) {
			data, _ := json.Marshal(player.Cards)
			go func() {
				room.MsgChan <- RoomMsg{
					Type: "showCard",
					Data: gconv.String(data),
				}
			}()

		}
		showPlayerCards(player)
	}

	// 确定谁先出牌，既方块3在谁那
	bidLandlord(room)

	// fmt.Println("\n底牌是:", showCards(room.Deck))
	// fmt.Println("地主是:", room.Landlord.Name)
	// fmt.Println("游戏开始!")

	// 显示地主的牌
	// showPlayerCards(room.Landlord)

	// 游戏主循环
	// passCount := 0
	room.Rgtimer.Add(context.Background(), 100*time.Millisecond, room.GameLoop)

}

var passCount = 0

// 房间循环定时器
func (room *Room) GameLoop(ctx context.Context) {
	// 检查游戏是否结束
	over, winner := isGameOver(room)
	if over {
		if winner.ID == room.Landlord.ID {
			fmt.Printf("\n游戏结束！地主%s获胜！\n", winner.Name)
		} else {
			fmt.Printf("\n游戏结束！农民获胜！恭喜%s！\n", winner.Name)
		}
		room.Rgtimer.Close()
		room.IsPlaying = false
	}

	currentPlayer := room.Players[room.Current]
	fmt.Printf("\n%s的回合 (当前手牌数: %d)\n", currentPlayer.Name, len(currentPlayer.Cards))
	showPlayerCards(currentPlayer)
	fmt.Printf("上一手牌: %s\n", showCards(room.LastCards))
	fmt.Print("请选择要出的牌 (输入牌的序号，用逗号分隔，0表示不出): ")

	var input string
	var indices []int
	var selectedCards []Card
	if currentPlayer.Type == AI {
		// AI决策
		time.Sleep(1 * time.Second) // 模拟思考时间
		indices = aiDecideCards(currentPlayer, room.LastCards)
		selectedCards = getSelectedCards(currentPlayer, indices)
		// 转换为输入格式
		if len(indices) == 0 {
			input = "0"
		} else {
			parts := make([]string, len(indices))
			for i, idx := range indices {
				parts[i] = strconv.Itoa(idx + 1) // 转换为1-based索引
			}
			input = strings.Join(parts, ",")
		}
		fmt.Println(input)
	} else {
		//记录开始出牌时间
		now := int(time.Now().Unix())
		if room.OutStarTime == 0 {
			room.OutStarTime = now
		}
		if len(currentPlayer.OutCards) <= 0 && now-room.OutStarTime < OutCard {
			return
		} else {
			if len(currentPlayer.OutCards) > 0 {
				g.Log().Debug(ctx, currentPlayer.OutCards)
				selectedCards = currentPlayer.OutCards
			}

			// selectedCards = getSelectedCards(currentPlayer, indices)
		}

		if now-room.OutStarTime >= OutCard {
			indices = make([]int, 0) //超过出牌时间过
		}

	}
	// 验证牌型
	cardType, valid := isValidCardType(selectedCards)
	if !valid {
		fmt.Println("牌型无效，请重新选择")
		return
	}
	// 验证是否能压过上一手牌
	if !compareCards(room.LastCards, selectedCards) {
		fmt.Println("不能压过上一手牌，请重新选择")
		return
	}

	// 如果是不出牌
	if len(selectedCards) == 0 {
		fmt.Printf("%s不出\n", currentPlayer.Name)
		//通知用户,某位机器人不出牌，
		if currentPlayer.Type == AI {
			go func() {
				room.MsgChan <- RoomMsg{
					Type: "pass",
					Data: gconv.String(currentPlayer.ID),
				}
			}()
		}
		passCount++
		// 如果三家都不出，重置上一手牌
		if passCount >= 2 {
			room.LastCards = []Card{}
			passCount = 0
		}
	} else {
		// 出牌
		fmt.Printf("%s出了: %s (%s)\n", currentPlayer.Name, showCards(selectedCards), cardType)
		//通知用户,某位机器人不出牌，
		if currentPlayer.Type == AI {
			go func() {
				data, _ := json.Marshal(struct {
					Pid      int    `json:"pid"`
					Cards    []Card `json:"cards"`
					CardType string `json:"card_type"`
				}{
					Pid:      currentPlayer.ID,
					Cards:    selectedCards,
					CardType: cardType,
				})
				room.MsgChan <- RoomMsg{
					Type: "pass",
					Data: gconv.String(data),
				}
			}()
		}
		removeCards(currentPlayer, indices)
		room.LastCards = selectedCards
		passCount = 0
	}

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
