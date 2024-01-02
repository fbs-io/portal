/*
 * @Author: reel
 * @Date: 2023-12-31 16:16:52
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 22:20:17
 * @Description: 票据管理
 */
package bill

import "github.com/fbs-io/core"

func New(route core.RouterGroup) {
	tx := route.Core().RDB()
	tx.Register(&InvoiceHeader{})

	billGroup := route.Group("bill", "票据管理").WithMeta("icon", "sc-icon-fin-bill")
	{
		invoice := billGroup.Group("invoice", "发票管理").WithMeta("icon", "sc-icon-fin-bill-invoice")
		invoice.GET("list", "发票列表", invoiceHeaderQueryParams{}, invoiceHeaderList())
		invoice.PUT("add", "新增发票", invoiceHeaderAddParams{}, invoiceHeaderAdd())
		invoice.POST("edit", "修改发票", invoiceHeaderEditParams{}, invoiceHeaderEdit())
		invoice.DELETE("delete", "删除发票", invoiceHeaderDeleteParams{}, invoiceHeaderDelete())
		// 测试用, 可以不用登陆和权限
		invoice.PUT("qradd", "二维码新增发票", invoiceHeaderAddWithQrParam{}, invoiceHeaderAddWithQr()).WithAllowSignature().WithPermission(core.SOURCE_TYPE_UNLIMITED)
	}
}
