package action

import (
	"context"
	v1 "gate-service/api/action/v1"
)

// ActionChan 操作信息通道，用于异步处理操作信息
var (
	actionChan   chan *v1.SendActionReq
	ctx          context.Context
	cancel       context.CancelFunc
	actionClient v1.ActionClient // 复用的 gRPC 客户端
)
