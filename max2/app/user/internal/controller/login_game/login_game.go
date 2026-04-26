package login_game

import (
	"context"
	v1 "user/api/login_game/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedLoginGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterLoginGameServer(s.Server, &Controller{})
}

func (*Controller) SendLogin(ctx context.Context, req *v1.SendLoginReq) (res *v1.SendLoginRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
