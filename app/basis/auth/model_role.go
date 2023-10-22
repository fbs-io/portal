/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-10-17 18:49:55
 * @Description: 请填写简介
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
)

type Role struct {
	Code        string           `gorm:"comment:code;unique" json:"code"`
	Label       string           `gorm:"comment:角色;unique" json:"label"`
	Sort        int              `gorm:"comment:排序" json:"sort"`
	Description string           `gorm:"comment:角色描述" json:"description"`
	Sources     rdb.ModeListJson `gorm:"comment:角色可用资源,使用::分割;type:varchar(1000)" json:"sources"`
	rdb.Model
	rdb.ShardingModel
}

func (r *Role) TableName() string {
	return r.ShardingModel.TableName(consts.TABLE_BASIS_AUTH_ROLE)
}
