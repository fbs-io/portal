<!--
 * @Author: reel
 * @Date: 2023-06-04 15:27:47
 * @LastEditors: reel
 * @LastEditTime: 2024-03-17 18:51:50
 * @Description: 登陆组件
-->
<template>
	<el-form ref="loginForm" :model="form" :rules="rules" label-width="0" size="large">
		<el-form-item prop="user">
			<el-input v-model="form.user" prefix-icon="el-icon-user" @blur="mouseLeave"  :placeholder="$t('login.userPlaceholder')">
			</el-input>
		</el-form-item>
		<el-form-item  prop="company" v-if="isCompany">
				<el-select v-model="form.company" style="width: 100%" clearable>
					<template #prefix>
						<sc-icon-company style="width: 15;height: 15;"></sc-icon-company>
					</template>
					<el-option v-for="item in companies" :key="item.company_code" :label="item.company_name" :value="item.company_code" />
				</el-select>
			</el-form-item>
		<el-form-item prop="password">
			<el-input v-model="form.password" prefix-icon="el-icon-lock" clearable show-password :placeholder="$t('login.PWPlaceholder')"></el-input>
		</el-form-item>
		<el-form-item>
			<el-button type="primary" style="width: 100%;" :loading="islogin" round @click="login">{{ $t('login.signIn') }}</el-button>
		</el-form-item>
	</el-form>
</template>

<script>
	export default {
		data() {
			return {
				form: {
					user: "",
					password: "",
					company:"",
					autologin: true
				},
				rules: {
					user: [
						{required: true, message: this.$t('login.userError'), trigger: 'blur'}
					],
					password: [
						{required: true, message: this.$t('login.PWError'), trigger: 'blur'}
					]
				},
				islogin: false,
				companies:[],
				isCompany:false
			}
		},
		watch:{
		},
		mounted() {

		},
		methods: {
			async login(){

				var validate = await this.$refs.loginForm.validate().catch(()=>{})
				if(!validate){ return false }

				this.islogin = true
				var data = {
					account: this.form.user,
					password: this.form.password,
					company_code: this.form.company
				}
				//获取token
				var res = await this.$API.basis_auth.token.post(data)
				if(res.errno == 0){
					this.$TOOL.cookie.set("TOKEN", res.details.token, {
						expires: this.form.autologin? this.$CONFIG.TOKEN_EXPIRED_TIME : 0
					})
					this.$TOOL.data.set("USER_INFO", res.details.userInfo)
					if (res.details.menu.length==0){
						this.islogin = false
						this.$alert("当前用户无任何菜单权限，请联系系统管理员", "无权限访问", {
							type: 'error',
							center: true
						})
						return false
					}
					this.$TOOL.data.set("MENU", res.details.menu)
					this.$TOOL.data.set("PERMISSIONS", res.details.permissions)
				}else{
					this.islogin = false
					// this.$message.warning(res.message)
					return false
				}

				this.$router.replace({
					path: '/',
					components:"overview",
				})
				// this.$message.success("Login Success 登录成功")
				this.islogin = false
			},
			async mouseLeave(){
				var res = await this.$API.basis_auth.user.getCompany.get({account:this.form.user})
				if (res.details){
					if (res.details.companies && res.details.companies.length>1){
						this.isCompany = true
					}else{
						this.isCompany = false
					}
					this.companies = res.details.companies
				}
			}
		}
	}
</script>

<style>
</style>
