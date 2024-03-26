/*
 * @Author: reel
 * @Date: 2023-07-18 21:46:02
 * @LastEditors: reel
 * @LastEditTime: 2024-03-26 20:06:16
 * @Description: 请填写简介
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/session"
)

type loginParams struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required" conditions:"-"`
	Company  string `json:"company_code" conditions:"-"`
}

// 登录时, 返回用户相关信息
// 如操作菜单, 权限等
func login() core.HandlerFunc {
	return func(ctx core.Context) {
		p := ctx.CtxGetParams().(*loginParams)
		// 获取公司code作为分区字段
		if p.Company == "" {
			// 如果公司code为空, 先尝试从缓存中获取上次存在的公司code
			p.Company = ctx.Core().Cache().Get(consts.GenUserCompanyKey(p.Account))
		}
		// 如果company不为空, 则写入上下文用于传递分区字段
		if p.Company != "" {
			ctx.CtxSet(core.CTX_SHARDING_KEY, p.Company)
		}

		user := UserService.GetByCode(p.Account)
		if user == nil {
			ctx.JSON(errno.ERRNO_AUTH_USER_OR_PWD)
			return
		}
		err := user.CheckPwd(p.Password)
		if err != nil {
			ctx.JSON(errno.ERRNO_AUTH_USER_OR_PWD)
			return
		}

		// session 设置
		sessionKey := session.GenSessionKey()
		ctx.Core().Session().SetWithCsrfToken(ctx.Ctx().Writer, sessionKey, user.Account)
		menu, _, permissions, _ := UserService.GetResourcePermission(ctx.TX(), user.Account)
		result := map[string]interface{}{
			"token":       sessionKey,
			"userInfo":    user.UserInfo(),
			"menu":        menu,
			"permissions": permissions,
		}
		if p.Company == "" && len(user.Company) > 0 {
			p.Company = user.Company[0]
		}
		user.Permissions[p.Company] = permissions
		result["company"] = p.Company
		// 设置公司code缓存
		ctx.Core().Cache().Set(consts.GenUserCompanyKey(user.Account), p.Company)
		ctx.JSON(errno.ERRNO_OK.WrapData(result).Notify())
	}
}

type userCompanyParams struct {
	Account string `form:"account"`
}

type setDefaultCompanyParams struct {
	Account string `json:"account" binding:"required"`
	Company string `json:"company" binding:"required"`
}

// 设置默认法人公司
func setDefaultCompany() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*setDefaultCompanyParams)
		ctx.Core().Cache().Set(consts.GenUserCompanyKey(param.Account), param.Company)
		ctx.CtxSet(core.CTX_SHARDING_KEY, param.Company)
		user := UserService.GetByCode(param.Account)
		if user == nil {
			ctx.JSON(errno.ERRNO_AUTH_NOT_LOGIN)
			return
		}
		menu, _, permissions, _ := UserService.GetResourcePermission(ctx.NewTX(), param.Account)
		result := map[string]interface{}{
			"menu":        menu,
			"permissions": permissions,
		}
		user.Permissions[param.Company] = permissions
		ctx.JSON(errno.ERRNO_OK.WrapData(result))
	}
}

type setDefaultPositionParams struct {
	Account  string `json:"account"`
	Position string `json:"position"`
}

// 设置默认岗位信息
func setDefaultPosition() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*setDefaultPositionParams)
		company_code := ctx.Core().Cache().Get(consts.GenUserCompanyKey(param.Account))

		ctx.Core().Cache().Set(consts.GenUserPositionKey(param.Account, company_code), param.Position)
		ctx.CtxSet(core.CTX_DATA_PERMISSION_KEY, param.Position)
		ctx.JSON(errno.ERRNO_OK)
	}
}

// 页面刷新后, 获取公司(分区)和岗位(数据权限)
//
// 不在登陆时返回前端并缓存以便及时获取信息,减少用户清缓存/重新登陆步骤
//
// 每次刷新整体页面或在服务端重启后时, 均可以获取最新的信息,而无需登陆
func getCompanyAndPosition() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userCompanyParams)
		res, err := UserService.getOrgPermission(ctx.TX(), param.Account)
		if err != nil {
			ctx.JSON(errno.ERRNO_CACHE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(res))
	}
}
