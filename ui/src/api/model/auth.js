/*
 * @Author: reel
 * @Date: 2023-06-24 08:35:55
 * @LastEditors: reel
 * @LastEditTime: 2023-08-23 06:47:37
 * @Description: 请填写简介
 */
import config from "@/config"
import http from "@/utils/request"

export default {
	token: {
		url: `${config.API_URL}/basis/user/login`,
		name: "登录获取TOKEN",
		post: async function(data={}){
			return await http.post(this.url, data);
		}
	}
}
