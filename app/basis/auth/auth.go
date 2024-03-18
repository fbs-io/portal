/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2024-03-17 16:22:49
 * @Description: 请填写简介
 */
package auth

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/env"
	"github.com/fbs-io/core/store/rdb"
)

func New(route core.RouterGroup) {
	// 角色code生成器
	roleSeq := sequence.New(route.Core(), "auth_role_sequence", sequence.SetDateFormat(""), sequence.SetPrefix("R"))
	// 注册超级管理员, 根据设置的后台管理员注册app管理员
	tx := route.Core().RDB()
	tx.Register(&User{}, func() error {
		return tx.DB().Create(&User{
			Account:  route.Core().Config().User,
			Password: route.Core().Config().Pwd,
			NickName: "超级管理员",
			Super:    "Y",
		}).Error
	})

	tx.Register(&Role{})
	tx.Register(&RlatUserRole{})
	tx.Register(&RlatUserCompany{})
	tx.Register(&RlatUserPosition{})

	auth := route.Group("auth", "用户中心").WithMeta("icon", "el-icon-stamp")
	// 用户个人信息操作
	//
	// 如登陆, 注销, 密码变更, 信息更改等
	info := auth.Group("user", "账户信息").WithPermission(core.SOURCE_TYPE_UNMENU).WithHidden()
	{

		// 允许鉴权例外
		info.POST("login", "登陆", loginParams{}, login()).WithAllowSignature()
		// 登陆前获取用户的公司列表
		info.GET("company", "获取公司列表", userCompanyParams{}, getCompanyAndPosition()).WithAllowSignature().WithPermission(core.SOURCE_TYPE_UNLIMITED)
		// 登陆后获取用户的公司列表和岗位列表
		info.GET("allowcompany", "获取公司列表", userCompanyParams{}, getCompanyAndPosition()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
		// 个人账户更改密码不受限
		info.POST("chpwd", "变更密码", userChPwdParams{}, chpwd()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		// 个人用户修改信息不受限
		info.POST("update", "更新账户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		// 设置/切换公司
		info.POST("default_company", "设置所在公司", setDefaultCompanyParams{}, setDefaultCompany()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		info.POST("default_position", "设置所在岗位", setDefaultPositionParams{}, setDefaultPosition()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
	}

	// 用户列表管理, 用于批量管理用户
	users := auth.Group("users", "用户管理").WithMeta("icon", "el-icon-user")
	{
		// 获取用户列表
		users.GET("list", "用户列表", usersQueryParams{}, usersQuery())
		// 添加用户
		users.PUT("add", "添加用户", userAddParams{}, userAdd())
		// 更新用户
		users.POST("edit", "更新用户", usersEditParams{}, usersUpdate())
		// 删除用户, 逻辑删除,
		users.DELETE("delete", "删除用户", rdb.DeleteParams{}, usersDelete())
	}

	roles := auth.Group("roles", "角色管理").WithMeta("icon", "el-icon-switch-filled")
	{
		// 角色列表
		roles.GET("list", "角色列表", rolesQueryParams{}, rolesQuery())
		// 单行添加
		roles.PUT("add", "添加角色", roleAddParams{}, roleAdd(roleSeq))
		// 更新用户
		roles.POST("edit", "编辑角色", roleEditParams{}, rolesUpdate())
		// 删除用户, 逻辑删除
		roles.DELETE("delete", "删除角色", rdb.DeleteParams{}, rolesDelete())
		// 权限菜单列表
		roles.GET("permission", "菜单列表", nil, menusQueryWithPermission())
	}

	//开发环境使用
	if env.Active().Value() == env.ENV_MODE_DEV {
		menus := auth.Group("menus", "菜单管理")
		{
			// 菜单列表
			menus.GET("list", "菜单列表", nil, menusQueryWithManager())
			// 添加菜单
			menus.PUT("add", "添加菜单", menusAddParams{}, menusAdd())
			// 更新菜单
			menus.POST("edit", "更新菜单", menusUpdateParams{}, menusUpdate())
			// 删除菜单, 逻辑删除
			menus.DELETE("delete", "添加角色", rdb.DeleteParams{}, menusDelete())
		}
	}

	route.Core().AddStartJobList(func() error {
		ResourceServiceInit(route.Core())
		RoleServiceInit(route.Core())
		UserServiceInit(route.Core())
		return nil
	})
}
