/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-09-01 06:27:12
 * @Description: 请填写简介
 */
package auth

import (
	"github.com/fbs-io/core/store/rdb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Code        string           `gorm:"comment:code;unique" json:"code"`
	Label       string           `gorm:"comment:角色;unique" json:"label"`
	Sort        int              `gorm:"comment:排序" json:"sort"`
	Description string           `gorm:"comment:角色描述" json:"description"`
	Sources     rdb.ModeListJson `gorm:"comment:角色可用资源,使用::分割;type:varchar(1000)" json:"sources"`
	rdb.Model
}

// type RoleList struct {
// 	Code        string           `json:"code"`
// 	ID          uint             `json:"id"`
// 	Label       string           `json:"label"`
// 	Sort        int              `json:"json"`
// 	Description string           `json:"description"`
// 	Sources     rdb.ModeListJson `json:"sources" gorm:"type:varchar(1000)"`
// 	CreatedAt   uint64           `json:"created_at"`
// 	Status      int8             `json:"status"`
// }

func (r *Role) TableName() string {
	return "e_auth_role"
}

// gorm 中间件操作
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	r.Code = uuid.New().String()
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

// 模型接口

func (r *Role) RoleInfo() map[string]interface{} {
	return map[string]interface{}{}
}
