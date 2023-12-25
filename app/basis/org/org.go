/*
 * @Author: reel
 * @Date: 2023-09-18 21:26:33
 * @LastEditors: reel
 * @LastEditTime: 2023-12-24 16:48:19
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
	tx.Register(&Position{})

	// 公司code生成器
	companySeq := sequence.New(route.Core(), "org_company_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("C"))
	// 组织code生成器
	departmentSeq := sequence.New(route.Core(), "org_company_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("D"))
	// 岗位code生成器
	positionSeq := sequence.New(route.Core(), "org_position_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("P"))

	orgGroup := route.Group("org", "组织管理").WithMeta("icon", "sc-icon-organization")

	// 可以作为帐套使用或作为环境隔离
	company := orgGroup.Group("company", "公司管理").WithMeta("icon", "sc-icon-company")
	{
		company.GET("list", "公司列表", companyQueryParams{}, companyList())
		company.PUT("add", "新增公司", companyAddParams{}, companyAdd(companySeq))
		company.POST("edit", "修改公司", companyEditParams{}, companyEdit())
		company.DELETE("delete", "删除公司", companyDeleteParams{}, companyDelete())
	}

	// 部门管理相关
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

	// 岗位管理相关
	position := orgGroup.Group("position", "岗位管理").WithMeta("icon", "sc-icon-position")
	{
		position.GET("list", "岗位列表", positionQueryParams{}, positionList())
		position.PUT("add", "新增岗位", positionAddParams{}, positionAdd(positionSeq))
		position.POST("edit", "修改岗位", positionEditParams{}, positionEdit())
		position.DELETE("delete", "删除岗位", positionDeleteParams{}, positionDelete())
	}
}
