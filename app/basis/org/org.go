/*
 * @Author: reel
 * @Date: 2023-09-18 21:26:33
 * @LastEditors: reel
 * @LastEditTime: 2023-10-28 21:49:46
 * @Description: 请填写简介
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
)

func New(route core.RouterGroup) {
	tx := route.Core().RDB()
	tx.Register(&Company{})
	tx.Register(&Department{})

	// 公司code生成器
	companySeq := sequence.New(route.Core(), "org_company_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("C"))
	// 组织code生成器
	departmentSeq := sequence.New(route.Core(), "org_company_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("D"))

	orgGroup := route.Group("org", "组织管理").WithMeta("icon", "sc-icon-organization")

	// 可以作为帐套使用或作为环境隔离
	company := orgGroup.Group("company", "公司管理").WithMeta("icon", "sc-icon-company")
	{
		company.GET("list", "公司列表", companyQueryParams{}, companyList())
		company.PUT("add", "新增公司", companyAddParams{}, companyAdd(companySeq))
		company.POST("edit", "修改公司", companyEditParams{}, companyEdit())
		company.DELETE("delete", "删除公司", companyDeleteParams{}, companyDelete())
	}

	department := orgGroup.Group("department", "部门管理").WithMeta("icon", "sc-icon-organization")
	{
		// 获取部门列表
		department.GET("list", "部门列表", departmentQueryParams{}, departmentList())
		// 获取部门树结构
		department.GET("tree", "部门树", departmentQueryParams{}, getDepartmentTree())
		department.PUT("add", "新增部门", departmentAddParams{}, departmentAdd(departmentSeq))
		department.POST("edit", "修改部门", departmentEditParams{}, departmentEdit())
		department.DELETE("delete", "删除部门", departmentDeleteParams{}, departmentDelete())
	}
}
