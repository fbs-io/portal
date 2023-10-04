<template>
	<!-- <el-container> -->
		<!-- <el-aside width="200px" v-loading="showGrouploading">
			<el-container>
				<el-header>
					<el-input placeholder="输入关键字进行过滤" v-model="groupFilterText" clearable></el-input>
				</el-header>
				<el-main class="nopadding">
					<el-tree ref="group" class="menu" node-key="id" :data="group" :current-node-key="''" :highlight-current="true" :expand-on-click-node="false" :filter-node-method="groupFilterNode" @node-click="groupClick"></el-tree>
				</el-main>
			</el-container>
		</el-aside> -->
		<el-container>
				<el-header>
					<div class="left-panel">
						<el-button type="primary" v-auth="auth.put" icon="el-icon-plus" @click="add"></el-button>
						<el-button type="danger" v-auth="auth.delete" plain icon="el-icon-delete" :disabled="selection.length==0" @click="batch_del"></el-button>
						<el-button type="primary" v-auth="auth.post" plain :disabled="selection.length==0" @click="batch_set_role">分配角色</el-button>
						<el-button type="primary" v-auth="auth.post" plain :disabled="selection.length==0" @click="batch_reset_pwd">密码重置</el-button>
					</div>
					<div class="right-panel">
						<div class="right-panel-search">
							<el-input v-model="search.nick_name" placeholder="姓名" clearable></el-input>
							<el-button type="primary" v-auth="auth.get" icon="el-icon-search" @click="upsearch"></el-button>
						</div>
					</div>
				</el-header>
				<el-main class="nopadding">
					<scTable ref="table" :apiObj="apiObj" @selection-change="selectionChange" stripe remoteSort remoteFilter>
						<el-table-column type="selection" width="50"></el-table-column>
						<el-table-column label="ID" prop="id" column-key="id" width="80" sortable='custom'></el-table-column>
						<!-- <el-table-column label="头像" width="80" column-key="Av" :filters="[{text: '已上传', value: '1'}, {text: '未上传', value: '0'}]">
							<template #default="scope">
								<el-avatar :src="scope.row.avatar" size="small"></el-avatar>
							</template>
						</el-table-column> -->
						
						<el-table-column label="登录账号" prop="account" width="130" sortable='custom' column-key="super" :filters="[{text: '管理账户', value: 'Y'}, {text: '普通账户', value: 'N'}]"></el-table-column>
						<el-table-column label="姓名" prop="nick_name" width="130" sortable='custom'></el-table-column>
						<el-table-column label="邮箱" prop="email" width="200" sortable='custom'></el-table-column>
						<el-table-column label="所属角色" prop="role" width="130" sortable='custom' :formatter="formatRoles"></el-table-column>
						<el-table-column label="状态" prop="status" width="130" sortable='custom'>
							<template #default="scope">
								<el-switch v-model="scope.row.status" @change="changeSwitch($event, scope.row)" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
							</template>
						</el-table-column>
						<el-table-column label="加入时间" prop="created_at"  :formatter="timestampToTime" width="170" sortable='custom'></el-table-column>
						<el-table-column label="操作" fixed="right" align="right" width="180">
							<template #default="scope">
								<el-button-group>
									<el-button text type="primary" v-auth="auth.get" size="small" @click="table_show(scope.row, scope.$index)">查看</el-button>
									<el-button text type="primary" v-auth="auth.post" size="small" @click="table_edit(scope.row, scope.$index)">编辑</el-button>
									<el-popconfirm title="确定删除吗？" @confirm="table_del(scope.row, scope.$index)">
										<template #reference>
											<el-button  v-auth="auth.delete" text type="primary" size="small">删除</el-button>
										</template>
									</el-popconfirm>
								</el-button-group>
							</template>
						</el-table-column>

					</scTable>
				</el-main>
		</el-container>
	<!-- </el-container> -->

	<save-dialog v-if="dialog.save" ref="saveDialog" @success="handleSuccess" @closed="dialog.save=false"></save-dialog>
	<role-dialog v-if="dialog.role" ref="roleDialog" @success="handleSuccess" @closed="dialog.role=false"></role-dialog>

</template>

<script>
	import saveDialog from './save'
	import roleDialog from './role'

	export default {
		name: 'user',
		components: {
			saveDialog,
			roleDialog,
		},
		data() {
			return {
				dialog: {
					save: false,
					role: false,
				},
				showGrouploading: false,
				groupFilterText: '',
				group: [],
				apiObj: this.$API.basis_auth.users.list,
				selection: [],
				search: {
					nick_name: null
				},
				roles:[],
				auth:{
					put: '',
					post: '',
					get: '',
					delete: '',
				}
			}
		},
		watch: {
			groupFilterText(val) {
				this.$refs.group.filter(val);
			}
		},
		mounted() {
			this.getGroup()
			this.getRoles()
			this.getUiPermission()
		},
		methods: {
			

			//添加
			add(){
				this.dialog.save = true
				this.$nextTick(() => {
					this.$refs.saveDialog.open()
				})
			},
			//编辑
			table_edit(row){
				this.dialog.save = true
				this.$nextTick(() => {
					this.$refs.saveDialog.open('edit').setData(row)
				})
			},
			//查看
			table_show(row){
				console.log(row)
				this.dialog.save = true
				this.$nextTick(() => {
					this.$refs.saveDialog.open('show').setData(row)
				})
			},
			select_id(){
				var reqData = []
				this.selection.forEach(item => {
					reqData.unshift(item.id)
				})
				return reqData
			},	
			//删除
			async table_del(row, index){
				var reqData = {id: [row.id]}
				await this.$API.basis_auth.users.delete(reqData);
				this.$refs.table.tableData.splice(index, 1)
				this.$refs.table.reload()
			},
			//批量删除
			async batch_del(){
				var reqData = []
				this.$confirm(`确定删除选中的 ${this.selection.length} 项吗？`, '提示', {
					type: 'warning'
				}).then(() => {
					this.selection.forEach(item => {
						this.$refs.table.tableData.forEach((itemI, indexI) => {
							if (item.id === itemI.id) {
								this.$refs.table.tableData.splice(indexI, 1)
							}
						})
						reqData.unshift(item.id)
					})
					const loading = this.$loading();
					this.$API.basis_auth.users.delete({id: reqData})
					loading.close();
					this.$refs.table.reload()
				}).catch(() => {
				})
			},
			// 批量重置密码
			async batch_reset_pwd(){
				this.$confirm(`确定删除选中的 ${this.selection.length} 项重置为密码: abc@123 吗？`, '提示', {
					type: 'warning'
				}).then(() => {
				 	this.$API.basis_auth.users.updates({password:"abc@123",id:this.select_id()})
				}).catch(() => {
				})
			},
			// 批量分配角色
			async batch_set_role(){
				this.dialog.role = true
				this.$nextTick(() => {
					// 把选择的id传个role组件
					this.$refs.roleDialog.open(this.select_id())
				})
			},
			//表格选择后回调事件
			selectionChange(selection){
				this.selection = selection;
			},
			//加载树数据
			async getGroup(){
				this.showGrouploading = true;
				var res = await this.$API.basis_auth.users.list();
				this.showGrouploading = false;
				var allNode ={id: '', label: '所有'}
				this.group = res.data;
			},
			//树过滤
			groupFilterNode(value, data){
				if (!value) return true;
				return data.label.indexOf(value) !== -1;
			},
			//树点击事件
			groupClick(data){
				var params = {
					groupId: data.id
				}
				this.$refs.table.reload(params)
			},
			//搜索
			upsearch(){
				this.$refs.table.upData(this.search)
			},
			//本地更新数据
			handleSuccess(data, mode){
					this.$refs.table.reload()
			},
			// 时间序列化
			timestampToTime (row, column) {
				var date = new Date(row[column.property]) //时间戳为10位需*1000，时间戳为13位的话不需乘1000
				var Y = date.getFullYear() + '-'
				var M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-'
				var D = date.getDate() + ' '
				var h = date.getHours() + ':'
				var m = date.getMinutes() + ':'
				var s = date.getSeconds()
				return Y+M+D+h+m+s
			},
			//表格内开关
			changeSwitch(val, row){
				row.status = row.status == 1 ? 1 : -1
				row.$switch_status = true;
				setTimeout(()=>{
					delete row.$switch_status;
					var res = this.$API.basis_auth.users.edit({id:[row.id],status:row.status});
				}, 1000)
			},

			async getRoles(){
				var res = await this.$API.basis_auth.roles.list();
				this.roles = res.details.rows;
			},
			formatRoles(row,column){
				var data = []
				var roles = row[column.property]
				if (!roles){
					return 
				}
				roles.forEach(item =>{
					var index = item
					var role = this.roles.filter(item => item.code == index)[0]
					if (role){
						data.unshift(role.label)
					}
					
				})
				return data.join(", ")
			},
			async getUiPermission(){
				var path = this.$router.currentRoute.value.fullPath
				var res = await this.$API.common.uiPermissions.get({path:path})
				this.auth = res.details
			},
		}
	}
</script>

<style>
</style>
