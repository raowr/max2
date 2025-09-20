// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package room

import (
	"context"

	"web/api/room/v1"
)

type IRoomV1 interface {
	Room(ctx context.Context, req *v1.RoomReq) (res *v1.RoomRes, err error)
}
