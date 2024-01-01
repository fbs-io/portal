/*
 * @Author: reel
 * @Date: 2023-09-18 19:26:52
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 20:11:29
 * @Description: 请填写简介
 */
package consts

const (
	// ↓↓↓↓↓↓↓↓↓↓基础配置↓↓↓↓↓↓↓↓↓↓↓↓
	// 组织表
	TABLE_BASIS_ORG_COMPANY    = "e_basis_org_company"
	TABLE_BASIS_ORG_DEPARTMENT = "e_basis_org_department"
	TABLE_BASIS_ORG_POSITION   = "e_basis_org_position"

	// 基础配置表
	TABLE_BASIS_AUTH_ROLE = "e_basis_auth_role"
	TABLE_BASIS_AUTH_USER = "e_basis_auth_user"
	// 关系表
	TABLE_BASIS_RLAT_USER_POSITION = "e_basis_rlat_user_position" // 岗位和用户的关系表
	// ↑↑↑↑↑↑↑↑↑↑基础配置↑↑↑↑↑↑↑↑↑↑↑↑

	// ↓↓↓↓↓↓↓↓↓↓财务管理↓↓↓↓↓↓↓↓↓↓↓↓
	// 票据相关
	// 发票
	TABLE_FIN_BILL_INVOICE_HEADER = "e_fin_bill_invoice_header"
	// ↑↑↑↑↑↑↑↑↑↑财务管理↑↑↑↑↑↑↑↑↑↑↑↑

)

var (
	TABLE_BASIS_AUTH_ROLE_TEST = "e_basis_auth_role"
)
