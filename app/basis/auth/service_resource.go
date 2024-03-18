/*
 * @Author: reel
 * @Date: 2024-01-17 21:49:09
 * @LastEditors: reel
 * @LastEditTime: 2024-03-19 06:33:00
 * @Description: 资源逻辑处理层
 */
package auth

import (
	"sync"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	QUERY_MENU_MODE_INFO       = "info"       // 用户登录时, 返回的菜单信息, 包含不受限菜单和资源
	QUERY_MENU_MODE_MANAGE     = "manage"     // 用于菜单管理返回的数据, 返回所有的菜单和资源用于前端设置, 查看
	QUERY_MENU_MODE_PERMISSION = "permission" // 用于菜单管理返回的数据, 仅返回受限的菜单和资源(api), 用于权限设置或修改
)

type resourceService struct {
	lock        *sync.RWMutex
	core        core.Core
	list        []*core.Sources // 存储所有资源列表
	listCode    []string        // 存储所有资源code, 用于和传入的权限值进行配合使用
	unLimitCode []string        // 存储不受限的code, 用于和传入的权限值进行配合使用
	idMap       map[uint]*core.Sources
	codeMap     map[string]*core.Sources
}

var ResourceService *resourceService

// 初始化
//
// 根据不同分区缓存
func ResourceServiceInit(c core.Core) {
	ResourceService = &resourceService{
		core:        c,
		lock:        &sync.RWMutex{},
		idMap:       make(map[uint]*core.Sources, 100),
		codeMap:     make(map[string]*core.Sources, 100),
		listCode:    make([]string, 0, 100),
		list:        make([]*core.Sources, 0, 100),
		unLimitCode: make([]string, 0, 100),
	}
	var list = make([]*core.Sources, 0, 100)
	tx := c.RDB().DB().Where("1=1")
	tx.Order("level, id").Find(&list)

	for _, item := range list {
		ResourceService.setCache(item)
	}

}

// 创建缓存
func (srv *resourceService) setCache(item *core.Sources) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	oItem := srv.codeMap[item.Code]

	// 如果无法查询到元素 则插入,否则更新
	if oItem == nil {
		srv.codeMap[item.Code] = item
		srv.idMap[item.ID] = item
		srv.list = append(srv.list, item)

		// 只获取可用资源
		if item.Type > 0 {
			srv.listCode = append(srv.listCode, item.Code)
			if item.Type == core.SOURCE_TYPE_UNMENU ||
				item.Type == core.SOURCE_TYPE_UNLIMITED ||
				item.Type == core.SOURCE_TYPE_UNPERMISSION {
				srv.unLimitCode = append(srv.unLimitCode, item.Code)
			}
		}
	} else {
		oItem.Desc = item.Desc
		oItem.Api = item.Api
		oItem.Type = item.Type
		oItem.Method = item.Method
		oItem.Params = item.Params
		oItem.AcceptType = item.AcceptType
		oItem.IsRouter = item.IsRouter
		oItem.Path = item.Path
		oItem.Component = item.Component
		oItem.Meta = item.Meta
		oItem.Status = item.Status
	}

}

// 通过code批量获取有效的资源
func (srv *resourceService) GetByCodes(codes []string) (list []*core.Sources) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()
	list = make([]*core.Sources, 0, 10)
	for _, code := range codes {
		item := srv.codeMap[code]
		if item != nil && item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 通过code获取有效的资源
func (srv *resourceService) GetByCode(code string) (itme *core.Sources) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	if srv.codeMap[code] != nil && srv.codeMap[code].Status == 1 {
		return srv.codeMap[code]
	}
	return
}

// 获取所有有效的资源列表
func (srv *resourceService) GetAllList() (list []*core.Sources) {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	list = make([]*core.Sources, 0, 100)
	for _, item := range srv.list {
		if srv.codeMap[item.Code] != nil && item.Status == 1 {
			list = append(list, item)
		}
	}
	return
}

// 按照id批量删除缓存
func (srv *resourceService) deleteByID(ids []uint) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	for _, id := range ids {
		item := srv.idMap[id]
		if item != nil {
			delete(srv.codeMap, item.Code)
			delete(srv.idMap, item.ID)
		}
	}
}

// 以下时对外服务的API

// 创建资源/菜单
//
// 以公司为分区, 不同公司的数据进行分开
//
// 同时创建缓存
type menusAddParams struct {
	Desc       string          `json:"desc" binding:"required"`
	PCode      string          `json:"pcode"`
	Level      int8            `json:"level"`
	Api        string          `json:"api"`
	Type       int8            `json:"type" binding:"required,gte=6,lte=9"`
	Method     string          `json:"method"`
	Params     string          `json:"params"`
	AcceptType string          `json:"accept_type"`
	IsRouter   int8            `json:"is_router" binding:"required"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Meta       rdb.ModeMapJson `json:"meta"` // 前端组件原信息
}

func (srv *resourceService) Create(tx *gorm.DB, param *menusAddParams) (err error) {
	model := &core.Sources{}
	model.Code = uuid.New().String()
	model.PCode = param.PCode
	model.Level = param.Level
	model.Desc = param.Desc
	model.Api = param.Api
	model.Type = param.Type
	model.Method = param.Method
	model.Params = param.Params
	model.AcceptType = param.AcceptType
	model.IsRouter = param.IsRouter
	model.Path = param.Path
	model.Component = param.Component
	model.Meta = param.Meta

	if model.PCode != "" {

		if srv.GetByCode(model.PCode) == nil {
			return errorx.New("无有效的父级code")
		}
		model.Level = model.Level + 1

	}

	err = tx.Create(model).Error
	if err != nil {
		return
	}
	srv.setCache(model)

	return
}

// 删除部门, 软删除
//
// 同时会删除部门缓存
func (srv *resourceService) Delete(tx *gorm.DB, param *rdb.DeleteParams) (err error) {
	model := &core.Sources{}

	err = tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
	if err != nil {
		return
	}

	srv.deleteByID(param.ID)

	return
}

// 批量更新参数
//
// id作为数组, 不适用于自动查询条件生成
type menusUpdateParams struct {
	ID         uint            `json:"id"  binding:"required"`
	Desc       string          `json:"desc" conditions:"-"`
	PCode      string          `json:"pcode" conditions:"-"`
	Level      int8            `json:"levle" conditions:"-"`
	Api        string          `json:"api" conditions:"-"`
	Type       int8            `json:"type" conditions:"-"`
	Method     string          `json:"method" conditions:"-"`
	Params     string          `json:"params" conditions:"-"`
	AcceptType string          `json:"accept_type" conditions:"-"`
	IsRouter   int8            `json:"is_router" conditions:"-"`
	Path       string          `json:"path" conditions:"-"`
	Component  string          `json:"component" conditions:"-"`
	Status     int8            `json:"status" conditions:"-"`
	Meta       rdb.ModeMapJson `json:"meta" conditions:"-"`
}

// 更新资源, 通过id批量更新
//
// 同时会更新缓存
func (srv *resourceService) Update(tx *gorm.DB, param *menusUpdateParams) (err error) {
	model := &core.Sources{}

	model.Desc = param.Desc
	model.Api = param.Api
	model.Type = param.Type
	model.Method = param.Method
	model.Params = param.Params
	model.AcceptType = param.AcceptType
	model.IsRouter = param.IsRouter
	model.Path = param.Path
	model.Component = param.Component
	model.Meta = param.Meta
	model.Status = param.Status
	// list := make([]*core.Sources, 0, 100)
	err = tx.Where("id = (?) ", param.ID).Updates(model).Find(model).Error
	if err != nil {
		return err
	}
	srv.setCache(model)

	return
}

// 通过记录行,生成菜单树和权限表以及授权表
//
// code为空时, 只给不受限内容
//
// mode控制生成的类型: 权限表, 角色用授权表, 菜单等
//
// 考虑到排序, 查询等情况, 菜单,权限等依旧通过直接查表获取,而不通过缓存构建
// func (srv *resourceService) GetSource(tx *gorm.DB, permissionList []string, mode string) (treeList []*core.Sources, permissions map[string]bool, err error) {
// 	// func getMenuTree(ctx core.Context, user *User, mode string) (tree []*core.Sources, permissions map[string]bool, err error) {
// 	source := &core.Sources{}

// 	// 构筑条件
// 	cond := rdb.NewCondition()
// 	cond.Orders = "level, id"
// 	cond.TableName = source.TableName()

// 	tx = tx.Where("type > (?)", 0)
// 	switch mode {
// 	// 登陆后返回的菜单信息
// 	case QUERY_MENU_MODE_INFO:
// 		tx = tx.Where("( code in (?) or type in (1,3,5) )", permissionList)

// 	// 菜单管理用的
// 	case QUERY_MENU_MODE_MANAGE:
// 		tx = tx.Where("code in (?)", permissionList)

// 	// 授权用
// 	case QUERY_MENU_MODE_PERMISSION:
// 		tx = tx.Where("type in (?)", []int8{core.SOURCE_TYPE_MENU, core.SOURCE_TYPE_PERMISSION})
// 		tx = tx.Where("code in (?)", permissionList)
// 	}

// 	permissions = make(map[string]bool, 100)

// 	menuList := make([]*core.Sources, 0, 1000)
// 	err = tx.Find(&menuList).Error
// 	if err != nil {
// 		return
// 	}
// 	// 构筑树表需要的变量
// 	menuMap := make(map[int8]map[string]*core.Sources, 100)
// 	treeList = make([]*core.Sources, 0, 100)
// 	allowTree := make([]*core.Sources, 0, 100)
// 	var level int8 = 0
// 	for i, item := range menuList {
// 		if i == 0 {
// 			level = item.Level
// 		}
// 		switch item.Type {
// 		case core.SOURCE_TYPE_UNLIMITED, core.SOURCE_TYPE_PERMISSION, core.SOURCE_TYPE_UNPERMISSION:
// 			permissions[item.Code] = true
// 		}
// 		itemTree := &core.Sources{
// 			ID:         item.ID,
// 			Code:       item.Code,
// 			PCode:      item.PCode,
// 			Name:       item.Name,
// 			Desc:       item.Desc,
// 			Level:      item.Level,
// 			Api:        item.Api,
// 			Type:       item.Type,
// 			Method:     item.Method,
// 			Params:     item.Params,
// 			AcceptType: item.AcceptType,
// 			IsRouter:   item.IsRouter,
// 			Path:       item.Path,
// 			Component:  item.Component,
// 			Meta:       item.Meta,
// 		}
// 		if menuMap[item.Level] == nil {
// 			menuMap[item.Level] = make(map[string]*core.Sources)
// 		}
// 		menuMap[item.Level][item.Code] = itemTree
// 		if item.Level == level && item.IsRouter == core.SOURCE_ROUTER_IS {
// 			treeList = append(treeList, itemTree)
// 		} else {
// 			pt := menuMap[item.Level-1][item.PCode]
// 			if pt != nil {
// 				pt.Children = append(pt.Children, itemTree)
// 			} else {
// 				if itemTree.IsRouter == core.SOURCE_ROUTER_IS && item.Type == core.SOURCE_TYPE_UNMENU {
// 					allowTree = append(allowTree, itemTree)
// 				}
// 			}
// 		}
// 	}
// 	treeList = append(treeList, allowTree...)

// 	return
// }

// 通过传入的codeMap获取对应的接口, 菜单, 管理菜单列表
//
// codeMap用于判断code是否存在
func (srv *resourceService) GetSource(codeMap map[string]bool) (menuList, manageList []*core.Sources, permissions map[string]bool, err error) {
	// 权限表
	permissions = make(map[string]bool, 100)
	menuTree := newTree()
	manageTree := newTree()

	for _, item := range srv.list {

		// 类型为0的数据为受限资源, 前端不可访问
		// 失效资源也被过滤
		// 过滤被删除的或不存在的数据数据
		if item == nil || item.Type == 0 || item.Status < 1 {
			continue
		}

		// 权限表单独构建
		switch item.Type {
		// 权限表构成, 取非受限和已有的权限
		case core.SOURCE_TYPE_UNLIMITED, core.SOURCE_TYPE_UNPERMISSION:
			permissions[item.Code] = true

		// 权限表构成, 取非受限和已有的权限
		case core.SOURCE_TYPE_PERMISSION:
			if codeMap != nil && codeMap[item.Code] {
				permissions[item.Code] = true
			}
		}
		// 前端页面菜单
		if (codeMap != nil && codeMap[item.Code] && item.IsRouter == core.SOURCE_ROUTER_IS) ||
			item.Type == core.SOURCE_TYPE_UNLIMITED ||
			item.Type == core.SOURCE_TYPE_UNMENU {
			menuTree.genMenuTree(item)
		}
		// 授权和管理菜单
		if codeMap != nil && codeMap[item.Code] && (item.Type == core.SOURCE_TYPE_MENU || item.Type == core.SOURCE_TYPE_PERMISSION) {
			manageTree.genMenuTree(item)
		}

	}
	menuList = append(menuTree.treeList, menuTree.allowList...)
	manageList = append(manageTree.treeList, manageTree.allowList...)

	return
}

// 用于构建树结构的菜单或数据
type tree struct {
	treeMap   map[int8]map[string]*core.Sources
	treeLevel int8
	treeList  []*core.Sources
	allowList []*core.Sources // 例外的菜单
}

func newTree() *tree {
	return &tree{
		treeMap:   map[int8]map[string]*core.Sources{},
		treeLevel: -1,
		treeList:  make([]*core.Sources, 0, 100),
		allowList: make([]*core.Sources, 0, 10),
	}
}

// 根据资源生成菜单
func (t *tree) genMenuTree(source *core.Sources) {
	item := &core.Sources{}

	item.Code = source.Code
	item.Name = source.Name
	item.Desc = source.Desc
	item.PCode = source.PCode
	item.Level = source.Level
	item.Api = source.Api
	item.Type = source.Type
	item.Sort = source.Sort
	item.Method = source.Method
	item.Params = source.Params
	item.AcceptType = source.AcceptType
	item.IsRouter = source.IsRouter
	item.Path = source.Path
	item.Component = source.Component
	item.Meta = source.Meta
	item.Children = make([]*core.Sources, 0, 100)
	item.ID = source.ID
	item.CreatedAT = source.CreatedAT
	item.CreatedBy = source.CreatedBy
	item.UpdatedAT = source.UpdatedAT
	item.UpdatedBy = source.UpdatedBy
	item.DeletedAT = source.DeletedAT
	item.DeletedBy = source.DeletedBy
	item.Status = source.Status
	if t.treeLevel == -1 {
		t.treeLevel = item.Level
	}

	if t.treeMap[item.Level] == nil {
		t.treeMap[item.Level] = make(map[string]*core.Sources, 100)
	}
	t.treeMap[item.Level][item.Code] = item
	if item.Level == t.treeLevel && item.IsRouter == core.SOURCE_ROUTER_IS {
		t.treeList = append(t.treeList, item)
	} else {
		pt := t.treeMap[item.Level-1][item.PCode]
		if pt != nil {
			pt.Children = append(pt.Children, item)
		} else {
			if item.IsRouter == core.SOURCE_ROUTER_IS && item.Type == core.SOURCE_TYPE_UNMENU {
				t.allowList = append(t.allowList, item)
			}
		}
	}

}
