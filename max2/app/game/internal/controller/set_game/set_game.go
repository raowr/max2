package set_game

import (
	"context"
	v1 "game/api/set_game/v1"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
)

type Controller struct {
	v1.UnimplementedSetGameServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterSetGameServer(s.Server, &Controller{})
}

func (*Controller) SendSet(ctx context.Context, req *v1.SendSetReq) (res *v1.SendSetRes, err error) {
	g.Log().Info(ctx, "SendSet", req)
	SetCache.Set(ctx, req.Key, req.Value, 0)
	return
}
