package room

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	"game_user/internal/message"
)

// Probability 判断指定百分比概率是否触发
// percent: 0-100 的百分比
func Probability(percent int) bool {
	if percent <= 0 {
		return false
	}
	if percent >= 100 {
		return true
	}
	return grand.Intn(100) < percent
}

// ParseRedisSubscribeMessage 解析 Redis 订阅消息
// msgString: Redis 订阅消息的字符串表示
// userName: 当前操作用户名（用于日志）
// ctx: 上下文
// 返回: 解析后的 ChatMsg，是否成功
func ParseRedisSubscribeMessage(msgString string, userName string, ctx context.Context) (*message.ChatMsg, bool) {
	// 解析消息
	msgJson, err := gjson.DecodeToJson(msgString)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s 消息解码失败: %v", userName, err)
		return nil, false
	}

	// 提取 Payload 字段（Redis 订阅消息格式）
	payloadStr := msgJson.Get("Payload").String()
	if payloadStr == "" {
		g.Log().Errorf(ctx, "用户 %s Payload 为空", userName)
		return nil, false
	}

	// 解析 Payload 中的实际消息内容
	payloadJson, err := gjson.DecodeToJson(payloadStr)
	if err != nil {
		g.Log().Errorf(ctx, "用户 %s Payload 解析失败: %v", userName, err)
		return nil, false
	}

	// 解析为 ChatMsg
	msgData := &message.ChatMsg{}
	payloadJson.Scan(msgData)
	return msgData, true
}
