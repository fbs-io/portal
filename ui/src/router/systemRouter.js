/*
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-20 07:05:25
 * @Description: 请填写简介
 */
import config from "@/config"

//系统路由
const routes = [
	{
		name: "layout",
		path: "/",
		component: () => import('@/layout'),
		meta:{
			title: "首页"
		},
		redirect: '/home',
		children: [

		]
	},
	{
		path: "/login",
		component: () => import('@/views/login'),
		meta: {
			title: "登录"
		}
	},
	{
		path: "/user_register",
		component: () => import('@/views/login/userRegister'),
		meta: {
			title: "注册"
		}
	},
	{
		path: "/reset_password",
		component: () => import('@/views/login/resetPassword'),
		meta: {
			title: "重置密码"
		}
	},
	{
		path: "/404",
		hidden: true,
		component: () => import('@/layout/other/404'),
	},
	{
		path:"/:catchAll(.*)",
		component: () => import('@/layout/other/404'),
		hidden: true
	}
]

export default routes;