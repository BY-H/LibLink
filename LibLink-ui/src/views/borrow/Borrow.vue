<template>
    <div class="borrow-wrapper">
        <el-row :gutter="20">
            <el-col :span="24">
                <el-card class="box-card">
                    <template #header>
                        <div class="card-header">
                            <el-icon><Reading /></el-icon>
                            <span>文献借阅</span>
                        </div>
                    </template>
                    <el-form @submit.prevent="handleBorrow" label-width="120px">
                        <el-form-item label="文献编号">
                            <el-input v-model="archiveID" placeholder="请输入文献编号" />
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="handleBorrow">
                                借阅
                            </el-button>
                        </el-form-item>
                    </el-form>
                </el-card>
            </el-col>
        </el-row>
        
        <el-row v-if="message" :gutter="20" class="message-row">
            <el-col :span="24">
                <el-card class="box-card">
                    <el-alert
                        :title="message"
                        type="success"
                        show-icon
                        :closable="false"
                    />
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { Reading, DocumentAdd } from '@element-plus/icons-vue'

// 响应式状态
const archiveID = ref('')
const message = ref('')

// 方法
const handleBorrow = () => {
    if (archiveID.value) {
        // 模拟借阅逻辑
        message.value = `文献编号 ${archiveID.value} 已成功借阅。`
        archiveID.value = ''
    } else {
        message.value = '请填写文献编号。'
    }
}
</script>

<style scoped>
.borrow-wrapper {
    padding: 20px;
    background-color: #f5f5f5;
}

.box-card {
    margin-bottom: 20px;
    border-radius: 10px;
    overflow: hidden;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
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

.message-row {
    margin-top: 20px;
}

:deep(.el-form-item__label) {
    font-weight: bold;
    color: #333;
}

:deep(.el-input) {
    max-width: 400px;
}

:deep(.el-button) {
    display: flex;
    align-items: center;
    gap: 5px;
}
</style>