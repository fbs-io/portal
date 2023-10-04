/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-04 16:56:10
 * @Description: 请填写简介
 */
import { permission } from '@/utils/permission'
import { watchEffect,ref } from "vue"
export default {
	mounted(el, binding) {
		const { value } = binding;
		var ok  = ref(false)
		watchEffect(()=>{
			if(Array.isArray(value)){
				value.forEach(item => {
					if(permission(item)){
						ok = ref(true);
					}
				})
			}else{
				ok = ref(permission(value))
			}
		})
		if (!ok){
			el.parentNode.removeChild(el)
		}
	},
	
};
