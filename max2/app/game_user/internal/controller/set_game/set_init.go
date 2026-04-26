package set_game

import (
	"context"
	"game_user/internal/dao"

	"game_user/internal/model/entity"

	"github.com/gogf/gf/v2/os/gcache"
)

var SetCache *gcache.Cache

func init() {
	SetCache = gcache.New()
	//初始化缓存
	var List []*entity.SetGame
	model := dao.SetGame.Ctx(context.Background())
	model.Scan(&List)
	for _, item := range List {
		SetCache.Set(context.Background(), item.Key, item.Value, 0)
	}
}
