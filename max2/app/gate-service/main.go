package main

import (
	_ "gate-service/internal/boot"
	_ "gate-service/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"gate-service/internal/cmd"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
