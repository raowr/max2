package room

import (
	"sort"
	"strings"
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
func (ps *Player) countSuits() map[int]int {
	count := make(map[int]int)
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

func (ps *Player) findStraightInSuit(cards []Card, cardSuit int) {
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
	// 外层循环：持续查找直到没有更多3带2组合
	for {
		// 关键修复：每次循环重新计算当前剩余牌的牌面数量（动态更新）
		rankCount := ps.countRanks()
		foundCombination := false // 标记是否找到组合

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

				// 步骤1：优先查找真正的对子（牌面数量为2的牌）
				foundPair := false
				for r, c := range rankCount {
					if c == 2 && r != rank { // 只匹配数量恰好为2的牌
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
							foundPair = true
							foundCombination = true // 标记找到组合
							break
						}
					}
				}

				// 步骤2：若没有真正的对子，从其他3张相同牌面的牌中取2张作为对子
				if !foundPair {
					for r, c := range rankCount {
						if c >= 3 && r != rank { // 从其他3+张的牌中取2张
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
								foundCombination = true // 标记找到组合
								break
							}
						}
					}
				}

				// 找到一个组合后，跳出当前循环，重新计算rankCount
				if foundCombination {
					break
				}
			}
		}

		// 若本轮未找到任何组合，退出外层循环
		if !foundCombination {
			break
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

// 查找顺子 - 2不参与任何顺子，顺子最大到A
func (ps *Player) findStraights() {
	if len(ps.Cards) < 5 {
		return
	}

	// 去重、排除2并排序
	uniqueRanks := make(map[int]bool)
	var sortedRanks []int
	twoRank := 15 // 2的Rank值

	for _, card := range ps.Cards {
		// 完全排除2，不参与顺子计算
		if card.Rank == twoRank {
			continue
		}

		if !uniqueRanks[card.Rank] {
			uniqueRanks[card.Rank] = true
			sortedRanks = append(sortedRanks, card.Rank)
		}
	}
	sort.Ints(sortedRanks)

	// 最大顺子只能到A(14)
	maxRank := 14

	// 检查顺子
	for i := 0; i <= len(sortedRanks)-5; i++ {
		// 如果顺子中最大的牌超过A，则跳过
		if sortedRanks[i+4] > maxRank {
			continue
		}

		isStraightValid := true
		// 检查是否连续递增
		for j := 0; j < 4; j++ {
			if sortedRanks[i+j+1]-sortedRanks[i+j] != 1 {
				isStraightValid = false
				break
			}
		}

		if isStraightValid {
			var straightCards []Card

			for rank := sortedRanks[i]; rank <= sortedRanks[i+4]; rank++ {
				// 找到对应rank的牌
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
	//所有剩余牌都是单牌
	for _, card := range ps.Cards {
		ps.handPattern[SINGLE] = append(ps.handPattern[SINGLE], []Card{card})
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
func countSuits(cards []Card) map[int]int {
	count := make(map[int]int)
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

// 判断是否是顺子 - 2不能参与顺子，顺子最大到A
func isStraight(cards []Card) bool {
	if len(cards) < 5 {
		return false
	}

	// 2的Rank值（根据原有定义，2的Rank是15）
	twoRank := 15
	// A的Rank值（根据原有定义，A的Rank是14）
	maxValidRank := 14

	// 获取唯一的牌面值并排除2，同时收集有效牌面
	uniqueRanks := make(map[int]bool)
	var validRanks []int
	for _, card := range cards {
		// 排除2，不参与顺子判断
		if card.Rank == twoRank {
			continue
		}
		// 只保留不超过A的牌面
		if card.Rank > maxValidRank {
			continue
		}
		if !uniqueRanks[card.Rank] {
			uniqueRanks[card.Rank] = true
			validRanks = append(validRanks, card.Rank)
		}
	}

	// 有效牌面不足5张，无法组成顺子
	if len(validRanks) < 5 {
		return false
	}

	// 对有效牌面排序
	sort.Ints(validRanks)

	// 检查是否存在连续的5张牌
	for i := 0; i <= len(validRanks)-5; i++ {
		isConsecutive := true
		for j := 0; j < 4; j++ {
			if validRanks[i+j+1]-validRanks[i+j] != 1 {
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

// 找到所有可能出现对子中最小的对子
func FindSmallestPossiblePair(lastPH int, lastCards, cards []Card, pairType PairType) []Card {
	if len(cards) < 2 {
		return nil
	}

	// 按牌的点数分组
	rankGroups := make(map[int][]Card)
	for _, card := range cards {
		rankGroups[card.Rank] = append(rankGroups[card.Rank], card)
	}

	// 收集所有可能的对子
	var allPairs [][]Card
	var gtlastCards [][]Card //大于上一手牌的对子

	for _, group := range rankGroups {
		if len(group) >= 2 {
			// 对每组牌按花色从大到小排序（黑桃3>红桃2>梅花1>方块0）
			sortCardsBySuitDesc(group)

			// 生成所有可能的对子组合
			for i := 0; i < len(group); i++ {
				for j := i + 1; j < len(group); j++ {
					pair := []Card{group[i], group[j]}
					allPairs = append(allPairs, pair)
				}
			}
		}
	}

	if len(allPairs) == 0 {
		return nil
	}

	for _, pair := range allPairs {
		if compareCards(lastPH, lastPH, lastCards, pair) {
			gtlastCards = append(gtlastCards, pair)
		}
	}
	if len(gtlastCards) == 0 {
		return nil
	}

	// 根据类型返回最大或最小的对子
	switch pairType {
	case SmallestPair:
		return findSmallestPair(allPairs)
	case LargestPair:
		return findLargestPair(allPairs)
	default:
		return nil
	}
}

// 按花色从大到小排序：黑桃(3) > 红桃(2) > 梅花(1) > 方块(0)
func sortCardsBySuitDesc(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Suit > cards[j].Suit
	})
}

// 比较两个对子的大小（用于找最小对子）
func comparePairsForSmallest(pair1, pair2 []Card) bool {
	// 先比较点数
	if pair1[0].Rank != pair2[0].Rank {
		return pair1[0].Rank < pair2[0].Rank
	}

	// 点数相同，比较第一张牌的花色
	if pair1[0].Suit != pair2[0].Suit {
		return pair1[0].Suit < pair2[0].Suit
	}

	// 第一张牌花色相同，比较第二张牌的花色
	return pair1[1].Suit < pair2[1].Suit
}

// 比较两个对子的大小（用于找最大对子）
func comparePairsForLargest(pair1, pair2 []Card) bool {
	// 先比较点数
	if pair1[0].Rank != pair2[0].Rank {
		return pair1[0].Rank > pair2[0].Rank
	}

	// 点数相同，比较第一张牌的花色
	if pair1[0].Suit != pair2[0].Suit {
		return pair1[0].Suit > pair2[0].Suit
	}

	// 第一张牌花色相同，比较第二张牌的花色
	return pair1[1].Suit > pair2[1].Suit
}

// 找到所有对子中最小的对子
func findSmallestPair(pairs [][]Card) []Card {
	if len(pairs) == 0 {
		return nil
	}

	smallest := pairs[0]
	for i := 1; i < len(pairs); i++ {
		if comparePairsForSmallest(pairs[i], smallest) {
			smallest = pairs[i]
		}
	}
	return smallest
}

// 找到所有对子中最大的对子
func findLargestPair(pairs [][]Card) []Card {
	if len(pairs) == 0 {
		return nil
	}

	largest := pairs[0]
	for i := 1; i < len(pairs); i++ {
		if comparePairsForLargest(pairs[i], largest) {
			largest = pairs[i]
		}
	}
	return largest
}
