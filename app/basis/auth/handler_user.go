/*
 * @Author: reel
 * @Date: 2023-08-19 17:38:01
 * @LastEditors: reel
 * @LastEditTime: 2024-03-15 20:42:54
 * @Description: 用户信息相关接口
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

func chpwd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userChPwdParams)

		err := UserService.ChangePwd(ctx.TX(), ctx.Auth(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 个人用户信息修改
func userUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userUpdateParams)

		user, err := UserService.UpdateByAccount(ctx.Auth(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(user.UserInfo()).Notify())
	}
}

func userAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userAddParams)

		err := UserService.Create(ctx.TX(), param)
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

func usersQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersQueryParams)

		data, err := UserService.Query(ctx.TX(core.SetTxMode(core.TX_QRY_MODE_SUBID)), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

// 用户更新
func usersUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*usersEditParams)

		err := UserService.UpdateByID(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 软删除
func usersDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rdb.DeleteParams)

		err := UserService.Delete(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
