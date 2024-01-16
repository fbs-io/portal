/*
 * @Author: reel
 * @Date: 2023-07-18 21:46:02
 * @LastEditors: reel
 * @LastEditTime: 2024-01-15 00:46:08
 * @Description: 请填写简介
 */
package auth

import (
	"fmt"
	"portal/app/basis/org"
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

		user := GetUser(p.Account, ctx, REFRESH_ALL)

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

		result := map[string]interface{}{
			"token":       sessionKey,
			"userInfo":    user.UserInfo(),
			"menu":        user.Menu,
			"permissions": user.Permissions,
		}
		SetUser(user.Account, user)
		if p.Company == "" && len(user.Company) > 0 {
			p.Company = user.Company[0].(string)
		}
		result["company"] = p.Company
		// 设置公司code缓存
		ctx.Core().Cache().Set(consts.GenUserCompanyKey(user.Account), p.Company)
		ctx.JSON(errno.ERRNO_OK.WrapData(result).Notify())
	}
}

type userCompanyParams struct {
	Account string `form:"account"`
}

func getCompany() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*userCompanyParams)
		user := GetUser(param.Account, ctx, REFRESH_USRE)
		if user == nil {
			ctx.JSON(errno.ERRNO_OK)
			return
		}

		var result = make([]map[string]interface{}, 0, 10)
		// tx := ctx.NewTX()
		companies := make([]string, 0, 10)

		for _, c := range user.Company {
			if c != nil {
				companies = append(companies, c.(string))
			}
		}
		if user.Super == "Y" {
			companies = nil

		}
		companyList := org.CompanySrvice.GetByCode(companies)
		for _, company := range companyList {
			result = append(result, map[string]interface{}{
				"company_code":       company.CompanyCode,
				"company_name":       company.CompanyName,
				"company_short_name": company.CompanyShortName,
			})
		}
		// 设置/更新默认公司
		company_code := ctx.Core().Cache().Get(consts.GenUserCompanyKey(user.Account))
		var isCheck = false
		if len(result) > 0 {
			for _, item := range result {
				if item["company_code"] == company_code {
					isCheck = true
				}
			}
			if !isCheck {
				company_code = result[0]["company_code"].(string)
			}
		}

		res := map[string]interface{}{
			"companies": result,
			"company":   company_code,
		}
		ctx.Core().Cache().Set(consts.GenUserCompanyKey(user.Account), company_code)
		SetUser(user.Account, user)

		// 岗位读取设置
		rlat := make([]*RlatUserPosition, 0, 10)
		ctx.NewTX().Where("account = ? ", user.Account).Find(&rlat)
		position_codes := make([]string, 0, 100)
		position_map := make(map[string]*RlatUserPosition, 10)
		for _, r := range rlat {
			position_codes = append(position_codes, r.PositionCode)
			position_map[r.PositionCode] = r
			if r.IsPosition == 1 {
				res["position"] = r.PositionCode
			}
		}
		position_code := ctx.Core().Cache().Get(consts.GenUserPositionKey(user.Account, company_code))
		if position_code != "" {
			res["position"] = position_code
		}

		positions := make([]*org.Position, 0, 100)

		ctx.NewTX().Where("position_code in ? ", position_codes).Find(&positions)

		department_codes := make([]string, 0, 10)
		deparment_map := make(map[string]*org.Position, 10)
		for _, position := range positions {
			department_codes = append(department_codes, position.DepartmentCode)
			deparment_map[position.DepartmentCode] = position
		}

		departments := make([]*org.Department, 0, 10)
		positionList := make([]map[string]string, 0, 10)
		ctx.NewTX().Where("department_code in ? ", department_codes).Find(&departments)
		for _, department := range departments {
			positionList = append(positionList, map[string]string{
				"position_code": deparment_map[department.DepartmentCode].PositionCode,
				"position_name": fmt.Sprintf("%s - %s", department.DepartmentName, deparment_map[department.DepartmentCode].PositionName),
			})
		}
		res["positions"] = positionList
		if res["position"] != nil {
			ctx.Core().Cache().Set(consts.GenUserPositionKey(user.Account, company_code), res["position"].(string))
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(res))
	}
}

type setDefaultCompanyParams struct {
	Account string `json:"account"`
	Company string `json:"company"`
}

func setDefaultCompany() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*setDefaultCompanyParams)
		ctx.Core().Cache().Set(consts.GenUserCompanyKey(param.Account), param.Company)
		ctx.CtxSet(core.CTX_SHARDING_KEY, param.Company)
		user := GetUser(param.Account, ctx, REFRESH_USRE)
		result := map[string]interface{}{
			"menu":        user.Menu,
			"permissions": user.Permissions,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(result))
	}
}

type setDefaultPositionParams struct {
	Account  string `json:"account"`
	Position string `json:"position"`
}

func setDefaultPosition() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*setDefaultPositionParams)
		company_code := ctx.Core().Cache().Get(consts.GenUserCompanyKey(param.Account))

		ctx.Core().Cache().Set(consts.GenUserPositionKey(param.Account, company_code), param.Position)
		ctx.CtxSet(core.CTX_DATA_PERMISSION_KEY, param.Position)
		ctx.JSON(errno.ERRNO_OK)
	}
}
