package api_service

import (
	"context"

	api "api-service/api"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	api.UnimplementedSetGameServer
}

func Register(s *grpcx.GrpcServer) {
	api.RegisterSetGameServer(s.Server, &Controller{})
}

func (*Controller) SendSet(ctx context.Context, req *api.SendSetReq) (res *api.SendSetRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
