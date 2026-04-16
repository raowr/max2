package log_game

import (
	"context"
	v1 "log-service/api/log_game/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct {
	v1.UnimplementedLogGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterLogGameServer(s.Server, &Controller{})
}

func (*Controller) SendLog(ctx context.Context, req *v1.SendLogReq) (res *v1.SendLogRes, err error) {
	g.Log().Infof(ctx, "收到日志信息: %v", req.GetRoomID())
	return nil, nil
}
