/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-07-18 22:10:11
 * @Description: 请填写简介
 */
package auth

import "github.com/fbs-io/core/store/rdb"

type RoleBase struct {
	Code        string   `gorm:"comment:角色编码" json:"code"`
	Description string   `gorm:"comment:角色描述" json:"description"`
	Source      string   `gorm:"comment:角色可用资源,使用::分割" json:"-" `
	Sources     []string `json:"sources" gorm:"-"`
}

type Role struct {
	RoleBase
	rdb.Model
}

func (r *Role) TableName() string {
	return "e_auth_role"
}
