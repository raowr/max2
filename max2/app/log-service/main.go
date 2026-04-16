package main

import (
	_ "log-service/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"log-service/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
