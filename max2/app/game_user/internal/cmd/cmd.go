package cmd

import (
	"context"
	"sync"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/gcmd"
	"google.golang.org/grpc"

	"game_user/internal/controller/action"
	"game_user/internal/controller/log_game"
	"game_user/internal/controller/set_game"
	"game_user/internal/controller/settle"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http and grpc servers",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			var wg sync.WaitGroup

			set_game.InitCache()

			// 2. gRPC 服务器
			wg.Add(1)
			go func() {
				defer wg.Done()
				c := grpcx.Server.NewConfig()
				if c.Address == "" {
					c.Address = ":8020"
				}
				grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))
				c.Options = append(c.Options, []grpc.ServerOption{
					grpcx.Server.ChainUnary(grpcx.Server.UnaryValidate),
				}...)
				gs := grpcx.Server.New(c)
				set_game.Register(gs)
				action.Register(gs)
				gs.Run()
			}()

			wg.Wait()

			log_game.ShutdownLog()
			settle.ShutdownSettle()
			return nil
		},
	}
)
