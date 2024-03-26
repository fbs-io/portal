/*
 * @Author: reel
 * @Date: 2024-01-14 22:23:07
 * @LastEditors: reel
 * @LastEditTime: 2024-03-18 22:00:10
 * @Description: 公司表相关逻辑处理
 */
package org

import (
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type companyService struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[uint]*Company
	codeMap map[string]*Company
	dimList []map[string]interface{}
}

var CompanyService *companyService

// 初始化
func CompanyServiceInit(c core.Core) {
	CompanyService = &companyService{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[uint]*Company, 100),
		codeMap: make(map[string]*Company, 100),
		dimList: make([]map[string]interface{}, 0, 100),
	}
	var companys = make([]*Company, 0, 100)
	c.RDB().DB().Order("id").Find(&companys)

	for _, company := range companys {
		CompanyService.setCache(company)
	}
	CompanyService.genDimList()
}

// 创建缓存
func (srv *companyService) setCache(item *Company) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	oItem := srv.codeMap[item.CompanyCode]

	// 插入或更新
	if oItem == nil {
		srv.codeMap[item.CompanyCode] = item
		srv.idMap[item.ID] = item
	} else {
		oItem.CompanyName = item.CompanyName
		oItem.CompanyShortName = item.CompanyShortName
		oItem.CompanyComment = item.CompanyComment
		oItem.CompanyBusiness = item.CompanyBusiness
		oItem.Status = item.Status
	}
	srv.codeMap[item.CompanyCode] = item
	srv.idMap[item.ID] = item
}

// 通过code列表获取有效公司
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *companyService) GetByCodes(codes []string) (list []*Company) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	list = make([]*Company, 0, 10)

	for _, code := range codes {
		company := srv.codeMap[code]
		if company != nil && company.Status == 1 {
			list = append(list, srv.codeMap[code])
		}
	}
	return
}

// 获取所有有效的公司
// 优化使用其他方法
func (srv *companyService) GetAll() (list []*Company) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	list = make([]*Company, 0, 10)
	for _, item := range srv.idMap {
		if item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 只获取有效的公司
func (srv *companyService) GetByCode(code string) (item *Company) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	if srv.codeMap[code].Status == 1 {
		return srv.codeMap[code]
	}
	return
}

// 按照id批量删除缓存
func (srv *companyService) deleteByID(ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, id := range ids {
		company := srv.idMap[id]
		if company != nil {
			delete(srv.codeMap, company.CompanyCode)
			delete(srv.idMap, company.ID)
		}
	}
}

// 以下时对外服务的API

// 创建公司
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
func (srv *companyService) Create(tx *gorm.DB, param *companyAddParams) (err error) {
	company := &Company{
		CompanyCode:      param.CompanyCode,
		CompanyName:      param.CompanyName,
		CompanyShortName: param.CompanyShortName,
		CompanyComment:   param.CompanyComment,
		CompanyBusiness:  param.CompanyBusiness,
	}
	err = tx.Create(company).Error
	if err != nil {
		return
	}
	srv.setCache(company)
	// 创建公司相关的库或者表
	err = srv.core.RDB().AddShardingSuffixs(company.CompanyCode)
	if err != nil {
		return errorx.Wrap(err, "迁移数据库/表发生错误")
	}
	srv.genDimList()
	return
}

// 删除公司, 软删除
//
// 同时会删除公司缓存
func (srv *companyService) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	company := &Company{}

	err = tx.Model(company).Where("id in (?)", param.ID).Delete(company).Error
	if err != nil {
		return
	}
	srv.deleteByID(param.ID)
	srv.genDimList()
	return
}

// 更新公司, 通过id批量更新
//
// 同时会更新缓存
func (srv *companyService) UpdateByID(tx *gorm.DB, param *companyEditParams) (err error) {
	company := &Company{
		CompanyName:      param.CompanyName,
		CompanyShortName: param.CompanyShortName,
		CompanyComment:   param.CompanyComment,
		CompanyBusiness:  param.CompanyBusiness,
	}
	company.Status = param.Status
	var companys = make([]*Company, 0, 100)
	err = tx.Where("id in (?) ", param.ID).Updates(company).Find(&companys).Error
	if err != nil {
		return
	}
	for _, co := range companys {
		srv.setCache(co)
	}
	srv.genDimList()
	return
}

// 查询
//
// TODO: 缓存结果
func (srv *companyService) Query(tx *gorm.DB, param *companyQueryParams) (data interface{}, err error) {
	model := &Company{}
	list := make([]*Company, 0, 100)

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
func (srv *companyService) DimList() (result []map[string]interface{}) {
	return srv.dimList
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *companyService) genDimList() {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	result := make([]map[string]interface{}, 0, 10)
	for _, company := range srv.idMap {
		result = append(result, map[string]interface{}{
			"code":   company.CompanyCode,
			"name":   company.CompanyName,
			"status": company.Status,
		})
	}
	srv.dimList = result
}
