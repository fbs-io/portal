/*
 * @Author: reel
 * @Date: 2023-08-23 06:20:02
 * @LastEditors: reel
 * @LastEditTime: 2023-08-23 21:27:17
 * @Description: 基础信息管理, 包含模块: 用户中心, 系统设置等信息
 */
package basis

import (
	"fbs-portal/app/basis/auth"

	"github.com/fbs-io/core"
)

func New(route core.RouterGroup) {
	basis := route.Group("basis", "基础信息管理")

	// 用户中心管理, 包含用户, 权限, 角色的管理
	auth.New(basis)

	// user := &User{}
	// user.Account = "root"
	// user.Password = "root123"
	// user.NickName = "超级管理员"
	// user.Super = "Y"
	// // 注册表
	// c.RDB().Register(&User{}, func() error { return c.RDB().DB().Create(user).Error })
	// c.RDB().Register(&Role{})

	// // 登陆
	// ajax := c.Group("ajax").WithPermission(core.SOURCE_TYPE_LIMITED)
	// {
	// 	ajax.POST("login", "登陆", loginParams{}, login()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
	// }

	// // 用户
	// userinfo := ajax.Group("user", "用户管理").WithPermission(core.SOURCE_TYPE_UNPERMISSION).WithRouter(core.SOURCE_ROUTER_IS)
	// {
	// 	// 用户登陆注销等操作
	// 	userinfo.PUT("chpwd", "变更密码", userChPwdParams{}, chpwd())
	// 	userinfo.PUT("update", "更新用户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
	// }

}
