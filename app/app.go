/*
 * @Author: reel
 * @Date: 2023-06-24 12:45:15
 * @LastEditors: reel
 * @LastEditTime: 2024-03-26 06:26:15
 * @Description: 业务代码，加载各个模块
 */
package app

import (
	"fmt"
	"net/http"
	"portal/app/basis"
	"portal/app/financial"
	"portal/pkg/consts"
	"portal/ui/dist"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
	"github.com/gin-gonic/gin"
)

func New(c core.Core) {
	// 使用数据分区模式
	c.RDB().AddMigrateList(func() error {
		result := make([]*dimResult, 0, 10)
		tx := c.RDB().DB()
		tx = tx.Table(consts.TABLE_BASIS_ORG_COMPANY).Select("company_code as code", "company_name as name", "status")
		tx.Where("status > 0 ").Find(&result)
		suffixList := make([]interface{}, 0, len(result))
		for _, item := range result {
			suffixList = append(suffixList, item.Code)
		}
		c.RDB().SetShardingModel(rdb.SHADING_MODEL_DB, suffixList)
		return nil
	})
	// 中间件使用
	core.STATIC_PATH_PREFIX = "/website/"

	core.AddAllowSource(fmt.Sprintf("%s:%s", "GET", "/"))
	core.AddAllowSource(fmt.Sprintf("%s:%s", "GET", ""))

	c.Engine().Use(core.LimiterMiddleware(c))
	c.Engine().Use(core.CorsMiddleware(c))
	c.Engine().Use(core.SignatureMiddleware(c))
	c.Engine().Use(permissionMiddleware(c))
	c.Engine().Use(core.LogMiddleware(c))
	c.Engine().Use(core.ParamsMiddleware(c))

	// 权限校验

	// 加载静态资源
	c.Engine().GET(core.STATIC_PATH_PREFIX+"*filepath", func(ctx *gin.Context) {
		staticSrv := http.FileServer(http.FS(dist.Static))
		staticSrv.ServeHTTP(ctx.Writer, ctx.Request)
	})
	c.Engine().GET("/", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html;charset=utf-8")
		ctx.String(200, string(dist.Index))
	})

	c.Engine().NoRoute(func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html;charset=utf-8")
		ctx.String(200, string(dist.Index))
		ctx.Writer.Flush()
	})
	// 加载
	ajax := c.Group("ajax").WithPermission(core.SOURCE_TYPE_LIMITED)
	{
		_ = ajax.Group("home", "首页").WithPermission(core.SOURCE_TYPE_UNMENU).WithMeta("icon", "el-icon-help-filled")

		ajax.GET("uipermission", "页面权限", uiPermissionParams{}, uiPermission()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
		ajax.GET("dimension", "维度列表", dimensionParams{}, dimList()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
	}

	// 初始化basis模块
	basis.New(ajax)

	// 初始化财务模块
	financial.New(ajax)
}
