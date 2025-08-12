<template>
    <div class="user-management">
        <el-card class="box-card">
            <div slot="header" class="clearfix">
                <span>用户管理</span>
                <el-button v-if="isAdmin" style="float: right" type="primary" @click="handleAddUser"> 添加用户 </el-button>
            </div>
            <el-table :data="filteredUsers" style="width: 100%">
                <el-table-column prop="id" label="用户ID" width="100"></el-table-column>
                <el-table-column prop="email" label="邮箱" width="250"></el-table-column>
                <el-table-column prop="role" label="角色" width="150"></el-table-column>
                <el-table-column prop="lastLogin" label="最近登录时间" width="200">
                    <template v-slot="scope">
                        <span>{{ formatDate(scope.row.lastLogin) }}</span>
                    </template>
                </el-table-column>
                <el-table-column v-if="isAdmin" label="操作" width="200">
                    <template v-slot="scope">
                        <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
                        <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
            <Pagination :pageObj="pageObj" :total="total" :onUpdate="fetchUsers" />
        </el-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Pagination from '@/components/Pagination.vue'

// 用户数据接口
interface User {
    id: number
    username: string
    email: string
    role: string
    lastLogin: string
}

// 当前登录用户信息（模拟从后端获取）
const currentUser = ref({
    id: 1,
    username: 'admin',
    role: 'admin' // 当前用户角色：'admin' 或 'user'
})

const users = ref<User[]>([]) // 用户列表
const total = ref(0) // 总用户数
const pageObj = ref({
    page: 1,
    page_size: 10
})

// 是否为管理员
const isAdmin = computed(() => currentUser.value.role === 'admin')

// 根据权限过滤用户列表
const filteredUsers = computed(() => {
    if (isAdmin.value) {
        return users.value // 管理员可以看到所有用户
    } else {
        return users.value.filter(user => user.id === currentUser.value.id) // 普通用户只能看到自己的信息
    }
})

// 格式化日期
function formatDate(date: string): string {
    const d = new Date(date)
    return d.toLocaleString()
}

// 获取用户列表
async function fetchUsers() {}

// 添加用户
function handleAddUser() {
    ElMessage.info('添加用户功能待实现')
}

// 编辑用户
function handleEdit(user: User) {
    ElMessage.info(`编辑用户：${user.username}`)
}

// 删除用户
function handleDelete(user: User) {
    ElMessageBox.confirm(`确定删除用户 ${user.username} 吗？`, '提示', {
        type: 'warning'
    })
        .then(() => {
            ElMessage.success('用户删除成功')
            // 调用删除 API 后刷新用户列表
            fetchUsers()
        })
        .catch(() => {
            ElMessage.info('已取消删除')
        })
}

// 初始化加载用户列表
onMounted(() => {
    fetchUsers()
})
</script>

<style scoped>
.user-management {
    padding: 20px;
}
</style>
