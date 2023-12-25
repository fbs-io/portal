/*
 * @Author: reel
 * @Date: 2023-10-28 10:29:01
 * @LastEditors: reel
 * @LastEditTime: 2023-12-24 17:56:03
 * @Description: 部门操作
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
)

type departmentAddParams struct {
	DepartmentCode        string `json:"department_code"`
	DepartmentName        string `json:"department_name"`
	DepartmentComment     string `json:"department_comment"`
	DepartmentParentCode  string `json:"department_parent_code"`
	DepartmentCustomLevel string `json:"department_custom_level"`
}

func departmentAdd(orgSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*departmentAddParams)
		if param.DepartmentCode == "" {
			param.DepartmentCode = orgSeq.Code()
		}
		if param.DepartmentCode == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无法获取或生成公司代码")))
			return
		}
		model := &Department{
			DepartmentCode:        param.DepartmentCode,
			DepartmentName:        param.DepartmentName,
			DepartmentComment:     param.DepartmentComment,
			DepartmentParentCode:  param.DepartmentParentCode,
			DepartmentCustomLevel: param.DepartmentCustomLevel,
		}
		if model.DepartmentParentCode != "" {
			tx := ctx.NewTX()
			pmodel := &Department{}
			err := tx.Table(model.TableName()).Where("department_code = ? and status = 1", model.DepartmentParentCode).First(pmodel).Error
			if err != nil || pmodel.DepartmentCode == "" {
				ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无有效的父级code")))
				return
			}
		}
		err := ctx.TX().Create(model).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type departmentEditParams struct {
	ID                    []uint `json:"id" binding:"required" conditions:"-"`
	DepartmentName        string `json:"department_name" conditions:"-"`
	DepartmentComment     string `json:"department_comment" conditions:"-"`
	DepartmentParentCode  string `json:"department_parent_code" conditions:"-"`
	DepartmentCustomLevel string `json:"department_custom_level" conditions:"-"`
	Status                int8   `json:"status" conditions:"-"`
}

func departmentEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*departmentEditParams)
		model := &Department{
			DepartmentName:        param.DepartmentName,
			DepartmentComment:     param.DepartmentComment,
			DepartmentParentCode:  param.DepartmentParentCode,
			DepartmentCustomLevel: param.DepartmentCustomLevel,
		}
		model.Status = param.Status

		if model.DepartmentParentCode != "" {
			tx := ctx.NewTX()
			pmodel := &Department{}
			err := tx.Table(model.TableName()).Where("department_code = ? and status = 1", model.DepartmentParentCode).First(pmodel).Error
			if err != nil || pmodel.DepartmentCode == "" {
				ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无有效的上级code")))
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

type departmentQueryParams struct {
	PageNum        int    `form:"page_num"`
	PageSize       int    `form:"page_size"`
	Orders         string `form:"orders"`
	DepartmentName string `form:"department_name" conditions:"like"`
}

// 列表的方式显示
func departmentList() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*departmentQueryParams)
		model := &Department{}

		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
			core.SetTxSubTable(model.TableName()),
		)
		modelList := make([]*Department, 0, 100)

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

// 树结构查询
func getDepartmentTree() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*departmentQueryParams)
		model := &Department{}

		modelList := make([]*Department, 0, 100)

		err := ctx.TX().Model(model).Offset(-1).Limit(-1).Order("id").Find(&modelList).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		modelTree, _ := genDepartmentTree(modelList)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"rows":      modelTree,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type departmentDeleteParams struct {
	ID []uint `json:"id" binding:"required" conditions:"-"`
}

// 软删除
func departmentDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*departmentDeleteParams)
		tx := ctx.TX()

		model := &Department{}

		err := tx.Model(model).Where("id in (?)", param.ID).Delete(model).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
