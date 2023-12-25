/*
 * @Author: reel
 * @Date: 2023-12-21 22:20:41
 * @LastEditors: reel
 * @LastEditTime: 2023-12-25 20:34:18
 * @Description: 岗位操作, 通过岗位定位用户权限, 如果没有权限, 则只能看到自己所属的内容
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
)

type positionAddParams struct {
	PositionCode         string           `json:"position_code"`
	PositionName         string           `json:"position_name" binding:"required"`
	PositionComment      string           `json:"position_comment"`
	PositionParentCode   string           `json:"position_parent_code"`
	DepartmentCode       string           `json:"department_code" binding:"required"`
	JobCode              string           `json:"job_code"`
	IsHead               int8             `json:"is_head"`
	IsApprove            int8             `json:"is_approve"`
	IsVritual            int8             `json:"is_vritual"`
	DataPermissionType   int8             `json:"data_permission_type"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom"`
}

func positionAdd(positionSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*positionAddParams)
		if param.PositionCode == "" {
			param.PositionCode = positionSeq.Code()
		}
		if param.DepartmentCode == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无法获取或生成公司代码")))
			return
		}
		model := &Position{
			PositionCode:         param.PositionCode,
			PositionName:         param.PositionName,
			PositionComment:      param.PositionComment,
			PositionParentCode:   param.PositionParentCode,
			DepartmentCode:       param.DepartmentCode,
			JobCode:              param.JobCode,
			IsHead:               param.IsHead,
			IsApprove:            param.IsApprove,
			IsVritual:            param.IsVritual,
			DataPermissionType:   param.DataPermissionType,
			DataPermissionCustom: param.DataPermissionCustom,
		}
		// 校验上级code
		if model.PositionParentCode != "" {
			tx := ctx.NewTX()
			pmodel := &Position{}
			err := tx.Table(model.TableName()).Where("department_code = ? and status = 1", model.PositionParentCode).First(pmodel).Error
			if err != nil || pmodel.DepartmentCode == "" {
				ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无有效的上级code")))
				return
			}
		}

		// 校验部门code
		tx := ctx.NewTX()
		dmodel := &Department{}
		err := tx.Table(dmodel.TableName()).Where("department_code = ? and status = 1", model.DepartmentCode).First(dmodel).Error
		if err != nil || dmodel.DepartmentCode == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无有效的部门code")))
			return
		}

		err = ctx.TX().Create(model).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type positionEditParams struct {
	ID                   []uint           `json:"id" binding:"required" conditions:"-"`
	PositionName         string           `json:"position_name" conditions:"-"`
	PositionComment      string           `json:"position_comment" conditions:"-"`
	PositionParentCode   string           `json:"position_parent_code" conditions:"-"`
	DepartmentCode       string           `json:"department_code" conditions:"-"`
	JobCode              string           `json:"job_code" conditions:"-"`
	IsHead               int8             `json:"is_head" conditions:"-"`
	IsApprove            int8             `json:"is_approve" conditions:"-"`
	IsVritual            int8             `json:"is_vritual" conditions:"-"`
	Status               int8             `json:"status" conditions:"-"`
	DataPermissionType   int8             `json:"data_permission_type" conditions:"-"`
	DataPermissionCustom rdb.ModeListJson `json:"data_permission_custom" conditions:"-"`
}

func positionEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*positionEditParams)
		model := &Position{
			PositionName:         param.PositionName,
			PositionComment:      param.PositionComment,
			PositionParentCode:   param.PositionParentCode,
			DepartmentCode:       param.DepartmentCode,
			JobCode:              param.JobCode,
			IsHead:               param.IsHead,
			IsApprove:            param.IsApprove,
			IsVritual:            param.IsVritual,
			DataPermissionType:   param.DataPermissionType,
			DataPermissionCustom: param.DataPermissionCustom,
		}

		model.Status = param.Status

		// 校验上级code
		if model.PositionParentCode != "" {
			tx := ctx.NewTX()
			pmodel := &Position{}
			err := tx.Table(model.TableName()).Where("department_code = ? and status = 1", model.PositionParentCode).First(pmodel).Error
			if err != nil || pmodel.DepartmentCode == "" {
				ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无有效的上级code")))
				return
			}
		}

		// 校验部门code
		if model.DepartmentCode != "" {
			tx := ctx.NewTX()
			dmodel := &Department{}
			err := tx.Table(dmodel.TableName()).Where("department_code = ? and status = 1", model.DepartmentCode).First(dmodel).Error
			if err != nil || dmodel.DepartmentCode == "" {
				ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(errorx.New("无有效的部门code")))
				return
			}
		}

		err := ctx.TX().Where("id in (?) ", param.ID).Updates(model).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type positionQueryParams struct {
	PageNum      int    `form:"page_num"`
	PageSize     int    `form:"page_size"`
	Orders       string `form:"orders"`
	PositionName string `form:"position_name" conditions:"like"`
	// DepartmentCode string `form:"department_name"`
}

// 列表的方式显示
func positionList() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*positionQueryParams)
		model := &Position{}

		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
			core.SetTxSubTable(model.TableName()),
		)
		modelList := make([]*Position, 0, 100)

		err := tx.Model(model).Find(&modelList).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		ctx.TX().Model(model).Offset(-1).Limit(-1).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      modelList,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type positionDeleteParams struct {
	ID []uint `json:"id" binding:"required" conditions:"-"`
}

func positionDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*positionDeleteParams)
		tx := ctx.TX()

		model := &Position{}

		err := tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
