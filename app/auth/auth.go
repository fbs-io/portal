/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2023-08-21 23:32:01
 * @Description: 请填写简介
 */
package auth

import (
	"sync"

	"github.com/fbs-io/core"
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

func GetUser(id uint) (user *User) {
	lock.RLock()
	defer lock.RUnlock()
	return userMap[id]
}

func New(c core.Core) {

	user := &User{}
	user.Account = "root"
	user.Password = "root123"
	user.NickName = "超级管理员"
	user.Super = "Y"
	// 注册表
	c.RDB().Register(&User{}, func() error { return c.RDB().DB().Create(user).Error })
	c.RDB().Register(&Role{})

	// 登陆
	ajax := c.Group("ajax").WithPermission(core.SOURCE_TYPE_LIMITED)
	{
		ajax.POST("login", "登陆", loginParams{}, login()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
	}

	// 用户
	userinfo := ajax.Group("user", "用户管理").WithPermission(core.SOURCE_TYPE_UNPERMISSION).WithRouter(core.SOURCE_ROUTER_IS)
	{
		// 用户登陆注销等操作
		userinfo.PUT("chpwd", "变更密码", userChPwdParams{}, chpwd())
		userinfo.PUT("update", "更新用户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
	}

}
