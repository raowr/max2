package set_game

import (
	"context"
	"game_user/internal/dao"
	"game_user/internal/model/entity"

	"game_user/internal/service"
)

func InitCache() {
	//初始化缓存
	var List []*entity.SetGame
	model := dao.SetGame.Ctx(context.Background())
	model.Scan(&List)
	for _, item := range List {
		service.Cache().Set(context.Background(), item.Key, item.Value, 0)
	}
}
