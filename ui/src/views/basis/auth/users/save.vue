<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="500" destroy-on-close @closed="$emit('closed')">
		<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="100px" label-position="left">
			<!-- <el-form-item label="头像" prop="avatar">
				<sc-upload v-model="form.avatar" title="上传头像"></sc-upload>
			</el-form-item> -->
			<el-form-item label="登录账号" prop="account">
				<el-input v-model="form.account" placeholder="用于登录系统" :disabled="mode=='edit'"></el-input>
			</el-form-item>
			<el-form-item label="姓名" prop="nick_name">
				<el-input v-model="form.nick_name" placeholder="请输入完整的真实姓名" clearable></el-input>
			</el-form-item>

			<template v-if="mode=='add'">
				<template v-if="mode=='add'">
					<el-form-item label="默认密码" prop="isDefaultPwd">
						<el-select v-model="isDefaultPwd" clearable >
							<el-option v-for="item in isPwds" :key="item.value" :label="item.label" :value="item.value"/>
						</el-select>
						<div class="el-form-item-msg">默认密码:账户+123</div>
					</el-form-item>
				</template>
				<template v-if="isDefaultPwd">
					<el-form-item label="登录密码" prop="password">
						<el-input type="password" v-model="form.password" clearable show-password></el-input>
					</el-form-item>
					<el-form-item label="确认密码" prop="password2">
						<el-input type="password" v-model="form.password2" clearable show-password></el-input>
					</el-form-item>
				</template>
			</template>
			<el-form-item label="邮箱" prop="email">
				<el-input v-model="form.email" placeholder="请输入完整的邮箱" clearable></el-input>
			</el-form-item>
			<el-form-item label="管理员" prop="super">
				<el-select v-model="form.super" clearable >
					 <el-option v-for="item in supers" :key="item.value" :label="item.label" :value="item.value"/>
				</el-select>
			</el-form-item>
			<el-form-item label="主要岗位" prop="position1">
				<template #default="scope">
					<el-select 
						ref="position" 
						filterable
						v-model="form.position1"
					>
						<el-option v-for="item in position.list" :label="item.name" :value="item.code" :disabled="form.position2.indexOf(item.code)>-1"/>
					</el-select>
				</template>
			</el-form-item>

			<el-form-item label="兼职岗位" prop="position2">
				<template #default="scope">
					<el-select 
						ref="position2" 
						filterable
						multiple
						v-model="form.position2"
					>
						<el-option v-for="item in position.list" :label="item.name" :value="item.code" :disabled="item.code==form.position1"/>
					</el-select>
				</template>
			</el-form-item>
			<el-form-item label="所属公司" prop="company">
				<el-select v-model="form.company" multiple filterable style="width: 100%">
					<el-option v-for="item in companies" :key="item.company_code" :label="item.company_name" :value="item.company_code"/>
				</el-select>
			</el-form-item>
			<el-form-item label="所属角色" prop="role">
				<el-select v-model="form.role" multiple filterable style="width: 100%">
					<el-option v-for="item in roles" :key="item.code" :label="item.label" :value="item.code"/>
				</el-select>
			</el-form-item>
		</el-form>
		<template #footer>
			<el-button @click="visible=false" >取 消</el-button>
			<el-button v-if="mode!='show'" type="primary" :loading="isSaveing" @click="submit(false)">保 存</el-button>
			<el-button v-if="mode=='add'" type="primary" :loading="isSaveing" @click="submit(true)">保存并继续</el-button>
		</template>
	</el-dialog>
</template>

<script>
	export default {
		emits: ['success', 'closed'],
		data() {
			return {
				mode: "add",
				titleMap: {
					add: '新增用户',
					edit: '编辑用户',
					show: '查看'
				},
				visible: false,
				isSaveing: false,
				isDefaultPwd:false,
				//表单数据
				form: {
					id:"",
					account: "",
					avatar: "",
					nick_name: "",
					dept: "",
					role: [],
					company: [],
					position1: "",
					position2: [],
					super:"N",
					password:"",
					password2:"",
					email:"",
				},
				form2:{},
				//验证规则
				rules: {
					// avatar:[
					// 	{required: true, message: '请上传头像'}
					// ],
					account: [
						{required: true, message: '请输入登录账号'}
					],
					nick_name: [
						{required: true, message: '请输入真实姓名'}
					],
					password: [
						{required: true, message: '请输入登录密码'},
						{validator: (rule, value, callback) => {
							if (this.form.password2 !== '') {
								this.$refs.dialogForm.validateField('password2');
							}
							callback();
						}}
					],
					password2: [
						{required: true, message: '请再次输入密码'},
						{validator: (rule, value, callback) => {
							if (value !== this.form.password) {
								callback(new Error('两次输入密码不一致!'));
							}else{
								callback();
							}
						}}
					],
					position1: [
						{required: true, message: '请选择所属岗位'}
					],
					role: [
						{required: true, message: '请选择所属角色', type: "Array"}
					],
					company: [
						{required: true, message: '请选择所属公司', type: "Array"}
					],
				},
				//所需数据选项
				companies: [],
				companiesProps: {
					value: "company_code",
					multiple: true,
					checkStrictly: true
				},
				roles: [],
				rolesProps: {
					value: "code",
					multiple: true,
					checkStrictly: true
				},
				position: {
					list: [],
					checked: [],
					props: {
						label: (data)=>{
							return data.position_name
						},
					}
				},
				supers:[
					{
						value:"N",
						label: "否"
					},
					{
						value:"Y",
						label:"是"
					}
				],
				isPwds:[
					{
						value:false,
						label: "是"
					},
					{
						value:true,
						label:"否"
					}
				]
			}
		},
		mounted() {
		},
		// watch: {
		// 	form:{
		// 		handler(val,oval){
		// 			console.log(val,oval)
		// 		},
		// 		deep:true
		// 	}
		// },
		methods: {
			//显示
			open(mode='add'){
				this.mode = mode;
				this.visible = true;
				
				// if (this.mode!="show"){
					this.getRoles()
					this.getCompanies()
					this.getPositions()
				// }
				if (mode=='add'){
					this.getAllowCompany()
				}
				return this
			},
			//加载树数据
			async getRoles(){
				var res = await this.$API.basis_auth.roles.list({page_num:-1,page_size:-1});
				this.roles = res.details.rows;
			},
			async getCompanies(){
				var res = await this.$API.basis_org.company.list({page_num:-1,page_size:-1});
				this.companies = res.details.rows;
			},
			async getPositions(){
				var res = await this.$API.common.dimension.get({dim_type:"position"});
				this.position.list = res.details
			},
			async getAllowCompany(){
				var user  = this.$TOOL.data.get("USER_INFO")
				var res = await this.$API.basis_auth.user.getAllowCompany.get({account:user.account})
				this.form.company.push(res.details.company)
			},
			//表单提交方法
			submit(isNext){
				this.$refs.dialogForm.validate(async (valid) => {
					if (valid) {
						this.isSaveing = true;
						if (this.mode=="add"){
							if (!this.isDefaultPwd){
								this.form.password = this.form.account+"123"
							}
							var res = await this.$API.basis_auth.users.add(this.form);
							this.isSaveing = false;
							if(res.errno == 0){								
								if (isNext) {
									this.form = {
										id:"",
										account: "",
										avatar: "",
										nick_name: "",
										role: "",
										company:[],
										super:"N",
										email:"",
										position1:"",
										position:"",
										position2:"",
									}
								}else{
									this.visible = false;
								}
							}
						}else if (this.mode=="edit") {
							var data = {
								id: [this.form.id],
								account: this.form.account,
								avatar: this.form.avatar,
								nick_name: this.form.nick_name,
								super: this.form.super,
								email: this.form.email,
								company: null,
								role: null,
								position1: "",
								position2: null,
							}
							// 如果一样重置为空, 不更新
							if (this.form.position1 != this.form2.position1 || this.form.position2 != this.form2.position2){
								data.position1 = this.form.position1
								data.position2 = this.form.position2
							}
							if (this.form.company != this.form2.company){
								data.company = this.form.company
							}
							if (this.form.role != this.form2.role){
								data.role = this.form.role
							}
							
							
							var res = await this.$API.basis_auth.users.edit(data);
							this.isSaveing = false;
							if (res.errno==0){
								this.visible = false;
							}
						}
						this.$emit("success",this.form,this.mode)
					}else{
						return false;
					}
				})
			},
			//表单注入数据
			setData(data){
				//可以和上面一样单个注入，也可以像下面一样直接合并进去
				Object.assign(this.form, data)
				this.form2["position1"] = data.position1
				this.form2["position2"] = data.position2
				this.form2["company"] = data.company
				this.form2["role"] = data.role
			},

		}
	}
</script>

<style>
</style>
