<!--
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-10-04 23:08:55
 * @Description: 请填写简介
-->
<template>
	<el-card shadow="hover" header="版本信息">
		<div style="height: 210px;text-align: center;">
			<img src="/website/img/ver.svg" style="height:140px"/>
			<h2 style="margin-top: 15px;">SCUI VITE ADMIN {{$CONFIG.CORE_VER}}</h2>
			<p style="margin-top: 5px;">最新版本 {{ver}}</p>
		</div>
		<div style="margin-top: 20px;">
			<el-button type="primary" plain round @click="golog">更新日志</el-button>
			<el-button type="primary" plain round @click="gogit">gitee</el-button>
		</div>
	</el-card>
</template>

<script>
	export default {
		title: "版本信息",
		icon: "el-icon-monitor",
		description: "当前项目版本信息",
		data() {
			return {
				ver: '0.1'
			}
		},
		mounted() {
			this.getVer()
		},
		methods: {
			async getVer(){
				// const ver = await this.$API.demo.ver.get()
				// this.ver = ver.data
			},
			golog(){
				window.open("https://gitee.com/qbiancheng/scui-vite-ui/releases")
			},
			gogit(){
				window.open("https://gitee.com/qbiancheng/scui-vite-ui")
			}
		}
	}
</script>
