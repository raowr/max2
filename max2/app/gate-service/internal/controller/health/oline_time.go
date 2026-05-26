package health

import (
	"context"
	"fmt"
	"gate-service/internal/consts"
	"gate-service/internal/service"
	"time"
)

// 假设每次心跳累加10秒后得到 newTotal
const (
	Warn8Hour      = 8 * 3600 // 8小时阈值
	RemindInterval = 30 * 60  // 30分钟间隔
)

// AddPlayTimeForUser 为用户增加今日在线时长（秒），返回最新累计值
func AddPlayTimeForUser(ctx context.Context, userName string, deltaSec int64) (bool, int64) {
	// key 格式：playtime:用户名:2026-05-26
	today := time.Now().Format("2006-01-02")
	key := fmt.Sprintf(consts.PlayTimePrefix+"%s:%s", userName, today)

	//先获取当前值
	valVal := service.Cache().Get(ctx, key).Int64()
	newTotalTime := valVal + deltaSec
	// 原子增加
	service.Cache().Set(ctx, key, newTotalTime, 48*time.Hour)

	// 获取当天最后提醒时间（Redis存储）
	lastRemindKey := fmt.Sprintf(consts.LastRemindPrefix+"%s:%s", userName, today)
	lastRemind := service.Cache().Get(ctx, lastRemindKey).Int64()

	if valVal >= Warn8Hour {
		// 如果从未提醒过（lastRemind==0）或距离上次提醒超过30分钟
		if lastRemind == 0 || time.Now().Unix()-lastRemind >= RemindInterval {
			// 更新提醒时间
			service.Cache().Set(ctx, lastRemindKey, time.Now().Unix(), 48*time.Hour)
			// 是否提醒
			return true, valVal / 3600
		}
	}
	return false, valVal / 3600
}
