package v1

import (
	"context"

	"gfast/api/v1/user"
	"gfast/internal/app/user/service"
)

func (c *ControllerUser) UserList(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error) {
	res, err = service.User().GetUserListSearch(ctx, req)
	return res, err
}
