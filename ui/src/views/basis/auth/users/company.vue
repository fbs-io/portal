<!--
 * @Author: reel
 * @Date: 2023-10-05 14:05:58
 * @LastEditors: reel
 * @LastEditTime: 2023-10-06 07:52:27
 * @Description: 分配给用户公司权限
-->
<template>
	<el-dialog :title="'分配公司'" v-model="visible" :width="500" destroy-on-close @closed="$emit('closed')">
		<el-form :model="form"  ref="dialogForm" label-width="100px" label-position="left">
			<el-form-item label="所属公司" prop="companies">
				<el-select v-model="form.company" multiple filterable style="width: 100%" :props="companiesProps">
					<el-option v-for="item in companies" :key="item.company_code" :label="item.company_name" :value="item.company_code" />
				</el-select>
			</el-form-item>
		</el-form>
		<template #footer>
			<el-button @click="visible=false" >取 消</el-button>
			<el-button type="primary" :loading="isSaveing" @click="submit(false)">保 存</el-button>
		</template>
	</el-dialog>
</template>

<script>
    export default {
        emits: ['success', 'closed'],
        data() {
            return {
                visible: false,
                isSaveing: false,
                //表单数据
                form: {
                    id: [],
                    company: [],
                },
                //所需数据选项
                companies: [],
                companiesProps:{
                    value:"company_code",
                    multiple: true,
                    checkStrictly: true
                }
            }
        },
        mounted() {
            this.getData()
        },
        methods: {
            //显示
            open(ids){
                this.form.id = ids;
                this.visible = true;
                return this
            },
            //加载树数据
            async getData(){
                var res = await this.$API.basis_org.company.list({page_num:-1,page_size:-1});
                this.companies = res.details.rows;
            },
            //表单提交方法
            async submit(){
                    this.isSaveing = true;
                    var res = await this.$API.basis_auth.users.edit(this.form);
                    this.isSaveing = false;
                    if(res.errno == 0){	
                        this.visible = false;
                        this.$emit("success",this.form,this.mode)
                    }else{
                        retrun
                    }
                
                }
        },
        //表单注入数据
        setData(data){
            //可以和上面一样单个注入，也可以像下面一样直接合并进去
            Object.assign(this.form, data)
        }
        
    }
</script>

<style>
</style>
