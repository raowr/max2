package settle

import (
	"context"
	v1 "game_user/api/settle/v1"
)

// SettleChan 结算信息通道，用于异步处理结算信息
var (
	SettleChan   chan *v1.SendSettleReq
	ctx          context.Context
	cancel       context.CancelFunc
	settleClient v1.SettleClient // 复用的 gRPC 客户端
)
