/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 09:32:14
 * @Description: 请填写简介
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type Role struct {
	CompanyCode string           `gorm:"index"`
	Code        string           `gorm:"comment:code;unique" json:"code"`
	Label       string           `gorm:"comment:角色;unique" json:"label"`
	Sort        int              `gorm:"comment:排序" json:"sort"`
	Description string           `gorm:"comment:角色描述" json:"description"`
	Sources     rdb.ModeListJson `gorm:"comment:角色可用资源,使用::分割;type:varchar(1000)" json:"sources"`
	rdb.Model
}

func (r *Role) TableName() string {
	return consts.TABLE_BASIS_AUTH_ROLE
}

// gorm 中间件操作
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	r.Model.BeforeCreate(tx)
	return nil
}

func (r *Role) BeforeUpdate(tx *gorm.DB) error {
	r.Model.BeforeUpdate(tx)
	return nil
}

func (r *Role) BeforeDelete(tx *gorm.DB) error {
	r.Model.BeforeDelete(tx)
	return nil
}
