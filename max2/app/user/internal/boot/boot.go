package boot

import (
	"user/internal/library/libBloom"

	_ "user/internal/logic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	ctx := gctx.New()
	if err := libBloom.BloomWarmUp(ctx); err != nil {
		g.Log().Error(ctx, "布隆过滤器预热失败:", err)
	} else {
		g.Log().Info(ctx, "布隆过滤器预热完成")
	}
}
