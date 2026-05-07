package user

import (
	"context"
	"encoding/json"
	"gfast/api/v1/user"
	cservice "gfast/internal/app/common/service"
	"gfast/internal/app/system/consts"
	"gfast/internal/app/user/dao"
	"gfast/internal/app/user/service"
	"gfast/library/liberr"
	"time"

	"gfast/internal/app/user/model/do"
	"gfast/internal/app/user/model/entity"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}

type sUser struct {
}

func (s *sUser) GetUserListSearch(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error) {
	res = new(user.UserListRes)
	g.Try(ctx, func(ctx context.Context) {
		model := dao.Users.Ctx(ctx)
		if req.Name != "" {
			model = model.Where("name like ?", "%"+req.Name+"%")
		}
		res.Total, err = model.Count()
		liberr.ErrIsNil(ctx, err, "获取用户总数失败")
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

func (s *sUser) UserUpdatePoint(ctx context.Context, req *user.UserUpdatePointReq) (res *user.UserUpdatePointRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = g.Try(ctx, func(ctx context.Context) {
			data := do.Users{
				Id:    req.Id,
				Point: req.Point,
			}
			_, e := dao.Users.Ctx(ctx).TX(tx).WherePri(req.Id).Update(data)
			liberr.ErrIsNil(ctx, e, "修改用户积分失败")
		})
		return err
	})

	// ✅ 将异步操作移到事务外部
	if err == nil {
		go func(userId uint, point int64) {
			// 使用独立的 context
			ctx := gctx.New()
			user, err := dao.Users.Ctx(ctx).WherePri(userId).One()
			if err != nil {
				g.Log().Errorf(ctx, "异步更新缓存失败: %v", err)
				return // 不要在 goroutine 中 panic
			}
			entUser := entity.Users{
				Id:       uint(user["id"].Uint()),
				Name:     user["name"].String(),
				Password: user["password"].String(),
				Token:    user["token"].String(),
				Point:    point,
			}
			jsonStr, err := json.Marshal(entUser)
			if err != nil {
				g.Log().Errorf(ctx, "序列化用户数据失败: %v", err)
				return
			}
			cservice.Cache().Set(ctx, "user:"+user["name"].String(), jsonStr, 7*24*time.Hour)
		}(uint(req.Id), int64(req.Point))
	}

	return
}

func (s *sUser) UserUpdatePassword(ctx context.Context, req *user.UserUpdatePasswordReq) (res *user.UserUpdatePasswordRes, err error) {
	var password string
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = g.Try(ctx, func(ctx context.Context) {
			var encryptErr error
			password, encryptErr = gmd5.Encrypt(req.Password)
			liberr.ErrIsNil(ctx, encryptErr, "密码加密失败")
			data := do.Users{
				Id:       req.Id,
				Password: password,
			}
			_, e := dao.Users.Ctx(ctx).TX(tx).WherePri(req.Id).Update(data)
			liberr.ErrIsNil(ctx, e, "修改失败")
		})
		return err
	})

	// ✅ 将异步操作移到事务外部
	if err == nil {
		go func(userId uint, encryptedPassword string) {
			// 使用独立的 context
			asyncCtx := gctx.New()
			userData, err := dao.Users.Ctx(asyncCtx).WherePri(userId).One()
			if err != nil || len(userData) == 0 {
				g.Log().Errorf(asyncCtx, "异步更新缓存失败，用户不存在: userId=%d, err=%v", userId, err)
				return
			}
			entUser := entity.Users{
				Id:       uint(userData["id"].Uint()),
				Name:     userData["name"].String(),
				Password: encryptedPassword,
				Token:    userData["token"].String(),
				Point:    int64(userData["point"].Int()),
			}
			jsonStr, err := json.Marshal(entUser)
			if err != nil {
				g.Log().Errorf(asyncCtx, "序列化用户数据失败: %v", err)
				return
			}
			// 修改redis缓存
			cservice.Cache().Set(asyncCtx, "user:"+userData["name"].String(), jsonStr, 7*24*time.Hour)
		}(uint(req.Id), password)
	}

	return
}
