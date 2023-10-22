/*
 * @Author: reel
 * @Date: 2023-10-05 19:59:11
 * @LastEditors: reel
 * @LastEditTime: 2023-10-17 19:09:12
 * @Description: 维度信息管理, 用于管理整个app的维度信息, 由多个表组合而成
 */
package app

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"gorm.io/gorm"
)

const (
	DIM_TYPE_COMPANY = "company"
	DIM_TYPE_ROLE    = "role"
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

		// TODO: 增加tx的上下文,
		tx := ctx.TX()
		result, err := queryDimList(tx, params.DimType)
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(result))

	}
}

func queryDimList(tx *gorm.DB, dimType string) (result []*companyResult, err error) {

	result = make([]*companyResult, 0, 1000)

	switch dimType {
	case DIM_TYPE_COMPANY:
		tx = tx.Table(consts.TABLE_BASIS_ORG_COMPANY).Select("company_code as code", "company_name as name")
		err = tx.Where("status > 0 ").Find(&result).Error
	case DIM_TYPE_ROLE:
		tx = tx.Table(consts.TABLE_BASIS_AUTH_ROLE).Select("code", "label as name")
		err = tx.Where("status > 0 ").Find(&result).Error
	default:
		err = errorx.New("dim_type不合法, 请输入正确的值")
	}
	return
}
