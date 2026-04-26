package main

import (
	_ "user/internal/boot"
	_ "user/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"user/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
