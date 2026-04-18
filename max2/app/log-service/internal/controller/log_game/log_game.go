package log_game

import (
	"context"
	v1 "log-service/api/log_game/v1"
	"log-service/internal/dao"
	"log-service/internal/model/do"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type Controller struct {
	v1.UnimplementedLogGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterLogGameServer(s.Server, &Controller{})
}

func (*Controller) SendLog(ctx context.Context, req *v1.SendLogReq) (res *v1.SendLogRes, err error) {
	g.Log().Infof(ctx, "收到日志信息: %v", req.GetRoomID())
	currentTime := gtime.Now()
	//写入数据库
	_, err = dao.LogGame.Ctx(ctx).Data(do.LogGame{
		RoomId:    req.GetRoomID(),
		Type:      req.GetType(),
		Status:    req.GetStatus(),
		UserId:    req.GetUserID(),
		Point:     req.GetPoint(),
		Action:    req.GetAction(),
		Remain:    req.GetRemain(),
		OutCards:  req.GetOutCardIds(),
		Text:      req.GetText(),
		CreatTime: currentTime,
	}).Insert()
	return nil, err
}
