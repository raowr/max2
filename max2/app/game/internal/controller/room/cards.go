package room

import (
	"sort"
)

// 按牌面排序
func (ps *Player) sortByRank() {
	sort.Slice(ps.Cards, func(i, j int) bool {
		return ps.Cards[i].Rank < ps.Cards[j].Rank
	})
}

// 按花色排序
func (ps *Player) sortBySuit() {
	sort.Slice(ps.Cards, func(i, j int) bool {
		if ps.Cards[i].Suit == ps.Cards[j].Suit {
			return ps.Cards[i].Rank < ps.Cards[j].Rank
		}
		return ps.Cards[i].Suit < ps.Cards[j].Suit
	})
}

// 统计牌面数量
func (ps *Player) countRanks() map[int]int {
	count := make(map[int]int)
	for _, card := range ps.Cards {
		count[card.Rank]++
	}
	return count
}

// 统计花色数量
func (ps *Player) countSuits() map[string]int {
	count := make(map[string]int)
	for _, card := range ps.Cards {
		count[card.Suit]++
	}
	return count
}

// 查找同花顺
func (ps *Player) findFlushes() {
	ps.sortBySuit()

	suits := ps.countSuits()
	for currentSuit, count := range suits {
		if count >= 5 {
			// 收集该花色的所有牌
			var suitedCards []Card
			for _, card := range ps.Cards {
				if card.Suit == currentSuit {
					suitedCards = append(suitedCards, card)
				}
			}

			// 查找顺子
			ps.findStraightInSuit(suitedCards, currentSuit)
		}
	}
}

func (ps *Player) findStraightInSuit(cards []Card, cardSuit string) {
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
func (ps *Player) findFours() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 4 {
			// 收集4张相同牌面的牌
			var fourCards []Card
			var otherCards []Card

			for i := 0; i < len(ps.Cards); i++ {
				if ps.Cards[i].Rank == rank {
					fourCards = append(fourCards, ps.Cards[i])
				} else {
					otherCards = append(otherCards, ps.Cards[i])
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
func (ps *Player) findThrees() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 3 {
			// 收集3张相同牌面的牌
			var threeCards []Card
			var pairCards []Card

			for _, card := range ps.Cards {
				if card.Rank == rank {
					threeCards = append(threeCards, card)
				}
			}

			// 查找对子
			for r, c := range rankCount {
				if c >= 2 && r != rank {
					pairCards = nil
					pairCount := 0
					for _, card := range ps.Cards {
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
func (ps *Player) findSuits() {
	suits := ps.countSuits()

	for currentSuit, count := range suits {
		if count >= 5 {
			var suitedCards []Card
			for _, card := range ps.Cards {
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
func (ps *Player) findStraights() {
	if len(ps.Cards) < 5 {
		return
	}

	// 去重并排序
	uniqueRanks := make(map[int]bool)
	var sortedRanks []int
	for _, card := range ps.Cards {
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
				for _, card := range ps.Cards {
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
func (ps *Player) findPairs() {
	rankCount := ps.countRanks()

	for rank, count := range rankCount {
		if count == 2 {
			var pairCards []Card
			for _, card := range ps.Cards {
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
func (ps *Player) findSingles() {
	// 统计2的数量
	twoCount := 0
	for _, card := range ps.Cards {
		if card.Rank == 15 { // 2的Rank为15
			twoCount++
		}
	}

	// 如果2的数量大于单牌数量，则2不参与单牌
	if twoCount > 0 {
		var singles []Card
		var twos []Card

		// 分离2和其他单牌
		for _, card := range ps.Cards {
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
		for _, card := range ps.Cards {
			ps.handPattern[SINGLE] = append(ps.handPattern[SINGLE], []Card{card})
		}
	}
}

// 移除已使用的牌
func (ps *Player) removeCards(toRemove []Card) {
	var newCards []Card
	removeMap := make(map[Card]bool)

	for _, card := range toRemove {
		removeMap[card] = true
	}

	for _, card := range ps.Cards {
		if !removeMap[card] {
			newCards = append(newCards, card)
		}
	}

	ps.Cards = newCards
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
func (ps *Player) Solve() map[int][][]Card {
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

// 判断两张牌的牌型
func judgeTwoCards(cards []Card) int {
	if len(cards) != 2 {
		return 0
	}

	// 检查是否是对子
	if cards[0].Rank == cards[1].Rank {
		return PAIR
	}

	return 0 // 两张不同牌面的牌不是有效牌型
}

// 判断五张牌的牌型
func judgeFiveCards(cards []Card) int {
	if len(cards) != 5 {
		return 0
	}

	// 检查同花顺
	if isFlush(cards) && isStraight(cards) {
		return FLUSH
	}

	// 检查四条
	if isFourOfAKind(cards) {
		return FOUR
	}

	// 检查葫芦（3带2）
	if isFullHouse(cards) {
		return THREE
	}

	// 检查同花
	if isFlush(cards) {
		return SUIT
	}

	// 检查顺子
	if isStraight(cards) {
		return STRAIGHT
	}

	return 0 // 不是有效牌型
}

// 统计牌面数量
func countRanks(cards []Card) map[int]int {
	count := make(map[int]int)
	for _, card := range cards {
		count[card.Rank]++
	}
	return count
}

// 统计花色数量
func countSuits(cards []Card) map[string]int {
	count := make(map[string]int)
	for _, card := range cards {
		count[card.Suit]++
	}
	return count
}

// 判断是否是同花
func isFlush(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}

	suits := countSuits(cards)
	for _, count := range suits {
		if count >= 5 {
			// 检查是否至少有5张同花色
			firstSuit := cards[0].Suit
			for _, card := range cards {
				if card.Suit != firstSuit {
					return false
				}
			}
			return true
		}
	}
	return false
}

// 判断是否是顺子
func isStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}

	// 获取唯一的牌面值并排序
	uniqueRanks := make(map[int]bool)
	var ranks []int
	for _, card := range cards {
		if !uniqueRanks[card.Rank] {
			uniqueRanks[card.Rank] = true
			ranks = append(ranks, card.Rank)
		}
	}

	if len(ranks) < 5 {
		return false
	}

	sort.Ints(ranks)

	// 检查连续5张牌
	for i := 0; i <= len(ranks)-5; i++ {
		isConsecutive := true
		for j := 0; j < 4; j++ {
			if ranks[i+j+1]-ranks[i+j] != 1 {
				isConsecutive = false
				break
			}
		}
		if isConsecutive {
			return true
		}
	}

	return false
}

// 判断是否是四条
func isFourOfAKind(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}

	rankCount := countRanks(cards)
	hasFour := false
	hasOne := false

	for _, count := range rankCount {
		if count == 4 {
			hasFour = true
		} else if count == 1 {
			hasOne = true
		}
	}

	return hasFour && hasOne
}

// 判断是否是葫芦（3带2）
func isFullHouse(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}

	rankCount := countRanks(cards)
	hasThree := false
	hasTwo := false

	for _, count := range rankCount {
		if count == 3 {
			hasThree = true
		} else if count == 2 {
			hasTwo = true
		}
	}

	return hasThree && hasTwo
}
