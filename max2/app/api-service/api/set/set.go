// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package set

import (
	"context"

	"api-service/api/set/v1"
)

type ISetV1 interface {
	Set(ctx context.Context, req *v1.SetReq) (res *v1.SetRes, err error)
}
