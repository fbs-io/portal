/*
 * @Author: reel
 * @Date: 2023-09-19 04:34:39
 * @LastEditors: reel
 * @LastEditTime: 2024-01-17 22:28:44
 * @Description: 公司管理, 系统按公司别进行数据隔离
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/fbs-io/core/store/rdb"
)

func companyAdd(orgSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*companyAddParams)
		if param.CompanyCode == "" {
			param.CompanyCode = orgSeq.Code()
		}
		if param.CompanyCode == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无法获取或生成公司代码")))
			return
		}

		err := CompanyService.Create(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

func companyEdit() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*companyEditParams)

		err := CompanyService.UpdateByID(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_UPDATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

func companyList() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*companyQueryParams)

		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
		)

		data, err := CompanyService.Query(tx, param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}

		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

// 软删除
func companyDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*rdb.DeleteParams)

		err := CompanyService.Delete(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}

		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
