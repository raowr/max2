package settle

import (
	"context"
	v1 "user/api/settle/v1"
	"user/internal/dao"
	"user/internal/library/liberr"
	"user/internal/model/do"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedSettleServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterSettleServer(s.Server, &Controller{})
}

func (*Controller) SendSettle(ctx context.Context, req *v1.SendSettleReq) (res *v1.SendSettleRes, err error) {
	//保存结算数据到数据库记录
	data := &do.Users{
		Point: req.Point,
	}
	_, err = dao.Users.Ctx(ctx).Where("id", req.Id).Update(data)
	if err != nil {
		liberr.ErrIsNil(ctx, err, "结算失败")
	}
	return nil, nil
}
