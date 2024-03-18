/*
 * @Author: reel
 * @Date: 2024-01-18 22:59:46
 * @LastEditors: reel
 * @LastEditTime: 2024-03-18 22:14:23
 * @Description: 用户处理逻辑
 */
package auth

import (
	"fmt"
	"portal/app/basis/org"
	"portal/pkg/consts"
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type userService struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[uint]*User
	codeMap map[string]*User
}

var UserService *userService

// 初始化
// toB 私有化部署设计, 默认加载所有用户
// TODO: 如果公网不知,需要考虑用户登录时再进行缓存设计
func UserServiceInit(c core.Core) {
	UserService = &userService{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[uint]*User, 100),
		codeMap: make(map[string]*User, 100),
	}
	var users = make([]*User, 0, 100)
	c.RDB().DB().Order("id").Find(&users)

	for _, user := range users {
		if user.Position1 == nil {
			user.Position1 = make(map[string]string, 100)
		}
		if user.Position2 == nil {
			user.Position2 = make(map[string][]string, 100)
		}
		UserService.setCache(user)
	}

	// 加载用户&公司关系
	var rucs = make([]*RlatUserCompany, 0, 100)
	c.RDB().DB().Where("1=1").Find(&rucs)
	for _, rlat := range rucs {
		user := UserService.codeMap[rlat.Account]
		if user == nil {
			continue
		}
		if user.Company == nil {
			user.Company = make([]string, 0, 10)
		}
		user.Company = append(user.Company, rlat.CompanyCode)
	}

	// 按分区加载角色和岗位
	for _, item := range org.CompanyService.GetAll() {
		// 加载用户position表
		var rups = make([]*RlatUserPosition, 0, 100)
		tx := c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, item.CompanyCode)
		tx.Order("id").Find(&rups)
		for _, item2 := range rups {
			user := UserService.codeMap[item2.Account]
			if user == nil {
				continue
			}

			if item2.IsPosition == 1 {
				user.Position1[item.CompanyCode] = item2.PositionCode
			} else {
				if user.Position2[item.CompanyCode] == nil {
					user.Position2[item.CompanyCode] = make([]string, 10)
				}
				user.Position2[item.CompanyCode] = append(user.Position2[item.CompanyCode], item2.PositionCode)
			}
		}

		// 加载用户角色表以及角色对应的权限
		var rurs = make([]*RlatUserRole, 0, 100)
		tx = c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, item.CompanyCode)
		tx.Order("id").Find(&rurs)
		for _, item2 := range rurs {
			user := UserService.codeMap[item2.Account]
			if user == nil {
				continue
			}
			if user.Roles == nil {
				user.Roles = make(map[string][]string, 100)
			}

			if user.Roles[item.CompanyCode] == nil {
				user.Roles[item.CompanyCode] = make([]string, 10)
			}
			user.Roles[item.CompanyCode] = append(user.Roles[item.CompanyCode], item2.RoleCode)

		}
	}

}

// 创建缓存
// TODO: 完善对公司, 角色的缓存
func (srv *userService) setCache(item *User) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	oItem := srv.codeMap[item.Account]
	// 插入或更新
	if oItem == nil {
		srv.codeMap[item.Account] = item
		srv.idMap[item.ID] = item
	} else {
		oItem.Password = item.Password
		oItem.NickName = item.NickName
		oItem.Email = item.Email
		oItem.Super = item.Super
		oItem.Status = item.Status
	}
	srv.codeMap[item.Account] = item
	srv.idMap[item.ID] = item
}

// 通过code批量获取公司列表
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *userService) GetByCodes(codes []string) (users []*User) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	users = make([]*User, 0, 10)
	if codes == nil {
		for _, user := range srv.codeMap {
			if user.Status == 1 {
				users = append(users, user)
			}
		}
		return
	}

	for _, code := range codes {
		user := srv.codeMap[code]
		if user != nil && user.Status == 1 {
			users = append(users, srv.codeMap[code])
		}
	}
	return
}

// 获取单个用户
func (srv *userService) GetByCode(code string) (user *User) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	if srv.codeMap[code].Status == 1 {
		return srv.codeMap[code]
	}
	return
}

// 通过ID获取单个用户
func (srv *userService) GetByID(id uint) (user *User) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	if srv.idMap[id].Status == 1 {
		return srv.idMap[id]
	}
	return
}

// 通过id批量获取公司列表
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *userService) GetByIDs(ids []uint) (users []*User) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	users = make([]*User, 0, 10)
	if ids == nil {
		for _, user := range srv.codeMap {
			if user.Status == 1 {
				users = append(users, user)
			}
		}
		return
	}

	for _, id := range ids {
		user := srv.idMap[id]
		if user != nil && user.Status == 1 {
			users = append(users, srv.idMap[id])
		}
	}
	return
}

// 按照id批量删除缓存
func (srv *userService) deleteByID(ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, id := range ids {
		user := srv.idMap[id]
		if user != nil {
			delete(srv.codeMap, user.Account)
			delete(srv.idMap, user.ID)
		}
	}
}

// 以下时对外服务的API

// 更改单用户昵称和邮箱, 其他信息不适合用户自己修改
func (srv *userService) UpdateByAccount(account string, param *userUpdateParams) (user *User, err error) {
	user = &User{
		NickName: param.NickName,
		Email:    param.Email,
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
		if RoleService.GetByCode(tx, companyCode) == nil {
			return errorx.Errorf("角色code:%s不存在或失效,请确认角色信息", companyCode)
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
	if err != nil {
		return
	}

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
		userData.Super = param.Super
		userData.NickName = param.NickName
		userData.Company = param.Company
		userData.Email = param.Email
		userData.Status = param.Status
		userData.Password = param.Pwd

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
		err = tx.Model(&User{}).Where("id in (?)", param.ID).Updates(userData).Find(&users).Error
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

	for _, item := range users {
		srv.setCache(item)
	}
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

		item.Position1 = srv.codeMap[item.Account].Position1[companyCode]
		item.Position2 = srv.codeMap[item.Account].Position2[companyCode]
		item.Departments = make([]string, 0, 10)
		positions = append(positions, org.PositionService.GetByCode(tx, item.Position1))
		if item.Position2 != nil {
			positions = append(positions, org.PositionService.GetByCodes(tx, item.Position2)...)
		}
		for _, position := range positions {
			if position == nil {
				continue
			}
			item.Departments = append(item.Departments, position.DepartmentCode)
		}
		item.Companies = srv.codeMap[item.Account].Company
		item.Roles = srv.codeMap[item.Account].Roles[companyCode]

	}
	data = map[string]interface{}{
		"page_num":  param.PageNum,
		"page_size": param.PageSize,
		"total":     count,
		"rows":      list,
	}
	return
}

// 设置用户&岗位表
func (srv *userService) SetRlatPosition(tx *gorm.DB, accounts []string, rups []*RlatUserPosition) (err error) {
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
	sk := rdb.GetShardingKey(tx)
	// 如果有更新, 则按用户全部更新
	if len(rups) > 0 {
		// 通过事物更新用户和岗位表
		err = ntx.Transaction(func(tx *gorm.DB) error {
			// 按照用户先删除所有的用户岗位关系数据
			err := ntx.Model(&RlatUserPosition{}).Unscoped().Where("account in (?) ", accounts).Delete(&RlatUserPosition{}).Error
			if err != nil {
				return errorx.Wrap(err, "更新岗位失败")
			}
			// 再更新用户用户/关系数据
			err = ntx.Model(&RlatUserPosition{}).Create(rups).Error
			if err != nil {
				return errorx.Wrap(err, "批量更新岗位失败")
			}
			for _, rup := range rups {
				user := srv.GetByCode(rup.Account)

				if rup.IsPosition == 1 {
					if user.Position1 == nil {
						user.Position1 = make(map[string]string, 10)
					}
					user.Position1[sk] = rup.PositionCode
				} else {
					if user.Position2 == nil {
						user.Position2 = make(map[string][]string)
					}
					if user.Position2[sk] == nil {
						user.Position2[sk] = make([]string, 0, 10)
					}
					user.Position2[sk] = append(user.Position2[sk], rup.PositionCode)
				}
			}
			return nil
		})
	}
	return
}

// 设置用户&角色表
func (srv *userService) SetRlatRole(tx *gorm.DB, accounts []string, rurs []*RlatUserRole) (err error) {
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
	sk := rdb.GetShardingKey(tx)
	// 如果有更新, 则按用户全部更新
	if len(rurs) > 0 {
		// 通过事物更新用户和角色表
		err = ntx.Transaction(func(tx *gorm.DB) error {
			// 按照用户先删除所有的用户岗位关系数据
			err := ntx.Unscoped().Model(&RlatUserRole{}).Unscoped().Where("account in (?) ", accounts).Delete(&RlatUserRole{}).Error
			if err != nil {
				return errorx.Wrap(err, "更新岗位失败")
			}
			// 再更新用户用户/关系数据
			err = ntx.Model(&RlatUserRole{}).Create(rurs).Error
			if err != nil {

				return errorx.Wrap(err, "批量更新岗位失败")
			}
			for _, rur := range rurs {
				user := srv.GetByCode(rur.Account)
				if user.Roles == nil {
					user.Roles = make(map[string][]string, 0)
				}
				if user.Roles[sk] == nil {
					user.Roles[sk] = make([]string, 0, 10)
				}
				user.Roles[sk] = append(user.Roles[sk], rur.RoleCode)
			}
			return nil
		})
	}
	return
}

// 设置用户&角色表
func (srv *userService) SetRlatCompany(tx *gorm.DB, accounts []string, rucs []*RlatUserCompany) (err error) {
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
	// 如果有更新, 则按用户全部更新
	if len(rucs) > 0 {
		// 通过事物更新用户和角色表
		err = ntx.Transaction(func(tx *gorm.DB) error {
			// 按照用户先删除所有的用户岗位关系数据
			err := ntx.Model(&RlatUserCompany{}).Unscoped().Where("account in (?) ", accounts).Delete(&RlatUserCompany{}).Error
			if err != nil {
				return errorx.Wrap(err, "批量更新公司失败")
			}
			// 再更新用户公司关系数据
			err = ntx.Model(&RlatUserCompany{}).Create(rucs).Error
			if err != nil {

				return errorx.Wrap(err, "批量更新公司失败")
			}
			for _, ruc := range rucs {
				user := srv.GetByCode(ruc.Account)
				if user.Company == nil {
					user.Company = make([]string, 0, 10)
				}

				user.Company = append(user.Company, ruc.CompanyCode)
			}
			return nil
		})
	}
	return
}

// 获取资源权限表
//
// 分为菜单管理权限, 接口权限, 设置菜单权限等
func (srv *userService) GetResourcePermission(tx *gorm.DB, account string) (menuList []*core.Sources, manageList []*core.Sources, permissions map[string]bool, err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
	sk := rdb.GetShardingKey(ntx)
	user := srv.codeMap[account]
	if user == nil || user.Status < 1 {
		return nil, nil, nil, errorx.Errorf("无有效用户:%s, 请检查用户是否存在或失效", account)
	}
	var permissionMap = make(map[string]bool, 100)
	// 从角色中获取用户的资源权限列表,如果角色列表为空, 则不过滤
	if user.Super == "Y" {
		for _, item := range ResourceService.GetAllList() {
			permissionMap[item.Code] = true
		}
	} else {
		for _, roleCode := range user.Roles[sk] {
			role := RoleService.GetByCode(ntx, roleCode)
			if role == nil || role.Status < 1 {
				continue
			}
			for _, codeI := range role.Sources {
				code, ok := codeI.(string)
				if ok {
					permissionMap[code] = true
				}
			}
		}
	}
	// 从权限列表中,获取用户对应的菜单权限,按钮(接口)权限, 菜单管理权限等
	return ResourceService.GetSource(permissionMap)
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
	if !isCheck && len(res) > 0 {
		company_code = res[0]["company_code"].(string)
	} else {
		company_code = ""
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
