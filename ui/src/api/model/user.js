/*
 * @Author: reel
 * @Date: 2023-08-21 06:01:52
 * @LastEditors: reel
 * @LastEditTime: 2023-08-29 21:31:21
 * @Description: 用户相关api, 如更新账户信息, 修改密码等
 */
import config from "@/config"
import http from "@/utils/request"

export default {
	updateUser: {
		url: `${config.API_URL}/basis/user/update`,
		name: "修改用户信息",
		put: async function(data={}){
			return await http.put(this.url, data);
		}
	},
	chPwd: {
		url: `${config.API_URL}/basis/user/chpwd`,
		name: "修改用户密码",
		put: async function(data={}){
			return await http.put(this.url, data);
		}
	},
	users:{
		url: `${config.API_URL}/basis/users`,
		name: "用户列表操作",
		// 查询
		list: async function(data={}){
			return await http.get(`${config.API_URL}/basis/users/list`, data);
		},
		// 新增
		add: async function(data={}){
			return await http.post(`${config.API_URL}/basis/user/add`, data);
		},
		// 更新
		edit: async function(data={}){
			return await http.put(`${config.API_URL}/basis/users/edit`, data);
		},
		// 删除
		delete: async function(data={}){
			return await http.delete(`${config.API_URL}/basis/users/delete`, data);
		},
		updates: async function(data={}){
			return await http.put(`${config.API_URL}/basis/users/edit`, data);
		},
		

	}
}
