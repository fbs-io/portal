/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 05:41:43
 * @Description: 请填写简介
 */
package auth

import (
	"sync"

	"github.com/fbs-io/core"
	"gorm.io/gorm"
)

var (
	userMap = make(map[string]*User, 100)
	roleMap = make(map[uint]*Role, 100)
	lock    = &sync.RWMutex{}
)

func SetUser(auth string, user *User) {
	lock.Lock()
	defer lock.Unlock()
	userMap[auth] = user
}

func GetUser(auth string, tx *gorm.DB) (user *User) {
	lock.RLock()
	defer lock.RUnlock()
	user = userMap[auth]
	if user == nil {
		user = &User{Account: auth}
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
	info := route.Group("user", "账户信息").WithPermission(core.SOURCE_TYPE_UNMENU).WithHidden()
	{

		// 允许鉴权例外
		info.POST("login", "登陆", loginParams{}, login()).WithAllowSignature()
		// 个人账户更改密码不受限
		info.PUT("chpwd", "变更密码", userChPwdParams{}, chpwd()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		// 个人用户修改信息不受限
		info.PUT("update", "更新账户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
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
