<template>
	<div class="system-dic-container">
		<el-card shadow="hover">
			<div class="system-user-search mb15">
        <el-form :model="tableData.param" ref="queryRef" :inline="true" label-width="68px">
          <el-form-item label="房间ID" prop="roomId">
            <el-input
                v-model="tableData.param.roomId"
                placeholder="请输入房间ID"
                clearable
                size="default"
                style="width: 240px"
                @keyup.enter.native="logGameList"
            />
          </el-form-item>
          <el-form-item label="用户ID" prop="userId">
            <el-input
                v-model="tableData.param.userId"
                placeholder="请输入用户ID"
                clearable
                size="default"
                style="width: 240px"
                @keyup.enter.native="logGameList"
            />
          </el-form-item>
          <el-form-item label="房间状态" prop="status" style="width: 200px;">
            <el-select
                v-model="tableData.param.status"
                placeholder="房间状态"
                clearable
                size="default"
                style="width: 240px"
            >
              <el-option label="未开始"  :value="0"/>
              <el-option label="进行中"  :value="1"/>
              <el-option label="已结算"  :value="2"/>
            </el-select>
          </el-form-item>
          <el-form-item label="房间类型" prop="status" style="width: 200px;">
            <el-select
                v-model="tableData.param.roomType"
                placeholder="房间类型"
                clearable
                size="default"
                style="width: 240px"
            >
              <el-option label="段位房"  :value="1"/>
              <el-option label="好友房"  :value="2"/>
            </el-select>
          </el-form-item>
          <el-form-item label="记录时间" prop="dateRange">
            <el-date-picker
                v-model="tableData.param.dateRange"
                size="default"
                style="width: 240px"
                value-format="YYYY-MM-DD"
                type="daterange"
                range-separator="-"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
            ></el-date-picker>
          </el-form-item>
          <el-form-item>
            <el-button size="default" type="primary" class="ml10" @click="logGameList">
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
        <el-table-column label="房间id" align="center" prop="roomId"/>
        <el-table-column label="房间类型" align="center" prop="type" >
        					<template #default="scope">
            <el-tag v-if="scope.row.type==1">段位房</el-tag>
						<el-tag v-else-if="scope.row.type==2">好友房</el-tag>
					</template>
        </el-table-column>
        <el-table-column label="房间状态" align="center" prop="status">
					<template #default="scope">
            <el-tag type="info" v-if="scope.row.status==0">未开始</el-tag>
						<el-tag type="success" v-else-if="scope.row.status==1">进行中</el-tag>
						<el-tag type="danger" v-else>已结算</el-tag>
					</template>
        </el-table-column>
				<el-table-column prop="userId" label="用户id" show-overflow-tooltip />
				<el-table-column prop="point" label="用户分数" show-overflow-tooltip></el-table-column>
        <el-table-column prop="action" label="行为" show-overflow-tooltip></el-table-column>
        <el-table-column prop="remain" label="剩余牌" show-overflow-tooltip></el-table-column>
        <el-table-column prop="outCards" label="出牌" show-overflow-tooltip></el-table-column>
				<el-table-column prop="creatTime" label="创建时间" show-overflow-tooltip></el-table-column>
			</el-table>
      <pagination
          v-show="tableData.total>0"
          :total="tableData.total"
          v-model:page="tableData.param.pageNum"
          v-model:limit="tableData.param.pageSize"
          @pagination="logGameList"
      />
		</el-card>
		<EditDic ref="editDicRef" @logGameList="logGameList"/>
	</div>
</template>

<script lang="ts">
import { toRefs, reactive, onMounted, ref, defineComponent } from 'vue';
import { FormInstance} from 'element-plus';
import EditDic from '/@/views/system/dict/component/editDic.vue';
import { getLogGameList} from "/@/api/log/log_game";
import { log } from 'console';


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
      roomId: string;
      userId: string;
      status: string;
      roomType: string;
      dateRange:string[];
		};
	};
}

export default defineComponent({
	name: 'systemDic',
	components: { EditDic },
	setup() {
		const addDicRef = ref();
		const editDicRef = ref();
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
          roomId:'',
          userId:'',
          status:'',
          roomType:'',
          dateRange:[],
				},
			},
		});
		// 初始化表格数据
		const initTableData = () => {
      logGameList()
		};
    const logGameList=()=>{
      getLogGameList(state.tableData.param).then((res:any)=>{
        console.log(res.data)
        state.tableData.data = res.data.list;
        state.tableData.total = res.data.total;
      });
    };

		// 页面加载时
		onMounted(() => {
			initTableData();
		});
    /** 重置按钮操作 */
    const resetQuery = (formEl: FormInstance | undefined) => {
      if (!formEl) return
      formEl.resetFields()
      logGameList()
    };
		return {
			addDicRef,
			editDicRef,
      queryRef,
      logGameList,
      resetQuery,
			...toRefs(state),
		};
	},
});
</script>
