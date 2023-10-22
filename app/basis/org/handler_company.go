/*
 * @Author: reel
 * @Date: 2023-09-19 04:34:39
 * @LastEditors: reel
 * @LastEditTime: 2023-10-15 23:12:22
 * @Description: 公司管理, 系统按公司别进行数据隔离
 */
package org

import (
	"portal/pkg/sequence"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
)

type companyAddParams struct {
	CompanyCode      string `json:"company_code"`
	CompanyName      string `json:"company_name"`
	CompanyShortName string `json:"company_shortname"`
	CompanyComment   string `json:"company_comment"`
	CompanyBusiness  string `json:"company_business"`
}

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

		company := &Company{
			CompanyCode:      param.CompanyCode,
			CompanyName:      param.CompanyName,
			CompanyShortName: param.CompanyShortName,
			CompanyComment:   param.CompanyComment,
			CompanyBusiness:  param.CompanyBusiness,
		}
		err := ctx.TX().Create(company).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		// 创建公司相关的库或者表
		ctx.Core().RDB().AddShardingSuffixs(company.CompanyCode)
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type companyEditParams struct {
	ID               []uint `json:"id" binding:"required"`
	CompanyName      string `json:"company_name"  conditions:"-"`
	CompanyShortName string `json:"company_shortname"  conditions:"-"`
	CompanyComment   string `json:"company_comment"  conditions:"-"`
	CompanyBusiness  string `json:"company_business"  conditions:"-"`
	Status           int8   `json:"status" conditions:"-"`
}

func companyEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*companyEditParams)
		company := &Company{
			CompanyName:      param.CompanyName,
			CompanyShortName: param.CompanyShortName,
			CompanyComment:   param.CompanyComment,
			CompanyBusiness:  param.CompanyBusiness,
		}
		company.Status = param.Status

		err := ctx.TX().Where("id in (?) ", param.ID).Updates(company).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type companyQueryParams struct {
	PageNum     int    `form:"page_num"`
	PageSize    int    `form:"page_size"`
	Orders      string `form:"orders"`
	CompanyName string `form:"company_name" conditions:"like"`
}

func companyList() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*companyQueryParams)
		comoany := &Company{}

		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
			core.SetTxSubTable(comoany.TableName()),
		)
		comoanys := make([]*Company, 0, 100)

		err := tx.Model(comoany).Find(&comoanys).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		ctx.TX().Model(comoany).Offset(-1).Limit(-1).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      comoanys,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type companyDeleteParams struct {
	ID []uint `json:"id" binding:"required" conditions:"-"`
}

// 软删除
func companyDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*companyDeleteParams)
		tx := ctx.TX()

		company := &Company{}

		err := tx.Model(company).Where("id in (?)", param.ID).Delete(company).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
