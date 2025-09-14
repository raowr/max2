package main

import (
	_ "web/app/job-template/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"web/app/job-template/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
