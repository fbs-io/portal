/*
 * @Author: reel
 * @Date: 2023-09-18 20:14:52
 * @LastEditors: reel
 * @LastEditTime: 2024-01-15 22:40:49
 * @Description: 法人管理
 */
package org

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
)

type Company struct {
	CompanyCode      string `json:"company_code" gorm:"comment:公司code;unique"`
	CompanyName      string `json:"company_name" gorm:"comment:公司名称"`
	CompanyShortName string `json:"company_short_name" gorm:"comment:公司简称"`
	CompanyComment   string `json:"company_comment" gorm:"comment:公司描述"`
	CompanyBusiness  string `json:"company_business" gorm:"comment:所在行业"`
	rdb.Model
}

func (o *Company) TableName() string {
	return consts.TABLE_BASIS_ORG_COMPANY
}

// 新增法人参数
type companyAddParams struct {
	CompanyCode      string `json:"company_code"`
	CompanyName      string `json:"company_name"`
	CompanyShortName string `json:"company_shortname"`
	CompanyComment   string `json:"company_comment"`
	CompanyBusiness  string `json:"company_business"`
}

// 更新公司参数
type companyEditParams struct {
	ID               []uint `json:"id" binding:"required"`
	CompanyName      string `json:"company_name"  conditions:"-"`
	CompanyShortName string `json:"company_shortname"  conditions:"-"`
	CompanyComment   string `json:"company_comment"  conditions:"-"`
	CompanyBusiness  string `json:"company_business"  conditions:"-"`
	Status           int8   `json:"status" conditions:"-"`
}

// 查询公司参数
type companyQueryParams struct {
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
	Orders      string `form:"orders"`
	CompanyName string `form:"company_name" conditions:"like"`
}
