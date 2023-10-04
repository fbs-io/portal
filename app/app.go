/*
 * @Author: reel
 * @Date: 2023-06-24 12:45:15
 * @LastEditors: reel
 * @LastEditTime: 2023-10-05 00:15:55
 * @Description: 业务代码，加载各个模块
 */
package app

import (
	"fmt"
	"net/http"
	"portal/app/basis"
	ui "portal/ui/dist"

	"github.com/fbs-io/core"
	"github.com/gin-gonic/gin"
)

func New(c core.Core) {
	// 中间件使用
	core.STATIC_PATH_PREFIX = "/website/"
	// c.Engine().GET(core.STATIC_PATH_PREFIX+"*filepath", func(ctx *gin.Context) {
	// 	staticSrv := http.FileServer(http.FS(ui.Static))
	// 	staticSrv.ServeHTTP(ctx.Writer, ctx.Request)
	// })
	core.AddAllowSource(fmt.Sprintf("%s:%s", "GET", "/"))
	core.AddAllowSource(fmt.Sprintf("%s:%s", "GET", ""))

	c.Engine().Use(core.LimiterMiddleware(c))
	c.Engine().Use(core.CorsMiddleware(c))
	c.Engine().Use(core.SignatureMiddleware(c, core.SINGULAR_TYPE_CSRF_TOKEN))
	c.Engine().Use(permissionMiddleware(c))
	c.Engine().Use(core.LogMiddleware(c))
	c.Engine().Use(core.ParamsMiddleware(c))

	// 权限校验

	// 加载静态资源
	// static, _ := fs.Sub(ui.Static, core.STATIC_PATH_PREFIX)
	// c.Engine().StaticFS("/website", http.FS(static))
	c.Engine().GET(core.STATIC_PATH_PREFIX+"*filepath", func(ctx *gin.Context) {
		staticSrv := http.FileServer(http.FS(ui.Static))
		staticSrv.ServeHTTP(ctx.Writer, ctx.Request)
	})
	// fs := http.FileServer(http.FS(ui.Static))
	c.Engine().GET("/", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html;charset=utf-8")
		ctx.String(200, string(ui.Index))
	})
	// http.Handle("/", http.StripPrefix("/", fs))
	// 加载
	ajax := c.Group("ajax").WithPermission(core.SOURCE_TYPE_LIMITED)
	{
		_ = ajax.Group("home", "首页").WithPermission(core.SOURCE_TYPE_UNMENU).WithMeta("icon", "el-icon-help-filled")

		ajax.GET("uipermission", "页面权限", uiPermissionParams{}, uiPermission()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
	}

	// 初始化basis模块
	basis.New(ajax)

}
