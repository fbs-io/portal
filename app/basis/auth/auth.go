/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2023-09-04 06:40:31
 * @Description: 请填写简介
 */
package auth

import (
	"sync"

	"github.com/fbs-io/core"
	"gorm.io/gorm"
)

var (
	userMap = make(map[uint]*User)
	roleMap = make(map[uint]*Role)
	lock    = &sync.RWMutex{}
)

func SetUser(id uint, user *User) {
	lock.Lock()
	defer lock.Unlock()
	userMap[id] = user
}

func GetUser(id uint, tx *gorm.DB) (user *User) {
	lock.RLock()
	defer lock.RUnlock()
	user = userMap[id]
	if user == nil {
		user = &User{}
		err := tx.Model(user).Find(user).Error
		if err != nil {
			return nil
		}
	}
	return
}

func SetRole(id uint, role *Role) {
	lock.Lock()
	defer lock.Unlock()
	roleMap[id] = role
}

func GetRole(id uint, tx *gorm.DB) (role *Role) {
	lock.RLock()
	defer lock.RUnlock()
	role = roleMap[id]
	if role == nil {
		role = &Role{}
		err := tx.Model(role).Where("id = (?)", id).Find(role).Error
		if err != nil {
			return nil
		}
	}
	return
}

func New(route core.RouterGroup) {

	user := &User{}
	user.Account = "root"
	user.Password = "root123"
	user.NickName = "超级管理员"
	user.Super = "Y"
	// 注册表
	tx := route.Core().RDB()
	tx.Register(&User{}, func() error { return tx.DB().Create(user).Error })
	tx.Register(&Role{})

	// 用户个人信息操作
	//
	// 如登陆, 注销, 密码变更, 信息更改等
	auth := route.Group("user", "账户信息").WithPermission(core.SOURCE_TYPE_UNMENU)
	{

		auth.POST("login", "登陆", loginParams{}, login()).WithAllowSignature()
		auth.PUT("chpwd", "变更密码", userChPwdParams{}, chpwd())
		auth.PUT("update", "更新账户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
	}

	// 用户列表管理, 用于批量管理用户
	users := route.Group("users", "用户管理")
	{

		users.PUT("add", "添加用户", userAddParams{}, userAdd())
		// 获取用户列表
		users.GET("list", "用户列表", usersQueryParams{}, usersQuery())
		// 批量更新用户
		users.POST("edit", "更新用户", usersUpdateParams{}, usersUpdate())
		// 批量更新, 也适用于单个用户
		users.DELETE("delete", "删除用户", usersDeleteParams{}, usersDelete())
	}

	// role := route.Group("role", "单角色管理").WithPermission(core.SOURCE_TYPE_LIMITED)

	// {
	// 	role.PUT("add", "添加角色", roleAddParams{}, roleAdd())
	// }

	roles := route.Group("roles", "角色管理")
	{
		// 单行添加
		roles.PUT("add", "添加角色", roleAddParams{}, roleAdd())

		// 批量操作
		roles.GET("list", "角色列表", rolesQueryParams{}, rolesQuery())
		roles.POST("edit", "编辑角色", rolesUpdateParams{}, rolesUpdate())
		roles.DELETE("delete", "删除角色", rolesDeleteParams{}, rolesDelete())
	}

	menus := route.Group("menus", "菜单管理")
	{
		// 单行添加
		menus.PUT("add", "添加角色", menusAddParams{}, menusAdd())
		menus.GET("list", "添加角色", nil, menusQuery())
		menus.POST("edit", "添加角色", menusUpdateParams{}, menusUpdate())
		menus.DELETE("delete", "添加角色", menusDeleteParams{}, menusDelete())
	}
}
