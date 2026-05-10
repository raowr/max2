package set

import (
	"context"
	"fmt"
	"log"
	"time"

	setGrpc "api-service/api"
	v1 "api-service/api/set/v1"
)

func (c *ControllerV1) Set(ctx context.Context, req *v1.SetReq) (res *v1.SetRes, err error) {
	log.Printf("Set req: %v", req)
	setReq := &setGrpc.SendSetReq{
		Key:   req.Key,
		Value: req.Value,
	}

	if client == nil {
		return nil, fmt.Errorf("gRPC client not initialized")
	}

	// 设置一个带超时的context，避免长时间等待导致context被取消
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = client.SendSet(timeoutCtx, setReq)
	if err != nil {
		log.Printf("Failed to call SendSet: %v", err)
		return nil, err
	}

	_, err = userClient.SendSet(timeoutCtx, setReq)
	if err != nil {
		log.Printf("Failed to call SendSet: %v", err)
		return nil, err
	}

	// 正确初始化返回值
	res = &v1.SetRes{}
	return res, nil
}
