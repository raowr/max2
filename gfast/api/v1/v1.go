// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package v1

import (
	"context"
	"gfast/api/v1/log"
)

type IV1Log interface {
	LogGameList(ctx context.Context, req *log.LogGameListReq) (res *log.LogGameListRes, err error)
}
