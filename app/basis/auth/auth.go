/*
 * @Author: reel
 * @Date: 2023-07-18 07:44:55
 * @LastEditors: reel
 * @LastEditTime: 2023-12-26 22:53:08
 * @Description: 请填写简介
 */
package auth

import (
	"portal/pkg/sequence"
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/env"
	"gorm.io/gorm"
)

var (
	userMap = make(map[string]*User, 100)
	roleMap = make(map[uint]*Role, 100)
	lock    = &sync.RWMutex{}
)

const (

	// 用户缓存刷新方式
	// 不刷新
	REFRESH_NOT int8 = iota
	// 刷新user
	REFRESH_USRE
	// 刷新user和菜单
	REFRESH_ALL
)

func SetUser(auth string, user *User) {
	lock.Lock()
	defer lock.Unlock()
	userMap[auth] = user
}

func GetUser(auth string, ctx core.Context, refresh int8) (user *User) {
	lock.RLock()
	defer lock.RUnlock()
	user = userMap[auth]
	if refresh == REFRESH_NOT {
		if user == nil {
			user = &User{Account: auth}
			err := ctx.NewTX().Model(user).Where("account = (?)", auth).Find(user).Error
			if err != nil {
				return nil
			}
			user.Menu, user.Permissions, _ = getMenuTree(ctx, user, QUERY_MENU_MODE_INFO)
		}
	}

	if refresh >= REFRESH_USRE {
		user = nil
	}
	if user == nil {
		user = &User{Account: auth}
		err := ctx.NewTX().Model(user).Where("account = (?)", auth).Find(user).Error
		if err != nil {
			return nil
		}

		// 继续判断是否有分区code, 如果无法从上下文获取code
		// 则尝试使用用户自带公司列表的第一个值作为公司code用于分区传递给上下文
		if ctx.CtxGet(core.CTX_SHARDING_KEY) == nil && user.Company != nil && len(user.Company) > 0 {
			ctx.CtxSet(core.CTX_SHARDING_KEY, user.Company[0].(string))
		}
		if refresh == REFRESH_ALL || user.Permissions == nil {
			user.Menu, user.Permissions, _ = getMenuTree(ctx, user, QUERY_MENU_MODE_INFO)
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
	r := &Role{}
	// tx.Register(r).AddShardingTable(r.TableName())
	tx.Register(r)
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
		info.GET("company", "获取公司列表", userCompanyParams{}, getCompany()).WithAllowSignature().WithPermission(core.SOURCE_TYPE_UNLIMITED)
		// 登陆后获取用户的公司列表
		info.GET("allowcompany", "获取公司列表", userCompanyParams{}, getCompany()).WithPermission(core.SOURCE_TYPE_UNLIMITED)
		// 个人账户更改密码不受限
		info.POST("chpwd", "变更密码", userChPwdParams{}, chpwd()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		// 个人用户修改信息不受限
		info.POST("update", "更新账户", userUpdateParams{}, userUpdate()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
		// 设置/切换公司
		info.POST("default_company", "更新账户", setDefaultCompanyParams{}, setDefaultCompany()).WithPermission(core.SOURCE_TYPE_UNPERMISSION)
	}

	// 用户列表管理, 用于批量管理用户
	users := auth.Group("users", "用户管理").WithMeta("icon", "el-icon-user")
	{
		// 获取用户列表
		users.GET("list", "用户列表", usersQueryParams{}, usersQuery())
		// 添加用户
		users.PUT("add", "添加用户", userAddParams{}, userAdd())
		// 更新用户
		users.POST("edit", "更新用户", usersUpdateParams{}, usersUpdate())
		// 删除用户, 逻辑删除,
		users.DELETE("delete", "删除用户", usersDeleteParams{}, usersDelete())
	}

	roles := auth.Group("roles", "角色管理").WithMeta("icon", "el-icon-switch-filled")
	{
		// 角色列表
		roles.GET("list", "角色列表", rolesQueryParams{}, rolesQuery())
		// 单行添加
		roles.PUT("add", "添加角色", roleAddParams{}, roleAdd(roleSeq))
		// 更新用户
		roles.POST("edit", "编辑角色", rolesUpdateParams{}, rolesUpdate())
		// 删除用户, 逻辑删除
		roles.DELETE("delete", "删除角色", rolesDeleteParams{}, rolesDelete())
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
			menus.DELETE("delete", "添加角色", menusDeleteParams{}, menusDelete())
		}
	}
}
