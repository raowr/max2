package login

import (
	"context"
	"encoding/json"
	"time"

	v1 "user/api/login/v1"
	"user/internal/library/liberr"
	"user/internal/model/entity"
	"user/internal/service"

	"user/internal/dao"

	"user/internal/consts"

	loginGameGrpc "user/api/login_game/v1"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Id       uint
	Username string
	jwt.RegisteredClaims
}

func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	res = new(v1.LoginRes)
	var password string
	var entUser entity.Users
	g.Try(ctx, func(ctx context.Context) {
		//登录用户
		//验证缓存中是否存在
		userInfo := service.Cache().Get(ctx, "user:"+req.Username)
		if userInfo == nil {
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

			entUser = entity.Users{
				Id:       uint(user["id"].Uint()),
				Name:     req.Username,
				Password: password,
				Token:    user["token"].String(),
				Point:    int64(user["point"].Int64()),
			}
		} else {
			//缓存有，验证密码正确
			if err := json.Unmarshal(userInfo.Bytes(), &entUser); err != nil {
				liberr.ErrIsNil(ctx, err, "用户名或密码错误")
			}
			password = entUser.Password
		}
		reqPassword, err := gmd5.Encrypt(req.Password)
		if err != nil {
			liberr.ErrIsNil(ctx, err, "用户名或密码错误")
		}
		if password != reqPassword {
			liberr.ErrIsNil(ctx, err, "用户名或密码错误")
		}
		//选择游戏服务器
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
		entUser.Token = res.Token

		//登录成功,缓存密码
		jsonStr, err := json.Marshal(entUser)
		if err != nil {
			liberr.ErrIsNil(ctx, err, "登录失败")
		}
		service.Cache().Set(ctx, "user:"+req.Username, jsonStr, 60*60)
		//修改数据库
		_, err = dao.Users.Ctx(ctx).Where("id", entUser.Id).Update("token", res.Token)
		if err != nil {
			liberr.ErrIsNil(ctx, err, "登录失败")
		}

		//发送token到游戏服务器
		_, err = client.SendLogin(ctx, &loginGameGrpc.SendLoginReq{
			UserId:   gconv.String(entUser.Id),
			UserName: entUser.Name,
			Token:    res.Token,
			Point:    gconv.String(entUser.Point),
		})
		if err != nil {
			liberr.ErrIsNil(ctx, err, "登录失败")
		}
	})
	return
}
