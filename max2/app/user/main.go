package main

import (
	_ "user/internal/boot"
	_ "user/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"user/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
