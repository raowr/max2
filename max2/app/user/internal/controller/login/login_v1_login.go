package login

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	v1 "user/api/login/v1"
	"user/internal/model/do"
	"user/internal/model/entity"
	"user/internal/service"

	"user/internal/dao"

	"user/internal/consts"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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
	var user gdb.Record
	//登录用户
	//验证缓存中是否存在
	userInfo := service.Cache().Get(ctx, "user:"+req.Username)
	if userInfo.IsEmpty() {
		//缓存没有，查数据库
		user, err = dao.Users.Ctx(ctx).Where("name", req.Username).One()
		if len(user) == 0 {
			err = gerror.New("用户不存在")
			return
		}
		if err != nil {
			g.Log().Error(ctx, err)
			err = gerror.New("用户名或密码错误")
			return
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
		if err = json.Unmarshal(userInfo.Bytes(), &entUser); err != nil {
			g.Log().Error(ctx, err)
			err = gerror.New("用户名或密码错误")
			return
		}
		password = entUser.Password
	}
	reqPassword, err := gmd5.Encrypt(req.Password)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("用户名或密码错误")
		return
	}
	if password != reqPassword {
		err = gerror.New("用户名或密码错误")
		return
	}
	res.Username = req.Username
	res.Point = entUser.Point
	//选择游戏服务器
	// 2. 获取全部节点
	serviceGsvc := GetGameUserService()
	if serviceGsvc == nil {
		err = gerror.New("无法获取 game_user 服务")
		return
	}
	endpoints := serviceGsvc.GetEndpoints()

	// 3. 随机选一个
	randomNode := endpoints[rand.Intn(len(endpoints))]

	// 4. 拿到地址（IP:端口）
	res.Node = randomNode.String()

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
		g.Log().Error(ctx, err)
		err = gerror.New("登录失败1")
		return
	}
	service.Cache().Set(ctx, "user:"+req.Username, jsonStr, 7*24*time.Hour)
	//修改数据库
	data := do.Users{
		Token: res.Token,
	}
	_, err = dao.Users.Ctx(ctx).Where("id", entUser.Id).Update(data)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("登录失败2")
		return
	}
	return
}
