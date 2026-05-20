package login

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// GetGameUserEndpoints 直接从 etcd 查询 game_user.svc 的端点列表
func GetGameUserEndpoints() gsvc.Endpoints {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		g.Log().Errorf(gctx.New(), "连接 etcd 失败: %v", err)
		return nil
	}
	defer cli.Close()

	// 查询 game_user.svc 的所有端点
	prefix := "/service/default/default/game_user.svc/latest/"
	resp, err := cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		g.Log().Errorf(gctx.New(), "查询 etcd 失败: %v", err)
		return nil
	}

	if len(resp.Kvs) == 0 {
		g.Log().Warning(gctx.New(), "未找到 game_user.svc 服务")
		return nil
	}

	// 提取所有端点地址
	endpoints := make(gsvc.Endpoints, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		addr := strings.TrimPrefix(string(kv.Key), prefix)
		endpoints = append(endpoints, gsvc.NewEndpoint(addr))
	}

	return endpoints
}
