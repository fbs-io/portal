/*
 * @Author: reel
 * @Date: 2023-12-20 23:37:54
 * @LastEditors: reel
 * @LastEditTime: 2024-01-20 15:30:48
 * @Description: 岗位管理, 岗位关联部门, 同时可以设置相关数据权限, 审批权限等
 */
package org

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"gorm.io/gorm"
)

type Position struct {
	PositionCode         string           `gorm:"comment:岗位code;unique" json:"position_code"`
	PositionName         string           `gorm:"comment:岗位名称" json:"position_name"`
	PositionComment      string           `gorm:"comment:岗位描述" json:"position_comment"`
	PositionParentCode   string           `gorm:"comment:上级岗位" json:"position_parent_code"`
	DepartmentCode       string           `gorm:"comment:部门code;index" json:"department_code"`
	DeptPositionName     string           `gorm:"comment:带部门的岗位名称" json:"dept_position_name"`
	JobCode              string           `gorm:"comment:职务code;index" json:"job_code"`
	IsHead               int8             `gorm:"comment:是否是部门领导" json:"is_head"`
	IsApprove            int8             `gorm:"comment:是否有审批权限" json:"is_approve"`
	IsVritual            int8             `gorm:"comment:是否是虚拟岗位" json:"is_vritual"`
	DataPermissionType   int8             `gorm:"comment:数据权限类型" json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `gorm:"comment:数据权限自定义列表;type:varchar(1024)" json:"data_permission_custom"`
	rdb.Model
	rdb.ShardingModel
}

func (model *Position) TableName() string {
	return consts.TABLE_BASIS_ORG_POSITION
}

func (model *Position) BeforeCreate(tx *gorm.DB) error {
	if model.IsHead == 0 {
		model.IsHead = -1
	}
	if model.IsApprove == 0 {
		model.IsApprove = -1
	}
	if model.IsVritual == 0 {
		model.IsVritual = -1
	}
	return nil
}

// 岗位新增参数
type positionAddParams struct {
	PositionCode         string           `json:"position_code"`
	PositionName         string           `json:"position_name" binding:"required"`
	PositionComment      string           `json:"position_comment"`
	PositionParentCode   string           `json:"position_parent_code"`
	DepartmentCode       string           `json:"department_code" binding:"required"`
	JobCode              string           `json:"job_code"`
	IsHead               int8             `json:"is_head"`
	IsApprove            int8             `json:"is_approve"`
	IsVritual            int8             `json:"is_vritual"`
	DataPermissionType   int8             `json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom"`
}

// 岗位修改参数
type positionEditParams struct {
	ID                   []uint           `json:"id" binding:"required" conditions:"-"`
	PositionName         string           `json:"position_name" conditions:"-"`
	PositionComment      string           `json:"position_comment" conditions:"-"`
	PositionParentCode   string           `json:"position_parent_code" conditions:"-"`
	DepartmentCode       string           `json:"department_code" conditions:"-"`
	JobCode              string           `json:"job_code" conditions:"-"`
	IsHead               int8             `json:"is_head" conditions:"-"`
	IsApprove            int8             `json:"is_approve" conditions:"-"`
	IsVritual            int8             `json:"is_vritual" conditions:"-"`
	Status               int8             `json:"status" conditions:"-"`
	DataPermissionType   int8             `json:"data_permission_type" conditions:"-"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom" conditions:"-"`
}

// 岗位查询参数
type positionQueryParams struct {
	PageNum      int    `form:"page_num"`
	PageSize     int    `form:"page_size"`
	Orders       string `form:"orders"`
	PositionName string `form:"position_name" conditions:"like"`
}
