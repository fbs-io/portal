/*
 * @Author: reel
 * @Date: 2023-09-18 21:26:33
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 09:31:44
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

	// 公司code生成器
	companySeq := sequence.New(route.Core(), "org_company_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("C"))

	orgGroup := route.Group("org", "组织管理").WithMeta("icon", "sc-icon-organization")

	// 可以作为帐套使用或作为环境隔离
	group := orgGroup.Group("company", "公司管理").WithMeta("icon", "sc-icon-company")
	{
		group.GET("list", "公司列表", companyQueryParams{}, companyList())
		group.PUT("add", "新增公司", companyAddParams{}, companyAdd(companySeq))
		group.POST("edit", "修改公司", companyEditParams{}, companyEdit())
		group.DELETE("delete", "删除公司", companyDeleteParams{}, companyDelete())
	}
}
