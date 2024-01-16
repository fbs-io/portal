/*
 * @Author: reel
 * @Date: 2024-01-16 22:33:20
 * @LastEditors: reel
 * @LastEditTime: 2024-01-16 23:45:31
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

type positionSrvice struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[string]map[uint]*Position
	codeMap map[string]map[string]*Position
	dimlist map[string][]map[string]interface{} // 维度化处理
}

var PositionSrvice *positionSrvice

// 初始化
//
// 根据不同分区缓存
func PositionSrviceInit(c core.Core) {
	PositionSrvice = &positionSrvice{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[string]map[uint]*Position, 100),
		codeMap: make(map[string]map[string]*Position, 100),
		dimlist: make(map[string][]map[string]interface{}, 100),
	}
	for sk := range CompanySrvice.codeMap {
		var list = make([]*Position, 0, 100)
		tx := c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, sk)
		tx.Where("deleted_at = 0").Order("id").Find(&list)

		for _, item := range list {
			PositionSrvice.setCache(tx, item)
		}

	}

}

// 创建缓存, 增加分区
func (srv *positionSrvice) setCache(tx *gorm.DB, item *Position) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] == nil {
		srv.codeMap[sk] = make(map[string]*Position, 100)
	}
	if srv.idMap[sk] == nil {
		srv.idMap[sk] = make(map[uint]*Position, 100)
	}

	oItem := srv.codeMap[sk][item.DepartmentCode]

	// 如果无法查询到元素 则插入,否则更新
	if oItem == nil {
		srv.codeMap[sk][item.DepartmentCode] = item
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
	}

}

// 通过id批量获取公司列表
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *positionSrvice) GetByCode(tx *gorm.DB, codes []string) (list []*Position) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	list = make([]*Position, 0, 10)
	if codes == nil {
		for _, item := range srv.codeMap[sk] {
			if item.Status == 1 {
				list = append(list, item)
			}
		}
		return
	}

	for _, code := range codes {
		item := srv.codeMap[sk][code]
		if item != nil && item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 按照id批量删除缓存
func (srv *positionSrvice) deleteByID(tx *gorm.DB, ids []uint) {
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

// 创建部门
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *positionSrvice) Create(tx *gorm.DB, param *positionAddParams) (err error) {
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
		pmodel := &Position{}
		err := tx.Table(model.TableName()).Where("position_code = ? and status = 1", model.PositionParentCode).First(pmodel).Error
		if err != nil || pmodel.PositionCode == "" {
			return errorx.New("无有效的上级code")
		}
	}

	// 校验部门code
	dmodel := DepartmentSrvice.GetByCode(tx, []string{model.DepartmentCode})
	if dmodel == nil || len(dmodel) == 0 {
		return errorx.New("无有效的部门code")
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
func (srv *positionSrvice) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	company := &Company{}

	err = tx.Model(company).Where("id in (?)", param.ID).Delete(company).Error
	if err != nil {
		return
	}
	// 创建公司相关的库或者表
	err = srv.core.RDB().AddShardingSuffixs(company.CompanyCode)
	if err != nil {
		return errorx.Wrap(err, "迁移数据库/表发生错误")
	}
	srv.deleteByID(tx, param.ID)
	srv.genDimList(tx)
	return
}

// 更新部门, 通过id批量更新
//
// 同时会更新缓存
func (srv *positionSrvice) UpdateByID(tx *gorm.DB, param *positionEditParams) (err error) {
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
		pmodel := &Position{}
		nTx := srv.core.RDB().DB().Where("1=1")
		rdb.CopyCtx(tx, nTx)

		err := nTx.Table(model.TableName()).Where("position_code = ? and status = 1", model.PositionParentCode).First(pmodel).Error
		if err != nil || pmodel.DepartmentCode == "" {

			return errorx.New("无有效的上级岗位code")
		}
	}

	// 如果参数中的部门code不为空, 则校验部门code
	if model.DepartmentCode != "" {
		dmodel := DepartmentSrvice.GetByCode(tx, []string{model.DepartmentCode})
		if dmodel == nil || len(dmodel) == 0 {
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

// 查询维度
func (srv *positionSrvice) DimList(tx *gorm.DB) (result []map[string]interface{}) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)

	return srv.dimlist[sk]
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *positionSrvice) genDimList(tx *gorm.DB) {
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
