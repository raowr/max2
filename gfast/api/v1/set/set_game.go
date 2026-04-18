package set

import (
	"gfast/internal/app/set/model/entity"

	commonApi "gfast/api/v1/common"

	"github.com/gogf/gf/v2/frame/g"
)

type SetGameListReq struct {
	g.Meta `path:"/setGame/list" tags:"游戏设置管理" method:"get" summary:"游戏设置列表"`
	Name   string `p:"name"` //名称
	commonApi.PageReq
}

type SetGameListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	List []*entity.SetGame `json:"list"`
}

type SetGameUpdateReq struct {
	g.Meta `path:"/setGame/update" tags:"游戏设置管理" method:"post" summary:"更新游戏设置"`
	Id     int    `p:"id"`    //id
	Key    string `p:"key"`   //键
	Value  string `p:"value"` //值
	commonApi.PageReq
}

type SetGameUpdateRes struct {
	g.Meta `mime:"application/json"`
}
