/*
 * @Author: reel
 * @Date: 2023-10-28 10:29:01
 * @LastEditors: reel
 * @LastEditTime: 2023-10-29 10:06:05
 * @Description: 部门模块相关处理逻辑
 */
package org

import "fmt"

type departmentTree struct {
	ID                    uint              `json:"id"`
	DepartmentCode        string            `json:"department_code"`
	DepartmentName        string            `json:"department_name"`
	DepartmentComment     string            `json:"department_comment"`
	DepartmentLevel       int8              `json:"department_level"`
	DepartmentFullPath    string            `json:"department_full_path"`
	DepartmentParentCode  string            `json:"department_parent_code"`
	DepartmentCustomLevel string            `json:"department_custom_level"`
	CreatedAT             uint              `json:"created_at"`
	CreatedBy             string            `json:"created_by"`
	UpdatedAT             uint              `json:"updated_at"`
	UpdatedBy             string            `json:"updated_by"`
	Status                int8              `json:"status"`
	Children              []*departmentTree `json:"children"`
}

// TODO: 做成通用方法
func genDepartmentTree(list []*Department) (tree []*departmentTree, err error) {
	// var level int8 = 0
	// var treeMap = make(map[int8]map[string]*departmentTree, 10)
	var treeMap = make(map[string]*departmentTree, 100)
	for _, item := range list {
		// if i == 0 {
		// 	level = item.DepartmentLevel
		// }

		itemTree := &departmentTree{
			ID:                    item.ID,
			DepartmentCode:        item.DepartmentCode,
			DepartmentName:        item.DepartmentName,
			DepartmentComment:     item.DepartmentComment,
			DepartmentLevel:       item.DepartmentLevel,
			DepartmentFullPath:    item.DepartmentFullPath,
			DepartmentParentCode:  item.DepartmentParentCode,
			DepartmentCustomLevel: item.DepartmentCustomLevel,
			CreatedAT:             item.CreatedAT,
			CreatedBy:             item.CreatedBy,
			UpdatedAT:             item.UpdatedAT,
			UpdatedBy:             item.UpdatedBy,
			Status:                item.Status,
			Children:              make([]*departmentTree, 0, 10),
		}

		// if treeMap[itemTree.DepartmentLevel] == nil {
		// 	treeMap[itemTree.DepartmentLevel] = make(map[string]*departmentTree)
		// }
		// 方便快速定位元素路径
		treeMap[itemTree.DepartmentCode] = itemTree

		// 根据层级判断, 如果时顶层, 写入list中, 如果是子级, 写入到对应的父级中
		// 如果父级不存在, 则子级也都不再录入
		if itemTree.DepartmentParentCode == "" {
			itemTree.DepartmentFullPath = itemTree.DepartmentName
			tree = append(tree, itemTree)
		} else {
			pt := treeMap[itemTree.DepartmentParentCode]
			if pt != nil {
				itemTree.DepartmentLevel = pt.DepartmentLevel + 1
				itemTree.DepartmentFullPath = fmt.Sprintf("%s-%s", pt.DepartmentFullPath, itemTree.DepartmentName)
				pt.Children = append(pt.Children, itemTree)
			}
		}
	}

	return
}
