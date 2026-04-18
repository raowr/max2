/*
* @desc:后台路由
* @company:max2
* @Author: raowenran
* @Date:   2026/4/18 17:34
 */

package router

import (
	"context"
	"gfast/internal/app/set/controller"
	"gfast/internal/app/system/service"
	"gfast/library/libRouter"

	"github.com/gogf/gf/v2/net/ghttp"
)

var R = new(Router)

type Router struct{}

func (router *Router) BindController(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/set", func(group *ghttp.RouterGroup) {
		//登录验证拦截
		service.GfToken().Middleware(group)
		//context拦截器
		group.Middleware(service.Middleware().Ctx, service.Middleware().Auth)
		//后台操作日志记录
		group.Hook("/*", ghttp.HookAfterOutput, service.OperateLog().OperationLog)
		group.Bind(
			controller.NewSet(),
		)
		//自动绑定定义的控制器
		if err := libRouter.RouterAutoBind(ctx, router, group); err != nil {
			panic(err)
		}
	})
}
