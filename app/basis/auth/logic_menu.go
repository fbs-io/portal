/*
 * @Author: reel
 * @Date: 2023-07-30 22:09:24
 * @LastEditors: reel
 * @LastEditTime: 2023-09-02 17:13:07
 * @Description: 临时生成菜单用
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
)

const (
	QUERY_MENU_MODE_INFO   = "info"   // 用户登录时, 返回的菜单信息, 包含不受限菜单和有权限的菜单
	QUERY_MENU_MODE_MANAGE = "manage" // 用于菜单管理返回的数据, 包含受限的菜单和资源(api), 不受限的菜单和资源并不需要显示
)

type menuTree struct {
	ID         uint            `json:"id"`
	Code       string          `json:"code"`
	PCode      string          `json:"pcode"`
	Level      int8            `json:"levle"`
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
	// *core.SourcesBase
}

// TODO: 增加角色判断
func getMenuTree(c core.Core, u *User, mode string) (tree []*menuTree, err error) {
	source := &core.Sources{}

	// 构筑条件
	cond := rdb.NewCondition()
	cond.PageSize = 1000
	cond.Orders = "level, id"
	cond.TableName = source.TableName()

	switch mode {

	// 登陆后返回的菜单信息
	case "info":
		cond.Where = map[string]interface{}{
			"is_router = ? ": core.SOURCE_ROUTER_IS,
		}
	// 菜单管理用的
	case "manage":
		cond.Where = map[string]interface{}{
			"type in (?)": []int{2, 4},
		}
	}

	menuList := make([]*core.Sources, 0, 1000)
	err = c.RDB().BuildQuery(cond).Offset(-1).Limit(-1).Find(&menuList).Error
	if err != nil {
		return
	}

	// 构筑树表需要的变量
	menuMap := make(map[int8]map[string]*menuTree, 100)
	tree = make([]*menuTree, 0, 100)

	var level int8 = 0
	for i, m := range menuList {
		if i == 0 {
			level = m.Level
		}
		mTree := &menuTree{
			ID:         m.ID,
			Code:       m.Code,
			PCode:      m.PCode,
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
			Children:   make([]*menuTree, 0, 100),
		}
		if menuMap[m.Level] == nil {
			menuMap[m.Level] = make(map[string]*menuTree)
		}
		menuMap[m.Level][m.Code] = mTree
		if m.Level == level {
			tree = append(tree, mTree)
		} else {
			pt := menuMap[m.Level-1][m.PCode]
			if pt != nil {
				pt.Children = append(pt.Children, mTree)
			}
		}
	}
	return
}
