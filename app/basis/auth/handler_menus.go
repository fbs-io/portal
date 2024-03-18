/*
 * @Author: reel
 * @Date: 2023-08-30 07:36:25
 * @LastEditors: reel
 * @LastEditTime: 2024-03-13 06:31:32
 * @Description: 菜单相关api逻辑
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

func menusAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*menusAddParams)
		tx := ctx.TX()

		err := ResourceService.Create(tx, param)
		if err != nil {
			if rdb.IsUniqueError(err) {
				ctx.JSON(errno.ERRNO_RDB_DUPLICATED_KEY.WrapError(err))
				return
			}
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())

	}
}

// 菜单查询, 返回树表结构
// 根据不同用户, 返回不同的菜单列表
func menusQueryWithManager() core.HandlerFunc {
	return func(ctx core.Context) {
		_, ManageMenus, _, err := UserService.GetResourcePermission(ctx.NewTX(), ctx.Auth())
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(ManageMenus))
	}
}

func menusUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*menusUpdateParams)
		tx := ctx.TX()
		err := ResourceService.Update(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 软删除
func menusDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rdb.DeleteParams)

		err := ResourceService.Delete(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
