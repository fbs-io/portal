<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="700" destroy-on-close @closed="$emit('closed')">
		<el-tabs tab-position="top">
			<el-tab-pane label="公司基本信息">
				<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="100px" label-position="left">
					<el-form-item label="公司代码" :disabled="mode=='edit'" prop="company_code">
						<el-input v-model="form.company_code" clearable placeholder="为空则自动生产代码"></el-input>
					</el-form-item>
					<el-form-item label="公司名称" prop="company_name">
						<el-input v-model="form.company_name" clearable></el-input>
					</el-form-item>
					<el-form-item label="公司简称" prop="company_shortname">
						<el-input v-model="form.company_shortname" clearable></el-input>
					</el-form-item>
					<el-form-item label="公司描述" prop="company_comment">
						<el-input v-model="form.company_comment" clearable></el-input>
					</el-form-item>
					<el-form-item label="所属行业" prop="company_business">
						<el-input v-model="form.company_business" clearable></el-input>
					</el-form-item>
					<el-form-item v-if="mode=='edit'" label="是否有效" prop="status">
						<el-switch v-model="form.status" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
				</el-form>
			</el-tab-pane>
		</el-tabs>
		<template #footer>
			<el-button @click="visible=false" >取 消</el-button>
			<el-button type="primary" :loading="isSaveing" @click="submit(false)">保 存</el-button>
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
					add: '新增',
					edit: '编辑',
					show: '查看'
				},
				visible: false,
				isSaveing: false,
				//表单数据
				form: {
					id: 0,
					company_code:"",
					company_name: "",
					company_shortname: "",
					company_comment: "",
					company_business: "",
					status:1
				},

				//验证规则
				rules: {
					company_name: [
						{required: true, message: '请输入公司名称'}
					],

				}
			}
		},
		mounted() {
		},
		methods: {
			//显示
			open(mode='add'){
				this.mode = mode;
				this.visible = true;
				return this
			},

			//表单提交方法
			submit(isNext){
				this.$refs.dialogForm.validate(async (valid) => {
					if (valid) {
						this.isSaveing = true;
						if (this.mode=="add"){
							var res = await this.$API.basis_org.company.add(this.form);
							this.isSaveing = false;
							if(res.errno == 0){								
								if (isNext) {
									this.form = {
										company_code:"",
										company_name: "",
										company_shortname: "",
										company_comment: "",
										company_business: "",
									}
								}else{
									this.visible = false;
								}
							}
						}else if (this.mode=="edit") {
							
							var data = {
								id:[this.form.id],
								company_code: this.form.company_code,
								company_name: this.form.company_name,
								company_shortname: this.form.company_shortname,
								company_comment: this.form.company_comment,
								company_business: this.form.company_business,
							}
							var res = await this.$API.basis_org.company.edit(data);
							this.isSaveing = false;
							if (res.errno==0){
								this.visible = false;
							}
						}
						this.$emit('success',this.form, this.mode)
					}else{
						return false;
					}
				})
			},
			//表单注入数据
			setData(data){
				//可以和上面一样单个注入，也可以像下面一样直接合并进去
				Object.assign(this.form, data)
			},
		}
	}
</script>

<style>
</style>
