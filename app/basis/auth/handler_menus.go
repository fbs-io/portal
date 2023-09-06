/*
 * @Author: reel
 * @Date: 2023-09-02 09:18:15
 * @LastEditors: reel
 * @LastEditTime: 2023-09-05 05:34:24
 * @Description: 请填写简介
 */
/*
 * @Author: reel
 * @Date: 2023-08-30 07:36:25
 * @LastEditors: reel
 * @LastEditTime: 2023-09-01 06:24:27
 * @Description: 菜单相关api逻辑
 */
package auth

import (
	"time"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

// TODO:补充其他用户信息, 如部门等
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type menusAddParams struct {
	Code       string          `json:"code" binding:"required"`
	Desc       string          `json:"desc"`
	PCode      string          `json:"pcode" binding:"required"`
	Level      int8            `json:"level" `
	Api        string          `json:"api"`
	Type       int8            `json:"type"`
	Method     string          `json:"method"`
	Params     string          `json:"params"`
	AcceptType string          `json:"accept_type"`
	IsRouter   int8            `json:"is_router"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Meta       rdb.ModeMapJson `json:"meta"` // 前端组件原信息
}

func menusAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*menusAddParams)
		tx := ctx.TX()

		menu := &core.Sources{}
		menu.Code = param.Code
		menu.PCode = param.PCode
		menu.Level = param.Level
		menu.Desc = param.Desc
		menu.Api = param.Api
		menu.Type = param.Type
		menu.Method = param.Method
		menu.Params = param.Params
		menu.AcceptType = param.AcceptType
		menu.IsRouter = param.IsRouter
		menu.Path = param.Path
		menu.Component = param.Component
		menu.Meta = param.Meta

		err := tx.Create(menu).Error
		if err != nil {
			if rdb.IsUniqueError(err) {
				ctx.JSON(errno.ERRNO_RDB_DUPLICATED_KEY.WrapError(err))
				return
			}
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())

	}
}

// 菜单查询, 返回树表结构
func menusQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		menus, _, err := getMenuTree(ctx.Core(), ctx.Auth(), "manage")
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(menus))
	}
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
	Meta       rdb.ModeMapJson `json:"meta" conditions:"-"`
}

func menusUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*menusUpdateParams)
		tx := ctx.TX()
		menu := &core.Sources{}
		menu.Desc = param.Desc
		menu.Api = param.Api
		menu.Type = param.Type
		menu.Method = param.Method
		menu.Params = param.Params
		menu.AcceptType = param.AcceptType
		menu.IsRouter = param.IsRouter
		menu.Path = param.Path
		menu.Component = param.Component
		menu.Meta = param.Meta
		// role.Status = param.Status
		err := tx.Model(menu).Updates(menu).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type menusDeleteParams struct {
	ID []uint `json:"id" conditions:"-"`
}

// 软删除
func menusDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rolesDeleteParams)
		tx := ctx.TX()

		menus := &core.Sources{}
		menus.DeletedBy = ctx.Auth()
		menus.DeletedAT = uint(time.Now().Unix())

		err := tx.Model(menus).Where("id in (?)", param.ID).Updates(menus).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
