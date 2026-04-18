<template>
	<div class="system-edit-dic-container">
		<el-dialog :title="'修改设置'" v-model="isShowDialog" width="769px">
			<el-form :model="ruleForm" ref="formRef" :rules="rules" size="default" label-width="90px">
        <el-form-item label="设置名称:" prop="name">
          {{ruleForm.name}}
        </el-form-item>
        <el-form-item label="值" prop="value">
          <el-input v-model="ruleForm.value" type="textarea" placeholder="请输入值"></el-input>
        </el-form-item>
			</el-form>
			<template #footer>
				<span class="dialog-footer">
					<el-button @click="onCancel" size="default">取 消</el-button>
					<el-button type="primary" @click="onSubmit" size="default">修 改</el-button>
				</span>
			</template>
		</el-dialog>
	</div>
</template>

<script lang="ts">
import { reactive, toRefs, defineComponent,ref, unref } from 'vue';
import { updateSetGame } from '/@/api/set/set_game';
import {ElMessage} from "element-plus";
interface RuleFormState {
  id:number;
  key:string;
  name:string;
  value:string;
}
interface DState {
	isShowDialog: boolean;
	ruleForm: RuleFormState;
  rules:{}
}

export default defineComponent({
	name: 'EditSet',
	setup(prop,{emit}) {
    const formRef = ref<HTMLElement | null>(null);
		const state = reactive<DState>({
			isShowDialog: false,
			ruleForm: {
        id:0,
        name:'',
        key:'',
        value:'',
			},
      rules: {
        value: [
          { required: true, message: "值不能为空", trigger: "blur" }
        ],
      }
		});
		// 打开弹窗
		const openDialog = (row: RuleFormState|null) => {
      resetForm();
      if (row){
        state.ruleForm = row;
      }
			state.isShowDialog = true;
		};
    const resetForm = ()=>{
      state.ruleForm = {
        id:0,
        name:'',
        key:'',
        value:'',
      }
    };
		// 关闭弹窗
		const closeDialog = () => {
			state.isShowDialog = false;
		};
		// 取消
		const onCancel = () => {
			closeDialog();
		};
		// 新增
		const onSubmit = () => {
      const formWrap = unref(formRef) as any;
      if (!formWrap) return;
      formWrap.validate((valid: boolean) => {
        if (valid) {
          if(state.ruleForm.id!==0){
            //修改
            updateSetGame(state.ruleForm).then(()=>{
              ElMessage.success('设置修改成功');
              closeDialog(); // 关闭弹窗
              emit('setList')
            })
          }
        }
      });
		};


		return {
			openDialog,
			closeDialog,
			onCancel,
			onSubmit,
      formRef,
			...toRefs(state),
		};
	},
});
</script>
