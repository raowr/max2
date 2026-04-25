package login

import (
	"context"
	"time"

	v1 "user/api/login/v1"
	"user/internal/library/liberr"
	"user/internal/service"

	"user/internal/dao"

	"user/internal/consts"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Id       uint
	Username string
	jwt.RegisteredClaims
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	res = new(v1.LoginRes)
	g.Try(ctx, func(ctx context.Context) {
		//登录用户
		//验证缓存中是否存在
		password := service.Cache().Get(ctx, req.Username).String()
		if password == "" {
			//缓存没有，查数据库
			user, err := dao.Users.Ctx(ctx).Where("name", req.Username).One()
			if err != nil {
				liberr.ErrIsNil(ctx, err, "用户名或密码错误")
			}
			//如果数据库也没有，返回错误
			if len(user) == 0 {
				liberr.ErrIsNil(ctx, err, "用户不存在")
			}
			//如果用户存在，验证密码
			password = user["password"].String()
			service.Cache().Set(ctx, req.Username, password, 60*60)
		}
		reqPassword, err := gmd5.Encrypt(req.Password)
		if err != nil {
			liberr.ErrIsNil(ctx, err, "用户名或密码错误")
		}
		if password != reqPassword {
			liberr.ErrIsNil(ctx, err, "用户名或密码错误")
		}
		//登录成功,缓存密码
		service.Cache().Set(ctx, req.Username, password, 60*60)
		//选择游戏服务器
		gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))

		gsel.SetBuilder(gsel.NewBuilderRoundRobin())
		selected, _, _ := gsel.GetBuilder().Build().Pick(ctx)
		targetNode := selected.Address()
		res.Node = targetNode

		// 生成token
		uc := &jwtClaims{
			Username: req.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
		res.Token, err = token.SignedString([]byte(consts.JwtKey))
	})
	return
}
