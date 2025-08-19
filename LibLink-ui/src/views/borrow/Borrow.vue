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
        <el-button type="primary" @click="openDrawer">
          新增
        </el-button>
      </div>

      <!-- 表格 -->
      <el-table :data="tableData" border style="width: 100%; margin-top: 15px;">
        <el-table-column prop="file_no" label="档案编号" />
        <el-table-column prop="arc_type" label="档案类型" />
        <el-table-column prop="contract_no" label="合同编号" />
        <el-table-column prop="name" label="姓名" />
        <el-table-column prop="id_card" label="身份证号" />
        <el-table-column prop="inst_no" label="网点编号" />
        <el-table-column prop="manager" label="客户经理" />
        <el-table-column prop="amount" label="合同金额" />
        <el-table-column prop="storage_date" label="入库日期" />

        <!-- 借阅状态 -->
        <el-table-column label="借阅状态">
          <template #default="scope">
            <el-tag
              :type="Number(scope.row.borrow_state) === 1 ? 'success' : 'info'"
              effect="light"
            >
              {{ Number(scope.row.borrow_state) === 1 ? '已借阅' : '未借阅' }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-link type="primary" @click="handleDelete(scope.row)">删除</el-link>
            <el-link
              type="primary"
              :disabled="Number(scope.row.borrow_state) === 1"
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

    <!-- 新增抽屉组件 -->
    <BorrowDraw ref="drawerRef" @submit="handleAdd" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { Reading } from "@element-plus/icons-vue"
import Pagination from "@/components/Pagination.vue"
import BorrowDraw from "./components/BorrowDraw.vue"
import { getArchives, addArchive } from "@/api/archives"

// 表格数据
const tableData = ref<any[]>([])

// 分页参数
const pageObj = ref({
  page: 1,
  page_size: 10,
})
const total = ref(0)

const fetchData = async () => {
  try {
    const response: any = await getArchives()
    // 处理 borrow_state 字段，确保是数字
    tableData.value = response.data.map((item: any) => ({
      ...item,
      borrow_state: Number(item.borrow_state),
    }))
    total.value = response.data.length // TODO: 后端分页后改成 response.total
    console.log("档案数据：", tableData.value)
  } catch (error) {
    console.error("获取数据失败：", error)
  }
}

const handleDelete = (row: any) => {
  console.log("删除", row)
}

const handleBorrow = (row: any) => {
  if (row.borrow_state === 1) return
  console.log("借阅", row)
}

const drawerRef = ref()

// 点击新增按钮，打开 Drawer
const openDrawer = () => {
  drawerRef.value.open()
}

// 接收新增数据
const handleAdd = async (data: any) => {
  try {
    await addArchive(data)
    console.log("新增档案成功：", data)
    fetchData()
  } catch (error) {
    console.error("新增档案失败：", error)
  }
}

onMounted(() => {
  // 初始化数据
  fetchData()
})
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
