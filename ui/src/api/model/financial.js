/*
 * @Author: reel
 * @Date: 2023-12-31 20:14:22
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 20:18:17
 * @Description: 财务模块接口
 */

import config from "@/config"
import http from "@/utils/request"

export default {

	// 发票批量操作
	invoice:{
		name: "发票批量操作",
		// 查询
		list: async function(data={}){
			return await http.get(`${config.API_URL}/financial/bill/invoice/list`, data);
		},
		// // 新增
		add: async function(data={}){
			return await http.put(`${config.API_URL}/financial/bill/invoice/add`, data);
		},
		// // 更新
		edit: async function(data={}){
			return await http.post(`${config.API_URL}/financial/bill/invoice/edit`, data);
		},
		// // 删除
		delete: async function(data={}){
			return await http.delete(`${config.API_URL}/financial/bill/invoice/delete`, data);
		},

	},
}
