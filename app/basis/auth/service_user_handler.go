/*
 * @Author: reel
 * @Date: 2024-03-21 05:43:50
 * @LastEditors: reel
 * @LastEditTime: 2024-03-23 05:46:53
 * @Description: 请填写简介
 */
/*
 * @Author: reel
 * @Date: 2024-01-18 22:59:46
 * @LastEditors: reel
 * @LastEditTime: 2024-03-21 05:40:54
 * @Description: 用户处理逻辑
 */
package auth

import (
	"fmt"
	"portal/app/basis/org"
	"portal/pkg/consts"

	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

// 以下时对外服务的API

// 更改单用户昵称和邮箱, 其他信息不适合用户自己修改
func (srv *userService) UpdateByAccount(account string, param *userUpdateParams) (user *User, err error) {
	user = &User{
		NickName: param.NickName,
		Email:    param.Email,
		Account:  account,
	}
	err = srv.core.RDB().DB().Where("account = (?)", account).Updates(user).Error
	if err != nil {
		return
	}
	srv.setCache(user)
	return
}

// 用户更改密码
func (srv *userService) ChangePwd(tx *gorm.DB, account string, param *userChPwdParams) (err error) {

	user := srv.codeMap[account]
	if user == nil {
		return errorx.New("无有效用户,请确认登陆用户是否正确")
	}
	err = user.chpwd(param)
	if err != nil {
		return errorx.Wrap(err, "更新密码发生失败")
	}
	err = tx.Save(user).Error
	if err != nil {
		return errorx.Wrap(err, "更新密码失败")
	}
	srv.setCache(user)
	return
}

// 创建用户
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *userService) Create(tx *gorm.DB, param *userAddParams) (err error) {
	companyCode := rdb.GetShardingKey(tx)
	rups := make([]*RlatUserPosition, 0, 10)
	user := &User{
		Account:   param.Account,
		Password:  param.Password,
		NickName:  param.NickName,
		Email:     param.Email,
		Super:     param.Super,
		Position1: make(map[string]string),
		Position2: make(map[string][]string),
		Roles:     make(map[string][]string),
	}

	// 岗位设置, 主岗和兼岗默认都会设置, 前端需要把参数都传过来
	// 主岗
	if param.Position1 != "" {
		if org.PositionService.GetByCode(tx, param.Position1) == nil {
			return errorx.Errorf("岗位code:%s不存在或失效,请确认岗位信息", param.Position1)
		}
		rups = append(rups, &RlatUserPosition{
			Account:      user.Account,
			PositionCode: param.Position1,
			IsPosition:   1,
		})
	}
	user.Position1[companyCode] = param.Position1

	//兼岗
	for _, positCode := range param.Position2 {
		if org.PositionService.GetByCode(tx, positCode) == nil {
			return errorx.Errorf("岗位code:%s不存在或失效,请确认岗位信息", positCode)
		}
		rup := &RlatUserPosition{
			Account:      user.Account,
			PositionCode: positCode,
			IsPosition:   -1,
		}
		rups = append(rups, rup)
	}
	user.Position2[companyCode] = param.Position2

	// 用户角色设置
	var rurs = make([]*RlatUserRole, 0, 100)
	for _, roleCode := range param.Role {
		if RoleService.GetByCode(tx, roleCode) == nil {
			return errorx.Errorf("角色code:%s不存在或失效,请确认角色信息", roleCode)
		}
		rur := &RlatUserRole{
			Account:  user.Account,
			RoleCode: roleCode,
		}
		rurs = append(rurs, rur)
	}
	user.Roles[companyCode] = param.Role

	// 用户公司设置
	var rucs = make([]*RlatUserCompany, 0, 100)
	for _, companyCode := range param.Company {
		if org.CompanyService.GetByCode(companyCode) == nil {
			return errorx.Errorf("公司code:%s不存在或失效,请确认公司信息", companyCode)
		}
		ruc := &RlatUserCompany{
			Account:     user.Account,
			CompanyCode: companyCode,
		}
		rucs = append(rucs, ruc)
	}
	user.Company = param.Company

	err = tx.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Create(user).Error
		if err != nil {
			return
		}
		err = srv.SetRlatPosition(tx, []string{user.Account}, rups)
		if err != nil {
			return
		}
		err = srv.SetRlatRole(tx, []string{user.Account}, rurs)
		if err != nil {
			return
		}
		err = srv.SetRlatCompany(tx, []string{user.Account}, rucs)
		if err != nil {
			return
		}
		return

	})
	srv.UpdateUsers(user.Account)
	return
}

// 删除用户, 软删除
//
// 同时会删除公司缓存
func (srv *userService) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {

	user := &User{}

	err = tx.Model(user).Where("id in (?)", param.ID).Delete(user).Error
	if err != nil {
		return
	}
	srv.deleteByID(param.ID)
	return
}

// 更新用户, 通过id批量更新
//
// 同时会更新缓存
func (srv *userService) UpdateByID(tx *gorm.DB, param *usersEditParams) (err error) {
	users := srv.GetByIDs(param.ID)
	rups := make([]*RlatUserPosition, 0, 10)
	rurs := make([]*RlatUserRole, 0, 100)
	rucs := make([]*RlatUserCompany, 0, 100)
	// 删除已存在用户的岗位用
	dus := make([]string, 0, 10)
	userData := &User{
		Super:    param.Super,
		NickName: param.NickName,
		Company:  param.Company,
		Email:    param.Email,
		Password: param.Pwd,
	}
	userData.Status = param.Status
	// 构建
	for _, user := range users {
		dus = append(dus, user.Account)
		// 主岗
		if param.Position1 != "" {
			if org.PositionService.GetByCode(tx, param.Position1) == nil {
				return errorx.Errorf("岗位code:%s不存在", param.Position1)
			}
			rups = append(rups, &RlatUserPosition{
				Account:      user.Account,
				PositionCode: param.Position1,
				IsPosition:   1,
			})
		}

		//兼岗
		for _, posit := range param.Position2 {
			if org.PositionService.GetByCode(tx, posit) == nil {
				return errorx.Errorf("岗位code:%s不存在", posit)
			}
			rup := &RlatUserPosition{
				Account:      user.Account,
				PositionCode: posit,
				IsPosition:   -1,
			}
			rups = append(rups, rup)
		}
		// 用户角色设置

		for _, roleCode := range param.Role {
			if RoleService.GetByCode(tx, roleCode) == nil {
				return errorx.Errorf("角色code:%s不存在或失效,请确认角色信息", roleCode)
			}
			rur := &RlatUserRole{
				Account:  user.Account,
				RoleCode: roleCode,
			}
			rurs = append(rurs, rur)
		}

		// 用户公司设置
		for _, companyCode := range param.Company {
			if org.CompanyService.GetByCode(companyCode) == nil {
				return errorx.Errorf("公司code:%s不存在或失效,请确认公司信息", companyCode)
			}
			ruc := &RlatUserCompany{
				Account:     user.Account,
				CompanyCode: companyCode,
			}
			rucs = append(rucs, ruc)
		}

	}

	err = tx.Transaction(func(tx *gorm.DB) (err error) {
		err = tx.Model(userData).Where("id in (?)", param.ID).Limit(-1).Offset(-1).Updates(userData).Find(&users).Error
		if err != nil {
			return
		}

		err = srv.SetRlatPosition(tx, dus, rups)
		if err != nil {
			return
		}

		err = srv.SetRlatRole(tx, dus, rurs)
		if err != nil {
			return
		}

		err = srv.SetRlatCompany(tx, dus, rucs)
		if err != nil {
			return
		}
		return
	})

	// 如果有更新, 则按用户全部更新
	srv.UpdateUsers(dus...)
	return
}

// 查询
//
// TODO: 缓存结果
func (srv *userService) Query(tx *gorm.DB, param *usersQueryParams) (data interface{}, err error) {
	model := &User{}
	list := make([]*UserList, 0, 100)
	companyCode := rdb.GetShardingKey(tx)

	err = tx.Model(model).Order("id").Find(&list).Error
	if err != nil {
		return
	}
	var count int64
	err = tx.Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return
	}
	// 处理员工岗位,部门信息
	for _, item := range list {
		// 获取岗位表
		positions := make([]*org.Position, 0, 100)
		item.Positions = make([]string, 0, 10)
		user := srv.GetByCode(item.Account)
		if user == nil {
			continue
		}
		item.Position1 = user.Position1[companyCode]
		item.Position2 = user.Position2[companyCode]
		item.Department = make([]string, 0, 10)
		positions = append(positions, org.PositionService.GetByCode(tx, item.Position1))
		if item.Position2 != nil {
			positions = append(positions, org.PositionService.GetByCodes(tx, item.Position2)...)
		}
		for _, position := range positions {
			if position == nil {
				continue
			}
			item.Department = append(item.Department, position.DepartmentCode)
			item.Positions = append(item.Positions, position.PositionCode)
		}

		item.Company = srv.codeMap[item.Account].Company
		item.Role = srv.codeMap[item.Account].Roles[companyCode]
	}
	data = map[string]interface{}{
		"page_num":  param.PageNum,
		"page_size": param.PageSize,
		"total":     count,
		"rows":      list,
	}
	return
}

// 获取组织权限(数据权限相关)
//
// 登陆后, 获取用户的公司(租户/分区)和岗位信息(数据权限)
//
// 主要用于前端顶部公司切换和岗位切换使用
func (srv *userService) getOrgPermission(tx *gorm.DB, account string) (result map[string]interface{}, err error) {

	user := srv.codeMap[account]
	if user == nil {
		return
	}
	res := make([]map[string]interface{}, 0, 10)
	companyList := org.CompanyService.GetByCodes(user.Company)
	if user.Super == "Y" {
		companyList = org.CompanyService.GetAll()

	}
	// 设置/更新默认公司
	var isCheck = false
	company_code := srv.core.Cache().Get(consts.GenUserCompanyKey(user.Account))
	for _, company := range companyList {
		res = append(res, map[string]interface{}{
			"company_code":       company.CompanyCode,
			"company_name":       company.CompanyName,
			"company_short_name": company.CompanyShortName,
		})

		if company_code == company.CompanyCode {
			isCheck = true

		}
	}
	if !isCheck {
		if len(res) > 0 {
			company_code = res[0]["company_code"].(string)
		} else {
			company_code = ""
		}
	}
	// 当前设置的公司(分区/租户)以及有权限的其他公司(分区/租户)
	result = map[string]interface{}{
		"companies": res,
		"company":   company_code,
	}
	err = srv.core.Cache().Set(consts.GenUserCompanyKey(user.Account), company_code)
	if err != nil {
		return

	}

	positions := make([]*org.Position, 0, 100)

	// 先从缓存获取默认岗位信息
	position_code := srv.core.Cache().Get(consts.GenUserPositionKey(user.Account, company_code))
	if position_code == "" {
		// 如果缓存获取不到,则获取默认主岗信息
		// 主岗获取不到,则从兼岗获取第一个code作为缓存
		position := org.PositionService.GetByCode(tx, user.Position1[company_code])
		position2 := org.PositionService.GetByCodes(tx, user.Position2[company_code])
		if position == nil && len(position2) > 0 {
			position = position2[0]
		}
		if position != nil {
			position_code = position.PositionCode

		}
	}
	err = srv.core.Cache().Set(consts.GenUserPositionKey(user.Account, company_code), position_code)
	if err != nil {
		return
	}

	// 岗位code不为空, 但无法获取有效的岗位, 则删除code
	if position_code != "" && org.PositionService.GetByCode(tx, position_code) == nil {
		err = srv.core.Cache().Del(consts.GenUserPositionKey(user.Account, company_code))
		if err != nil {
			return
		}
		position_code = ""
	}

	result["position"] = position_code

	positions = append(positions, org.PositionService.GetByCode(tx, user.Position1[company_code]))
	positions = append(positions, org.PositionService.GetByCodes(tx, user.Position2[company_code])...)

	positionList := make([]map[string]string, 0, 10)
	for _, position := range positions {
		if position == nil {
			continue
		}
		// 如果岗位上的部门不存在, 则不合法, 岗位也无需显示
		if org.DepartmentService.GetByCode(tx, position.DepartmentCode) != nil {
			positionList = append(positionList, map[string]string{
				"position_code": position.PositionCode,
				"position_name": fmt.Sprintf("%s - %s", org.DepartmentService.GetByCode(tx, position.DepartmentCode).DepartmentName, position.PositionName),
			})
		}
	}
	result["positions"] = positionList
	return
}
