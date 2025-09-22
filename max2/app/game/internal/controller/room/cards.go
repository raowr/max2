package room

import (
	"sort"
	"strconv"
)

func rankToNumber(rank string) int {
	switch rank {
	case "A":
		return 1
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	default:
		if num, err := strconv.Atoi(rank); err == nil {
			return num
		}
		return 0
	}
}

// 顺子(5张点数连续的牌,花色不同)
func findAllStraights(cards []Card) [][]Card {
	// 创建点数到牌的映射
	rankMap := make(map[int][]Card)
	for _, card := range cards {
		num := rankToNumber(card.Rank)
		rankMap[num] = append(rankMap[num], card)
		// A也可以作为14
		if num == 1 {
			rankMap[14] = append(rankMap[14], card)
		}
	}

	// 获取所有点数并排序
	var allRanks []int
	for rank := range rankMap {
		allRanks = append(allRanks, rank)
	}
	sort.Ints(allRanks)

	var results [][]Card

	// 遍历所有可能的起始点数
	for i := 0; i <= len(allRanks)-5; i++ {
		start := allRanks[i]

		// 检查是否构成顺子
		valid := true
		for j := 0; j < 5; j++ {
			if len(rankMap[start+j]) == 0 {
				valid = false
				break
			}
		}

		if valid {
			// 生成这个顺子的所有可能组合
			generateCombinations(start, rankMap, []Card{}, &results)
		}
	}

	return results
}

func generateCombinations(start int, rankMap map[int][]Card, current []Card, results *[][]Card) {
	if len(current) == 5 {
		straight := make([]Card, 5)
		copy(straight, current)
		*results = append(*results, straight)
		return
	}

	currentRank := start + len(current)
	for _, card := range rankMap[currentRank] {
		// 检查这张牌是否已经在当前组合中（避免同一张牌重复使用）
		if !isCardInSlice(current, card) {
			generateCombinations(start, rankMap, append(current, card), results)
		}
	}
}

func isCardInSlice(cards []Card, target Card) bool {
	for _, card := range cards {
		if card.Value == target.Value && card.Suit == target.Suit {
			return true
		}
	}
	return false
}

// 花色
// 同花(5张相同的花色)
func findAllFlushes(cards []Card) [][]Card {
	// 按花色分组并排序（为了输出美观）
	suitGroups := make(map[int][]Card)
	for _, card := range cards {
		suitGroups[card.Suit] = append(suitGroups[card.Suit], card)
	}

	var results [][]Card

	for _, suitedCards := range suitGroups {
		if len(suitedCards) >= 5 {
			// 对同花色的牌按点数排序（为了输出美观）
			sort.Slice(suitedCards, func(i, j int) bool {
				return rankToNumber(suitedCards[i].Rank) < rankToNumber(suitedCards[j].Rank)
			})

			// 生成组合
			combinations := generateCombinationsIterative(suitedCards, 5)
			results = append(results, combinations...)
		}
	}

	return results
}

// 迭代方式生成组合，避免递归深度问题
func generateCombinationsIterative(cards []Card, k int) [][]Card {
	if k == 0 || len(cards) < k {
		return [][]Card{}
	}

	var results [][]Card
	stack := [][]int{{0}} // 存储索引组合

	for len(stack) > 0 {
		// 取出最后一个组合
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(last) == k {
			// 找到一个完整组合
			combination := make([]Card, k)
			for i, idx := range last {
				combination[i] = cards[idx]
			}
			results = append(results, combination)
			continue
		}

		// 扩展组合
		start := 0
		if len(last) > 0 {
			start = last[len(last)-1] + 1
		}

		for i := start; i <= len(cards)-k+len(last); i++ {
			newComb := make([]int, len(last)+1)
			copy(newComb, last)
			newComb[len(last)] = i
			stack = append(stack, newComb)
		}
	}

	return results
}

// 福禄(3带2)

// 使用迭代方法生成组合，更高效
func findAllFlushesEfficient(cards []Card) [][]Card {
	// 按花色分组并排序（为了输出美观）
	suitGroups := make(map[int][]Card)
	for _, card := range cards {
		suitGroups[card.Suit] = append(suitGroups[card.Suit], card)
	}

	var results [][]Card

	for _, suitedCards := range suitGroups {
		if len(suitedCards) >= 5 {
			// 对同花色的牌按点数排序（为了输出美观）
			sort.Slice(suitedCards, func(i, j int) bool {
				return rankToNumber(suitedCards[i].Rank) < rankToNumber(suitedCards[j].Rank)
			})

			// 生成组合
			combinations := generateCombinationsIterative(suitedCards, 5)
			results = append(results, combinations...)
		}
	}

	return results
}
