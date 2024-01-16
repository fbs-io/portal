/*
 * @Author: reel
 * @Date: 2023-08-20 15:42:11
 * @LastEditors: reel
 * @LastEditTime: 2024-01-17 00:12:16
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

		// 鉴权中间件已经校验过session, 并通过session获取用户名保存至上下文中
		// 如果无法获取上下文中的用户, 则返回签名错误, 前端调转到用户登陆
		// TODO:校验用户是否存在
		accountI, ok := ctx.Get(core.CTX_AUTH)
		if !ok {
			ctx.JSON(200, errno.ERRNO_AUTH_PERMISSION.ToMap())
			ctx.Abort()
			return
		}
		account := accountI.(string)

		// 通过用户名获取缓存的用户所属公司
		// TODO:如果后去不到,则增加默认公司
		company := c.Cache().Get(consts.GenUserCompanyKey(account))

		// 公司code作为关键的数据分区值, 缓存到上下文, 同时给orm层使用
		ctx.Set(core.CTX_SHARDING_KEY, company)

		// 刷新用户, 获取最新的用户信息
		// TODO:用户重新创建缓存
		user := auth.GetUser(account, core.NewCtx(c, ctx), auth.REFRESH_NOT)
		auth.SetUser(user.Account, user)
		permissionKey := genPermissionKey(ctx)

		// 处理数据权限
		// 根据用户和当前缓存的公司, 获取已缓存的岗位信息
		position_code := c.Cache().Get(consts.GenUserPositionKey(user.Account, company))
		// 定义的数据权限数据结构, 默认只读个人创建的
		dp := &rdb.DataPermissionStringCtx{
			DataPermissionType: rdb.DATA_PERMISSION_ONESELF,
		}
		tx := c.RDB().DB().Where("1=1").Set(core.CTX_SHARDING_KEY, company)

		// 如果岗位存在且有效, 则获取岗位的权限并更新dp
		position := org.PositionSrvice.GetByCode(tx, []string{position_code})
		if position != nil && len(position) > 0 {
			// 设置用户所在岗位的部门(数据权限)以及设置数据权限类型
			dp.DataPermission = position[0].DepartmentCode
			dp.DataPermissionType = position[0].DataPermissionType
			dp.DataPermissionScope = make([]string, 0, 100)
			switch position[0].DataPermissionType {

			// 数据权限: 当前部门及子部门可见
			case rdb.DATA_PERMISSION_ONLY_DEPT_ALL:
				depts := org.DepartmentSrvice.GetAllChildren(tx, position[0].DepartmentCode)
				for key := range depts {
					dp.DataPermissionScope = append(dp.DataPermissionScope, key)
				}
			case rdb.DATA_PERMISSION_ONLY_CUSTOM:
				for _, dept_code := range position[0].DataPermissionCustom {
					dp.DataPermissionScope = append(dp.DataPermissionScope, dept_code.(string))
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
