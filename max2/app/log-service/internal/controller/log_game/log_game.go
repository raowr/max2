package log_game

import (
	"context"
	v1 "log-service/api/log_game/v1"
	"log-service/internal/dao"
	"log-service/internal/model/do"
	"time"

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

	// 创建一个新的上下文，不受原始 ctx 取消的影响
	logCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 复制请求数据，避免引用问题
	reqData := req
	go func() {
		// 在 goroutine 中写入数据库
		_, err := dao.LogGame.Ctx(logCtx).Data(do.LogGame{
			RoomId:    reqData.GetRoomID(),
			Type:      reqData.GetType(),
			Status:    reqData.GetStatus(),
			UserId:    reqData.GetUserID(),
			Point:     reqData.GetPoint(),
			Action:    reqData.GetAction(),
			Remain:    reqData.GetRemain(),
			OutCards:  reqData.GetOutCardIds(),
			Text:      reqData.GetText(),
			CreatTime: currentTime,
		}).Insert()

		if err != nil {
			g.Log().Error(logCtx, "日志写入失败:", err)
		}
	}()

	return nil, nil
}
