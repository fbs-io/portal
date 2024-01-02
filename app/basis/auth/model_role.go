/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2023-11-03 06:42:55
 * @Description: 角色信息管理
 */
package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
)

type Role struct {
	Code                 string           `gorm:"comment:code;unique" json:"code"`
	Label                string           `gorm:"comment:角色;unique" json:"label"`
	Sort                 int              `gorm:"comment:排序" json:"sort"`
	Description          string           `gorm:"comment:角色描述" json:"description"`
	DataPermissionType   int8             `gorm:"comment:角色数据权限" json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `gorm:"comment:自定义选择的数据权限;type:varchar(1000)" json:"data_permission_custom"`
	Sources              rdb.ModeListJson `gorm:"comment:角色可用资源;type:varchar(1000)" json:"sources"`
	rdb.Model
	rdb.ShardingModel
}

func (r *Role) TableName() string {
	return r.ShardingModel.TableName(consts.TABLE_BASIS_AUTH_ROLE)
}
