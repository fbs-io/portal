/*
 * @Author: reel
 * @Date: 2023-09-18 20:14:52
 * @LastEditors: reel
 * @LastEditTime: 2023-09-18 21:39:42
 * @Description: 组织管理
 */
package org

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type Org struct {
	OrgCode      string `gorm:"comment:组织code"`
	OrgName      string `gorm:"comment:组织名称"`
	OrgShortName string `gorm:"comment:组织简称"`
	OrgComment   string `gorm:"comment:组织描述"`
	OrgBusiness  string `gorm:"comment:所在行业"`
	rdb.Model
}

func (o *Org) TableName() string {
	return consts.TABLE_BASIS_ORG
}

func (o *Org) BeforeCreate(tx *gorm.DB) error {
	o.Model.BeforeCreate(tx)
	return nil
}

func (o *Org) BeforeUpdate(tx *gorm.DB) error {
	o.Model.BeforeUpdate(tx)
	return nil
}
