package log_game

import (
	"context"
	v1 "game_user/api/log_game/v1"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
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

// 全局变量
var (
	logClient v1.LogGameClient // 复用的 gRPC 客户端
)

// 初始化
func init() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
	ctx, cancel = context.WithCancel(context.Background())
	LogChan = make(chan *v1.SendLogReq, 100)

	// 创建复用的 gRPC 连接
	conn := grpcx.Client.MustNewGrpcClientConn("log-service")
	logClient = v1.NewLogGameClient(conn)

	// 启动异步任务
	go processLogChan()
}

func processLogChan() {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "日志异步任务通道收到退出信号，正在清理...")
			return
		case logInfo := <-LogChan: // 直接从通道读取，非阻塞时会等待
			g.Log().Infof(ctx, "执行日志异步任务: %v", logInfo)
			sendLogToService(logInfo)
		}
	}
}

func SendLog(logInfo *v1.SendLogReq) {
	select {
	case LogChan <- logInfo:
		// 发送成功
	default:
		// 通道满，丢弃日志
		g.Log().Warning(ctx, "日志通道已满，丢弃日志")
	}
}

func ShutdownLog() {
	if cancel != nil {
		cancel()
	}
}

func sendLogToService(logInfo *v1.SendLogReq) {
	_, err := logClient.SendLog(gctx.New(), logInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
}
