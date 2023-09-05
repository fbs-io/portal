/*
 * @Author: reel
 * @Date: 2023-08-20 15:42:11
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 19:59:43
 * @Description: 权限校验中间件
 */
package app

import (
	"fbs-portal/app/basis/auth"
	"fmt"
	"strings"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/gin-gonic/gin"
)

func permissionMiddleware(c core.Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if core.GetAllowResource(ctx) {
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
