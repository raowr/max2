package logGame

import (
	"context"
	"gfast/api/v1/log"
	"gfast/internal/app/log/dao"
	"gfast/internal/app/log/service"
	"gfast/internal/app/system/consts"
	"gfast/library/liberr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func init() {
	service.RegisterLogGame(New())
}

func New() *sLogGame {
	return &sLogGame{}
}

type sLogGame struct {
}

func (s *sLogGame) GetLogGameListSearch(ctx context.Context, req *log.LogGameListReq) (res *log.LogGameListRes, err error) {
	res = new(log.LogGameListRes)
	g.Try(ctx, func(ctx context.Context) {
		model := dao.LogGame.Ctx(ctx)
		if len(req.DateRange) > 0 {
			model = model.Where("creat_time between ? and ?", req.DateRange[0], req.DateRange[1])
		}
		if req.RoomId != "" {
			model = model.Where("room_id = ?", req.RoomId)
		}
		if req.Status != 0 {
			model = model.Where("status", gconv.Int(req.Status))
		}
		if req.Type != 0 {
			model = model.Where("type", gconv.Int(req.Type))
		}
		if req.UserId != "" {
			model = model.Where("user_id = ?", req.UserId)
		}
		res.Total, err = model.Count()
		liberr.ErrIsNil(ctx, err, "获取游戏日志总数失败")
		if req.PageNum == 0 {
			req.PageNum = 1
		}
		res.CurrentPage = req.PageNum
		if req.PageSize == 0 {
			req.PageSize = consts.PageSize
		}
		err = model.Page(res.CurrentPage, req.PageSize).Order("id desc").Scan(&res.List)
		liberr.ErrIsNil(ctx, err, "获取游戏日志数据失败")
	})
	return
}
