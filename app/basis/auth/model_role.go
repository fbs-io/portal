/*
 * @Author: reel
 * @Date: 2023-07-18 06:41:27
 * @LastEditors: reel
 * @LastEditTime: 2024-01-21 18:33:12
 * @Description: 角色信息管理
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
	Sources     rdb.ModeListJson `gorm:"comment:角色可用资源;type:varchar(1000)" json:"sources"`
	rdb.Model
	rdb.ShardingModel
}

func (r *Role) TableName() string {
	return r.ShardingModel.TableName(consts.TABLE_BASIS_AUTH_ROLE)
}

type roleAddParams struct {
	Code                 string           `json:"code"`
	Label                string           `json:"label"`
	Sort                 int              `json:"sort"`
	Description          string           `json:"description"`
	Sources              rdb.ModeListJson `json:"sources"`
	DataPermissionType   int8             `json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom"`
}

// orders, page_num, page_size 作为保留字段用于条件生成
type rolesQueryParams struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Orders   string `form:"orders"`
	Label    string `form:"label" conditions:"like"`
}

// 批量更新参数
//
// id作为数组, 不适用于自动查询条件生成
type roleEditParams struct {
	ID                   []uint           `json:"id"  binding:"required" conditions:"-"`
	Label                string           `json:"label" conditions:"-"`
	Sort                 int              `json:"json" conditions:"-"`
	Description          string           `json:"description" conditions:"-"`
	Sources              rdb.ModeListJson `json:"sources" conditions:"-"`
	Status               int8             `json:"status" conditions:"-"`
	DataPermissionType   int8             `json:"data_permission_type" conditions:"-"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom" conditions:"-"`
}
