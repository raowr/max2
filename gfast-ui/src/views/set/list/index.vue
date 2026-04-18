<template>
	<div class="system-dic-container">
		<el-card shadow="hover">
			<div class="system-user-search mb15">
        <el-form :model="tableData.param" ref="queryRef" :inline="true" label-width="68px">
          <el-form-item label="设置名称" prop="name">
            <el-input
                v-model="tableData.param.name"
                placeholder="请输入设置名称"
                clearable
                size="default"
                style="width: 240px"
                @keyup.enter.native="setGameList"
            />
          </el-form-item>
          <el-form-item>
            <el-button size="default" type="primary" class="ml10" @click="setGameList">
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
        <el-table-column label="设置名称" align="center" prop="name"/>
        <el-table-column label="key" align="center" prop="key" />
        <el-table-column label="值" align="center" prop="value"/>
        <el-table-column label="操作" width="200">
					<template #default="scope">
						<el-button size="small" text type="primary" @click="onOpenEditSet(scope.row)">修改</el-button>
					</template>
				</el-table-column>
			</el-table>
      <pagination
          v-show="tableData.total>0"
          :total="tableData.total"
          v-model:page="tableData.param.pageNum"
          v-model:limit="tableData.param.pageSize"
          @pagination="setGameList"
      />
		</el-card>
		<EditSet ref="editSetRef" @setGameList="setGameList"/>
	</div>
</template>

<script lang="ts">
import { toRefs, reactive, onMounted, ref, defineComponent } from 'vue';
import { FormInstance} from 'element-plus';
import EditSet from '/@/views/set/list/component/editSet.vue';
import { getSetGameList} from "/@/api/set/set_game";


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
      key: string;
      value: string;
		};
	};
}

export default defineComponent({
	name: 'systemSet',
	components: { EditSet },
	setup() {
		const addDicRef = ref();
		const editSetRef = ref();
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
          key:'',
          value:'',
			},
    }
		});
		// 初始化表格数据
		const initTableData = () => {
      setGameList()
		};
    const setGameList=()=>{
      getSetGameList(state.tableData.param).then((res:any)=>{
        console.log(res.data)
        state.tableData.data = res.data.list;
        state.tableData.total = res.data.total;
      });
    };

		// 页面加载时
		onMounted(() => {
			initTableData();
		});
    // 打开修改设置弹窗
		const onOpenEditSet = (row: TableDataRow) => {
			editSetRef.value.openDialog(row);
		};
    /** 重置按钮操作 */
    const resetQuery = (formEl: FormInstance | undefined) => {
      if (!formEl) return
      formEl.resetFields()
      setGameList()
    };
		return {
			addDicRef,
			editSetRef,
      queryRef,
      setGameList,
      onOpenEditSet,
      resetQuery,
			...toRefs(state),
		};
	},
});
</script>
