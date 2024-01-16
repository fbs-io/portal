/*
 * @Author: reel
 * @Date: 2023-07-30 22:09:24
 * @LastEditors: reel
 * @LastEditTime: 2024-01-14 20:40:53
 * @Description: 临时生成菜单用, 菜单功能主要思路: 由后端完成菜单的生成, 前端主要用于查看
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
)

const (
	QUERY_MENU_MODE_INFO       = "info"       // 用户登录时, 返回的菜单信息, 包含不受限菜单和资源
	QUERY_MENU_MODE_MANAGE     = "manage"     // 用于菜单管理返回的数据, 返回所有的菜单和资源用于前端设置, 查看
	QUERY_MENU_MODE_PERMISSION = "permission" // 用于菜单管理返回的数据, 仅返回受限的菜单和资源(api), 用于权限设置或修改
)

type menuTree struct {
	ID         uint            `json:"id"`
	Code       string          `json:"code"`
	PCode      string          `json:"pcode"`
	Level      int8            `json:"levle"`
	Name       string          `json:"name"`
	Desc       string          `json:"desc"`
	Api        string          `json:"api"`
	Type       int8            `json:"type"`
	Method     string          `json:"method"`
	Params     string          `json:"params"`
	AcceptType string          `json:"accept_type"`
	IsRouter   int8            `json:"is_router"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Meta       rdb.ModeMapJson `json:"meta"` // 前端组件原信息
	Children   []*menuTree     `json:"children"`
}

func getMenuTree(ctx core.Context, user *User, mode string) (tree []*menuTree, permissions map[string]bool, err error) {
	source := &core.Sources{}

	// 构筑条件
	cond := rdb.NewCondition()
	cond.Orders = "level, id"
	cond.TableName = source.TableName()

	permissionList := []string{""}
	// dataPermissionCtx := &rdb.DataPermissionStringCtx{
	// 	DataPermissionType:  rdb.DATA_PERMISSION_ONESELF,
	// 	DataPermission:      user.Department,
	// 	DataPermissionScope: make([]string, 0, 100),
	// }
	if user.Super != "Y" {
		// 获取用户关联的角色
		roles := make([]*Role, 0, 100)
		userRoles := []string{}
		for _, roleCode := range user.Role {
			userRoles = append(userRoles, roleCode.(string))
		}
		// 设置把分区字段放入上下文,
		// 用户登陆逻辑复杂, 本处特别处理
		err = ctx.NewTX().Where("code in (?) ", userRoles).Find(&roles).Error
		if err != nil {
			return
		}

		for _, role := range roles {
			for _, source := range role.Sources {
				permissionList = append(permissionList, source.(string))
			}
			// // 处理role上的部门
			// switch role.DataPermissionType {
			// case rdb.DATA_PERMISSION_ONLY_CUSTOM:
			// 	for _, cust := range role.DataPermissionCustom {
			// 		dataPermissionCtx.DataPermissionScope = append(dataPermissionCtx.DataPermissionScope, cust.(string))
			// 	}
			// case rdb.DATA_PERMISSION_ONLY_DEPT:
			// 	dataPermissionCtx.DataPermissionScope = append(dataPermissionCtx.DataPermissionScope, user.Department)
			// case rdb.DATA_PERMISSION_ONLY_DEPT_ALL:
			// 	depts := make([]*org.Department, 0, 100)
			// 	err = ctx.NewTX().Where("code in (?) ", userRoles).Find(&depts).Error
			// 	if err != nil {
			// 		return
			// 	}
		}
		// }
	}
	switch mode {

	// 登陆后返回的菜单信息
	case QUERY_MENU_MODE_INFO:
		cond.Where["type > (?)"] = 0
		if user.Super != "Y" {
			cond.Where["code in (?) or type in (1,3,5)"] = permissionList
		}
	// 菜单管理用的
	case QUERY_MENU_MODE_MANAGE:
		cond.Where["type > (?)"] = 0
		if user.Super != "Y" {
			cond.Where["code in (?) "] = permissionList
		}
	case QUERY_MENU_MODE_PERMISSION:
		cond.Where["type in (?)"] = []int8{core.SOURCE_TYPE_MENU, core.SOURCE_TYPE_PERMISSION}
		if user.Super != "Y" {
			cond.Where["code in (?) "] = permissionList
		}
	}

	permissions = make(map[string]bool, 100)
	// 通过角色关联权限和菜单
	menuList := make([]*core.Sources, 0, 1000)
	tx := ctx.Core().RDB().BuildQuery(cond).Offset(-1).Limit(-1)
	err = tx.Find(&menuList).Error
	if err != nil {
		return
	}
	// 构筑树表需要的变量
	menuMap := make(map[int8]map[string]*menuTree, 100)
	tree = make([]*menuTree, 0, 100)
	allowTree := make([]*menuTree, 0, 100)
	var level int8 = 0
	for i, m := range menuList {
		if i == 0 {
			level = m.Level
		}
		switch m.Type {
		case core.SOURCE_TYPE_UNLIMITED, core.SOURCE_TYPE_PERMISSION, core.SOURCE_TYPE_UNPERMISSION:
			permissions[m.Code] = true
		}
		mTree := &menuTree{
			ID:         m.ID,
			Code:       m.Code,
			PCode:      m.PCode,
			Name:       m.Name,
			Desc:       m.Desc,
			Level:      m.Level,
			Api:        m.Api,
			Type:       m.Type,
			Method:     m.Method,
			Params:     m.Params,
			AcceptType: m.AcceptType,
			IsRouter:   m.IsRouter,
			Path:       m.Path,
			Component:  m.Component,
			Meta:       m.Meta,
		}
		if menuMap[m.Level] == nil {
			menuMap[m.Level] = make(map[string]*menuTree)
		}
		menuMap[m.Level][m.Code] = mTree
		if m.Level == level && m.IsRouter == core.SOURCE_ROUTER_IS {
			tree = append(tree, mTree)
		} else {
			pt := menuMap[m.Level-1][m.PCode]
			if pt != nil {
				pt.Children = append(pt.Children, mTree)
			} else {
				if mTree.IsRouter == core.SOURCE_ROUTER_IS && mTree.Type == core.SOURCE_TYPE_UNMENU {
					allowTree = append(allowTree, mTree)
				}
			}
		}
	}
	tree = append(tree, allowTree...)

	return
}
