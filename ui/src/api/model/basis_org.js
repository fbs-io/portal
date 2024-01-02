
/*
 * @Author: reel
 * @Date: 2023-06-24 08:35:55
 * @LastEditors: reel
 * @LastEditTime: 2024-01-02 22:56:19
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

	department:{
		url: `${config.API_URL}/basis/org/department/`,
		name: "用户列表操作",
		// 查询
		list: async function(data={}){
			return await http.get(`${config.API_URL}/basis/org/department/list`, data);
		},
		// 查询
		tree: async function(data={}){
			return await http.get(`${config.API_URL}/basis/org/department/tree`, data);
		},
		// // 新增
		add: async function(data={}){
			return await http.put(`${config.API_URL}/basis/org/department/add`, data);
		},
		// // 更新
		edit: async function(data={}){
			return await http.post(`${config.API_URL}/basis/org/department/edit`, data);
		},
		// // 删除
		delete: async function(data={}){
			return await http.delete(`${config.API_URL}/basis/org/department/delete`, data);
		},

	},

	//
	position:{
		url: `${config.API_URL}/basis/org/position/`,
		name: "用户列表操作",
		// 查询
		list: async function(data={}){
			return await http.get(`${config.API_URL}/basis/org/position/list`, data);
		},
		// // 查询
		// tree: async function(data={}){
		// 	return await http.get(`${config.API_URL}/basis/org/department/tree`, data);
		// },
		// // 新增
		add: async function(data={}){
			return await http.put(`${config.API_URL}/basis/org/position/add`, data);
		},
		// // 更新
		edit: async function(data={}){
			return await http.post(`${config.API_URL}/basis/org/position/edit`, data);
		},
		// // 删除
		delete: async function(data={}){
			return await http.delete(`${config.API_URL}/basis/org/position/delete`, data);
		},

	},
}
