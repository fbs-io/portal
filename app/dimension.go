/*
 * @Author: reel
 * @Date: 2023-10-05 19:59:11
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:57:05
 * @Description: 维度信息管理, 用于管理整个app的维度信息, 由多个表组合而成
 */
package app

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
)

type dimensionParams struct {
	DimType string `form:"dim_type" binding:"required" conditions:"-"`
}

type companyResult struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	PCode string `json:"pcode"`
}

func dimList() core.HandlerFunc {
	return func(ctx core.Context) {
		params := ctx.CtxGetParams().(*dimensionParams)
		tx := ctx.NewTX()
		var (
			err    error
			result = make([]*companyResult, 0, 1000)
		)
		switch params.DimType {
		case "company":
			// result = make([]*companyResult, 0, 1000)
			err = tx.Table(consts.TABLE_BASIS_ORG_COMPANY).Select("company_code as code", "company_name as name").Find(&result).Error
		case "role":
			// result = make([]*roleResult, 0, 1000)
			err = tx.Table(consts.TABLE_BASIS_AUTH_ROLE).Select("code", "label as name").Find(&result).Error
		default:
			err = errorx.New("dim_type不合法, 请输入正确的值")
		}

		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(result))

	}
}
