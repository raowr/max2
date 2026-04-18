// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gfast/api/v1/set"
)

type (
	ISetGame interface {
		GetSetGameListSearch(ctx context.Context, req *set.SetGameListReq) (res *set.SetGameListRes, err error)
		Update(ctx context.Context, req *set.SetGameUpdateReq) (err error)
	}
)

var (
	localSetGame ISetGame
)

func SetGame() ISetGame {
	if localSetGame == nil {
		panic("implement not found for interface ISetGame, forgot register?")
	}
	return localSetGame
}

func RegisterSetGame(i ISetGame) {
	localSetGame = i
}
