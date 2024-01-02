/*
 * @Author: reel
 * @Date: 2023-12-31 14:06:55
 * @LastEditors: reel
 * @LastEditTime: 2024-01-01 17:14:51
 * @Description: 财务管理
 */
package financial

import (
	"portal/app/financial/bill"

	"github.com/fbs-io/core"
)

func New(route core.RouterGroup) {
	fin := route.Group("financial", "财务").WithMeta("icon", "sc-icon-financial")

	// 票据相关
	bill.New(fin)
}
