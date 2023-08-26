/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2023-08-24 07:23:31
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
	userList := route.Group("users", "用户管理")
	{
		userList.GET("list", "用户列表", userListQueryParams{}, userListQuery())
		userList.POST("add", "用户列表", nil, nil)
		userList.PUT("edit", "用户列表", nil, nil)
		userList.DELETE("delete", "用户列表", nil, nil)
	}

}
