/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-04 17:43:12
 * @Description: 请填写简介
 */
import store from '@/stores'
import { nextTick } from 'vue'

export function beforeEach(to, from){
	var adminMain = document.querySelector('#adminui-main')
	if(!adminMain){return false}
	store.commit("updateViewTags", {
		fullPath: from.fullPath,
		scrollTop: adminMain.scrollTop
	})
}

export function afterEach(to){
	var adminMain = document.querySelector('#adminui-main')
	if(!adminMain){return false}
	nextTick(()=>{
		var beforeRoute = store.state.viewTags.viewTags.filter(v => v.fullPath == to.fullPath)[0]
		if(beforeRoute){
			adminMain.scrollTop = beforeRoute.scrollTop || 0
		}
	})
}