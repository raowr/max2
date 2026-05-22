package login

import (
	"context"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GetRandomGameUserNode 直接从 etcd 查询 game_user.svc 的随机节点地址
func GetRandomGameUserNode() string {
	rand.Seed(time.Now().UnixNano())

	// 连接 etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		g.Log().Errorf(gctx.New(), "连接 etcd 失败: %v", err)
		return ""
	}
	defer cli.Close()

	// 查询服务（go-zero 写入的 key 格式是 /discov/service_name/address/lease_id）
	prefix := "/discov/game_user.svc/"
	resp, err := cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		g.Log().Errorf(gctx.New(), "查询 etcd 失败: %v", err)
		return ""
	}

	if len(resp.Kvs) == 0 {
		g.Log().Warning(gctx.New(), "未找到 game_user.svc 服务")
		return ""
	}

	// 解析节点地址（去重）
	var nodes []string
	seen := make(map[string]bool)

	for _, kv := range resp.Kvs {
		// key := string(kv.Key)
		value := string(kv.Value)

		// key 格式: /discov/game_user.svc/8.155.147.137:8010/lease_id
		// value 就是地址: 8.155.147.137:8010
		if !seen[value] {
			seen[value] = true
			nodes = append(nodes, value)
			g.Log().Debugf(gctx.New(), "发现服务节点: %s", value)
		}
	}

	if len(nodes) == 0 {
		g.Log().Warning(gctx.New(), "game_user.svc 服务没有可用节点")
		return ""
	}

	return nodes[rand.Intn(len(nodes))]
}
