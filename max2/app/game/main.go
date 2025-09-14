package main

import (
	_ "game/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"game/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
