package controller

import (
	"context"

	"gfast/api/v1/set"
	"gfast/internal/app/set/service"
)

func (c *ControllerSet) SetGameList(ctx context.Context, req *set.SetGameListReq) (res *set.SetGameListRes, err error) {
	return service.SetGame().GetSetGameListSearch(ctx, req)
}
