<!--
 * @Author: reel
 * @Date: 2023-09-06 20:38:09
 * @LastEditors: reel
 * @LastEditTime: 2023-09-07 06:15:26
 * @Description: 请填写简介
-->
<template>
	<el-dialog :title="'分配角色'" v-model="visible" :width="500" destroy-on-close @closed="$emit('closed')">
		<el-form :model="form"  ref="dialogForm" label-width="100px" label-position="left">
			<el-form-item label="所属角色" prop="role">
				<el-select v-model="form.role" multiple filterable style="width: 100%">
					<el-option v-for="item in roles" :key="item.code" :label="item.label" :value="item.code"/>
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
                isDefaultPwd:false,
                //表单数据
                form: {
                    id: [],
                    role: [],
                },
                //所需数据选项
                roles: [],
            }
        },
        mounted() {
            this.getRoles()
        },
        methods: {
            //显示
            open(ids){
                this.form.id = ids;
                this.visible = true;
                return this
            },
            //加载树数据
            async getRoles(){
                var res = await this.$API.basis_auth.roles.list();
                this.roles = res.details.rows;
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
