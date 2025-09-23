package room

import (
	"fmt"
	"sort"
)

type Card struct {
	Suit string // 花色: ♠, ♥, ♦, ♣
	Rank int    // 牌面: 3-14 (3-10, 11=J, 12=Q, 13=K, 14=A, 15=2)
}

const (
	SINGLE   = iota + 1 // 单牌
	PAIR                // 对子
	STRAIGHT            // 顺子
	SUIT                // 花色
	THREE               // 3带2
	FOUR                // 4带1
	FLUSH               // 同花顺
)

type PokerSolver struct {
	cards       []Card
	handPattern map[int][][]Card
}

func NewPokerSolver(cards []Card) *PokerSolver {
	return &PokerSolver{
		cards:       cards,
		handPattern: make(map[int][][]Card),
	}
}

// 按牌面排序
func (ps *PokerSolver) sortByRank() {
	sort.Slice(ps.cards, func(i, j int) bool {
		return ps.cards[i].Rank < ps.cards[j].Rank
	})
}

// 按花色排序
func (ps *PokerSolver) sortBySuit() {
	sort.Slice(ps.cards, func(i, j int) bool {
		if ps.cards[i].Suit == ps.cards[j].Suit {
			return ps.cards[i].Rank < ps.cards[j].Rank
		}
		return ps.cards[i].Suit < ps.cards[j].Suit
	})
}

// 统计牌面数量
func (ps *PokerSolver) countRanks() map[int]int {
	count := make(map[int]int)
	for _, card := range ps.cards {
		count[card.Rank]++
	}
	return count
}

// 统计花色数量
func (ps *PokerSolver) countSuits() map[string]int {
	count := make(map[string]int)
	for _, card := range ps.cards {
		count[card.Suit]++
	}
	return count
}

// 查找同花顺
func (ps *PokerSolver) findFlushes() {
	ps.sortBySuit()

	suits := ps.countSuits()
	for currentSuit, count := range suits {
		if count >= 5 {
			// 收集该花色的所有牌
			var suitedCards []Card
			for _, card := range ps.cards {
				if card.Suit == currentSuit {
					suitedCards = append(suitedCards, card)
				}
			}

			// 查找顺子
			ps.findStraightInSuit(suitedCards, currentSuit)
		}
	}
}

func (ps *PokerSolver) findStraightInSuit(cards []Card, cardSuit string) {
	if len(cards) < 5 {
		return
	}

	// 去重并排序
	uniqueRanks := make(map[int]bool)
	for _, card := range cards {
		uniqueRanks[card.Rank] = true
	}

	// 检查顺子
	for start := 3; start <= 10; start++ {
		isStraightValid := true
		var straightCards []Card

		for i := 0; i < 5; i++ {
			rank := start + i
			if !uniqueRanks[rank] {
				isStraightValid = false
				break
			}
			// 找到对应花色的牌
			for _, card := range cards {
				if card.Rank == rank && card.Suit == cardSuit {
					straightCards = append(straightCards, card)
					break
				}
			}
		}

		if isStraightValid && len(straightCards) == 5 {
			ps.handPattern[FLUSH] = append(ps.handPattern[FLUSH], straightCards)
			// 移除已使用的牌
			ps.removeCards(straightCards)
		}
	}
}

// 查找4带1
func (ps *PokerSolver) findFours() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 4 {
			// 收集4张相同牌面的牌
			var fourCards []Card
			var otherCards []Card

			for i := 0; i < len(ps.cards); i++ {
				if ps.cards[i].Rank == rank {
					fourCards = append(fourCards, ps.cards[i])
				} else {
					otherCards = append(otherCards, ps.cards[i])
				}
			}

			if len(otherCards) >= 1 {
				// 取一张其他牌作为搭配
				combination := append(fourCards, otherCards[0])
				ps.handPattern[FOUR] = append(ps.handPattern[FOUR], combination)
				ps.removeCards(combination)
			}
		}
	}
}

// 查找3带2
func (ps *PokerSolver) findThrees() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 3 {
			// 收集3张相同牌面的牌
			var threeCards []Card
			var pairCards []Card

			for _, card := range ps.cards {
				if card.Rank == rank {
					threeCards = append(threeCards, card)
				}
			}

			// 查找对子
			for r, c := range rankCount {
				if c >= 2 && r != rank {
					pairCards = nil
					pairCount := 0
					for _, card := range ps.cards {
						if card.Rank == r && pairCount < 2 {
							pairCards = append(pairCards, card)
							pairCount++
						}
					}
					if len(pairCards) == 2 {
						combination := append(threeCards, pairCards...)
						ps.handPattern[THREE] = append(ps.handPattern[THREE], combination)
						ps.removeCards(combination)
						break
					}
				}
			}
		}
	}
}

// 查找同花
func (ps *PokerSolver) findSuits() {
	suits := ps.countSuits()

	for currentSuit, count := range suits {
		if count >= 5 {
			var suitedCards []Card
			for _, card := range ps.cards {
				if card.Suit == currentSuit {
					suitedCards = append(suitedCards, card)
				}
			}

			if len(suitedCards) >= 5 {
				// 取最大的5张同花
				sort.Slice(suitedCards, func(i, j int) bool {
					return suitedCards[i].Rank > suitedCards[j].Rank
				})
				combination := suitedCards[:5]
				ps.handPattern[SUIT] = append(ps.handPattern[SUIT], combination)
				ps.removeCards(combination)
			}
		}
	}
}

// 查找顺子
func (ps *PokerSolver) findStraights() {
	if len(ps.cards) < 5 {
		return
	}

	// 去重并排序
	uniqueRanks := make(map[int]bool)
	var sortedRanks []int
	for _, card := range ps.cards {
		if !uniqueRanks[card.Rank] {
			uniqueRanks[card.Rank] = true
			sortedRanks = append(sortedRanks, card.Rank)
		}
	}
	sort.Ints(sortedRanks)

	// 检查顺子
	for i := 0; i <= len(sortedRanks)-5; i++ {
		isStraightValid := true
		for j := 0; j < 4; j++ {
			if sortedRanks[i+j+1]-sortedRanks[i+j] != 1 {
				isStraightValid = false
				break
			}
		}

		if isStraightValid {
			var straightCards []Card
			usedRanks := make(map[int]bool)

			for rank := sortedRanks[i]; rank <= sortedRanks[i+4]; rank++ {
				usedRanks[rank] = true
				// 找到第一张该牌面的牌
				for _, card := range ps.cards {
					if card.Rank == rank && !containsCard(straightCards, card) {
						straightCards = append(straightCards, card)
						break
					}
				}
			}

			if len(straightCards) == 5 {
				ps.handPattern[STRAIGHT] = append(ps.handPattern[STRAIGHT], straightCards)
				ps.removeCards(straightCards)
			}
		}
	}
}

// 查找对子
func (ps *PokerSolver) findPairs() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 2 {
			var pairCards []Card
			for _, card := range ps.cards {
				if card.Rank == rank {
					pairCards = append(pairCards, card)
				}
			}

			if len(pairCards) == 2 {
				ps.handPattern[PAIR] = append(ps.handPattern[PAIR], pairCards)
				ps.removeCards(pairCards)
			}
		}
	}
}

// 处理单牌（考虑2的数量）
func (ps *PokerSolver) findSingles() {
	// 统计2的数量
	twoCount := 0
	for _, card := range ps.cards {
		if card.Rank == 15 { // 2的Rank为15
			twoCount++
		}
	}

	// 如果2的数量大于单牌数量，则2不参与单牌
	if twoCount > 0 {
		var singles []Card
		var twos []Card

		// 分离2和其他单牌
		for _, card := range ps.cards {
			if card.Rank == 15 {
				twos = append(twos, card)
			} else {
				singles = append(singles, card)
			}
		}

		// 添加单牌
		for _, card := range singles {
			ps.handPattern[SINGLE] = append(ps.handPattern[SINGLE], []Card{card})
		}

		// 只有当2的数量小于等于单牌数量时，2才作为单牌
		if len(twos) <= len(singles) {
			for _, card := range twos {
				ps.handPattern[SINGLE] = append(ps.handPattern[SINGLE], []Card{card})
			}
		}
	} else {
		// 没有2，所有剩余牌都是单牌
		for _, card := range ps.cards {
			ps.handPattern[SINGLE] = append(ps.handPattern[SINGLE], []Card{card})
		}
	}
}

// 移除已使用的牌
func (ps *PokerSolver) removeCards(toRemove []Card) {
	var newCards []Card
	removeMap := make(map[Card]bool)

	for _, card := range toRemove {
		removeMap[card] = true
	}

	for _, card := range ps.cards {
		if !removeMap[card] {
			newCards = append(newCards, card)
		}
	}

	ps.cards = newCards
}

// 检查牌是否在数组中
func containsCard(cards []Card, target Card) bool {
	for _, card := range cards {
		if card.Suit == target.Suit && card.Rank == target.Rank {
			return true
		}
	}
	return false
}

// 主要处理函数
func (ps *PokerSolver) Solve() map[int][][]Card {
	// 按照优先级顺序查找牌型
	ps.findFlushes()   // 同花顺
	ps.findFours()     // 4带1
	ps.findThrees()    // 3带2
	ps.findSuits()     // 同花
	ps.findStraights() // 顺子
	ps.findPairs()     // 对子
	ps.findSingles()   // 单牌

	return ps.handPattern
}

// 打印结果
func (ps *PokerSolver) PrintResult() {
	patterns := []struct {
		name string
		code int
	}{
		{"同花顺", FLUSH},
		{"4带1", FOUR},
		{"3带2", THREE},
		{"同花", SUIT},
		{"顺子", STRAIGHT},
		{"对子", PAIR},
		{"单牌", SINGLE},
	}

	for _, pattern := range patterns {
		if hands, exists := ps.handPattern[pattern.code]; exists && len(hands) > 0 {
			fmt.Printf("%s (%d组):\n", pattern.name, len(hands))
			for i, hand := range hands {
				fmt.Printf("  第%d组: ", i+1)
				for _, card := range hand {
					fmt.Printf("%s%s ", card.Suit, rankToString(card.Rank))
				}
				fmt.Println()
			}
		}
	}

	fmt.Printf("\n剩余牌数: %d\n", len(ps.cards))
}

// 辅助函数：将数字转换为牌面显示
func rankToString(rank int) string {
	switch rank {
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	case 14:
		return "A"
	case 15:
		return "2"
	default:
		return fmt.Sprintf("%d", rank)
	}
}

// 示例使用
func main() {
	// 创建示例牌组
	cards := []Card{
		{"♠", 3}, {"♥", 3}, {"♦", 3}, {"♣", 4},
		{"♠", 4}, {"♥", 5}, {"♦", 5}, {"♣", 6},
		{"♠", 6}, {"♥", 7}, {"♦", 7}, {"♣", 8},
		{"♠", 15}, // 2
	}

	solver := NewPokerSolver(cards)
	result := solver.Solve()

	fmt.Println("扑克牌型分析结果:")
	solver.PrintResult()

	// 输出整理后的handPattern
	fmt.Println("\nhandPattern内容:")
	for patternType, hands := range result {
		fmt.Printf("牌型%d: ", patternType)
		for i, hand := range hands {
			fmt.Printf("第%d组[", i+1)
			for j, card := range hand {
				fmt.Printf("%s%s", card.Suit, rankToString(card.Rank))
				if j < len(hand)-1 {
					fmt.Printf(",")
				}
			}
			fmt.Printf("] ")
		}
		fmt.Println()
	}
}
