package settle

import (
	"context"
	v1 "game_user/api/settle/v1"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type Controller struct {
	v1.UnimplementedSettleServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterSettleServer(s.Server, &Controller{})
}

func (*Controller) SendSettle(ctx context.Context, req *v1.SendSettleReq) (res *v1.SendSettleRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

// 初始化
func init() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
	ctx, cancel = context.WithCancel(context.Background())
	SettleChan = make(chan *v1.SendSettleReq, 100)

	// 创建复用的 gRPC 连接
	conn := grpcx.Client.MustNewGrpcClientConn("user")
	settleClient = v1.NewSettleClient(conn)

	// 启动异步任务
	go processSettleChan()
}

func processSettleChan() {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "结算异步任务通道收到退出信号，正在清理...")
			return
		case settleInfo := <-SettleChan:
			g.Log().Infof(ctx, "执行结算异步任务: %v", settleInfo)
			sendService(settleInfo)
		}
	}
}

func SendSettle(settleInfo *v1.SendSettleReq) {
	select {
	case SettleChan <- settleInfo:
		// 发送成功
	default:
		// 通道满，丢弃结算信息
		g.Log().Warning(ctx, "结算通道已满，丢弃结算信息")
	}
}

func ShutdownSettle() {
	if cancel != nil {
		cancel()
	}
}

func sendService(settleInfo *v1.SendSettleReq) {
	_, err := settleClient.SendSettle(gctx.New(), settleInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
}
