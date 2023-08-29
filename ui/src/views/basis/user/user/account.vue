<!--
 * @Author: reel
 * @Date: 2023-07-30 22:36:55
 * @LastEditors: reel
 * @LastEditTime: 2023-08-26 22:06:52
 * @Description: 请填写简介
-->
<template>
	<el-card shadow="never" header="账户信息">
		<el-form ref="account" :model="account" :rules="rules" label-width="120px" style="margin-top:20px;">
			<el-form-item label="账号">
				<el-input v-model="account.account" disabled></el-input>
				<div class="el-form-item-msg">账号信息用于登录，系统不允许修改</div>
			</el-form-item>
			<el-form-item label="姓名">
				<el-input v-model="account.nick_name"></el-input>
			</el-form-item>
			<el-form-item label="邮箱" prop="email">
				<el-input v-model="account.email"></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" @click="updateAccount">保存</el-button>
			</el-form-item>
		</el-form>
	</el-card>
</template>

<script>

	export default {
		data() {
			return {
				account: {},
				rules: {
					email: [
						{ message: '请输入邮箱地址'}
					],
				}
			}
		},
		created(){
			this.account = this.$TOOL.data.get("USER_INFO");
		},
		methods:{
			async updateAccount(){			
				// 校验参数
				var validate = await this.$refs.account.validate().catch(()=>{})
				if(!validate){ return false }

				var data = {
					id: this.account.id,
					nick_name: this.account.nick_name,
					email: this.account.email
				}

				// 更新信息
				var res = await this.$API.user.updateUser.put(data)
				if(res.errno == 0){
					this.$TOOL.data.set("USER_INFO", res.details)
				// }else{
				// 	this.$message.warning(res.message)
				// 	return false
				}

			}
		}
	}
</script>

<style>
</style>
