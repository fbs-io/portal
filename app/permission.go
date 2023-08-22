/*
 * @Author: reel
 * @Date: 2023-08-20 15:42:11
 * @LastEditors: reel
 * @LastEditTime: 2023-08-20 23:35:53
 * @Description: 权限校验中间件
 */
package app

import (
	"github.com/fbs-io/core"
	"github.com/gin-gonic/gin"
)

func permissionMiddleware(c core.Core) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
