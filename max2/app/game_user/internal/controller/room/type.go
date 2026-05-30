package room

import (
	"context"
	v1 "game_user/api/action/v1"
	"sync"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gtimer"
)

var ctx = context.Background()

type MsgChan chan RoomMsg

type RoomMsg struct {
	Type string
	Data string
}

// 牌的结构体
type Card struct {
	Value    string // 牌值 3-15，分别对应3-A
	Suit     int    // 花色 0-3 方块、梅花、红桃、黑桃
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
	UserName    string           //用户名
}

// 房间结构体
type Room struct {
	ID           string
	Players      []*Player
	Deck         []Card
	Landlord     *Player
	Farmers      []*Player
	Current      int
	LastCards    []Card
	LastPH       int
	Turn         int
	IsPlaying    bool
	Rgtimer      *gtimer.Timer `json:"-"` // 忽略
	OutStarTime  int
	passCount    int
	NextPlayerID int
	Status       int
	mutex        sync.RWMutex `json:"-"` // 忽略
	Type         int
	subConn      gredis.Conn            `json:"-"` // 忽略
	subClient    *gredis.Redis          `json:"-"` // 忽略
	pubClient    *gredis.Redis          `json:"-"` // 忽略
	MsgQueue     chan *v1.SendActionReq `json:"-"` // 忽略
	MsgCancel    context.CancelFunc     `json:"-"` // 忽略 ← 核心问题
}

// 房间管理器
type RoomManager struct {
	Rooms      map[string]*Room
	PlayerList map[string]*Player
}

// 对子查找类型
type PairType int

const (
	SmallestPair PairType = iota // 最小的对子
	LargestPair                  // 最大的对子
)

type overMsg struct {
	WinName string `json:"winName"` //玩家名称
	Winner  int    `json:"winner"`  //玩家ID
	Point   int64  `json:"point"`   //玩家总瓜子数
	Win     int64  `json:"win"`     //当次赢分,正数为赢，负数为输
}

type Users struct {
	Id       uint   `json:"id"       orm:"id"       description:"自增ID"`
	Name     string `json:"name"     orm:"name"     description:"用户名称"`
	Password string `json:"password" orm:"password" description:"用户密码"`
	Point    int64  `json:"point"    orm:"point"    description:"用户分数"`
	Token    string `json:"token"    orm:"token"    description:"登录token"`
}
