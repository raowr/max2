package main

import (
	_ "api-service/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"api-service/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
