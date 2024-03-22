/*
 * @Author: reel
 * @Date: 2024-01-18 22:59:46
 * @LastEditors: reel
 * @LastEditTime: 2024-03-22 07:43:37
 * @Description: 用户处理逻辑
 */
package auth

import (
	"portal/app/basis/org"
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
	UserService.UpdateUsers()
}

func (srv *userService) UpdateUsers(accounts ...string) {
	var users = make([]*User, 0, 100)
	tx := srv.core.RDB().DB()
	if accounts != nil {
		tx = tx.Where("account in ?", accounts)
	}
	tx.Order("id").Find(&users)

	srv.setCache(users...)

	// 加载用户&公司关系
	var rucs = make([]*RlatUserCompany, 0, 100)
	tx = srv.core.RDB().DB().Where("1=1")
	if accounts != nil {
		tx = tx.Where("account in ?", accounts)
	}
	tx.Order("id").Find(&rucs)
	for _, rlat := range rucs {
		user := srv.GetByCode(rlat.Account)
		if user == nil {
			continue
		}
		user.UpdateCompany(rlat.CompanyCode)
	}

	// 按分区加载角色和岗位
	for _, item := range org.CompanyService.GetAll() {
		// 加载用户position表
		var rups = make([]*RlatUserPosition, 0, 100)
		tx := srv.core.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, item.CompanyCode)
		if accounts != nil {
			tx = tx.Where("account in ?", accounts)
		}
		tx.Order("id").Find(&rups)
		for _, item2 := range rups {
			user := srv.GetByCode(item2.Account)
			if user == nil {
				continue
			}
			user.UpdatePosition(item.CompanyCode, item2)
		}

		// 加载用户角色表以及角色对应的权限
		var rurs = make([]*RlatUserRole, 0, 100)
		tx = srv.core.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, item.CompanyCode)
		if accounts != nil {
			tx = tx.Where("account in ?", accounts)
		}
		tx.Order("id").Find(&rurs)
		for _, item2 := range rurs {
			user := srv.GetByCode(item2.Account)
			if user == nil {
				continue
			}
			user.UpdateRole(item.CompanyCode, item2.RoleCode)
		}
	}
}

// 创建缓存
// TODO: 完善对公司, 角色的缓存
func (srv *userService) setCache(items ...*User) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, item := range items {
		if item == nil {
			continue
		}
		srv.codeMap[item.Account] = item
		srv.idMap[item.ID] = item
	}

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

	if srv.codeMap[code] != nil && srv.codeMap[code].Status == 1 {
		return srv.codeMap[code]
	}
	return
}

// 通过ID获取单个用户
func (srv *userService) GetByID(id uint) (user *User) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	if srv.idMap[id] != nil && srv.idMap[id].Status == 1 {
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

// 设置用户&岗位表
func (srv *userService) SetRlatPosition(tx *gorm.DB, accounts []string, rups []*RlatUserPosition) (err error) {
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
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
			return nil
		})
	}
	return
}

// 设置用户&角色表
func (srv *userService) SetRlatRole(tx *gorm.DB, accounts []string, rurs []*RlatUserRole) (err error) {
	ntx := srv.core.RDB().DB().Where("1=1")
	rdb.CopyCtx(tx, ntx)
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

			return nil
		})
	}
	return
}

// 设置用户&公司表
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
			return nil
		})
	}
	return
}
