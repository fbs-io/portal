/*
 * @Author: reel
 * @Date: 2023-10-28 10:29:01
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 20:47:42
 * @Description: 部门模块相关处理逻辑
 */
package org

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fbs-io/core"
	"gorm.io/gorm"
)

var (
	// 用于缓存各个法人整个组织
	deptAllList = make(map[string][]*Department, 100)
	// 用于缓存各个法人公司的下的部门及子部门
	deptAndAllChildrenMap = make(map[string]map[string]map[string]*departmentTree, 100)

	// 岗位list
	positionAllList = make(map[string]map[string]*Position, 100)

	// 用于map的锁
	lock = &sync.RWMutex{}
)

type departmentTree struct {
	ID                    uint              `json:"id"`
	DepartmentCode        string            `json:"department_code"`
	DepartmentName        string            `json:"department_name"`
	DepartmentComment     string            `json:"department_comment"`
	DepartmentLevel       int8              `json:"department_level"`
	DepartmentFullPath    string            `json:"department_full_path"`
	DepartmentFullPath2   string            `json:"-"` // 用于记录组织code全路径
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
// 通过list生成树结构
func GenDepartmentTree(list []*Department) (tree []*departmentTree, treeMap map[string]*departmentTree, treeMap2 map[string]map[string]*departmentTree, err error) {

	treeMap = make(map[string]*departmentTree, 100)
	treeMap2 = make(map[string]map[string]*departmentTree, 100)
	for _, item := range list {
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

		// 方便快速定位元素路径
		treeMap[itemTree.DepartmentCode] = itemTree
		if treeMap2[itemTree.DepartmentCode] == nil {
			treeMap2[itemTree.DepartmentCode] = make(map[string]*departmentTree, 100)
			treeMap2[itemTree.DepartmentCode][itemTree.DepartmentCode] = itemTree
		}
		// 根据层级判断, 如果时顶层, 写入list中, 如果是子级, 写入到对应的父级中
		// 如果父级不存在, 则子级也都不再录入
		if itemTree.DepartmentParentCode == "" {
			itemTree.DepartmentFullPath = itemTree.DepartmentName
			itemTree.DepartmentFullPath2 = itemTree.DepartmentCode
			tree = append(tree, itemTree)
		} else {
			pt := treeMap[itemTree.DepartmentParentCode]
			if pt != nil {
				itemTree.DepartmentLevel = pt.DepartmentLevel + 1
				itemTree.DepartmentFullPath = fmt.Sprintf("%s-%s", pt.DepartmentFullPath, itemTree.DepartmentName)
				itemTree.DepartmentFullPath2 = fmt.Sprintf("%s-%s", pt.DepartmentFullPath2, itemTree.DepartmentCode)
				pt.Children = append(pt.Children, itemTree)
				for _, pcode := range strings.Split(pt.DepartmentFullPath2, "-") {
					treeMap2[pcode][itemTree.DepartmentCode] = itemTree
				}
			}
		}
	}

	return
}

func SetDeptAllList(company_code string, deptList []*Department) {
	lock.Lock()
	defer lock.Unlock()
	deptAllList[company_code] = deptList
}

func GetDeptAndAllChildren(company_code, department_code string, tx *gorm.DB) (result map[string]*departmentTree) {
	var deptList []*Department
	var deptMap map[string]map[string]*departmentTree
	lock.RLock()
	defer func() {
		lock.Lock()
		if deptList != nil {
			deptAllList[company_code] = deptList
			// fmt.Println("完成部门list缓存设置")
		}
		lock.Unlock()
	}()

	defer func() {
		lock.Lock()
		if deptMap != nil {
			if deptAndAllChildrenMap[company_code] == nil {
				deptAndAllChildrenMap[company_code] = deptMap
				// fmt.Println("完成部门及子部门的map缓存设置")
			}
		}
		lock.Unlock()
	}()
	defer lock.RUnlock()
	// fmt.Println("查询开始")
	tx.Set(core.CTX_SHARDING_KEY, company_code)
	// 先判断法人公司下是否有部门list, 如果法人公司下没有组织list, 则重新生成组织list和部门map
	// 当对部门进行增删改时, 将删除 部门list缓存, 重新生成
	deptList = deptAllList[company_code]
	if deptList == nil {
		deptList = make([]*Department, 0, 100)
		// fmt.Println("没有部门list, 进行查表操作")
		// 失效部门和删除部门的数据也允许查看
		tx.Offset(-1).Limit(-1).Where("sk = ?", company_code).Order("id").Find(&deptList)
		// fmt.Println("开始生成组织数")
		_, _, deptMap, _ = GenDepartmentTree(deptList)
		if deptMap != nil {
			return deptMap[department_code]
		}
		// fmt.Println("返回部门及所有子部门数据")

	}

	if deptAndAllChildrenMap[company_code] != nil && deptAndAllChildrenMap[company_code][department_code] != nil {
		return deptAndAllChildrenMap[company_code][department_code]
	}
	_, _, deptMap, _ = GenDepartmentTree(deptList)

	if deptMap != nil {
		result = deptMap[department_code]
		return
	}

	return
}

func GetPosition(company_code, position_code string, tx gorm.DB) (position *Position) {
	lock.RLock()
	// 设置position
	defer func() {
		lock.Lock()
		if position.PositionCode != "" {
			if positionAllList[company_code] == nil {
				positionAllList[company_code] = make(map[string]*Position, 100)
			}
			positionAllList[company_code][position.PositionCode] = position
			// fmt.Println("完成岗位的map缓存设置")
		}
		lock.Unlock()
	}()
	defer lock.RUnlock()
	// fmt.Println("开始查询岗位缓存")
	tx.Set(core.CTX_SHARDING_KEY, company_code)
	if positionAllList[company_code] == nil {
		// fmt.Println("没有查询到岗位缓存,开始进行查表操作")
		position = &Position{}
		tx.Where("position_code = ? and status = 1 ", position_code).Find(position)
		// fmt.Println("完成查表操作, 返回岗位信息")
		return
	}
	// fmt.Println("有岗位缓存,开始进行岗位code查询")
	if positionAllList[company_code][position_code] != nil {
		// fmt.Println("有岗位code缓存,将直接返回")
		return positionAllList[company_code][position_code]
	}
	// fmt.Println("没有有岗位code缓存,开始进行岗位code查表操作")
	position = &Position{}
	tx.Where("position_code = ? and status = 1 and (deleted_at>0 or deleted_at is null) ", position_code).Find(position)
	// fmt.Println("完成岗位code查表操作")
	return
}
