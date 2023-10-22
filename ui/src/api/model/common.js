/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:59:12
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"

export default {
	uiPermissions: {
		url: `${config.API_URL}/uipermission`,
		name: "页面按钮权限",
		get: function(data, config={}){
			return http.get(this.url, data, config);
		}
	},
	dimension: {
		url: `${config.API_URL}/dimension`,
		name: "维度列表",
		get: async function(data, config={}){
			return await http.get(this.url, data, config);
		}
	},
}
