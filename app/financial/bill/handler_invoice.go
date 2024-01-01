/*
 * @Author: reel
 * @Date: 2023-12-31 16:41:30
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 21:35:53
 * @Description: 发票管理
 */

package bill

import (
	"strings"

	"github.com/fbs-io/core"
	"github.com/fbs-io/core/pkg/errno"
	"github.com/fbs-io/core/pkg/errorx"
	"github.com/shopspring/decimal"
)

type invoiceHeaderAddParams struct {
	InvoiceNo            string `json:"invoice_no" binding:"numeric,required"`           //发票号码
	InvoiceType          string `json:"invoice_type" binding:"numeric,required"`         //发票类型
	InvoiceDate          string `json:"invoice_date" binding:"numeric,required"`         //开票时间
	InvoiceOrgCode       string `json:"invoice_org_code" binding:"numeric,required"`     //开票机构代码
	InvoiceVerifyCode    string `json:"invoice_verify_code" binding:"omitempty,numeric"` //发票校验码
	InvoiceEncryptCode   string `json:"invoice_encrypt_code"`                            //发票加密字符
	InvoiceComment       string `json:"invoice_comment"`                                 //发票描述
	InvoicePurchaserName string `json:"invoice_purchas_name"`                            //发票购买方姓名
	InvoicePurchaserCode string `json:"invoice_purchas_code"`                            //发票购买方代码
	InvoiceNum           string `json:"invoice_num" binding:"omitempty,numeric"`         //开票数量
	InvoiceAmount        string `json:"invoice_amount" binding:"numeric,required"`       //开票金额
	InvoiceTaxes         string `json:"invoice_taxes" binding:"omitempty,numeric"`       //开票税金
	InvoiceAllAmount     string `json:"invoice_all_amount" binding:"omitempty,numeric"`  //价税合计
}

func invoiceHeaderAdd() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*invoiceHeaderAddParams)

		invoice := &InvoiceHeader{
			InvoiceNo:            param.InvoiceNo,
			InvoiceType:          param.InvoiceType,
			InvoiceDate:          param.InvoiceDate,
			InvoiceOrgCode:       param.InvoiceOrgCode,
			InvoiceVerifyCode:    param.InvoiceVerifyCode,
			InvoiceComment:       param.InvoiceComment,
			InvoicePurchaserName: param.InvoicePurchaserName,
			InvoicePurchaserCode: param.InvoicePurchaserCode,
		}
		invoice.InvoiceNum, _ = decimal.NewFromString(param.InvoiceNum)
		invoice.InvoiceAmount, _ = decimal.NewFromString(param.InvoiceAmount)
		invoice.InvoiceTaxes, _ = decimal.NewFromString(param.InvoiceTaxes)
		invoice.InvoiceAllAmount, _ = decimal.NewFromString(param.InvoiceAllAmount)
		err := ctx.TX().Create(invoice).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type invoiceHeaderAddWithQrParam struct {
	Info string `json:"info"` // 发票信息
	Auth string `json:"auth"` // 用户
	SK   string `json:"sk"`   // 数据分析
	DP   string `json:"dp"`   // 数据权限
}

// 通过扫描二维码录入信息
// 测试用,不经过用户认证
func invoiceHeaderAddWithQr() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*invoiceHeaderAddWithQrParam)
		infors := strings.Split(param.Info, ",")
		if len(infors) != 8 {
			ctx.JSON(errno.ERRNO_PARAMS_INVALID.WrapError(errorx.New("发票二维码不正确")))
			return
		}
		invoice := &InvoiceHeader{
			InvoiceNo:          infors[3],
			InvoiceType:        infors[1],
			InvoiceDate:        infors[5],
			InvoiceOrgCode:     infors[2],
			InvoiceVerifyCode:  infors[6],
			InvoiceEncryptCode: infors[7],
		}
		invoice.InvoiceAmount, _ = decimal.NewFromString(infors[4])

		err := ctx.TX().Create(invoice).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type invoiceHeaderEditParams struct {
	ID                   []uint `json:"id" binding:"required"`
	InvoiceNo            string `json:"invoice_no" binding:"numeric,required"`       //发票号码
	InvoiceType          string `json:"invoice_type" binding:"numeric,required"`     //发票类型
	InvoiceDate          string `json:"invoice_date" binding:"numeric,required"`     //开票时间
	InvoiceOrgCode       string `json:"invoice_orgc_ode" binding:"numeric,required"` //开票机构代码
	InvoiceVerifyCode    string `json:"invoice_verify_code" binding:"numeric"`       //发票校验码
	InvoiceEncryptCode   string `json:"invoice_encrypt_code"`                        //发票加密字符
	InvoiceComment       string `json:"invoice_comment"`                             //发票描述
	InvoicePurchaserName string `json:"invoice_purchas_name"`                        //发票购买方姓名
	InvoicePurchaserCode string `json:"invoice_purchas_code"`                        //发票购买方代码
	InvoiceNum           string `json:"invoice_num" binding:"numeric"`               //开票数量
	InvoiceAmount        string `json:"invoice_amount" binding:"numeric,required"`   //开票金额
	InvoiceTaxes         string `json:"invoice_taxes" binding:"numeric"`             //开票税金
	InvoiceAllAmount     string `json:"invoice_all_amount" binding:"numeric"`        //价税合计
	Status               int8   `json:"status"`
}

func invoiceHeaderEdit() core.HandlerFunc {
	return func(ctx core.Context) {

		param := ctx.CtxGetParams().(*invoiceHeaderEditParams)
		invoice := &InvoiceHeader{
			InvoiceNo:            param.InvoiceNo,
			InvoiceType:          param.InvoiceType,
			InvoiceDate:          param.InvoiceDate,
			InvoiceOrgCode:       param.InvoiceOrgCode,
			InvoiceVerifyCode:    param.InvoiceVerifyCode,
			InvoiceComment:       param.InvoiceComment,
			InvoicePurchaserName: param.InvoicePurchaserName,
			InvoicePurchaserCode: param.InvoicePurchaserCode,
		}
		invoice.Status = param.Status
		invoice.InvoiceNum, _ = decimal.NewFromString(param.InvoiceNum)
		invoice.InvoiceAmount, _ = decimal.NewFromString(param.InvoiceAmount)
		invoice.InvoiceTaxes, _ = decimal.NewFromString(param.InvoiceTaxes)
		invoice.InvoiceAllAmount, _ = decimal.NewFromString(param.InvoiceAllAmount)

		err := ctx.TX().Where("id in (?) ", param.ID).Updates(invoice).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_CREATE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}

type invoiceHeaderQueryParams struct {
	PageNum   int    `form:"page_num"`
	PageSize  int    `form:"page_size"`
	Orders    string `form:"orders"`
	InvoiceNo string `form:"invoice_no" conditions:"like"`
}

func invoiceHeaderList() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*invoiceHeaderQueryParams)
		invoice := &InvoiceHeader{}

		// 使用子查询, 优化分页查询
		tx := ctx.TX(
			core.SetTxMode(core.TX_QRY_MODE_SUBID),
			core.SetTxSubTable(invoice.TableName()),
		)
		invoices := make([]*InvoiceHeader, 0, 100)

		err := tx.Model(invoice).Find(&invoices).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_QUERY.WrapError(err))
			return
		}
		var count int64
		ctx.TX().Model(invoice).Offset(-1).Limit(-1).Count(&count)
		data := map[string]interface{}{
			"page_num":  param.PageNum,
			"page_size": param.PageSize,
			"total":     count,
			"rows":      invoices,
		}
		ctx.JSON(errno.ERRNO_OK.WrapData(data))
	}
}

type invoiceHeaderDeleteParams struct {
	ID []uint `json:"id" binding:"required" conditions:"-"`
}

// 软删除
func invoiceHeaderDelete() core.HandlerFunc {
	return func(ctx core.Context) {
		param := ctx.CtxGetParams().(*invoiceHeaderDeleteParams)
		tx := ctx.TX()

		invoiceHeader := &InvoiceHeader{}

		err := tx.Model(invoiceHeader).Where("id in (?)", param.ID).Delete(invoiceHeader).Error
		if err != nil {
			ctx.JSON(errno.ERRNO_RDB_DELETE.WrapError(err))
			return
		}
		ctx.JSON(errno.ERRNO_OK.Notify())
	}
}
