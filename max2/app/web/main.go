package main

import (
	_ "web/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	_ "web/config"
	"web/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
