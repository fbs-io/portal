/*
 * @Author: reel
 * @Date: 2024-01-20 20:24:52
 * @LastEditors: reel
 * @LastEditTime: 2024-03-21 06:59:26
 * @Description: 用户相关关系表
 */

package auth

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
)

// 用户和岗位关系表
type RlatUserPosition struct {
	Account      string `gorm:"comment:用户code;index"`
	PositionCode string `gorm:"comment:岗位code;index"`
	IsPosition   int8   `gorm:"comment:是否主岗"`
	rdb.Model
	rdb.ShardingModel
}

func (model *RlatUserPosition) TableName() string {
	return consts.TABLE_BASIS_RLAT_USER_POSITION
}

// 用户和角色关系表
type RlatUserRole struct {
	Account  string `gorm:"comment:用户code;index"`
	RoleCode string `gorm:"comment:角色code;index"`
	rdb.Model
	rdb.ShardingModel
}

func (model *RlatUserRole) TableName() string {
	return consts.TABLE_BASIS_RLAT_USER_ROLE
}

// 角色和资源关系表
type RlatRoleResource struct {
	RoleCode     string `gorm:"comment:角色code;index"`
	ResourceCode string `gorm:"comment:资源code;index"`
	rdb.Model
	rdb.ShardingModel
}

func (model *RlatRoleResource) TableName() string {
	return consts.TABLE_BASIS_RLAT_ROLE_RESOURCE
}

// 用户和公司关系表
// 不设置分区表
type RlatUserCompany struct {
	Account     string `gorm:"comment:用户code;index"`
	CompanyCode string `gorm:"comment:公司code;index"`
	rdb.Model
}

func (model *RlatUserCompany) TableName() string {
	return consts.TABLE_BASIS_RLAT_USER_COMPANY
}
