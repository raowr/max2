package log

import (
	"gfast/internal/app/log/model/entity"

	commonApi "gfast/api/v1/common"

	"github.com/gogf/gf/v2/frame/g"
)

type LogGameListReq struct {
	g.Meta    `path:"/logGame/list" tags:"游戏日志管理" method:"get" summary:"游戏日志列表"`
	RoomId    string   `p:"roomId"`    //房间id
	Type      int      `p:"type"`      //类型
	Status    int      `p:"status"`    //状态
	UserId    string   `p:"userId"`    //用户id
	DateRange []string `p:"dateRange"` //日期范围
	commonApi.PageReq
}

type LogGameListRes struct {
	g.Meta `mime:"application/json"`
	commonApi.ListRes
	List []*entity.LogGame `json:"list"`
}
