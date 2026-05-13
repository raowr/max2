package register

import (
	"context"
	"encoding/json"
	"time"
	"unicode/utf8"
	"user/internal/dao"
	"user/internal/library/liberr"
	"user/internal/model/do"
	"user/internal/model/entity"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"

	v1 "user/api/register/v1"
	"user/internal/service"
)

func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	//校验用户名长度
	if utf8.RuneCountInString(req.Username) < 1 || utf8.RuneCountInString(req.Username) > 12 {
		liberr.ErrIsNil(ctx, err, "用户名长度必须在1-12之间")
	}
	// 校验密码是否一致
	if req.Password != req.Password2 {
		liberr.ErrIsNil(ctx, err, "密码不一致")
	}
	// 校验用户名是否存在
	if req.Username == "admin" {
		liberr.ErrIsNil(ctx, err, "用户名已存在")
	}
	//判断缓存中是否有这个用户名
	if !service.Cache().Get(ctx, req.Username).IsEmpty() {
		liberr.ErrIsNil(ctx, err, "用户名已存在")
	}
	// 校验用户名是否为空
	if req.Username == "" {
		liberr.ErrIsNil(ctx, err, "用户名不能为空")
	}
	// 校验密码是否为空
	if req.Password == "" {
		liberr.ErrIsNil(ctx, err, "密码不能为空")
	}
	// 校验密码是否为空
	if req.Password2 == "" {
		liberr.ErrIsNil(ctx, err, "确认密码不能为空")
	}
	//判断mysql中是否有这个用户名
	model := dao.Users.Ctx(ctx).Where("name = ?", req.Username)
	record, err := model.Fields("id").One()
	if err != nil {
		g.Log().Error(ctx, err)
		liberr.ErrIsNil(ctx, err, "用户名已存在")
	}
	if len(record) > 0 {
		liberr.ErrIsNil(ctx, err, "用户名已存在")
	}
	//如果注册成功，返回成功
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		err = g.Try(ctx, func(ctx context.Context) {
			//注册用户
			password, err := gmd5.Encrypt(req.Password)
			if err != nil {
				panic(err)
			}
			result, err := dao.Users.Ctx(ctx).TX(tx).Insert(do.Users{
				Name:     req.Username,
				Password: password,
				Point:    100,
				Token:    "",
			})
			if err != nil {
				panic(err)
			}
			userId, err := result.LastInsertId()
			if err != nil {
				panic(err)
			}
			entUser := entity.Users{
				Id:       uint(userId),
				Name:     req.Username,
				Password: password,
				Token:    "",
				Point:    100,
			}
			jsonStr, err := json.Marshal(entUser)
			if err != nil {
				liberr.ErrIsNil(ctx, err, "登录失败")
			}
			//缓存用户名
			service.Cache().Set(ctx, "user:"+req.Username, jsonStr, 7*24*time.Hour)
		})
		return err
	})
	return
}
