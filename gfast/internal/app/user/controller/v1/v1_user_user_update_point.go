package v1

import (
	"context"

	"gfast/api/v1/user"
	"gfast/internal/app/user/service"
)

func (c *ControllerUser) UserUpdatePoint(ctx context.Context, req *user.UserUpdatePointReq) (res *user.UserUpdatePointRes, err error) {
	res, err = service.User().UserUpdatePoint(ctx, req)
	return res, err
}
