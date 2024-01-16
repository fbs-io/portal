/*
 * @Author: reel
 * @Date: 2023-12-21 22:20:41
 * @LastEditors: reel
 * @LastEditTime: 2024-01-16 23:09:12
 * @Description: 岗位操作, 通过岗位定位用户权限, 如果没有权限, 则只能看到自己所属的内容
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
)

func positionAdd(positionSeq sequence.Sequence) core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*positionAddParams)
		if param.PositionCode == "" {
			param.PositionCode = positionSeq.Code()
		}
		if param.PositionCode == "" {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(errorx.New("无法获取或生成岗位代码")))
			return
		}
		err := PositionSrvice.Create(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

func positionEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*positionEditParams)

		err := PositionSrvice.UpdateByID(ctx.TX(), param)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
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
