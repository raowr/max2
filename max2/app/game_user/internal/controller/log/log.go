package log

import (
	"context"
	v1 "game_user/api/log_game/v1"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func SendLog(logInfo LogInfo) {
	LogChan <- logInfo
}
func GetLogChan() {
	ctx, cancel = context.WithCancel(context.Background())
	LogChan = make(chan LogInfo, 100)
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
	// 启动异步任务
	go func() {
		for {
			select {
			case <-ctx.Done():
				g.Log().Infof(ctx, "日志异步任务通道收到退出信号，正在清理...")
				return
			default:
				// 执行异步操作
				logInfo := <-LogChan
				g.Log().Infof(ctx, "执行日志异步任务: %v", logInfo)
				sendLogToService(logInfo)
			}
		}
	}()
}

// ShutdownLog 关闭日志系统
func ShutdownLog() {
	if cancel != nil {
		cancel()
	}
}

// sendLogToService 发送日志到 log-service
func sendLogToService(logInfo LogInfo) {
	ctx := gctx.New()
	conn := grpcx.Client.MustNewGrpcClientConn("log-service")
	client := v1.NewLogGameClient(conn)

	req := &v1.SendLogReq{
		RoomID:     logInfo.RoomID,
		Type:       int32(logInfo.Type),
		Status:     int32(logInfo.Status),
		UserID:     logInfo.UserID,
		Point:      logInfo.Point,
		Action:     logInfo.Action,
		Remain:     logInfo.Remain,
		OutCardIds: logInfo.OutCardIds,
		Text:       logInfo.Text,
	}

	_, err := client.SendLog(ctx, req)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Debug(ctx, "Log sent successfully")
}
