/*
 * @Author: reel
 * @Date: 2023-09-18 21:26:33
 * @LastEditors: reel
 * @LastEditTime: 2023-09-27 21:48:25
 * @Description: 请填写简介
 */
package org

import "github.com/fbs-io/core"

func New(route core.RouterGroup) {
	tx := route.Core().RDB()
	tx.Register(&Org{})

	orgGroup := route.Group("org", "组织管理")
	{
		orgGroup.GET("add", "新增组织", OrgAddParams{}, add()).WithPermission(core.SOURCE_TYPE_UNLIMITED).WithAllowSignature()
	}
}
