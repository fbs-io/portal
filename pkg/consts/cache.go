/*
 * @Author: reel
 * @Date: 2023-10-06 15:40:40
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 15:56:31
 * @Description: 请填写简介
 */
package consts

import "fmt"

const (
	// 用户所属公司缓存前缀
	CACHE_USER_COMPANY_PREFIX = "user::company::%s"

	// 上下文使用的缓存
	CTX_COMPANY = "ctx_company"
)

func GenUserCompanyKey(account string) string {
	return fmt.Sprintf(CACHE_USER_COMPANY_PREFIX, account)
}
