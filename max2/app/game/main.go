package main

import (
	_ "game/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"game/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
