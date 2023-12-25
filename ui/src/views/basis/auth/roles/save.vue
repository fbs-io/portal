<template>
	<el-dialog :title="titleMap[mode]" v-model="visible" :width="700" destroy-on-close @closed="$emit('closed')">
		<el-tabs tab-position="top" @tab-click="tabClick">
			<el-tab-pane :label="tablePane.roleInfo">
				<el-form :model="form" :rules="rules" :disabled="mode=='show'" ref="dialogForm" label-width="100px" label-position="left">
					<el-form-item label="角色名称" prop="label">
						<el-input v-model="form.label" clearable></el-input>
					</el-form-item>
					<el-form-item label="角色描述" prop="description">
						<el-input v-model="form.description" clearable></el-input>
					</el-form-item>
					<el-form-item label="排序" prop="sort">
						<el-input-number v-model="form.sort" controls-position="right" :min="1" style="width: 100%;"></el-input-number>
					</el-form-item>
					<el-form-item v-if="mode=='edit'" label="是否有效" prop="status">
						<el-switch v-model="form.status" :active-value="1" :inactive-value="-1"></el-switch>
					</el-form-item>
				</el-form>
			</el-tab-pane>
			<el-tab-pane :label="tablePane.menuPermission">
				<div class="treeMain" >
					<el-tree ref="menu" node-key="code" :data="menu.list" :props="menu.props" show-checkbox></el-tree>
				</div>
			</el-tab-pane>
			<!-- TODO:完善数据权限, 首页画面, 首页组件 -->
			<!-- <el-tab-pane :label="tablePane.dataPermission">
				<el-form label-width="100px" label-position="left">
					<el-form-item label="规则类型">
						<el-select v-model="form.data_permission_type" placeholder="请选择">
							<el-option label="本人可见" value="1"></el-option>
							<el-option label="全部可见" value="2"></el-option>
							<el-option label="所在部门可见" value="3"></el-option>
							<el-option label="所在部门及子级可见" value="4"></el-option>
							<el-option label="选择的部门可见" value="5"></el-option>
							<el-option label="自定义" value="6"></el-option>
						</el-select>
					</el-form-item>
					<el-form-item label="选择部门" v-show="form.data_permission_type==5">
						<div class="treeMain" style="width: 100%;">
							<el-tree ref="dept" node-key="department_code" :data="department.list" :props="department.props" show-checkbox></el-tree>
						</div>
					</el-form-item>
					<el-form-item label="规则值" v-show="data.dataType=='6'">
						<el-input v-model="data.rule" clearable type="textarea" :rows="6" placeholder="请输入自定义规则代码"></el-input>
					</el-form-item>
				</el-form> 
			</el-tab-pane> -->
		</el-tabs>
		<template #footer>
			<el-button @click="visible=false" >取 消</el-button>
			<el-button type="primary" v-if="mode!='show'" :loading="isSaveing" @click="submit(false)">保 存</el-button>
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
				tablePane:{
					roleInfo:"角色基本信息",
					menuPermission:"菜单权限",
					dataPermission:"数据权限",
				},
				visible: false,
				isSaveing: false,
				//表单数据
				form: {
					id:"",
					label: "",
					description: "",
					sort: 1,
					status: 1,
					sources:[],
					data_permission_type:"1",
					data_permission_custom:[],
				},
				menu: {
					list: [],
					checked: [],
					props: {
						label: (data)=>{
							return data.meta.title
						}
					}
				},
				//验证规则
				rules: {
					label: [
						{required: true, message: '请输入角色名称'}
					],
				},
				department: {
					list: [],
					props: {
						label: (data)=>{
							return data.department_name
						},
					},
					checked: [],

				},
			}
		},
		mounted() {
			// this.getMenu()
		},
		watch:{
			form:{
				handler(val){
					if (val.data_permission_type=='5'){
						this.getDepartmentTree()
					}
				},
				immediate: true,
				deep: true,
			}
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
						this.form.sources = this.$refs.menu.getCheckedKeys().concat(this.$refs.menu.getHalfCheckedKeys())
						this.form.data_permission_custom = this.$refs.dept.getCheckedKeys().concat(this.$refs.dept.getHalfCheckedKeys())
						if (this.mode=="add"){
							this.form.data_permission_type=Number(this.form.data_permission_type)
							var res = await this.$API.basis_auth.roles.add(this.form);
							this.isSaveing = false;
							if(res.errno == 0){
								if (isNext) {
									this.form = {
										id:"",
										label: "",
										description: "",
										sort: 1,
										status: "",
										sources: [],
										data_permission_custom:[]
									}
								}else{
									this.visible = false;
								}
							}
						}else if (this.mode=="edit") {
							
							var data = {
								id:[this.form.id],
								label: this.form.label,
								description: this.form.description,
								sort: this.form.sort,
								status: this.form.status,
								sources: this.form.sources,
								data_permission_type:Number(this.form.data_permission_type),
								data_permission_custom: this.form.data_permission_custom,
							}
						
							var res = await this.$API.basis_auth.roles.edit(data);
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
				this.form.data_permission_type = String(this.form.data_permission_type)
			},
			// 菜单列表
			async getMenu(){
				var res = await this.$API.basis_auth.roles.permission()
				this.menu.list = res.details

				// //获取接口返回的之前选中的和半选的合并，处理过滤掉有叶子节点的key
				this.menu.checked = this.form.sources
				this.$nextTick(() => {
					if (this.menu.checked ){
						let filterKeys = this.menu.checked.filter(key => this.$refs.menu.getNode(key).isLeaf)
						this.$refs.menu.setCheckedKeys(filterKeys, true)
					}
				})
			},
			async getDepartmentTree(){
				let res = await this.$API.basis_org.department.tree()
				this.department.list = res.details.rows

				this.department.checked = this.form.data_permission_custom
				this.$nextTick(() => {
					if (this.department.checked ){
						let filterKeys = this.department.checked.filter(key => this.$refs.dept.getNode(key).isLeaf)
						this.$refs.dept.setCheckedKeys(filterKeys, true)
					}
				})
			},
			tabClick(table,event){
				if (table.props.label==this.tablePane.menuPermission){
					this.getMenu()
				}
			}
		}
	}
</script>

<style>
</style>
