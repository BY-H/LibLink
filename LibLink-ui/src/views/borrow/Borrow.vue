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
                    <el-form 
                        ref="formRef"
                        :model="borrowForm"
                        :rules="rules"
                        label-width="120px"
                    >
                        <el-form-item label="文献分类" prop="category">
                            <el-tree-select
                                v-model="borrowForm.category"
                                :data="treeData"
                                placeholder="请选择文献分类"
                                check-strictly
                                :render-after-expand="false"
                            />
                        </el-form-item>
                        <el-form-item label="借阅备注" prop="remark">
                            <el-input 
                                v-model="borrowForm.remark"
                                type="textarea"
                                placeholder="请输入借阅备注（选填）"
                                :rows="3"
                            />
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="submitForm(formRef)">
                                确认借阅
                            </el-button>
                            <el-button @click="resetForm(formRef)">
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
                        :closable="true"
                    />
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Reading, DocumentAdd, RefreshRight } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

interface TreeNode {
    value: string
    label: string
    children?: TreeNode[]
}

const formRef = ref<FormInstance>()
const message = ref('')
const messageType = ref('success')

const borrowForm = reactive({
    category: '',
    archiveId: '',
    duration: '',
    remark: ''
})

const treeData: TreeNode[] = [
    {
        value: '1',
        label: '计算机科学',
        children: [
            {
                value: '1-1',
                label: '软件工程',
                children: [
                    { value: '1-1-1', label: '软件测试' },
                    { value: '1-1-2', label: '软件架构' }
                ]
            },
            {
                value: '1-2',
                label: '人工智能',
                children: [
                    { value: '1-2-1', label: '机器学习' },
                    { value: '1-2-2', label: '深度学习' }
                ]
            }
        ]
    },
    {
        value: '2',
        label: '数学',
        children: [
            {
                value: '2-1',
                label: '基础数学',
                children: [
                    { value: '2-1-1', label: '微积分' },
                    { value: '2-1-2', label: '线性代数' }
                ]
            },
            {
                value: '2-2',
                label: '应用数学',
                children: [
                    { value: '2-2-1', label: '概率论' },
                    { value: '2-2-2', label: '数理统计' }
                ]
            }
        ]
    }
]

const rules = reactive<FormRules>({
    category: [
        { required: true, message: '请选择文献分类', trigger: 'change' }
    ],
    archiveId: [
        { required: true, message: '请输入文献编号', trigger: 'blur' },
        { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
    ],
    duration: [
        { required: true, message: '请选择借阅时长', trigger: 'change' }
    ]
})

const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return
    await formEl.validate((valid) => {
        if (valid) {
            // TODO: 调用借阅 API
            message.value = `已成功借阅 ${getCategoryLabel(borrowForm.category)} 分类下的文献：${borrowForm.archiveId}，借阅时长：${borrowForm.duration}天`
            messageType.value = 'success'
            resetForm(formEl)
        } else {
            message.value = '请完善借阅信息'
            messageType.value = 'error'
        }
    })
}

const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
    message.value = ''
}

const getCategoryLabel = (value: string): string => {
    const findLabel = (nodes: TreeNode[], targetValue: string): string => {
        for (const node of nodes) {
            if (node.value === targetValue) return node.label
            if (node.children) {
                const label = findLabel(node.children, targetValue)
                if (label) return label
            }
        }
        return ''
    }
    return findLabel(treeData, value)
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
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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

:deep(.el-input), :deep(.el-tree-select), :deep(.el-select) {
    max-width: 400px;
    width: 100%;
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