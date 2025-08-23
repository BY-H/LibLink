<template>
    <div class="borrow-wrapper">
        <el-card class="box-card">
            <template #header>
                <div class="card-header">
                    <el-icon><Reading /></el-icon>
                    <span>文献借阅</span>
                </div>
            </template>

            <!-- 搜索 + 新增操作区域 -->
            <div class="table-actions">
                <div class="search-form">
                    <el-input v-model="searchContractNo" placeholder="请输入合同编号进行搜索" clearable style="width: 300px" @keyup.enter="handleSearch">
                        <template #prefix>
                            <el-icon><Search /></el-icon>
                        </template>
                    </el-input>
                    <el-button type="primary" @click="handleSearch" :icon="Search"> 搜索 </el-button>
                    <el-button @click="resetSearch" :icon="Refresh"> 重置 </el-button>
                </div>

                <div class="right-actions">
                    <el-button type="primary" @click="fetchData"> 批量新增 </el-button>
                    <el-button type="primary" @click="openDrawer"> 新增 </el-button>
                </div>
            </div>

            <!-- 表格 -->
            <el-table :data="tableData" border style="width: 100%; margin-top: 15px">
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
                        <el-tag :type="Number(scope.row.borrow_state) == 1 ? 'success' : 'info'" effect="light">
                            {{ Number(scope.row.borrow_state) == 1 ? '已借阅' : '未借阅' }}
                        </el-tag>
                    </template>
                </el-table-column>

                <!-- 操作列 -->
                <el-table-column label="操作" width="180">
                    <template #default="scope">
                        <!-- <el-link type="primary" @click="handleDelete(scope.row)">删除</el-link> -->
                        <el-link v-if="scope.row.borrow_state == 0" type="primary" @click="handleBorrow(scope.row)" style="margin-left: 10px"> 借阅 </el-link>
                        <el-link v-else type="warning" @click="handleReturn(scope.row)" style="margin-left: 10px"> 归还 </el-link>
                    </template>
                </el-table-column>
            </el-table>

            <!-- 分页 -->
            <Pagination v-model:pageObj="pageObj" v-model:page-size="pageObj.page_size" :total="total" :onUpdate="fetchData" @update:page-size="handlePageSizeChange" />
        </el-card>

        <!-- 新增抽屉组件 -->
        <BorrowDraw ref="drawerRef" @submit="handleAdd" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed } from 'vue'
import { Reading, Search, Refresh } from '@element-plus/icons-vue'
import Pagination from '@/components/Pagination.vue'
import BorrowDraw from './components/BorrowDraw.vue'
import { getArchives, addArchive, borrowArchive, returnArchive } from '@/api/archives'

// 搜索关键词
const searchContractNo = ref('')

// 搜索表单
const searchForm = reactive({
    file_no: '',
    contract_no: '',
    arc_type: '',
    inst_no: '',
    name: '',
    borrow_state: ''
})

// 表格数据
const tableData = ref<any[]>([])

// 分页参数
const pageObj = ref({
    page: 1,
    page_size: 10
})
const total = ref(0)

const filteredData = computed(() => {
    if (!searchContractNo.value) {
        return tableData.value
    }

    return tableData.value.filter(item => item.contract_no && item.contract_no.includes(searchContractNo.value))
})

const fetchData = async () => {
    try {
        // 构建查询参数
        const params: any = {
            page: pageObj.value.page,
            page_size: pageObj.value.page_size
        }

        // 添加搜索参数
        if (searchContractNo.value.trim()) {
            params.contract_no = searchContractNo.value.trim()
        }

        console.log('请求参数:', params)

        const response: any = await getArchives(params)

        if (response && response.data) {
            tableData.value = response.data.map((item: any) => ({
                ...item,
                borrow_state: Number(item.borrow_state) || 0
            }))
            total.value = response.total || response.data.length
            console.log('获取数据成功:', tableData.value)
        } else {
            console.error('响应数据格式错误:', response)
        }
    } catch (error: any) {
        console.error('获取数据失败：', error)
    }
}

// 搜索处理
const handleSearch = () => {
    pageObj.value.page = 1 // 重置到第一页
    fetchData()
}

// 重置搜索
const resetSearch = () => {
    searchContractNo.value = ''
    pageObj.value.page = 1
    fetchData()
}

// 每页条数变化
const handlePageSizeChange = (newSize: number) => {
    pageObj.value.page_size = newSize
    pageObj.value.page = 1
    fetchData()
}

const handleDelete = (row: any) => {
    console.log('删除', row)
}

const handleBorrow = (row: any) => {
    if (row.borrow_state === 1) return
    try {
        borrowArchive({
            contract_no: row.contract_no
        })
        row.borrow_state = '1' // 更新状态为已借阅
    } catch (error) {
        console.error('借阅失败：', error)
    }
}

const handleReturn = (row: any) => {
    if (row.borrow_state === 0) return
    try {
        returnArchive({
            contract_no: row.contract_no
        })
        row.borrow_state = '0' // 更新状态为未借阅
    } catch (error) {
        console.error('归还失败：', error)
    }
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
        console.log('新增档案成功：', data)
        fetchData()
    } catch (error) {
        console.error('新增档案失败：', error)
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

.table-actions {
    display: flex;
    align-items: center;
    justify-content: space-between; /* 左右分开 */
    margin-bottom: 15px;
}

.right-actions {
    display: flex;
    gap: 10px; /* 两个按钮之间的间距 */
}

.search-form {
    display: flex;
    align-items: center;
    gap: 10px; /* 按钮之间的间距 */
}
</style>
