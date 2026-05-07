<template>
	<div class="system-dic-container">
		<el-card shadow="hover">
			<div class="system-user-search mb15">
        <el-form :model="tableData.param" ref="queryRef" :inline="true" label-width="68px">
          <el-form-item label="用户名称" prop="name">
            <el-input
                v-model="tableData.param.name"
                placeholder="用户名称"
                clearable
                size="default"
                style="width: 240px"
                @keyup.enter.native="userList"
            />
          </el-form-item>
          <el-form-item>
            <el-button size="default" type="primary" class="ml10" @click="userList">
              <el-icon>
                <ele-Search />
              </el-icon>
              查询
            </el-button>
            <el-button size="default" @click="resetQuery(queryRef)">
              <el-icon>
                <ele-Refresh />
              </el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
			</div>
			<el-table :data="tableData.data" style="width: 100%" >
        <el-table-column label="id" align="center" prop="id"/>
        <el-table-column label="用户名称" align="center" prop="name" />
        <el-table-column label="积分" align="center" prop="point"/>
        <el-table-column label="操作" width="200">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="onOpenEditPoint(scope.row)">修改积分</el-button>
            <el-button size="small" text type="primary" @click="onOpenEditPassword(scope.row)">重置密码</el-button>
					</template>
				</el-table-column>
			</el-table>
      <pagination
          v-show="tableData.total>0"
          :total="tableData.total"
          v-model:page="tableData.param.pageNum"
          v-model:limit="tableData.param.pageSize"
          @pagination="getUserList"
      />
		</el-card>
		<EditPoint ref="editPointRef" @getUserList="getUserList"/>
    <EditPassword ref="editPasswordRef" @getUserList="getUserList"/>
	</div>
</template>

<script lang="ts">
import { toRefs, reactive, onMounted, ref, defineComponent } from 'vue';
import { FormInstance} from 'element-plus';
import EditPoint from '/@/views/user/list/component/editPoint.vue';
import EditPassword from '/@/views/user/list/component/editPassword.vue';
import { getUserList, updateUserPoint, updateUserPassword} from "/@/api/user/user";


// 定义接口来定义对象的类型
interface TableDataRow {
  id:number;
  roomId: string;
  roomType: string;
  status: number;
  userId: string;
  point: number;
  action: string;
  remain: number;
  output: string;
  creatTime: string;
  password:string;
}
interface TableDataState {
  ids:number[];
	tableData: {
		data: Array<TableDataRow>;
		total: number;
		loading: boolean;
		param: {
			pageNum: number;
			pageSize: number;
      name: string;
		};
	};
}

export default defineComponent({
	name: 'User',
	components: { EditPoint, EditPassword },
	setup() {
		const addDicRef = ref();
		const editPointRef = ref();
    const editPasswordRef = ref();
    const queryRef = ref();
		const state = reactive<TableDataState>({
      ids:[],
			tableData: {
				data: [],
				total: 0,
				loading: false,
				param: {
					pageNum: 1,
					pageSize: 10,
          name:'',
			},
    }
		});
		// 初始化表格数据
		const initTableData = () => {
      userList()
		};
    const userList=()=>{
      getUserList(state.tableData.param).then((res:any)=>{
        console.log(res.data)
        state.tableData.data = res.data.list;
        state.tableData.total = res.data.total;
      });
    };

		// 页面加载时
		onMounted(() => {
			initTableData();
		});
    // 打开修改积分弹窗
		const onOpenEditPoint = (row: TableDataRow) => {
			editPointRef.value.openDialog(row);
		};
    // 打开重置密码弹窗
		const onOpenEditPassword = (row: TableDataRow) => {
      row.password = ''
			editPasswordRef.value.openDialog(row);
		};
    /** 重置按钮操作 */
    const resetQuery = (formEl: FormInstance | undefined) => {
      if (!formEl) return
      formEl.resetFields()
      userList()
    };
		return {
			addDicRef,
      editPointRef,
      editPasswordRef,
      queryRef,
      userList,
      resetQuery,
      onOpenEditPoint,
      onOpenEditPassword,
			...toRefs(state),
		};
	},
});
</script>
