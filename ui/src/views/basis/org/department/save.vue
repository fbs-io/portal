<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="700" destroy-on-close @closed="$emit('closed')">
		<el-tabs tab-position="top">
			<el-tab-pane label="角色基本信息">
				<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="100px" label-position="left">
					<el-form-item label="部门代码" prop="department_code">
						<el-input v-model="form.department_code" clearable placeholder="为空则自动生产代码"></el-input>
					</el-form-item>
					<el-form-item label="部门名称" prop="department_name">
						<el-input v-model="form.department_name" clearable></el-input>
					</el-form-item>
					<el-form-item label="部门描述" prop="department_comment">
						<el-input v-model="form.department_comment" clearable></el-input>
					</el-form-item>
					<el-form-item label="父级部门" prop="department_parent_code">
						<!-- <el-input v-model="form.department_parent_code" clearable></el-input> -->
						<!-- <div class="treeMain" > -->
						<template #default="scope">

							<el-tree-select 
								ref="department" 
								v-model="form.department_parent_code"
								node-key="department_code"
								:data="department.list" 
								:props="department.props" 
								check-strictly
								:render-after-expand="false"
							></el-tree-select>
						</template>
						<!-- </div> -->
					</el-form-item>
					<el-form-item label="自定义层级" prop="department_custom_level">
						<el-input v-model="form.department_custom_level" clearable></el-input>
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
	import { ref } from 'vue';
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
					department_code:"",
					department_name: "",
					department_comment: "",
					department_custom_level: "",
					department_parent_code:"",
					status:1
				},

				//验证规则
				rules: {
					department_name: [
						{required: true, message: '请输入部门名称'}
					],

				},				
				department: {
					list: [],
					checked: [],
					props: {
						label: (data)=>{
							return data.department_name
						},
						// value: (data)=>{
						// 	return data.department_code
						// }
					}
				},
			}
		},
		mounted() {
			this.getDepartmentTree()
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
							var res = await this.$API.basis_org.department.add(this.form);
							this.isSaveing = false;
							if(res.errno == 0){
								if (isNext) {
									this.form = {
										department_code:"",
										department_name: "",
										department_comment: "",
										department_custom_level: "",
										department_parent_code:"",
									}
								}else{
									this.visible = false;
								}
							}
						}else if (this.mode=="edit") {
							
							var data = {
								id:[this.form.id],
								department_code: this.form.department_code,
								department_name: this.form.department_name,
								department_comment: this.form.department_comment,
								department_custom_level: this.form.department_custom_level,
								department_parent_code: this.form.department_parent_code,
								status:this.status,

							}
							var res = await this.$API.basis_org.department.edit(data);
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
			async getDepartmentTree(){
				let res = await this.$API.basis_org.department.tree()
				this.department.list = res.details.rows
			}
		}
	}
</script>

<style>
</style>
