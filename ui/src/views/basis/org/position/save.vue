<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="700" destroy-on-close @closed="$emit('closed')">
		<el-tabs tab-position="top">
			<el-tab-pane label="岗位基本信息">
				<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="130px" label-position="left">
					<el-form-item label="岗位代码" prop="position_code">
						<el-input v-model="form.position_code" :disabled="mode=='edit'" clearable placeholder="为空则自动生产代码"></el-input>
					</el-form-item>
					<el-form-item label="岗位名称" prop="position_name">
						<el-input v-model="form.position_name" clearable></el-input>
					</el-form-item>
					<el-form-item label="岗位描述" prop="position_comment">
						<el-input v-model="form.position_comment" clearable></el-input>
					</el-form-item>
					<el-form-item label="上级岗位" prop="position_parent_code">
						<template #default="scope">
							<el-select 
								ref="position" 
								filterable
								v-model="form.position_parent_code"
							>
								<el-option v-for="item in position.list" :label="item.name" :value="item.code"/>
							</el-select>
						</template>
					
					</el-form-item>
					<el-form-item label="所属部门" prop="department_code">
						<template #default="scope">
							<el-tree-select 
								ref="department" 
								v-model="form.department_code"
								node-key="department_code"
								:data="department.list" 
								:props="department.props" 
								check-strictly
								:render-after-expand="false"
							></el-tree-select>
						</template>
					</el-form-item>
					<!-- TODO: 添加职务-->
					<!-- <el-form-item label="所属职务" prop="job_code">
						<el-select v-model="form.job_code" multiple filterable style="width: 100%">
							<el-option v-for="item in jobs" :key="item.job_code" :label="item.job_name" :value="item.job_code"/>
						</el-select>
					</el-form-item> -->
					<el-form-item  label="是否部门负责人" prop="is_head">
						<el-switch v-model="form.is_head" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
					<el-form-item  label="是否有审批权限" prop="is_approve">
						<el-switch v-model="form.is_approve" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
					<el-form-item  label="是否为虚拟岗位" prop="is_vritual">
						<el-switch v-model="form.is_vritual" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
					<el-form-item v-if="mode=='edit'" label="是否有效" prop="status">
						<el-switch v-model="form.status" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
				</el-form>
			</el-tab-pane>
			<el-tab-pane label="数据权限">
				<el-form label-width="100px" label-position="left">
					<el-form-item label="规则类型">
						<el-select v-model="form.data_permission_type" placeholder="请选择">
							<el-option v-for="item in data_permission_types" :label="item.label" :value="item.value"></el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="选择部门" v-show="form.data_permission_type==5">
						<div class="treeMain" style="width: 100%;">
							<el-tree ref="department" node-key="department_code" :data="department.list" :props="department.props" show-checkbox></el-tree>
						</div>
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
					position_code:"",
					position_name:"",
					position_comment:"",
					position_parent_code:"",
					department_code:"",
					job_code:"",
					is_head:-1,
					is_approve:-1,
					is_vritual:-1,
					data_permission_type:1,
					data_permission_custom:[],
					status:1
				},
				data_permission_types:[
					{label:"本人可见", value:1},
					{label:"全部可见", value:2},
					{label:"所在部门可见", value:3},
					{label:"所在部门及子级可见", value:4},
					{label:"选择的部门可见", value:5},
				],
				//验证规则
				rules: {
					department_name: [
						{required: true, message: '请输入部门名称'}
					],

				},
				position: {
					list: [],
					checked: [],
					props: {
						label: (data)=>{
							return data.name
						},
					}
				},
				department: {
					list: [],
					checked: [],
					props: {
						label: (data)=>{
							return data.department_name
						},
					}
				},
			}
		},
		mounted() {
			this.getPosition()
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
						this.form.data_permission_custom = this.$refs.department.getCheckedKeys().concat(this.$refs.department.getHalfCheckedKeys())
						if (this.mode=="add"){
							var data = {
								position_code: this.form.position_code,
								position_name: this.form.position_name,
								position_comment: this.form.position_comment,
								position_parent_code: this.form.position_parent_code,
								department_code: this.form.department_code,
								job_code: this.form.job_code,
								is_head: this.form.is_head,
								is_approve: this.form.is_approve,
								is_vritual: this.form.is_vritual,
								data_permission_type: Number(this.form.data_permission_type),
								data_permission_custom: this.form.data_permission_custom,
								status: this.form.status,
							}
							var res = await this.$API.basis_org.position.add(data);
							this.isSaveing = false;
							if(res.errno == 0){
								if (isNext) {
									this.form = {
										position_code:"",
										position_name:"",
										position_comment:"",
										position_parent_code:"",
										department_code:"",
										job_code:"",
										is_head:-1,
										is_approve:-1,
										is_vritual:-1,
										data_permission_type:"1",
										data_permission_custom:"",
										status:1
									}
								}else{
									this.visible = false;
								}
							}
						}else if (this.mode=="edit") {
							var data = {
								id:[this.form.id],
								position_code: this.form.position_code,
								position_name: this.form.position_name,
								position_comment: this.form.position_comment,
								position_parent_code: this.form.position_parent_code,
								department_code: this.form.department_code,
								job_code: this.form.job_code,
								is_head: this.form.is_head,
								is_approve: this.form.is_approve,
								is_vritual: this.form.is_vritual,
								data_permission_type: Number(this.form.data_permission_type),
								data_permission_custom: this.form.data_permission_custom,
								status: this.form.status,
							}
							var res = await this.$API.basis_org.position.edit(data);
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
			async getPosition(){
				var res = await this.$API.common.dimension.get({dim_type:"position"});
				this.position.list = res.details
			},
			//表单注入数据
			setData(data){
				//可以和上面一样单个注入，也可以像下面一样直接合并进去
				Object.assign(this.form, data)
			},
			async getDepartmentTree(){
				let res = await this.$API.basis_org.department.tree()
				this.department.list = res.details.rows

				this.department.checked = this.form.data_permission_custom
				this.$nextTick(() => {
					if (this.department.checked ){
						let filterKeys = this.department.checked.filter(key => this.$refs.department.getNode(key).isLeaf)
						this.$refs.department.setCheckedKeys(filterKeys, true)
					}
				})
			}
		}
	}
</script>

<style>
</style>
