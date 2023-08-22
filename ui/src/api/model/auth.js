/*
 * @Author: reel
 * @Date: 2023-06-24 08:35:55
 * @LastEditors: reel
 * @LastEditTime: 2023-08-19 14:40:22
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"

export default {
	token: {
		url: `${config.API_URL}/login`,
		name: "登录获取TOKEN",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	}
}
