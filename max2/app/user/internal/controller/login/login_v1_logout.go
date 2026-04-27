package login

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "user/api/login/v1"
	"user/internal/dao"
	"user/internal/model/entity"
	"user/internal/service"
)

func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	res = new(v1.LogoutRes)

	g.Try(ctx, func(ctx context.Context) {
		// 1. 从缓存中获取用户信息，验证token是否有效
		cacheKey := "user:" + req.Username
		cacheData := service.Cache().Get(ctx, cacheKey)
		if cacheData.IsEmpty() {
			gerror.New("用户未登录或登录已过期")
		}

		// 2. 解析缓存数据，验证token
		var userInfo entity.Users
		if err := cacheData.Scan(&userInfo); err != nil {
			gerror.New("解析缓存数据失败")
		}

		if userInfo.Token != req.Token {
			gerror.New("无效的token")
		}

		// 3. 从缓存中删除用户信息
		service.Cache().Remove(ctx, cacheKey)

		// 4. 从数据库中清除用户的token
		_, err = dao.Users.Ctx(ctx).Where("name", req.Username).Update("token", "")
		if err != nil {
			gerror.New("退出登录失败")
		}
	})
	return
}
