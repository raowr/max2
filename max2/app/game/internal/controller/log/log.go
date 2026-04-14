package log

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func SendLog(logInfo LogInfo) {
	LogChan <- logInfo
}
func GetLogChan() {
	ctx, cancel = context.WithCancel(context.Background())
	LogChan = make(chan LogInfo, 100)
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
