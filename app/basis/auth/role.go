/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-08-28 22:50:51
 * @Description: 请填写简介
 */
package auth

import "github.com/fbs-io/core/store/rdb"

type Role struct {
	Code        string           `gorm:"comment:角色编码" json:"code"`
	Description string           `gorm:"comment:角色描述" json:"description"`
	Sort        int              `gorm:"comment:排序" json:"json"`
	Sources     rdb.ModeListJson `gorm:"comment:角色可用资源,使用::分割;type:varchar(1000)" json:"-"`
	rdb.Model
}

func (r *Role) TableName() string {
	return "e_auth_role"
}
