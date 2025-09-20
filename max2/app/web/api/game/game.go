// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package game

import (
	"context"

	"web/api/game/v1"
)

type IGameV1 interface {
	Game(ctx context.Context, req *v1.GameReq) (res *v1.GameRes, err error)
}
