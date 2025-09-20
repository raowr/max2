package main

import (
	_ "max2/app/job-template/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"max2/app/job-template/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
