<template>
  <div class="borrow-wrapper">
    <el-card class="box-card">
      <!-- 卡片头部 -->
      <template #header>
        <div class="card-header">
          <el-icon><Reading /></el-icon>
          <span>文献借阅</span>
        </div>
      </template>

      <div class="table-actions">
        <el-button type="primary" @click="handleAdd">
          新增
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border style="width: 100%; margin-top: 15px;">
        <el-table-column prop="fileNo" label="档案编号" />
        <el-table-column prop="contractNo" label="合同编号" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="idCard" label="身份证号" />
        <el-table-column prop="branchNo" label="网点编号" />
        <el-table-column prop="manager" label="客户经理" />
        <el-table-column prop="amount" label="合同金额" />
        <el-table-column prop="storageDate" label="入库日期" />

        <!-- 借阅状态 -->
        <el-table-column label="借阅状态">
          <template #default="scope">
            <el-tag
              :type="scope.row.borrowStatus === 1 ? 'success' : 'info'"
              effect="light"
            >
              {{ scope.row.borrowStatus === 1 ? '已借阅' : '未借阅' }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-link type="primary" @click="handleDelete(scope.row)">删除</el-link>
            <el-link
              type="primary"
              :disabled="scope.row.borrowStatus === 1"
              @click="handleBorrow(scope.row)"
              style="margin-left: 10px"
            >
              借阅
            </el-link>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <Pagination
        v-model:pageObj="pageObj"
        :total="total"
        :onUpdate="fetchData"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { Plus } from "@element-plus/icons-vue"
import Pagination from "@/components/Pagination.vue"

// 表格数据
const tableData = ref([
  {
    fileNo: "FLOANA11",
    contractNo: "HT12345678",
    name: "张三",
    idCard: "350427200001013527",
    branchNo: "903091001",
    manager: "张国强",
    amount: 10000,
    storageDate: "20250801",
    borrowStatus: 1,
  },
  {
    fileNo: "FLOANA12",
    contractNo: "HT87654321",
    name: "李四",
    idCard: "350427199901015050",
    branchNo: "903091101",
    manager: "张志伟",
    amount: 100000,
    storageDate: "20241221",
    borrowStatus: 1,
  },
  {
    fileNo: "FLOANA13",
    contractNo: "HT13572468",
    name: "王五",
    idCard: "3504271980010102010",
    branchNo: "903091102",
    manager: "肖磊",
    amount: 1000000,
    storageDate: "20250731",
    borrowStatus: 0,
  },
])

// 分页参数
const pageObj = ref({
  page: 1,
  page_size: 10,
})
const total = ref(3)

const fetchData = () => {
  console.log("分页参数：", pageObj.value)
}

const handleDelete = (row: any) => {
  console.log("删除", row)
}

const handleBorrow = (row: any) => {
  if (row.borrowStatus === 1) return
  console.log("借阅", row)
}

const handleAdd = () => {
  console.log("点击新增")
  // TODO: 打开新增对话框 / 跳转新增页面
}
</script>

<style scoped>
.borrow-wrapper {
  padding: 20px;
  background-color: #f5f5f5;
}

.box-card {
  border-radius: 10px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.el-icon {
  margin-right: 10px;
  color: #409eff;
}

.table-actions {
  justify-content: flex-end;
  display: flex;
  margin-bottom: 10px;
}
</style>
