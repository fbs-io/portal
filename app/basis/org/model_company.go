/*
 * @Author: reel
 * @Date: 2023-09-18 20:14:52
 * @LastEditors: reel
 * @LastEditTime: 2023-10-05 21:04:15
 * @Description: 组织管理
 */
package org

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type Company struct {
	CompanyCode      string `json:"company_code" gorm:"comment:公司code;unique"`
	CompanyName      string `json:"company_name" gorm:"comment:公司名称"`
	CompanyShortName string `json:"company_shortname" gorm:"comment:公司简称"`
	CompanyComment   string `json:"company_comment" gorm:"comment:公司描述"`
	CompanyBusiness  string `json:"company_business" gorm:"comment:所在行业"`
	rdb.Model
}

func (o *Company) TableName() string {
	return consts.TABLE_BASIS_ORG_COMPANY
}

func (o *Company) BeforeCreate(tx *gorm.DB) error {
	o.Model.BeforeCreate(tx)
	return nil
}

func (o *Company) BeforeUpdate(tx *gorm.DB) error {
	o.Model.BeforeUpdate(tx)
	return nil
}
