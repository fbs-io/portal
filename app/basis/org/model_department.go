/*
 * @Author: reel
 * @Date: 2023-10-28 07:41:57
 * @LastEditors: reel
 * @LastEditTime: 2024-01-16 22:41:46
 * @Description: 部门设置
 */
package org

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
)

type Department struct {
	DepartmentCode        string `gorm:"comment:部门code;unique"`
	DepartmentName        string `gorm:"comment:部门名称"`
	DepartmentComment     string `gorm:"comment:部门描述"`
	DepartmentLevel       int8   `gorm:"comment:部门层级"`
	DepartmentFullPath    string `gorm:"comment:部门全路径"`
	DepartmentParentCode  string `gorm:"comment:部门父级code"`
	DepartmentCustomLevel string `gorm:"comment:部门自定义层级"`
	rdb.Model
	rdb.ShardingModel
}

func (model *Department) TableName() string {
	return consts.TABLE_BASIS_ORG_DEPARTMENT
}

// 新增部门参数
type departmentAddParams struct {
	DepartmentCode        string `json:"department_code"`
	DepartmentName        string `json:"department_name"`
	DepartmentComment     string `json:"department_comment"`
	DepartmentParentCode  string `json:"department_parent_code"`
	DepartmentCustomLevel string `json:"department_custom_level"`
}

// 更新部门参数
type departmentEditParams struct {
	ID                    []uint `json:"id" binding:"required" conditions:"-"`
	DepartmentName        string `json:"department_name" conditions:"-"`
	DepartmentComment     string `json:"department_comment" conditions:"-"`
	DepartmentParentCode  string `json:"department_parent_code" conditions:"-"`
	DepartmentCustomLevel string `json:"department_custom_level" conditions:"-"`
	Status                int8   `json:"status" conditions:"-"`
}

// 查询公司参数
type departmentQueryParams struct {
	PageNum        int    `form:"page_num"`
	PageSize       int    `form:"page_size"`
	Orders         string `form:"orders"`
	DepartmentName string `form:"department_name" conditions:"like"`
}

type departmentTree struct {
	ID                    uint              `json:"id"`
	DepartmentCode        string            `json:"department_code"`
	DepartmentName        string            `json:"department_name"`
	DepartmentComment     string            `json:"department_comment"`
	DepartmentLevel       int8              `json:"department_level"`
	DepartmentFullPath    string            `json:"department_full_path"`
	DepartmentFullPath2   string            `json:"-"` // 用于记录组织code全路径
	DepartmentParentCode  string            `json:"department_parent_code"`
	DepartmentCustomLevel string            `json:"department_custom_level"`
	CreatedAT             uint              `json:"created_at"`
	CreatedBy             string            `json:"created_by"`
	UpdatedAT             uint              `json:"updated_at"`
	UpdatedBy             string            `json:"updated_by"`
	Status                int8              `json:"status"`
	Children              []*departmentTree `json:"children"`
}
