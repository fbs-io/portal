<!--
 * @Author: reel
 * @Date: 2023-08-31 21:51:57
 * @LastEditors: reel
 * @LastEditTime: 2023-12-31 21:22:34
 * @Description: 请填写简介
-->
<template>
	<el-container>
		<el-header>
			<div class="left-panel">
				<el-button type="primary" v-auth="auth?auth.put:auth" icon="el-icon-plus" @click="add"></el-button>
				<el-button type="danger" v-auth="auth?auth.delete:auth" plain icon="el-icon-delete" :disabled="selection.length==0" @click="batch_del"></el-button>
			</div>
			<div class="right-panel">
				<div class="right-panel-search">
					<el-input v-model="search.invoice_code" placeholder="发票代码" clearable></el-input>
					<el-button type="primary"  v-auth="auth?auth.get:auth"  icon="el-icon-search" @click="upsearch"></el-button>
				</div>
			</div>
		</el-header>
		<el-main class="nopadding">
			<scTable ref="table" :apiObj="apiObj" row-key="id" @selection-change="selectionChange" stripe>
				<el-table-column type="selection" width="50"></el-table-column>
				<el-table-column label="#" porp="id" type="index" width="50"></el-table-column>
				<el-table-column label="发票号码" sortable='custom' prop="invoice_no" width="120"></el-table-column>
				<el-table-column label="发票代码" sortable='custom' prop="invoice_org_code" width="120"></el-table-column>
				<el-table-column label="发票类型" sortable='custom' prop="invoice_type" width="200"></el-table-column>
				<el-table-column label="开票时间" sortable='custom' prop="invoice_date" width="120"></el-table-column>
				<el-table-column label="发票金额" sortable='custom' prop="invoice_amount" width="150"></el-table-column>
				<el-table-column label="发票描述" sortable='custom' prop="invoice_commnet" width="120"></el-table-column>
				<el-table-column label="状态" prop="status" width="80">
					<template #default="scope">
						<el-switch v-model="scope.row.status" @change="changeSwitch($event, scope.row)" :loading="scope.row.$switch_status" :active-value="1" :inactive-value="-1"></el-switch>
					</template>
				</el-table-column>
				<el-table-column label="创建时间" prop="created_at" :formatter="timestampToTime" width="180"></el-table-column>
				<el-table-column label="创建人" prop="created_by"  width="100"></el-table-column>
				<el-table-column label="更新时间" prop="updated_at" :formatter="timestampToTime" width="180"></el-table-column>
				<el-table-column label="更新人" prop="updated_by"  width="100"></el-table-column>
				<el-table-column label="操作" fixed="right" align="right" width="170">
					<template #default="scope">
						<el-button-group>
							<el-button text type="primary"  v-auth="auth?auth.get:auth" size="small" @click="table_show(scope.row, scope.$index)">查看</el-button>
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
				groupFilterText: '',
				group: [],
				apiObj: this.$API.financial.invoice.list,
				selection: [],
				search: {
					company_name: null
				},
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
				await this.$API.financial.invoice.delete(reqData);
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
					this.$API.financial.invoice.delete({id: reqData})
					loading.close();
					this.$refs.table.reload()
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
				if(mode=='add'){
					this.$refs.table.reload()
				}else if(mode=='edit'){
					this.$refs.table.tableData.filter(item => item.id===data.id ).forEach(item => {
						Object.assign(item, data)
					})
				}
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
			changeSwitch(val, row){
				row.status = row.status == 1 ? 1 : -1
				row.$switch_status = true;
				setTimeout(()=>{
					delete row.$switch_status;
					this.$API.financial.invoice.edit({id:[row.id],status:row.status});
				}, 1000)
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
