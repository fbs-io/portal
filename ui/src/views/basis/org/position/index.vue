<template>
	<el-container>
		<el-header>
			<div class="left-panel">
				<el-button type="primary" v-auth="auth?auth.put:auth" icon="el-icon-plus" @click="add"></el-button>
				<el-button type="danger" v-auth="auth?auth.delete:auth" plain icon="el-icon-delete" :disabled="selection.length==0" @click="batch_del"></el-button>
			</div>
			<div class="right-panel">
				<div class="right-panel-search">
					<el-input v-model="search.position_name" placeholder="岗位名称" clearable></el-input>
					<el-button type="primary"  v-auth="auth?auth.get:auth"  icon="el-icon-search" @click="upsearch"></el-button>
				</div>
			</div>
		</el-header>
		<el-main class="nopadding">
			<scTable ref="table" :apiObj="apiObj" row-key="id" @selection-change="selectionChange" stripe>
				<el-table-column type="selection" width="50"></el-table-column>
				<el-table-column label="岗位编码"  prop="position_code" width="200"></el-table-column>
				<el-table-column label="岗位名称"  prop="position_name" width="150"></el-table-column>
				<el-table-column label="岗位描述"  prop="position_comment" width="150"></el-table-column>
				<el-table-column label="部门名称"  prop="department_code" width="120" :formatter="formatter"></el-table-column>
				<el-table-column label="职务名称"  prop="job_name" width="120"></el-table-column>
				<el-table-column label="部门负责人" prop="is_head" width="100">
					<template #default="scope">
						<el-switch v-model="scope.row.is_head" @change="changeSwitch($event, scope.row, 'is_head')" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
					</template>
				</el-table-column>
				<el-table-column label="审批权限" prop="is_approve" width="100">
					<template #default="scope">
						<el-switch v-model="scope.row.is_approve" @change="changeSwitch($event, scope.row,'is_approve')" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
					</template>
				</el-table-column>
				<el-table-column label="虚拟岗位" prop="is_vritual" width="100">
					<template #default="scope">
						<el-switch v-model="scope.row.is_vritual" @change="changeSwitch($event, scope.row,'is_vritual')" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
					</template>
				</el-table-column>
				<el-table-column label="状态" prop="status" width="80">
					<template #default="scope">
						<el-switch v-model="scope.row.status" @change="changeSwitch($event, scope.row,'status')" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
					</template>
				</el-table-column>
				<el-table-column label="创建时间" prop="created_at" :formatter="timestampToTime" width="180"></el-table-column>
				<el-table-column label="创建人" prop="created_by"  width="100"></el-table-column>
				<el-table-column label="更新时间" prop="updated_at" :formatter="timestampToTime" width="180"></el-table-column>
				<el-table-column label="更新人" prop="updated_by"  width="100"></el-table-column>
				<el-table-column label="操作" fixed="right" align="right" width="170">
					<template #default="scope">
						<el-button-group>
							<el-button text type="primary"  v-auth="auth?auth.post:auth"  size="small" @click="table_edit(scope.row, scope.$index)">编辑</el-button>
							<el-popconfirm title="确定删除吗？"  @confirm="table_del(scope.row, scope.$index)">
								<template #reference>
									<el-button text type="primary"  v-auth="auth?auth.delete:auth"  size="small">删除</el-button>
								</template>
							</el-popconfirm>
						</el-button-group>
					</template>
				</el-table-column>

			</scTable>
		</el-main>
	</el-container>

	<save-dialog v-if="dialog.save" ref="saveDialog" @success="handleSuccess" @closed="dialog.save=false"></save-dialog>

</template>


<script>
	import saveDialog from './save'

	export default {
		name: 'user',
		components: {
			saveDialog
		},
		data() {
			return {
				dialog: {
					save: false
				},
				showGrouploading: false,
				group: [],
				apiObj: this.$API.basis_org.position.list,
				selection: [],
				search: {
					company_name: null
				},
				auth:{
					put: '',
					post: '',
					get: '',
					delete: '',
				},
				formatData:{
					department_code:{}
				}
			}
		},
		// watch: {

		// },
		mounted() {
			this.getUiPermission()
			this.getDepartments()
		},
		methods: {
			//添加
			add(){
				this.dialog.save = true
				this.$nextTick(() => {
					this.$refs.saveDialog.open()
					this.$refs.saveDialog.getDepartmentTree()
				})
			},
			//编辑
			table_edit(row){
				this.dialog.save = true
				this.$nextTick(() => {
					this.$refs.saveDialog.open('edit').setData(row)
					this.$refs.saveDialog.getDepartmentTree()
				})
			},
			//查看
			table_show(row){
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
				await this.$API.basis_org.position.delete(reqData);
				this.$refs.table.tableData.splice(index, 1)
				setTimeout(()=>{
						this.$refs.table.reload()
					}, 500)
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
					this.$API.basis_org.position.delete({id: reqData})
					loading.close();
					setTimeout(()=>{
						this.$refs.table.reload()
					}, 500)
					
				}).catch(() => {
				})
			},

			//表格选择后回调事件
			selectionChange(selection){
				this.selection = selection;
			},

			//搜索
			upsearch(){
				this.$refs.table.upData(this.search)
			},
			// //本地更新数据
			handleSuccess(data, mode){
				setTimeout(()=>{
						this.$refs.table.reload()
					}, 500)
			},
			// 时间序列化
			timestampToTime (row, column) {
				var date = new Date(row[column.property]) //时间戳为10位需*1000，时间戳为13位的话不需乘1000
				var Y = date.getFullYear() + '-'
				var M =  String(date.getMonth()+1).padStart(2,"0") + '-'
				var D = String(date.getDate()).padStart(2,"0") + ' '
				var h = String(date.getHours()).padStart(2,"0") + ':'
				var m = String(date.getMinutes()).padStart(2,"0") + ':'
				var s = String(date.getSeconds()).padStart(2,"0")
				return Y+M+D+h+m+s
			},
			//表格内开关
			changeSwitch(val,row,column){
				console.log(column)
				console.log(row[column])
				// console.log(column)
				row[column] = row[column] == 1 ? 1 : -1
				row.$switch_status = true;
				// setTimeout(()=>{
					var data = {id:[row.id]}
					data[column] = row[column]
					this.$API.basis_org.position.edit(data);
					// }, 1000)
					delete row.$switch_status;
			},
			async getUiPermission(){
				var path = this.$router.currentRoute.value.fullPath
				var res = await this.$API.common.uiPermissions.get({path:path})
				this.auth = res.details
			},
			async getDepartments(){
				var res = await this.$API.common.dimension.get({dim_type:"department"});
				res.details.forEach(item=>{
					this.formatData.department_code[item.code] =item.name
				})
			},
			formatter(row,column,cols){
				var key = column.property
				var data = []
				var filterData = this.formatData[key]
				if (!cols || !filterData){
					return 
				}
				if (typeof cols == "string"){
					return filterData[cols]
				}
				cols.forEach(item =>{
					var dim = filterData[item]
					if (dim){
						data.unshift(dim)
					}
					
				})
				return data.join(", ")
			},

		}
	}
</script>

<style>
</style>
