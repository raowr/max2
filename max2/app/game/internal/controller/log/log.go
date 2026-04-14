package log

import (
	"context"
)

func SendLog(logInfo LogInfo) {
	LogChan <- logInfo
}
func GetLogChan() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	LogChan = make(chan LogInfo, 100)
	// 启动异步任务
	go func() {
		for {
			select {
			case <-ctx.Done():
				println("日志异步任务通道收到退出信号，正在清理...")
				return
			default:
				// 执行异步操作
				logInfo := <-LogChan
				println("执行日志异步任务", logInfo)
			}
		}
	}()
}
