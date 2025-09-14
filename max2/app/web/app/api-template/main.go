package main

import (
	_ "web/app/api-template/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"web/app/api-template/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
