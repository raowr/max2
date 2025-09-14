package main

import (
	_ "web/app/svc-template/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"web/app/svc-template/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
