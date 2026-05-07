// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package v1

import (
	"context"

	"gfast/api/v1/log"
	"gfast/api/v1/set"
	"gfast/api/v1/user"
)

type IV1Log interface {
	LogGameList(ctx context.Context, req *log.LogGameListReq) (res *log.LogGameListRes, err error)
}

type IV1Set interface {
	SetGameList(ctx context.Context, req *set.SetGameListReq) (res *set.SetGameListRes, err error)
	SetGameUpdate(ctx context.Context, req *set.SetGameUpdateReq) (res *set.SetGameUpdateRes, err error)
}

type IV1User interface {
	UserList(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error)
	UserUpdatePoint(ctx context.Context, req *user.UserUpdatePointReq) (res *user.UserUpdatePointRes, err error)
	UserUpdatePassword(ctx context.Context, req *user.UserUpdatePasswordReq) (res *user.UserUpdatePasswordRes, err error)
}
