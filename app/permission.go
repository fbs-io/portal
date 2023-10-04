/*
 * @Author: reel
 * @Date: 2023-08-20 15:42:11
 * @LastEditors: reel
 * @LastEditTime: 2023-10-04 21:48:26
 * @Description: 权限校验中间件
 */
package app

import (
	"fmt"
	"portal/app/basis/auth"
	"strings"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/gin-gonic/gin"
)

func permissionMiddleware(c core.Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if core.GetAllowSource(ctx) {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.RequestURI, core.STATIC_PATH_PREFIX) {
			ctx.Next()
			return
		}
		accountI, ok := ctx.Get(core.CTX_AUTH)
		if !ok {
			ctx.JSON(200, errno.ERRNO_AUTH_PERMISSION.ToMap())
			ctx.Abort()
			return
		}
		account := accountI.(string)
		user := auth.GetUser(account, c.RDB().DB())
		permissionKey := fmt.Sprintf("%s%s", strings.ToLower(ctx.Request.Method), strings.ReplaceAll(ctx.FullPath(), "/", ":"))
		if user.Permissions[permissionKey] {
			ctx.Next()
			return
		}
		ctx.JSON(200, errno.ERRNO_AUTH_PERMISSION.ToMap())
		ctx.Abort()

	}
}

type uiPermissionParams struct {
	Path string `form:"path" binding:"required" conditions:"like%"`
}

func uiPermission() core.HandlerFunc {
	return func(ctx core.Context) {
		tx := ctx.TX()
		soureces := make([]core.Sources, 0, 10)
		err := tx.Find(&soureces).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		result := make(map[string]interface{}, 10)
		for _, item := range soureces {
			if item.Method == "" {
				continue
			}
			result[item.Method] = item.Code
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(result))

	}
}
