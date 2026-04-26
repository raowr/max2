package login_game

import (
	"context"
	"encoding/json"
	v1 "game_user/api/login_game/v1"
	"game_user/internal/service"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedLoginGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterLoginGameServer(s.Server, &Controller{})
}

func (*Controller) SendLogin(ctx context.Context, req *v1.SendLoginReq) (res *v1.SendLoginRes, err error) {
	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	service.Cache().Set(ctx, req.Token, reqJson, 7*24*60*60)
	return nil, nil
}
