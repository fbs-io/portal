/*
 * @Author: reel
 * @Date: 2023-06-24 12:45:15
 * @LastEditors: reel
 * @LastEditTime: 2023-08-21 20:51:37
 * @Description: 业务代码，加载各个模块
 */
package app

import (
	"fbs-portal/app/auth"
	"fmt"

	"github.com/fbs-io/core"
)

func New(c core.Core) {
	// 中间件使用
	core.AddAllowResource(fmt.Sprintf("%s:%s", "POST", "/ajax/login"))
	
	c.Engine().Use(core.LimiterMiddleware(c))
	c.Engine().Use(core.CorsMiddleware(c))
	c.Engine().Use(core.SignatureMiddleware(c, core.SINGULAR_TYPE_CSRF_TOKEN))
	// c.Engine().Use(permissionMiddleware(c))
	c.Engine().Use(core.LogMiddleware(c))
	c.Engine().Use(core.ParamsMiddleware(c))

	auth.New(c)

	ajax := c.Group("ajax")
	{
		_ = ajax.Group("home", "首页").WithPermission(core.SOURCE_TYPE_UNMENU)

	}

}
