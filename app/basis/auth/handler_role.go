/*
 * @Author: reel
 * @Date: 2023-08-30 07:36:25
 * @LastEditors: reel
 * @LastEditTime: 2023-11-19 20:57:52
 * @Description: 角色相关api逻辑
 */
package auth

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/store/rdb"
)

// const (
// 	DATA_PERMISSION_ONESELF       int8 = iota + 1 //本人可见
// 	DATA_PERMISSION_ALL                           //全部可见
// 	DATA_PERMISSION_ONLY_DEPT                     //所在部门可见
// 	DATA_PERMISSION_ONLY_DEPT_ALL                 //所在部门及子级可见
// 	DATA_PERMISSION_ONLY_CUSTOM                   //选择的部门可见
// )

// TODO:补充其他用户信息, 如部门等
// 设计思路: 用户和员工分开, 用户可以绑定员工, 但员工不一定有登陆账户
type roleAddParams struct {
	Label                string           `json:"label"`
	Sort                 int              `json:"sort"`
	Description          string           `json:"description"`
	Sources              rdb.ModeListJson `json:"sources"`
	DataPermissionType   int8             `json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom"`
}

func roleAdd(roleSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*roleAddParams)
		tx := ctx.TX()
		role := &Role{
			Code:               roleSeq.Code(),
			Label:              param.Label,
			Sort:               param.Sort,
			Description:        param.Description,
			Sources:            param.Sources,
			DataPermissionType: param.DataPermissionType,
		}
		if role.DataPermissionType == rdb.DATA_PERMISSION_ONLY_CUSTOM {
			role.DataPermissionCustom = param.DataPermissionCustom
		}
		err := tx.Create(role).Error
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

// orders, page_num, page_size 作为保留字段用于条件生成
type rolesQueryParams struct {
	PageNum  int    `form:"page_num"`
	PageSize int    `form:"page_size"`
	Orders   string `form:"orders"`
	Label    string `form:"label" conditions:"like"`
}

func rolesQuery() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rolesQueryParams)
		role := &Role{}
		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
		)
		roles := make([]*Role, 0, 100)

		err := tx.Model(role).Find(&roles).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		ctx.TX().Model(role).Offset(-1).Limit(-1).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      roles,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

// 批量更新参数
//
// id作为数组, 不适用于自动查询条件生成
type rolesUpdateParams struct {
	ID                   []uint           `json:"id"  binding:"required" conditions:"-"`
	Label                string           `json:"label" conditions:"-"`
	Sort                 int              `json:"json" conditions:"-"`
	Description          string           `json:"description" conditions:"-"`
	Sources              rdb.ModeListJson `json:"sources" conditions:"-"`
	Status               int8             `json:"status" conditions:"-"`
	DataPermissionType   int8             `json:"data_permission_type" conditions:"-"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom" conditions:"-"`
}

func rolesUpdate() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rolesUpdateParams)
		tx := ctx.TX()
		role := &Role{
			Label:              param.Label,
			Sort:               param.Sort,
			Description:        param.Description,
			Sources:            param.Sources,
			DataPermissionType: param.DataPermissionType,
		}
		if param.DataPermissionType == rdb.DATA_PERMISSION_ONLY_CUSTOM {
			role.DataPermissionCustom = param.DataPermissionCustom
		}
		role.Status = param.Status
		err := tx.Model(role).Where("id in (?)", param.ID).Updates(role).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type rolesDeleteParams struct {
	ID []uint `json:"id" conditions:"-"`
}

// 软删除
func rolesDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rolesDeleteParams)
		tx := ctx.TX()

		role := &Role{}

		err := tx.Model(role).Where("id in (?)", param.ID).Delete(role).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 菜单和权限查询, 返回树表结构
func menusQueryWithPermission() core.HandlerFunc {
	return func(ctx core.Context) {
		user := GetUser(ctx.Auth(), ctx, REFRESH_NOT)
		menus, _, err := getMenuTree(ctx, user, QUERY_MENU_MODE_MANAGE)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(menus))
	}
}
