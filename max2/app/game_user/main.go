package main

import (
	_ "game_user/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"game_user/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
