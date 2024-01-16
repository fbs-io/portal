/*
 * @Author: reel
 * @Date: 2024-01-14 22:23:07
 * @LastEditors: reel
 * @LastEditTime: 2024-01-16 23:41:57
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

type companySrvice struct {
	lock    *sync.RWMutex
	core    core.Core
	idMap   map[uint]*Company
	codeMap map[string]*Company
	dimList []map[string]interface{}
}

var CompanySrvice *companySrvice

// 初始化
func CompanySrviceInit(c core.Core) {
	CompanySrvice = &companySrvice{
		core:    c,
		lock:    &sync.RWMutex{},
		idMap:   make(map[uint]*Company, 100),
		codeMap: make(map[string]*Company, 100),
		dimList: make([]map[string]interface{}, 0, 100),
	}
	var companys = make([]*Company, 0, 100)
	c.RDB().DB().Where("deleted_at = 0").Order("id").Find(&companys)

	for _, company := range companys {
		CompanySrvice.setCache(company)
	}
}

// 创建缓存
func (srv *companySrvice) setCache(item *Company) {
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

// 通过id批量获取公司列表
//
// 当入参为nil时, 获取所有且有效的数据
func (srv *companySrvice) GetByCode(codes []string) (companys []*Company) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	companys = make([]*Company, 0, 10)
	if codes == nil {
		for _, company := range srv.codeMap {
			if company.Status == 1 {
				companys = append(companys, company)
			}
		}
		return
	}

	for _, code := range codes {
		company := srv.codeMap[code]
		if company != nil && company.Status == 1 {
			companys = append(companys, srv.codeMap[code])
		}
	}
	return
}

// 按照id批量删除缓存
func (srv *companySrvice) deleteByID(ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	// companys = make([]*Company, 10)
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
func (srv *companySrvice) Create(tx *gorm.DB, param *companyAddParams) (err error) {
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
func (srv *companySrvice) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
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
func (srv *companySrvice) UpdateByID(tx *gorm.DB, param *companyEditParams) (err error) {
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

// 查询维度
func (srv *companySrvice) DimList() (result []map[string]interface{}) {
	return srv.dimList
}

// 生成维度信息, 简化变量生成的次数, 减少垃圾回收次数
func (srv *companySrvice) genDimList() {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	result := make([]map[string]interface{}, 10)
	for _, company := range srv.idMap {
		result = append(result, map[string]interface{}{
			"code":   company.CompanyCode,
			"name":   company.CompanyName,
			"status": company.Status,
		})
	}
	srv.dimList = result
}
