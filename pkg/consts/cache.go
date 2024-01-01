/*
 * @Author: reel
 * @Date: 2023-10-06 15:40:40
 * @LastEditors: reel
 * @LastEditTime: 2023-12-30 22:48:58
 * @Description: 请填写简介
 */
package consts

import "fmt"

const (
	// 用户所属公司缓存前缀
	CACHE_USER_COMPANY_PREFIX = "user::company::%s"
	// 用户岗位缓存前缀
	CACHE_USER_POSITION_PREFIX = "user::company::position::%s-%s"

	// 上下文使用公司的缓存
	CTX_COMPANY = "ctx_company"
	// 上下文使用的岗位缓存
	CTX_POSITION = "ctx_company"
)

func GenUserCompanyKey(account string) string {
	return fmt.Sprintf(CACHE_USER_COMPANY_PREFIX, account)
}

func GenUserPositionKey(account, company string) string {
	return fmt.Sprintf(CACHE_USER_POSITION_PREFIX, account, company)
}
