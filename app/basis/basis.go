/*
 * @Author: reel
 * @Date: 2023-08-23 06:20:02
 * @LastEditors: reel
 * @LastEditTime: 2023-09-27 22:38:04
 * @Description: 基础信息管理, 包含模块: 用户中心, 系统设置等信息
 */
package basis

import (
	"portal/app/basis/auth"
	"portal/app/basis/org"

	"github.com/fbs-io/core"
)

func New(route core.RouterGroup) {
	basis := route.Group("basis", "配置").WithMeta("icon", "el-icon-aim")

	// 用户中心管理, 包含用户, 权限, 角色的管理
	auth.New(basis)

	// 包含组织, 部门管理, 未来增加法人公司管理, 用于复杂组织架构的管理
	org.New(basis)
}
