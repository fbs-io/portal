<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="700" destroy-on-close @closed="$emit('closed')">
		<el-tabs tab-position="top">
			<el-tab-pane label="角色基本信息">
				<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="100px" label-position="left">

					<el-form-item label="发票号码" :disabled="mode=='edit'" prop="invoice_no">
						<el-input v-model="form.invoice_no" clearable placeholder="发票号码"></el-input>
					</el-form-item>
					<el-form-item label="发票代码" prop="invoice_org_code">
						<el-input v-model="form.invoice_org_code" clearable></el-input>
					</el-form-item>
					<el-form-item label="发票类型" prop="invoice_type">
						<el-input v-model="form.invoice_type" clearable></el-input>
					</el-form-item>
					<el-form-item label="开票时间" prop="invoice_date">
						<el-input v-model="form.invoice_date" clearable></el-input>
					</el-form-item>
					<el-form-item label="发票金额" prop="invoice_amount">
						<el-input v-model="form.invoice_amount" clearable></el-input>
					</el-form-item>
					<el-form-item label="发票校验码" prop="invoice_verify_code">
						<el-input v-model="form.invoice_verify_code" clearable></el-input>
					</el-form-item>
					<el-form-item label="加密字符" prop="invoice_encrypt_code">
						<el-input v-model="form.invoice_encrypt_code" clearable></el-input>
					</el-form-item>
					<el-form-item label="购方名称" prop="invoice_purchas_name">
						<el-input v-model="form.invoice_purchas_name" clearable></el-input>
					</el-form-item>
					<el-form-item label="购方税号" prop="invoice_purchas_code">
						<el-input v-model="form.invoice_purchas_code" clearable></el-input>
					</el-form-item>
					<el-form-item label="发票描述" prop="invoice_commnet">
						<el-input v-model="form.invoice_comment" clearable></el-input>
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
					invoice_no:"",
					invoice_type:"",
					invoice_date:"",
					invoice_org_code:"",
					invoice_verify_code:"",
					invoice_encrypt_code:"",
					invoice_comment:"",
					invoice_purchas_name:"",
					invoice_purchas_code:"",
					invoice_num:"",
					invoice_amount:"",
					invoice_taxes:"",
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
							var res = await this.$API.financial.invoice.add(this.form);
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
							var res = await this.$API.financial.invoice.edit(data);
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
