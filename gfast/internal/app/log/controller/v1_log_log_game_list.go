package controller

import (
	"context"

	"gfast/api/v1/log"
	"gfast/internal/app/log/service"
)

func (c *ControllerLog) LogGameList(ctx context.Context, req *log.LogGameListReq) (res *log.LogGameListRes, err error) {
	res, err = service.LogGame().GetLogGameListSearch(ctx, req)
	return res, err
}
