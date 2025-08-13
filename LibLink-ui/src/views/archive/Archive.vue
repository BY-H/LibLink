<template>
    <div class="archive-wrapper">
        <el-row :gutter="20">
            <el-col :span="24">
                <el-card class="box-card">
                    <template #header>
                        <div class="card-header">
                            <el-icon><Folder /></el-icon>
                            <span>文献入库</span>
                        </div>
                    </template>
                    <el-form 
                        ref="formRef"
                        :model="archiveForm"
                        :rules="rules"
                        label-width="120px"
                        @submit.prevent="handleSubmit"
                    >
                        <el-form-item label="档案编号" prop="fileNo">
                            <el-input 
                                v-model="archiveForm.fileNo"
                                placeholder="请输入档案编号"
                            />
                        </el-form-item>
                        
                        <el-form-item label="合同编号" prop="contractNo">
                            <el-input 
                                v-model="archiveForm.contractNo"
                                placeholder="请输入合同编号"
                            />
                        </el-form-item>

                        <el-form-item label="归属网点" prop="instNo">
                            <el-input 
                                v-model="archiveForm.instNo"
                                placeholder="请输入归属网点"
                            />
                        </el-form-item>

                        <el-form-item label="档案类型" prop="arcType">
                            <el-select 
                                v-model="archiveForm.arcType"
                                placeholder="请选择档案类型"
                                style="width: 100%"
                            >
                                <el-option label="图书" value="BOOK" />
                                <el-option label="期刊" value="JOURNAL" />
                                <el-option label="论文" value="THESIS" />
                                <el-option label="其他" value="OTHER" />
                            </el-select>
                        </el-form-item>

                        <el-form-item>
                            <el-button type="primary" @click="handleSubmit">
                                入库
                            </el-button>
                            <el-button @click="resetForm">
                                <el-icon><RefreshRight /></el-icon>
                                重置
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
                        :type="messageType"
                        show-icon
                        :closable="false"
                    />
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Folder, DocumentAdd, RefreshRight } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const formRef = ref<FormInstance>()
const message = ref('')
const messageType = ref('success')

const archiveForm = reactive({
    fileNo: '',
    contractNo: '',
    instNo: '',
    arcType: ''
})

const rules = reactive<FormRules>({
    fileNo: [
        { required: true, message: '请输入档案编号', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
    ],
    contractNo: [
        { required: true, message: '请输入合同编号', trigger: 'blur' },
        { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
    ],
    instNo: [
        { required: true, message: '请输入归属网点', trigger: 'blur' },
        { min: 3, max: 50, message: '长度在 3 到 50 个字符', trigger: 'blur' }
    ],
    arcType: [
        { required: true, message: '请选择档案类型', trigger: 'change' }
    ]
})

// 将原来的 handleSubmit 函数替换为：
const handleSubmit = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    try {
        await formEl.validate((valid) => {
            if (valid) {
                // TODO: 调用API保存数据
                message.value = '文献入库成功！'
                messageType.value = 'success'
                resetForm()
            } else {
                message.value = '请填写完整的表单信息！'
                messageType.value = 'error'
            }
        })
    } catch (error) {
        console.error('表单验证失败:', error)
        message.value = '表单验证失败'
        messageType.value = 'error'
    }
}

const resetForm = () => {
    if (!formRef.value) return
    formRef.value.resetFields()
    message.value = ''
}
</script>

<style scoped>
.archive-wrapper {
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

:deep(.el-input), :deep(.el-select) {
    max-width: 400px;
}

:deep(.el-button) {
    display: flex;
    align-items: center;
    gap: 5px;
}

:deep(.el-button + .el-button) {
    margin-left: 10px;
}
</style>