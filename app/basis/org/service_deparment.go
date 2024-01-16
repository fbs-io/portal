/*
 * @Author: reel
 * @Date: 2024-01-15 20:03:41
 * @LastEditors: reel
 * @LastEditTime: 2024-01-17 00:10:58
 * @Description: 部门表相关逻辑处理
 */
package org

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type departmentSrvice struct {
	lock        *sync.RWMutex
	core        core.Core
	list        map[string][]*Department
	idMap       map[string]map[uint]*Department
	codeMap     map[string]map[string]*Department
	treeList    map[string][]*departmentTree                     // 组织树list, 前端展示
	allChildren map[string]map[string]map[string]*departmentTree // map下包含所有下级
	dimList     map[string][]map[string]interface{}              // 维度表
}

var DepartmentSrvice *departmentSrvice

// 初始化
//
// 根据不同分区缓存
func DepartmentSrviceInit(c core.Core) {
	DepartmentSrvice = &departmentSrvice{
		core:        c,
		lock:        &sync.RWMutex{},
		list:        make(map[string][]*Department, 100),
		idMap:       make(map[string]map[uint]*Department, 100),
		codeMap:     make(map[string]map[string]*Department, 100),
		treeList:    make(map[string][]*departmentTree, 100),
		allChildren: make(map[string]map[string]map[string]*departmentTree, 100),
		dimList:     make(map[string][]map[string]interface{}, 100),
	}
	for sk := range CompanySrvice.codeMap {
		var deparments = make([]*Department, 0, 100)
		tx := c.RDB().DB().Where("1=1")
		tx.Set(core.CTX_SHARDING_KEY, sk)
		tx.Where("deleted_at = 0").Order("id").Find(&deparments)

		for _, item := range deparments {
			DepartmentSrvice.setCache(tx, item)
		}
		DepartmentSrvice.list[sk] = deparments
		DepartmentSrvice.GenDepartmentTree(tx)

	}

}

// 创建缓存, 增加分区
func (srv *departmentSrvice) setCache(tx *gorm.DB, item *Department) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	if srv.codeMap[sk] == nil {
		srv.codeMap[sk] = make(map[string]*Department, 100)
	}
	if srv.idMap[sk] == nil {
		srv.idMap[sk] = make(map[uint]*Department, 100)
	}
	oItem := srv.codeMap[sk][item.DepartmentCode]

	// 插入或更新
	if oItem == nil {
		srv.list[sk] = append(srv.list[sk], item)
		srv.codeMap[sk][item.DepartmentCode] = item
		srv.idMap[sk][item.ID] = item
	} else {
		oItem.DepartmentName = item.DepartmentName
		oItem.DepartmentComment = item.DepartmentComment
		oItem.DepartmentLevel = item.DepartmentLevel
		oItem.DepartmentFullPath = item.DepartmentFullPath
		oItem.DepartmentParentCode = item.DepartmentParentCode
		oItem.DepartmentCustomLevel = item.DepartmentCustomLevel
		oItem.Status = item.Status
	}

}

// 通过id批量获取公司列表
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *departmentSrvice) GetByCode(tx *gorm.DB, codes []string) (items []*Department) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	items = make([]*Department, 0, 10)
	if codes == nil {
		for _, item := range srv.codeMap[sk] {
			if item.Status == 1 {
				items = append(items, item)
			}
		}
		return
	}

	for _, code := range codes {
		item := srv.codeMap[sk][code]
		if item != nil && item.Status == 1 {
			items = append(items, item)
		}
	}
	return
}

// 按照id批量删除缓存
func (srv *departmentSrvice) deleteByID(tx *gorm.DB, ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	sk := rdb.GetShardingKey(tx)
	for _, id := range ids {
		model := srv.idMap[sk][id]
		if model != nil {
			delete(srv.codeMap[sk], model.DepartmentCode)
			delete(srv.idMap[sk], model.ID)
		}
	}
}

// 以下时对外服务的API

// 创建部门
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *departmentSrvice) Create(tx *gorm.DB, param *departmentAddParams) (err error) {
	model := &Department{
		DepartmentCode:        param.DepartmentCode,
		DepartmentName:        param.DepartmentName,
		DepartmentComment:     param.DepartmentComment,
		DepartmentParentCode:  param.DepartmentParentCode,
		DepartmentCustomLevel: param.DepartmentCustomLevel,
	}

	if model.DepartmentParentCode != "" {
		pmodel := &Department{}
		err = tx.Table(model.TableName()).Where("department_code = ? and status = 1", model.DepartmentParentCode).First(pmodel).Error
		if err != nil || pmodel.DepartmentCode == "" {
			return errorx.New("无有效的父级code")
		}
	}
	err = tx.Create(model).Error
	srv.setCache(tx, model)
	srv.GenDepartmentTree(tx)
	srv.genDimList(tx)
	return
}

// 删除部门, 软删除
//
// 同时会删除部门缓存
func (srv *departmentSrvice) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	model := &Department{}

	err = tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
	if err != nil {
		return
	}

	srv.deleteByID(tx, param.ID)
	srv.GenDepartmentTree(tx)
	srv.genDimList(tx)
	return
}

// 更新部门, 通过id批量更新
//
// 同时会更新缓存
func (srv *departmentSrvice) UpdateByID(tx *gorm.DB, param *departmentEditParams) (err error) {
	model := &Department{
		DepartmentName:        param.DepartmentName,
		DepartmentComment:     param.DepartmentComment,
		DepartmentParentCode:  param.DepartmentParentCode,
		DepartmentCustomLevel: param.DepartmentCustomLevel,
	}
	model.Status = param.Status

	if model.DepartmentParentCode != "" {
		pmodel := &Department{}
		nTx := srv.core.RDB().DB().Where("1=1")
		rdb.CopyCtx(tx, nTx)
		err = nTx.Table(model.TableName()).Where("department_code = ? and status = 1", model.DepartmentParentCode).First(pmodel).Error
		if err != nil || pmodel.DepartmentCode == "" {
			return errorx.New("无有效的上级code")
		}
	}
	items := make([]*Department, 0, 100)
	err = tx.Where("id in (?) ", param.ID).Updates(model).Find(&items).Error
	if err != nil {
		return err
	}
	for _, item := range items {
		srv.setCache(tx, item)
	}
	srv.GenDepartmentTree(tx)
	srv.genDimList(tx)
	return
}

// 通过记录行,生成各种树表
//
// 分别生成list结构的数据表(前端展示用), 包含所有子节点的map(查询子级时使用)
//
// 当更新,新增,删除时, 该函数均需要重新执行以更新缓存
func (srv *departmentSrvice) GenDepartmentTree(tx *gorm.DB) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	treeMap := make(map[string]*departmentTree, 100)
	treeList := make([]*departmentTree, 0, 100)
	allChildrenMap := make(map[string]map[string]*departmentTree, 100)

	sk := rdb.GetShardingKey(tx)
	for _, item := range srv.list[sk] {

		// 如果无法从codeMap中查询到数据,表示数据被删除, 直接跳过
		if srv.codeMap[sk][item.DepartmentCode] == nil {
			continue
		}

		itemTree := &departmentTree{
			ID:                    item.ID,
			DepartmentCode:        item.DepartmentCode,
			DepartmentName:        item.DepartmentName,
			DepartmentComment:     item.DepartmentComment,
			DepartmentLevel:       item.DepartmentLevel,
			DepartmentFullPath:    item.DepartmentFullPath,
			DepartmentParentCode:  item.DepartmentParentCode,
			DepartmentCustomLevel: item.DepartmentCustomLevel,
			CreatedAT:             item.CreatedAT,
			CreatedBy:             item.CreatedBy,
			UpdatedAT:             item.UpdatedAT,
			UpdatedBy:             item.UpdatedBy,
			Status:                item.Status,
			Children:              make([]*departmentTree, 0, 10),
		}

		// 方便快速定位元素路径
		treeMap[itemTree.DepartmentCode] = itemTree
		if allChildrenMap[itemTree.DepartmentCode] == nil {
			allChildrenMap[itemTree.DepartmentCode] = make(map[string]*departmentTree, 100)
			allChildrenMap[itemTree.DepartmentCode][itemTree.DepartmentCode] = itemTree
		}
		// 根据层级判断, 如果时顶层, 写入list中, 如果是子级, 写入到对应的父级中
		// 如果父级不存在, 则子级也都不再录入
		if itemTree.DepartmentParentCode == "" {
			itemTree.DepartmentFullPath = itemTree.DepartmentName
			itemTree.DepartmentFullPath2 = itemTree.DepartmentCode
			treeList = append(treeList, itemTree)
		} else {
			pt := treeMap[itemTree.DepartmentParentCode]
			if pt != nil {
				itemTree.DepartmentLevel = pt.DepartmentLevel + 1
				itemTree.DepartmentFullPath = fmt.Sprintf("%s-%s", pt.DepartmentFullPath, itemTree.DepartmentName)
				itemTree.DepartmentFullPath2 = fmt.Sprintf("%s-%s", pt.DepartmentFullPath2, itemTree.DepartmentCode)
				// 把子部门依次添加到父级部门的列表中
				pt.Children = append(pt.Children, itemTree)
				for _, pcode := range strings.Split(pt.DepartmentFullPath2, "-") {
					allChildrenMap[pcode][itemTree.DepartmentCode] = itemTree
				}
			}
		}
	}

	DepartmentSrvice.allChildren[sk] = allChildrenMap
	DepartmentSrvice.treeList[sk] = treeList
}

func (srv *departmentSrvice) TreeList(tx *gorm.DB) []*departmentTree {
	return srv.treeList[rdb.GetShardingKey(tx)]
}

// 查询维度
func (srv *departmentSrvice) DimList(tx *gorm.DB) (result []map[string]interface{}) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)

	return srv.dimList[sk]
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *departmentSrvice) genDimList(tx *gorm.DB) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	result := make([]map[string]interface{}, 10)
	for _, item := range srv.idMap[sk] {
		result = append(result, map[string]interface{}{
			"code":   item.DepartmentCode,
			"name":   item.DepartmentName,
			"status": item.Status,
		})
	}
	srv.dimList[sk] = result
}

// 获取所有子级部门
func (srv *departmentSrvice) GetAllChildren(tx *gorm.DB, dept_code string) map[string]*departmentTree {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	sk := rdb.GetShardingKey(tx)
	return srv.allChildren[sk][dept_code]
}

// 列表的方式显示
// func departmentList() core.HandlerFunc {
// 	return func(ctx core.Context) {
// 		param := ctx.CtxGetParams().(*departmentQueryParams)
// 		model := &Department{}

// 		// 使用子查询, 优化分页查询
// 		tx := ctx.TX(
// 			core.SetTxMode(core.TX_QRY_MODE_SUBID),
// 			core.SetTxSubTable(model.TableName()),
// 		)
// 		modelList := make([]*Department, 0, 100)

// 		err := tx.Model(model).Find(&modelList).Error
// 		if err != nil {
// 			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
// 			return
// 		}
// 		var count int64
// 		ctx.TX().Model(model).Offset(-1).Limit(-1).Count(&count)
// 		data := map[string]interface{}{
// 			"page_num":  param.PageNum,
// 			"page_size": param.PageSize,
// 			"total":     count,
// 			"rows":      modelList,
// 		}
// 		ctx.JSON(errno.ERRNO_OK.WrapData(data))
// 	}
// }
