package logic

import (
	"context"
	"gfast/api/v1/set"
	"gfast/internal/app/set/dao"
	"gfast/internal/app/set/service"
	"gfast/internal/app/system/consts"
	"gfast/library/liberr"

	"gfast/internal/app/set/model/do"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	service.RegisterSetGame(New())
}

func New() *sSetGame {
	return &sSetGame{}
}

type sSetGame struct {
}

func (s *sSetGame) GetSetGameListSearch(ctx context.Context, req *set.SetGameListReq) (res *set.SetGameListRes, err error) {
	res = new(set.SetGameListRes)
	g.Try(ctx, func(ctx context.Context) {
		model := dao.SetGame.Ctx(ctx)
		if req.Name != "" {
			model = model.Where("name like ?", req.Name)
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

func (s *sSetGame) Update(ctx context.Context, req *set.SetGameUpdateReq) (err error) {

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = g.Try(ctx, func(ctx context.Context) {
			//菜单数据
			data := do.SetGame{
				Id:    req.Id,
				Value: req.Value,
			}
			_, e := dao.SetGame.Ctx(ctx).TX(tx).WherePri(req.Id).Update(data)
			go func() {
				r, _ := g.Client().Post(
					ctx,
					"http://127.0.0.1:8002/set",
					g.Map{
						"key":   req.Key,
						"value": req.Value,
					},
				)
				r.Close()
			}()

			if err != nil {
				panic(err)
			}
			liberr.ErrIsNil(ctx, e, "修改失败")

		})
		return err
	})
	return
}
