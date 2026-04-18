// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"gfast/api/v1/log"
)

type (
	ILogGame interface {
		GetLogGameListSearch(ctx context.Context, req *log.LogGameListReq) (res *log.LogGameListRes, err error)
	}
)

var (
	localLogGame ILogGame
)

func LogGame() ILogGame {
	if localLogGame == nil {
		panic("implement not found for interface ILogGame, forgot register?")
	}
	return localLogGame
}

func RegisterLogGame(i ILogGame) {
	localLogGame = i
}
