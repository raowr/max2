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

	// ========== 1. 密码错误限流：防止暴力猜密码 ==========
	pwdErrCount := service.Cache().Get(ctx, "pwd_err:"+req.Username)
	if !pwdErrCount.IsEmpty() {
		errCount := pwdErrCount.Int()
		// 连续5次密码错误，锁定1分钟
		if errCount >= 5 {
			err = gerror.New("密码错误次数过多，请1分钟后再试")
			return
		}
	}

	// ========== 2. 缓存穿透防护：用户信息缓存 ==========
	userInfo := service.Cache().Get(ctx, "user:"+req.Username)
	if userInfo.IsEmpty() {
		//缓存没有，查数据库
		user, err = dao.Users.Ctx(ctx).Where("name", req.Username).One()
		if len(user) == 0 {
			// 空值缓存：用户不存在也缓存10秒（+随机抖动防雪崩）
			nullTTL := 10*time.Second + time.Duration(rand.Intn(5))*time.Second
			service.Cache().Set(ctx, "user:"+req.Username, []byte("null"), nullTTL)
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
		//缓存有，检测是否为空值缓存（防止缓存穿透）
		if string(userInfo.Bytes()) == "null" {
			err = gerror.New("用户不存在")
			return
		}
		//正常用户缓存，验证密码正确
		if err = json.Unmarshal(userInfo.Bytes(), &entUser); err != nil {
			g.Log().Error(ctx, err)
			err = gerror.New("用户名或密码错误")
			return
		}
		password = entUser.Password
	}

	// ========== 3. 密码验证 ==========
	reqPassword, err := gmd5.Encrypt(req.Password)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("用户名或密码错误")
		return
	}
	if password != reqPassword {
		// 密码错误计数 +1
		newCount := pwdErrCount.Int() + 1
		service.Cache().Set(ctx, "pwd_err:"+req.Username, newCount, 5*time.Minute)
		err = gerror.New("用户名或密码错误")
		return
	}

	// ========== 4. 登录成功，清除密码错误计数 ==========
	service.Cache().Remove(ctx, "pwd_err:"+req.Username)

	res.Username = req.Username
	res.Point = entUser.Point

	// 拿到地址（IP:端口）
	res.Node = GetRandomGameUserNode()
	g.Log().Debugf(ctx, "发现服务节点: %s", res.Node)

	// 随机抖动300秒，防止缓存雪崩（大量缓存同时过期）
	SuffixTime := time.Duration(rand.Intn(300)) * time.Second

	// 生成token
	uc := &jwtClaims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7*24*time.Hour + SuffixTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	res.Token, err = token.SignedString([]byte(consts.JwtKey))
	entUser.Token = res.Token

	//登录成功,缓存用户信息（过期时间带随机抖动，防雪崩）
	jsonStr, err := json.Marshal(entUser)
	if err != nil {
		g.Log().Error(ctx, err)
		err = gerror.New("登录失败1")
		return
	}
	userTTL := 7*24*time.Hour + SuffixTime
	service.Cache().Set(ctx, "user:"+req.Username, jsonStr, userTTL)

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
