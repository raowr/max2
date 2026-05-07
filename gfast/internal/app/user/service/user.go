// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gfast/api/v1/user"
)

type (
	IUser interface {
		GetUserListSearch(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error)
		UserUpdatePoint(ctx context.Context, req *user.UserUpdatePointReq) (res *user.UserUpdatePointRes, err error)
		UserUpdatePassword(ctx context.Context, req *user.UserUpdatePasswordReq) (res *user.UserUpdatePasswordRes, err error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
