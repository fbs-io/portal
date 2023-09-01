/*
 * @Author: reel
 * @Date: 2023-07-30 22:09:24
 * @LastEditors: reel
 * @LastEditTime: 2023-08-30 07:33:17
 * @Description: 临时生成菜单用
 */
package auth

import (
	"github.com/fbs-io/core"
	"github.com/fbs-io/core/store/rdb"
)

type menuTree struct {
	Name      string          `json:"name"`
	Path      string          `json:"path"`      // 前端用组件方法
	Component string          `json:"component"` // 前端组件名称
	Meta      rdb.ModeMapJson `json:"meta"`      // 前端组件原信息
	Children  []*menuTree     `json:"children"`
}

// TODO: 增加角色判断
func getMenuTree(c core.Core, u *User) (tree []*menuTree, err error) {
	source := &core.Sources{}

	// 构筑条件
	cond := rdb.NewCondition()
	cond.PageSize = 1000
	cond.Orders = "level, id"
	cond.TableName = source.TableName()
	cond.Where = map[string]interface{}{
		"is_router = ? ": core.SOURCE_ROUTER_IS,
	}

	menuList := make([]*core.Sources, 0, 1000)
	err = c.RDB().BuildQuery(cond).Find(&menuList).Error
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
			Name:      m.Name,
			Path:      m.Path,
			Component: m.Component,
			Meta:      m.Meta,
			Children:  make([]*menuTree, 0, 100),
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
