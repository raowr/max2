// util/card.js
export const CARD_TYPE = {
    NONE: 0,          // 无效牌型
    SINGLE: 1,        // 单张（1张）
    PAIR: 2,          // 对子（2张）
    STRAIGHT: 3,      // 顺子（5张）
    SUIT: 4,          // 同花（5张）
    FULL_HOUSE: 5,    // 葫芦（3带2）
    FOUR_OF_A_KIND: 6,// 四条（4带1）
    STRAIGHT_FLUSH: 7 // 同花顺（5张）
};

/**
 * 判断是否可以出牌
 * @param {Array} currentCards 当前要出的牌数组
 * @param {Object} lastHand 上一手牌信息 {cards: [], isSelf: boolean, type: number}
 * @returns {boolean} 是否可以出牌
 */
export function canPlayCards(currentCards, lastHand) {
    const currentLen = currentCards.length;
    const currentType = getCardType(currentCards);

    if (currentType === CARD_TYPE.NONE) return false;

    if (lastHand?.isSelf) return true;

    if (!lastHand || !lastHand.cards || lastHand.cards.length === 0) return true;

    const lastLen = lastHand.cards.length;
    if (currentLen !== lastLen) return false;

    if (lastHand.type === CARD_TYPE.NONE) return true;

    return isCurrentTypeGreater(currentType, currentCards, lastHand.type, lastHand.cards);
}

/**
 * 获取牌型
 * @param {Array} cards 牌数组
 * @returns {number} 牌型常量
 */
export function getCardType(cards) {
    const len = cards.length;
    switch (len) {
        case 1: return CARD_TYPE.SINGLE;
        case 2: return judgeTwoCards(cards);
        case 5: return judgeFiveCards(cards);
        default: return CARD_TYPE.NONE;
    }
}

/**
 * 比较当前牌型是否大于上一手
 */
function isCurrentTypeGreater(currentType, currentCards, lastType, lastCards) {
    if (currentType > lastType) return true;
    if (currentType !== lastType) return false;

    // 先按牌型规则比较核心大小
    const baseCompare = compareByRankRule(currentType, currentCards, lastCards);
    if (baseCompare !== 0) return baseCompare > 0;

    // 核心大小相同时，按牌型特定规则比较最大ID
    const currentKeyId = getKeyIdByType(currentType, currentCards);
    const lastKeyId = getKeyIdByType(lastType, lastCards);
    return currentKeyId > lastKeyId;
}

/**
 * 根据牌型获取关键ID（核心修改点）
 * - 单张/对子/同花顺/同花/顺子：取牌组中最大ID
 * - 葫芦（3带2）：取3张相同牌中的最大ID
 * - 四条（4带1）：取4张相同牌中的最大ID
 */
function getKeyIdByType(type, cards) {
    switch (type) {
        case CARD_TYPE.FULL_HOUSE: {
            // 葫芦：获取3张相同牌中的最大ID
            const rankCount = countRanks(cards);
            const threeRank = getRankByCount(rankCount, 3);
            const threeCards = cards.filter(card => card.rank === threeRank);
            return Math.max(...threeCards.map(c => c.id));
        }
        case CARD_TYPE.FOUR_OF_A_KIND: {
            // 四条：获取4张相同牌中的最大ID
            const rankCount = countRanks(cards);
            const fourRank = getRankByCount(rankCount, 4);
            const fourCards = cards.filter(card => card.rank === fourRank);
            return Math.max(...fourCards.map(c => c.id));
        }
        default:
            // 其他牌型：获取牌组中最大ID
            return Math.max(...cards.map(c => c.id));
    }
}

/**
 * 按牌型规则比较核心大小（返回1:当前大, -1:上一手大, 0:相同）
 */
function compareByRankRule(type, current, last) {
    switch (type) {
        case CARD_TYPE.SINGLE:
            return current[0].rank > last[0].rank ? 1 : (current[0].rank < last[0].rank ? -1 : 0);
        case CARD_TYPE.PAIR:
            return current[0].rank > last[0].rank ? 1 : (current[0].rank < last[0].rank ? -1 : 0);
        case CARD_TYPE.STRAIGHT_FLUSH:
        case CARD_TYPE.SUIT:
        case CARD_TYPE.STRAIGHT: {
            const currentMax = getMaxRank(current);
            const lastMax = getMaxRank(last);
            return currentMax > lastMax ? 1 : (currentMax < lastMax ? -1 : 0);
        }
        case CARD_TYPE.FULL_HOUSE: {
            const currentThree = getRankByCount(countRanks(current), 3);
            const lastThree = getRankByCount(countRanks(last), 3);
            return currentThree > lastThree ? 1 : (currentThree < lastThree ? -1 : 0);
        }
        case CARD_TYPE.FOUR_OF_A_KIND: {
            const currentFour = getRankByCount(countRanks(current), 4);
            const lastFour = getRankByCount(countRanks(last), 4);
            return currentFour > lastFour ? 1 : (currentFour < lastFour ? -1 : 0);
        }
        default: return 0;
    }
}

// 辅助函数：判断两张牌是否为对子
function judgeTwoCards(cards) {
    return cards.length === 2 && cards[0].rank === cards[1].rank
        ? CARD_TYPE.PAIR
        : CARD_TYPE.NONE;
}

// 辅助函数：判断五张牌的牌型
function judgeFiveCards(cards) {
    if (cards.length !== 5) return CARD_TYPE.NONE;
    if (isFlush(cards) && isStraight(cards)) return CARD_TYPE.STRAIGHT_FLUSH;
    if (isFourOfAKind(cards)) return CARD_TYPE.FOUR_OF_A_KIND;
    if (isFullHouse(cards)) return CARD_TYPE.FULL_HOUSE;
    if (isFlush(cards)) return CARD_TYPE.SUIT;
    if (isStraight(cards)) return CARD_TYPE.STRAIGHT;
    return CARD_TYPE.NONE;
}

// 辅助函数：统计牌面数量
export function countRanks(cards) {
    return cards.reduce((count, card) => {
        count[card.rank] = (count[card.rank] || 0) + 1;
        return count;
    }, {});
}

// 辅助函数：判断是否为同花
function isFlush(cards) {
    const suit = cards[0].suit;
    return cards.every(card => card.suit === suit);
}

// 辅助函数：判断是否为顺子
function isStraight(cards) {
    const twoRank = 15;
    const maxValid = 14;
    const uniqueRanks = [...new Set(cards.map(c => c.rank))]
        .filter(rank => rank !== twoRank && rank <= maxValid)
        .sort((a, b) => a - b);
    if (uniqueRanks.length < 5) return false;
    for (let i = 0; i <= uniqueRanks.length - 5; i++) {
        if (uniqueRanks[i + 4] - uniqueRanks[i] === 4) return true;
    }
    return false;
}

// 辅助函数：判断是否为四条
function isFourOfAKind(cards) {
    const counts = countRanks(cards);
    const values = Object.values(counts);
    // 确保只有两个不同的牌点，且一个出现4次，另一个出现1次
    return values.length === 2 && values.includes(4) && values.includes(1);
}

// 辅助函数：判断是否为葫芦
function isFullHouse(cards) {
    const counts = countRanks(cards);
    const values = Object.values(counts);
     // 确保只有两个不同的牌点，且一个出现3次，另一个出现2次
    return values.length === 2 && values.includes(3) && values.includes(2);
}

// 辅助函数：获取最大牌点
function getMaxRank(cards) {
    return Math.max(...cards.map(card => card.rank));
}

// 辅助函数：根据数量获取牌点
export function getRankByCount(countMap, targetCount) {
    for (const [rank, count] of Object.entries(countMap)) {
        if (count === targetCount) return parseInt(rank, 10);
    }
    return 0;
}

// 辅助函数：获取排序后的牌点数组
function getSortedRanks(cards, descending = false) {
    return cards.map(card => card.rank).sort((a, b) =>
        descending ? b - a : a - b
    );
}

export const cardUtil = {
    canPlayCards,
    getCardType,
    CARD_TYPE,
    getKeyIdByType, // 暴露按牌型获取关键ID的函数
    countRanks,
    getRankByCount,
};

export default cardUtil;