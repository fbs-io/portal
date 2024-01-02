/*
 * @Author: reel
 * @Date: 2023-08-20 15:42:11
 * @LastEditors: reel
 * @LastEditTime: 2024-01-02 22:49:42
 * @Description: 权限校验中间件
 */
package app

import (
	"fmt"
	"portal/app/basis/auth"
	"portal/app/basis/org"
	"portal/pkg/consts"
	"strings"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
	"github.com/gin-gonic/gin"
)

func genPermissionKey(ctx *gin.Context) string {
	return fmt.Sprintf("%s%s", strings.ToLower(ctx.Request.Method), strings.ReplaceAll(ctx.FullPath(), "/", ":"))
}

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
		company := c.Cache().Get(consts.GenUserCompanyKey(account))
		ctx.Set(core.CTX_SHARDING_KEY, company)
		user := auth.GetUser(account, core.NewCtx(c, ctx), auth.REFRESH_NOT)
		auth.SetUser(user.Account, user)
		permissionKey := genPermissionKey(ctx)

		// 处理数据权限
		position_code := c.Cache().Get(consts.GenUserPositionKey(user.Account, company))
		dp := &rdb.DataPermissionStringCtx{
			DataPermissionType: rdb.DATA_PERMISSION_ONESELF,
		}

		if position_code != "" {
			position := org.GetPosition(company, position_code, *c.RDB().DB())
			if position != nil {
				// 设置用户所在岗位的部门(数据权限)以及设置数据权限类型
				dp.DataPermission = position.DepartmentCode
				dp.DataPermissionType = position.DataPermissionType
				dp.DataPermissionScope = make([]string, 0, 100)
				switch position.DataPermissionType {

				// 数据权限: 当前部门及子部门可见
				case rdb.DATA_PERMISSION_ONLY_DEPT_ALL:
					depts := org.GetDeptAndAllChildren(company, position.DepartmentCode, c.RDB().DB())
					for key := range depts {
						dp.DataPermissionScope = append(dp.DataPermissionScope, key)
					}
				case rdb.DATA_PERMISSION_ONLY_CUSTOM:
					for _, dept_code := range position.DataPermissionCustom {
						dp.DataPermissionScope = append(dp.DataPermissionScope, dept_code.(string))
					}
				}
			}
		}
		ctx.Set(core.CTX_DATA_PERMISSION_KEY, dp)
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
