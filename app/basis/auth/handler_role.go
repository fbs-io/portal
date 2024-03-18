/*
 * @Author: reel
 * @Date: 2023-08-30 07:36:25
 * @LastEditors: reel
 * @LastEditTime: 2024-01-20 21:24:07
 * @Description: 角色相关api逻辑
 */
package auth

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
)

// const (
// 	DATA_PERMISSION_ONESELF       int8 = iota + 1 //本人可见
// 	DATA_PERMISSION_ALL                           //全部可见
// 	DATA_PERMISSION_ONLY_DEPT                     //所在部门可见
// 	DATA_PERMISSION_ONLY_DEPT_ALL                 //所在部门及子级可见
// 	DATA_PERMISSION_ONLY_CUSTOM                   //选择的部门可见
// )

// TODO:补充其他用户信息, 如部门等
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户

func roleAdd(roleSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*roleAddParams)
		if param.Code == "" {
			param.Code = roleSeq.Code()
		}
		if param.Code == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无法获取或生成岗位代码")))
			return
		}

		err := RoleService.Create(ctx.TX(), param)
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

func rolesQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rolesQueryParams)
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
		)
		data, err := RoleService.Query(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

func rolesUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*roleEditParams)

		err := RoleService.UpdateByID(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 软删除
func rolesDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rdb.DeleteParams)

		err := RoleService.Delete(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 菜单和权限查询, 返回树表结构
func menusQueryWithPermission() core.HandlerFunc {
	return func(ctx core.Context) {
		// user := GetUser(ctx.Auth(), ctx, REFRESH_NOT)

		// menus, _, err := ResourceService.GetSource(ctx, user.Roles, QUERY_MENU_MODE_MANAGE)
		// if err != nil {
		// 	ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
		// 	return
		// }
		// ctx.JSON(errno.ERRNO_OK.WrapData(menus))
	}
}
