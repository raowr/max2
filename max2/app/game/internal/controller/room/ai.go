package room

import "sort"

func DecideCards(player, landlord *Player, lastPH int, lastCards []Card) (pokerHandType int, playCards []Card) {
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

		//如果时3号玩家，人类玩家剩余1张牌
		//上次出牌是单牌，玩家剩余1张牌，选择比上家大的最大的单牌出
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
		//如果打的是单张,玩家3没有单张，需要拆出单张，如果打的是对，需要从3条和4条中拆
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

		//如果是3号玩家，人类玩家剩余2张牌，并且牌组上没有并且上一手是人类打的，必须拆牌
		if player.ID == 3 && len(landlord.Cards) == 2 && (lastPH == SINGLE || lastPH == PAIR) && landlord.Must && len(playCards) <= 0 {
			//如果是单牌
			if lastPH == SINGLE {
				gtLandlordCardIds := make([]Card, 0) //所有比人类玩家大的牌的Id
				if player.CardNum > 0 {
					for _, v := range player.handPattern {
						for _, v1 := range v {
							for _, v2 := range v1 {
								singleCards := make([]Card, 0)
								singleCards = append(singleCards, v2)
								//如果比上手大，就需要重新整理牌
								if compareCards(lastPH, lastPH, lastCards, singleCards) {
									gtLandlordCardIds = append(gtLandlordCardIds, v2)
								}
							}
						}
					}
				}
				if len(gtLandlordCardIds) > 0 {
					//取大于玩家且最小的打
					sort.Slice(gtLandlordCardIds, func(i, j int) bool {
						return gtLandlordCardIds[i].Id < gtLandlordCardIds[j].Id
					})
					playCards = append(playCards, gtLandlordCardIds[0])
					player.ReCard = true
				}
			}
			//如果是对子
			if lastPH == PAIR {
				allCards := make([]Card, 0)
				//查找所有的牌
				if player.CardNum > 0 {
					for _, v := range player.handPattern {
						for _, v1 := range v {
							allCards = append(allCards, v1...)
						}
					}
				}
				pairs := FindSmallestPossiblePair(lastPH, lastCards, allCards, LargestPair)
				if pairs != nil {
					playCards = pairs
					player.ReCard = true
				}
			}
		}
		//如果是3号玩家，人类玩家剩余5张牌，50%概率拆并且牌组上没有并且上一手是人类打的，必须拆牌
		if player.ID == 3 && len(landlord.Cards) == 5 && (lastPH == SINGLE || lastPH == PAIR) && landlord.Must && Probability(40) && len(playCards) <= 0 {
			//如果是单牌
			if lastPH == SINGLE {
				gtLandlordCardIds := make([]Card, 0) //所有比人类玩家大的牌的Id
				if player.CardNum > 0 {
					for _, v := range player.handPattern {
						for _, v1 := range v {
							for _, v2 := range v1 {
								singleCards := make([]Card, 0)
								singleCards = append(singleCards, v2)
								//如果比上手大，就需要重新整理牌
								if compareCards(lastPH, lastPH, lastCards, singleCards) {
									gtLandlordCardIds = append(gtLandlordCardIds, v2)
								}
							}
						}
					}
				}
				if len(gtLandlordCardIds) > 0 {
					//取大于玩家且最小的打
					sort.Slice(gtLandlordCardIds, func(i, j int) bool {
						return gtLandlordCardIds[i].Id < gtLandlordCardIds[j].Id
					})
					playCards = append(playCards, gtLandlordCardIds[0])
					player.ReCard = true
				}
			}
			//如果是对子
			if lastPH == PAIR {
				allCards := make([]Card, 0)
				//查找所有的牌
				if player.CardNum > 0 {
					for _, v := range player.handPattern {
						for _, v1 := range v {
							allCards = append(allCards, v1...)
						}
					}
				}
				pairs := FindSmallestPossiblePair(lastPH, lastCards, allCards, SmallestPair)
				if pairs != nil {
					playCards = pairs
					player.ReCard = true
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
