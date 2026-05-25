package action

import (
	"context"
	v1 "gate-service/api/action/v1"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type Controller struct {
	v1.UnimplementedActionServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterActionServer(s.Server, &Controller{})
}

func (*Controller) SendAction(ctx context.Context, req *v1.SendActionReq) (res *v1.SendActionRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

// 初始化
func init() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
	ctx, cancel = context.WithCancel(context.Background())
	actionChan = make(chan *v1.SendActionReq, 100)

	// 创建复用的 gRPC 连接
	conn := grpcx.Client.MustNewGrpcClientConn("game_user")
	actionClient = v1.NewActionClient(conn)

	// 启动异步任务
	go processActionChan()
}

func processActionChan() {
	for {
		select {
		case <-ctx.Done():
			g.Log().Infof(ctx, "结算异步任务通道收到退出信号，正在清理...")
			return
		case actionInfo := <-actionChan:
			g.Log().Infof(ctx, "执行操作异步任务: %v", actionInfo)
			sendService(actionInfo)
		}
	}
}
func sendService(actionInfo *v1.SendActionReq) {
	_, err := actionClient.SendAction(gctx.New(), actionInfo)
	if err != nil {
		g.Log().Error(ctx, err)
	}
}

func SendAction(actionInfo *v1.SendActionReq) {
	select {
	case actionChan <- actionInfo:
		// 发送成功
	default:
		// 通道满，丢弃操作信息
		g.Log().Warning(ctx, "操作通道已满，丢弃操作信息")
	}
}

func ShutdownAction() {
	if cancel != nil {
		cancel()
	}
}
