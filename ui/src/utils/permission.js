/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:58:32
 * @Description: 按钮权限校验
 */
import tool from '@/utils/tool';

export function permission(data) {
	let permissions = tool.data.get("PERMISSIONS");
	if(!permissions){
		return false;
	}
	let isHave = permissions[data];
	return isHave;
}

export function rolePermission(data) {
	let userInfo = tool.data.get("USER_INFO");
	if(!userInfo){
		return false;
	}
	let role = userInfo.role;
	if(!role){
		return false;
	}
	let isHave = role.includes(data);
	return isHave;
}
