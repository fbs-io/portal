/*
 * @Author: reel
 * @Date: 2023-12-31 14:09:18
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 21:25:06
 * @Description: 发票管理
 */
package bill

import (
	"portal/pkg/consts"

	"github.com/fbs-io/core/store/rdb"
	"github.com/shopspring/decimal"
)

type InvoiceHeader struct {
	InvoiceNo            string          `gorm:"comment:发票号码;unique" json:"invoice_no"`
	InvoiceType          string          `gorm:"comment:发票类型" json:"invoice_type"`
	InvoiceDate          string          `gorm:"comment:开票时间" json:"invoice_date"`
	InvoiceOrgCode       string          `gorm:"comment:开票机构代码;index" json:"invoice_org_code"`
	InvoiceVerifyCode    string          `gorm:"comment:发票校验码" json:"invoice_verify_code"`
	InvoiceEncryptCode   string          `gorm:"comment:发票加密字符" json:"invoice_encrypt_code"`
	InvoiceComment       string          `gorm:"comment:发票说明" json:"invoice_comment"`
	InvoicePurchaserName string          `gorm:"comment:购买名称" json:"invoice_purchas_name"`
	InvoicePurchaserCode string          `gorm:"comment:购方税号" json:"invoice_purchas_code"`
	InvoiceNum           decimal.Decimal `gorm:"comment:开票数量;type:decimal(18,4)" json:"invoice_num"`
	InvoiceAmount        decimal.Decimal `gorm:"comment:开票金额(不含税);type:decimal(18,4)" json:"invoice_amount"`
	InvoiceTaxes         decimal.Decimal `gorm:"comment:开票税金;type:decimal(18,4)" json:"invoice_taxes"`
	InvoiceAllAmount     decimal.Decimal `gorm:"comment:价税合计;type:decimal(18,4)" json:"invoice_all_amount"`
	rdb.Model
	rdb.ShardingModel
	rdb.DataPermissionStringModel
	// TODO: 维护主数据信息
	// Invoice                 string          `gorm:"comment:发票加密字符"`
	// InvoicePurchaserTelAddr string          `gorm:"comment:发票购买方地址电话"`
	// InvoicePurchaserTelAddr string          `gorm:"comment:发票地址电话"`
}

func (i *InvoiceHeader) TableName() string {
	return consts.TABLE_FIN_BILL_INVOICE_HEADER
}
