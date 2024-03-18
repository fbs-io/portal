/*
 * @Author: reel
 * @Date: 2024-01-16 22:33:20
 * @LastEditors: reel
 * @LastEditTime: 2024-03-15 21:34:24
 * @Description: 岗位表相关处理逻辑
 */

package org

import (
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type positionService struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[string]map[uint]*Position
	codeMap map[string]map[string]*Position
	dimlist map[string][]map[string]interface{} // 维度化处理
}

var PositionService *positionService

// 初始化
//
// 根据不同分区缓存
func PositionServiceInit(c core.Core) {
	PositionService = &positionService{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[string]map[uint]*Position, 100),
		codeMap: make(map[string]map[string]*Position, 100),
		dimlist: make(map[string][]map[string]interface{}, 100),
	}
	for sk := range CompanyService.codeMap {
		var list = make([]*Position, 0, 100)
		tx := c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, sk)
		tx.Order("id").Find(&list)

		for _, item := range list {
			PositionService.setCache(tx, item)
		}
		PositionService.genDimList(tx)

	}

}

// 创建缓存, 增加分区
func (srv *positionService) setCache(tx *gorm.DB, item *Position) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] == nil {
		srv.codeMap[sk] = make(map[string]*Position, 100)
	}
	if srv.idMap[sk] == nil {
		srv.idMap[sk] = make(map[uint]*Position, 100)
	}

	oItem := srv.codeMap[sk][item.PositionCode]

	// 如果无法查询到元素 则插入,否则更新
	if oItem == nil {
		srv.codeMap[sk][item.PositionCode] = item
		srv.idMap[sk][item.ID] = item
	} else {
		oItem.PositionName = item.PositionName
		oItem.PositionComment = item.PositionComment
		oItem.PositionParentCode = item.PositionParentCode
		oItem.DepartmentCode = item.DepartmentCode
		oItem.JobCode = item.JobCode
		oItem.IsHead = item.IsHead
		oItem.IsApprove = item.IsApprove
		oItem.IsVritual = item.IsVritual
		oItem.DataPermissionType = item.DataPermissionType
		oItem.DataPermissionCustom = item.DataPermissionCustom
		oItem.Status = item.Status
	}

}

// 通过code批量获取有效岗位列表
func (srv *positionService) GetByCodes(tx *gorm.DB, codes []string) (list []*Position) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	list = make([]*Position, 0, 10)
	for _, code := range codes {
		if srv.codeMap[sk] != nil {
			item := srv.codeMap[sk][code]
			if item != nil && item.Status == 1 {
				list = append(list, item)
			}
		}
	}
	return
}

// 查询所有有效的岗位
func (srv *positionService) GetAll(tx *gorm.DB) (list []*Position) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	list = make([]*Position, 0, 10)
	for _, item := range srv.codeMap[sk] {
		if item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 通过code获取对应分区的单个且有效的值
func (srv *positionService) GetByCode(tx *gorm.DB, code string) (item *Position) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] != nil && srv.codeMap[sk][code] != nil && srv.codeMap[sk][code].Status == 1 {
		return srv.codeMap[sk][code]
	}
	return
}

// 按照id批量删除缓存
func (srv *positionService) deleteByID(tx *gorm.DB, ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	for _, id := range ids {
		item := srv.idMap[sk][id]
		if item != nil {
			delete(srv.codeMap[sk], item.DepartmentCode)
			delete(srv.idMap[sk], item.ID)
		}
	}
}

// 以下时对外服务的API

// 创建岗位
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *positionService) Create(tx *gorm.DB, param *positionAddParams) (err error) {
	model := &Position{
		PositionCode:         param.PositionCode,
		PositionName:         param.PositionName,
		PositionComment:      param.PositionComment,
		PositionParentCode:   param.PositionParentCode,
		DepartmentCode:       param.DepartmentCode,
		JobCode:              param.JobCode,
		IsHead:               param.IsHead,
		IsApprove:            param.IsApprove,
		IsVritual:            param.IsVritual,
		DataPermissionType:   param.DataPermissionType,
		DataPermissionCustom: param.DataPermissionCustom,
	}
	// 校验上级code
	if model.PositionParentCode != "" {
		if srv.GetByCode(tx, model.PositionParentCode) == nil {
			return errorx.New("无有效的上级code")
		}
	}

	// 如果参数中的部门code不为空, 则校验部门code
	if model.DepartmentCode != "" {
		dmodel := DepartmentService.GetByCode(tx, model.DepartmentCode)
		if dmodel == nil {
			return errorx.New("无有效的部门code")
		}
	}

	err = tx.Create(model).Error
	if err != nil {
		return
	}
	srv.setCache(tx, model)
	srv.genDimList(tx)
	return
}

// 删除岗位, 软删除
//
// 同时会删除岗位缓存
func (srv *positionService) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	model := &Position{}

	err = tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
	if err != nil {
		return
	}

	srv.deleteByID(tx, param.ID)
	srv.genDimList(tx)
	return
}

// 更新岗位, 通过id批量更新
//
// 同时会更新缓存
func (srv *positionService) UpdateByID(tx *gorm.DB, param *positionEditParams) (err error) {
	model := &Position{
		PositionName:         param.PositionName,
		PositionComment:      param.PositionComment,
		PositionParentCode:   param.PositionParentCode,
		DepartmentCode:       param.DepartmentCode,
		JobCode:              param.JobCode,
		IsHead:               param.IsHead,
		IsApprove:            param.IsApprove,
		IsVritual:            param.IsVritual,
		DataPermissionType:   param.DataPermissionType,
		DataPermissionCustom: param.DataPermissionCustom,
	}

	model.Status = param.Status

	// 校验上级code
	if model.PositionParentCode != "" {
		if srv.GetByCode(tx, model.PositionParentCode) == nil {

			return errorx.New("无有效的上级岗位code")
		}
	}

	// 如果参数中的部门code不为空, 则校验部门code
	if model.DepartmentCode != "" {
		dmodel := DepartmentService.GetByCode(tx, model.DepartmentCode)
		if dmodel == nil {
			return errorx.New("无有效的部门code")
		}
	}

	list := make([]*Position, 0, 100)
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
func (srv *positionService) Query(tx *gorm.DB, param *positionQueryParams) (data interface{}, err error) {
	model := &Position{}
	list := make([]*Position, 0, 100)
	var count int64

	err = tx.Model(model).Find(&list).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		return
	}

	data = map[string]interface{}{
		"page_num":  param.PageNum,
		"page_size": param.PageSize,
		"total":     count,
		"rows":      list,
	}
	return
}

// 查询维度
func (srv *positionService) DimList(tx *gorm.DB) (result []map[string]interface{}) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)

	return srv.dimlist[sk]
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *positionService) genDimList(tx *gorm.DB) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	result := make([]map[string]interface{}, 0, 100)
	for _, item := range srv.idMap[sk] {
		result = append(result, map[string]interface{}{
			"code":   item.PositionCode,
			"name":   item.PositionName,
			"status": item.Status,
		})
	}
	srv.dimlist[sk] = result
}
