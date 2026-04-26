package log_game

import (
	"context"
	v1 "game_user/api/log_game/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedLogGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterLogGameServer(s.Server, &Controller{})
}

func (*Controller) SendLog(ctx context.Context, req *v1.SendLogReq) (res *v1.SendLogRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
