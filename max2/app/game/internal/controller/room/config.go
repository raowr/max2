package room

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

// 将全局变量改为函数，每次需要时重新从配置中获取
func GetOutCardTimeout() int {
	return g.Cfg().MustGet(context.Background(), "game.outCard").Int() //出牌最大时间
}

func GetGameCommission() int64 {
	return g.Cfg().MustGet(context.Background(), "game.comm").Int64() //每局抽水2
}
