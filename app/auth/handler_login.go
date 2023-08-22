/*
 * @Author: reel
 * @Date: 2023-07-18 21:46:02
 * @LastEditors: reel
 * @LastEditTime: 2023-08-21 23:48:11
 * @Description: 请填写简介
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/session"
)

type loginParams struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required" conditions:"-"`
}

// 登录时, 返回用户相关信息
// 如操作菜单, 权限等
func login() core.HandlerFunc {
	return func(ctx core.Context) {

		p := ctx.CtxGetParams().(*loginParams)
		tx := ctx.TX()
		user := &User{}
		err := tx.Find(user).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_AUTH_USER_OR_PWD.ToMap())
			return
		}
		err = user.CheckPwd(p.Password)

		if err != nil {
			ctx.JSON(errno.ERRNO_AUTH_USER_OR_PWD.ToMap())
			return
		}
		// session 设置
		sessionKey := session.GenSessionKey()
		ctx.Core().Session().SetWithCsrfToken(ctx.Ctx().Writer, sessionKey, user.Account)

		// 菜单获取
		menu, _ := getMenuTree(ctx.Core(), user)

		// 权限获取及绑定
		permissionMap, permissionList, _ := getPermission(ctx.Core(), user)
		user.Permissions = permissionMap

		result := map[string]interface{}{
			"token":       sessionKey,
			"userInfo":    user.UserInfo(),
			"menu":        menu,
			"permissions": permissionList,
		}

		SetUser(user.ID, user)
		ctx.JSON(errno.ERRNO_OK.ToMapWithData(result))
	}
}
