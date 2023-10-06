/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:58:06
 * @Description: 自定义v-auth指令
 */
import { permission } from '@/utils/permission'
import { ref } from "vue"

export default {
	updated(el,binding){
		var ok  = ref(false)
		const { value } = binding
		if(Array.isArray(value)){
			value.forEach(item => {
				if(permission(item)){
					ok = ref(true);
				}
			})
		}else{
			ok = ref(permission(value))
		}
		el.style.display = ok.value ? "block":"none"
	}
	
};
