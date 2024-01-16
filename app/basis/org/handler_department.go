/*
 * @Author: reel
 * @Date: 2023-10-28 10:29:01
 * @LastEditors: reel
 * @LastEditTime: 2024-01-16 23:45:10
 * @Description: 部门操作
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
)

// 增加部门
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
		err := DepartmentSrvice.Create(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 编辑部门
func departmentEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*departmentEditParams)

		err := DepartmentSrvice.UpdateByID(ctx.TX(), param)

		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

// 树结构查询用于数据展示, 不需要过滤等
func getDepartmentTree() core.HandlerFunc {
	return func(ctx core.Context) {

		modelTree := DepartmentSrvice.TreeList(ctx.TX())
		data := map[string]interface{}{
			"rows": modelTree,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

// 软删除
func departmentDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rdb.DeleteParams)
		tx := ctx.TX()

		err := DepartmentSrvice.Delete(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
