/*
 * @Author: reel
 * @Date: 2023-10-28 07:41:57
 * @LastEditors: reel
 * @LastEditTime: 2023-11-15 06:27:52
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
