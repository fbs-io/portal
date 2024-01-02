/*
 * @Author: reel
 * @Date: 2023-10-05 19:59:11
 * @LastEditors: reel
 * @LastEditTime: 2023-12-25 23:04:19
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
	DIM_TYPE_ROLE = "role" // 角色

	// 组织相关
	DIM_TYPE_COMPANY    = "company"    // 法人
	DIM_TYPE_POSITION   = "position"   // 岗位
	DIM_TYPE_DEPARTMENT = "department" // 部门
)

type dimensionParams struct {
	DimType string `form:"dim_type" binding:"required" conditions:"-"`
	DimName string `form:"dim_name" conditions:"like%"`
}

type dimResult struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	PCode  string `json:"pcode"`
	Status int8   `json:"status"`
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

// 获取维度表
//
// 用于前端的数据展示等
func queryDimList(tx *gorm.DB, dimType string) (result []*dimResult, err error) {

	result = make([]*dimResult, 0, 1000)

	switch dimType {
	case DIM_TYPE_COMPANY:
		tx = tx.Table(consts.TABLE_BASIS_ORG_COMPANY).Select("company_code as code", "company_name as name", "status")
		err = tx.Where("status > 0 ").Find(&result).Error
	case DIM_TYPE_DEPARTMENT:
		tx = tx.Table(consts.TABLE_BASIS_ORG_DEPARTMENT).Select("department_code as code", "department_name as name", "department_parent_code as pcode", "status")
		err = tx.Where("status > 0 ").Find(&result).Error
	case DIM_TYPE_ROLE:
		tx = tx.Table(consts.TABLE_BASIS_AUTH_ROLE).Select("code", "label as name", "status")
		err = tx.Where("status > 0 ").Find(&result).Error
	case DIM_TYPE_POSITION:
		tx = tx.Table(consts.TABLE_BASIS_ORG_POSITION).Select("position_code as code", "position_name as name", "position_parent_code as pcode", "status")
		err = tx.Where("status > 0 ").Find(&result).Error
	default:
		err = errorx.New("dim_type不合法, 请输入正确的值")
	}
	return
}
