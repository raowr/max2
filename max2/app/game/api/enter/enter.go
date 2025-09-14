// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package enter

import (
	"context"

	"game/api/enter/v1"
)

type IEnterV1 interface {
	Enter(ctx context.Context, req *v1.EnterReq) (res *v1.EnterRes, err error)
}
