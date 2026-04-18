package controller

import (
	"context"

	"gfast/api/v1/set"
	"gfast/internal/app/set/service"
)

func (c *ControllerSet) SetGameUpdate(ctx context.Context, req *set.SetGameUpdateReq) (res *set.SetGameUpdateRes, err error) {
	err = service.SetGame().Update(ctx, req)
	return
}
