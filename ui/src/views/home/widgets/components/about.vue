<template>
	<el-card shadow="hover" header="关于项目" class="item-background">
		<p>高性能 / 精致 / 优雅。基于Vue3 + Element-Plus + Vite 的中后台前端解决方案，如果喜欢就点个星星支持一下。</p>
		<p>
			<a href='https://gitee.com/qbiancheng/scui-vite-ui/stargazers'>
				<img src='https://gitee.com/qbiancheng/scui-vite-ui/badge/star.svg?theme=dark' alt='star'  style="vertical-align: middle">
			</a>
		</p>
	</el-card>
</template>

<script>
	export default {
		title: "关于项目",
		icon: "el-icon-setting",
		description: "点个星星支持一下",
		data() {
			return {
				
			}
		}
	}
</script>

<style scoped>
	.item-background p {color: #999;margin-top:10px;line-height: 1.8;}
</style>
