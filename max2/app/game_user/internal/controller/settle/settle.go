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

func SendSettle(settleInfo *v1.SendSettleReq) {
	SettleChan <- settleInfo
}
func GetSettleChan() {
	ctx, cancel = context.WithCancel(context.Background())
	SettleChan = make(chan *v1.SendSettleReq, 100)
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
				settleInfo := <-SettleChan
				g.Log().Infof(ctx, "执行结算异步任务: %v", settleInfo)
				sendService(settleInfo)
			}
		}
	}()
}

// ShutdownLogtle 关闭结算信息系统
func ShutdownSettle() {
	if cancel != nil {
		cancel()
	}
}

// sendLogtleToService 发送结算信息到 settle-service
func sendService(settleInfo *v1.SendSettleReq) {
	ctx := gctx.New()
	conn := grpcx.Client.MustNewGrpcClientConn("user")
	client := v1.NewSettleClient(conn)
	_, err := client.SendSettle(ctx, settleInfo)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	g.Log().Debug(ctx, "Log sent successfully")
}
