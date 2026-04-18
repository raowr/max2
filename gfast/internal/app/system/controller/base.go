/*
* @desc:system base controller
* @company:云南奇讯科技有限公司
* @Author: yixiaohu
* @Date:   2022/3/4 18:12
 */

package controller

import (
	commonController "gfast/internal/app/common/controller"

	"github.com/gogf/gf/v2/net/ghttp"
)

type BaseController struct {
	commonController.BaseController
}

// Init 自动执行的初始化方法
func (c *BaseController) Init(r *ghttp.Request) {
	c.BaseController.Init(r)
}
