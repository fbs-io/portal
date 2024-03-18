/*
 * @Author: reel
 * @Date: 2024-01-17 07:04:34
 * @LastEditors: reel
 * @LastEditTime: 2024-03-13 07:03:41
 * @Description: 角色表逻辑处理层
 */
package auth

import (
	"portal/app/basis/org"
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type roleService struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[string]map[uint]*Role
	codeMap map[string]map[string]*Role
	dimlist map[string][]map[string]interface{} // 维度化处理
}

var RoleService *roleService

// 初始化
//
// 根据不同分区缓存
func RoleServiceInit(c core.Core) {
	RoleService = &roleService{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[string]map[uint]*Role, 100),
		codeMap: make(map[string]map[string]*Role, 100),
		dimlist: make(map[string][]map[string]interface{}, 100),
	}
	// 通过CompanyCode作为分区并初始化不同分区的数据
	for _, item := range org.CompanyService.GetAll() {
		var list = make([]*Role, 0, 100)
		tx := c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, item.CompanyCode)

		tx.Order("id").Find(&list)

		for _, item := range list {
			RoleService.setCache(tx, item)
		}
		RoleService.genDimList(tx)
	}

}

// 创建缓存, 增加分区
func (srv *roleService) setCache(tx *gorm.DB, item *Role) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] == nil {
		srv.codeMap[sk] = make(map[string]*Role, 100)
	}
	if srv.idMap[sk] == nil {
		srv.idMap[sk] = make(map[uint]*Role, 100)
	}

	oItem := srv.codeMap[sk][item.Code]

	// 如果无法查询到元素 则插入,否则更新
	if oItem == nil {
		srv.codeMap[sk][item.Code] = item
		srv.idMap[sk][item.ID] = item
	} else {
		oItem.Label = item.Label
		oItem.Sort = item.Sort
		oItem.Description = item.Description
		oItem.Sources = item.Sources
		oItem.Status = item.Status
	}

}

// 通过gorm和code获取对应分区有效的角色
func (srv *roleService) GetByCodes(tx *gorm.DB, codes []string) (list []*Role) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	list = make([]*Role, 0, 10)
	for _, code := range codes {
		item := srv.codeMap[sk][code]
		if item != nil && item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 通过gorm上下文获取对应分区内所有有效的角色
func (srv *roleService) GetAll(tx *gorm.DB, codes []string) (list []*Role) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	list = make([]*Role, 0, 10)
	for _, item := range srv.codeMap[sk] {
		if item.Status == 1 {
			list = append(list, item)
		}
	}
	return

}

// 通过code和gorm上下文获取对应分区的有效的角色
func (srv *roleService) GetByCode(tx *gorm.DB, code string) (item *Role) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] != nil && srv.codeMap[sk][code].Status == 1 {
		return srv.codeMap[sk][code]
	}

	return
}

// 按照id批量删除缓存
func (srv *roleService) deleteByID(tx *gorm.DB, ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	for _, id := range ids {
		item := srv.idMap[sk][id]
		if item != nil {
			delete(srv.codeMap[sk], item.Code)
			delete(srv.idMap[sk], item.ID)
		}
	}
}

// 以下时对外服务的API

// 创建部门
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *roleService) Create(tx *gorm.DB, param *roleAddParams) (err error) {
	model := &Role{
		Code:        param.Code,
		Label:       param.Label,
		Sort:        param.Sort,
		Description: param.Description,
		Sources:     param.Sources,
	}

	err = tx.Create(model).Error
	if err != nil {
		return
	}
	srv.setCache(tx, model)
	srv.genDimList(tx)
	return
}

// 删除部门, 软删除
//
// 同时会删除部门缓存
func (srv *roleService) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	model := &Role{}

	err = tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
	if err != nil {
		return
	}

	srv.deleteByID(tx, param.ID)
	srv.genDimList(tx)
	return
}

// 更新, 通过id批量更新
//
// 同时会更新缓存
func (srv *roleService) UpdateByID(tx *gorm.DB, param *roleEditParams) (err error) {
	model := &Role{
		Label:       param.Label,
		Sort:        param.Sort,
		Description: param.Description,
		Sources:     param.Sources,
	}

	model.Status = param.Status

	list := make([]*Role, 0, 100)
	err = tx.Where("id in (?) ", param.ID).Updates(model).Find(&list).Error
	if err != nil {
		return err
	}
	for _, item := range list {
		srv.setCache(tx, item)
	}
	srv.genDimList(tx)
	return
}

// 查询
//
// TODO: 缓存结果
func (srv *roleService) Query(tx *gorm.DB, param *rolesQueryParams) (data interface{}, err error) {
	role := &Role{}
	roles := make([]*Role, 0, 100)

	var count int64
	err = tx.Model(role).Find(&roles).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return
	}

	data = map[string]interface{}{
		"page_num":  param.PageNum,
		"page_size": param.PageSize,
		"total":     count,
		"rows":      roles,
	}
	return
}

// 查询维度
func (srv *roleService) DimList(tx *gorm.DB) (result []map[string]interface{}) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)

	return srv.dimlist[sk]
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *roleService) genDimList(tx *gorm.DB) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	result := make([]map[string]interface{}, 0, 100)
	for _, item := range srv.idMap[sk] {
		result = append(result, map[string]interface{}{
			"code":   item.Code,
			"name":   item.Label,
			"status": item.Status,
		})
	}
	srv.dimlist[sk] = result
}
