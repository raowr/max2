package v1

import (
	"context"

	"gfast/api/v1/user"
	"gfast/internal/app/user/service"
)

func (c *ControllerUser) UserUpdatePassword(ctx context.Context, req *user.UserUpdatePasswordReq) (res *user.UserUpdatePasswordRes, err error) {
	res, err = service.User().UserUpdatePassword(ctx, req)
	return res, err
}
