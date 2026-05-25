package action

import (
	"context"
	"game_user/internal/controller/room"
	"sync"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gorilla/websocket"
)

var (
	clientCache *gcache.Cache
	cacheCtx    = context.Background()
	// 缓存键常量
	cacheKeyPrefix = "ws:client:"

	// fake *faker.Faker

	err error
)

// 客户端连接结构体
type Client struct {
	conn      *websocket.Conn    // WebSocket 连接
	userName  string             // 用户标识（用于重连）
	heartbeat time.Time          // 最后心跳时间
	pid       int                // 玩家id
	cancel    context.CancelFunc // 添加取消函数
	mutex     sync.RWMutex       // 新增：保护 conn 等字段的并发访问
	closed    int32              // 添加原子关闭标记
	subClient *gredis.Redis      // 订阅客户端
	pubClient *gredis.Redis      // 发布客户端
}

// 全局房间管理器及并发安全锁（核心优化：解决全局资源竞争）
var (
	clients       = make(map[string]*Client)
	clientsMu     sync.RWMutex // 保护clients的读写锁
	rm            *room.RoomManager
	rmMu          sync.RWMutex              // 保护roomManager的读写锁
	allowedOrigin = room.GetAllowedOrigin() // 生产环境需替换为实际域名
)

// 推送玩家列表（处理JSON错误）
// 创建一个临时结构体，只包含可序列化的字段
type PlayerDTO struct {
	ID       int    `json:"ID"`
	Name     string `json:"Name"`
	RoomID   string `json:"RoomID"`
	Type     int    `json:"Type"`
	CardNum  int    `json:"CardNum"`
	Point    int64  `json:"Point"`
	UserId   string `json:"UserId"`
	RoomType int    `json:"RoomType"`
	// 只包含需要序列化的字段，排除MsgChan等不可序列化字段
}

func init() {
	// 初始化全局房间管理器（带锁保护）
	rmMu.Lock()
	if rm == nil {
		rm = room.NewRoomManager()
	}
	rmMu.Unlock()

	clientCache = gcache.New()

	// 先尝试简体中文
	// fake, err = faker.New("zh-CN")
	// if err != nil {
	// 	// 如果失败，尝试 zh
	// 	fake, err = faker.New("zh")
	// 	if err != nil {
	// 		// 使用默认语言（英文）
	// 		fake, _ = faker.New("en")
	// 	}
	// }

}

// 生成缓存键
func getClientCacheKey(userID string) string {
	return cacheKeyPrefix + userID
}

// 添加客户端到缓存
func addClient(client *Client) error {
	cacheKey := getClientCacheKey(client.userName)
	return clientCache.Set(cacheCtx, cacheKey, client, time.Hour)
}

// 根据userID获取客户端
func getClient(userID string) (*Client, error) {
	cacheKey := getClientCacheKey(userID)

	value, err := clientCache.Get(cacheCtx, cacheKey)
	if err != nil {
		return nil, err
	}

	if value == nil || value.IsNil() {
		return nil, nil // 客户端不存在
	}

	client, ok := value.Val().(*Client)
	if !ok {
		// 类型断言失败，删除无效缓存
		clientCache.Remove(cacheCtx, cacheKey)
		return nil, nil
	}

	return client, nil
}

// 删除客户端
func removeClient(userID string) error {
	cacheKey := getClientCacheKey(userID)
	_, err := clientCache.Remove(cacheCtx, cacheKey)
	return err
}

// 批量删除（用于清理操作）
func removeClients(userIDs []string) error {
	for _, userID := range userIDs {
		if err := removeClient(userID); err != nil {
			return err
		}
	}
	return nil
}

func getPlayers(roomInfo *room.Room) []*PlayerDTO {
	// 转换玩家列表为可序列化的DTO列表
	playerDTOs := make([]*PlayerDTO, len(roomInfo.Players))
	for i, p := range roomInfo.Players {
		playerDTOs[i] = &PlayerDTO{
			ID:       p.ID,
			Name:     p.Name,
			RoomID:   p.RoomID,
			Type:     int(p.Type),
			CardNum:  p.CardNum,
			Point:    p.Point,
			UserId:   p.UserName,
			RoomType: roomInfo.Type,
		}
	}
	return playerDTOs
}
