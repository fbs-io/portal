
/*
 * @Author: reel
 * @Date: 2023-06-24 08:35:55
 * @LastEditors: reel
 * @LastEditTime: 2023-09-29 21:57:53
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"

export default {

	// 公司批量操作
	company:{
		url: `${config.API_URL}/basis/org/company/`,
		name: "用户列表操作",
		// 查询
		list: async function(data={}){
			return await http.get(`${config.API_URL}/basis/org/company/list`, data);
		},
		// // 新增
		add: async function(data={}){
			return await http.put(`${config.API_URL}/basis/org/company/add`, data);
		},
		// // 更新
		edit: async function(data={}){
			return await http.post(`${config.API_URL}/basis/org/company/edit`, data);
		},
		// // 删除
		delete: async function(data={}){
			return await http.delete(`${config.API_URL}/basis/org/company/delete`, data);
		},

	},

	// // 角色操作
	// roles:{
	// 	url: `${config.API_URL}/basis/roles`,
	// 	name: "角色列表操作",

	// 	// 单行新增
	// 	add: async function(data={}){
	// 		return await http.put(`${config.API_URL}/basis/auth/roles/add`, data);
	// 	},
		
	// 	// 查询列表
	// 	list: async function(data={}){
	// 		return await http.get(`${config.API_URL}/basis/auth/roles/list`, data);
	// 	},
		
	// 	// 批量更新, 
	// 	// 但根据业务, 不允许赋予多个角色相同的权限
	// 	// 但允许批量停用, 删除等
	// 	edit: async function(data={}){
	// 		return await http.post(`${config.API_URL}/basis/auth/roles/edit`, data);
	// 	},

	// 	// 批量删除
	// 	delete: async function(data={}){
	// 		return await http.delete(`${config.API_URL}/basis/auth/roles/delete`, data);
	// 	},

	// 	// 获取权限菜单树表
	// 	permission: async function(data={}){
	// 		return await http.get(`${config.API_URL}/basis/roles/auth/permission`, data);
	// 	},
	// },

	// // 角色操作
	// menus:{
	// 	url: `${config.API_URL}/basis/menus`,
	// 	name: "菜单列表操作",

	// 	// 单行新增
	// 	add: async function(data={}){
	// 		return await http.put(`${config.API_URL}/basis/auth/menus/add`, data);
	// 	},
		
	// 	// 查询列表
	// 	list: async function(data={}){
	// 		return await http.get(`${config.API_URL}/basis/auth/menus/list`, data);
	// 	},
		
	// 	// 批量更新, 
	// 	// 但根据业务, 不允许赋予多个角色相同的权限
	// 	// 但允许批量停用, 删除等
	// 	edit: async function(data={}){
	// 		return await http.post(`${config.API_URL}/basis/auth/menus/edit`, data);
	// 	},

	// 	// 批量删除
	// 	delete: async function(data={}){
	// 		return await http.delete(`${config.API_URL}/basis/auth/menus/delete`, data);
	// 	},
	// }
}
