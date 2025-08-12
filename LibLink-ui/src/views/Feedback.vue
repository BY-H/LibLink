<template>
    <div class="feedback-wrapper">
        <el-card class="box-card">
            <div slot="header" class="clearfix">
                <span>用户反馈</span>
            </div>
            <el-form :model="feedbackForm" :rules="rules" ref="feedbackFormRef" label-width="100px" class="feedback-form">
                <el-form-item label="反馈内容" prop="content">
                    <el-input type="textarea" v-model="feedbackForm.content" placeholder="请输入反馈内容"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="submitFeedback">提交反馈</el-button>
                    <el-button @click="resetForm">重置</el-button>
                </el-form-item>
            </el-form>
        </el-card>

        <el-card class="box-card" style="margin-top: 20px">
            <div slot="header" class="clearfix">
                <span>反馈记录</span>
            </div>
            <el-table :data="feedbackList" style="width: 100%">
                <el-table-column prop="content" label="反馈内容"></el-table-column>
                <el-table-column prop="created_at" label="提交时间" width="220"></el-table-column>
            </el-table>
        </el-card>
    </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { onMounted } from 'vue'
import { getFeedbacks, addFeedback } from '@/api/system'

interface Feedback {
    id: number
    content: string
    created_at: string
}

const feedbackForm = ref({
    content: ''
})

const feedbackList = ref<Feedback[]>([
    { id: 1, content: '系统很好用！', created_at: '2025-03-18 10:00' },
    { id: 2, content: '希望增加更多功能。', created_at: '2025-03-19 12:00' }
])

const getFeedbackList = async () => {
    try {
        const response: any = await getFeedbacks()
        feedbackList.value = response.list
    } catch (error) {
        console.error('获取反馈列表失败:', error)
    }
}

const AddFeedback = async (feedback: any) => {
    try {
        await addFeedback(feedback)
    } catch (error) {
        console.error('提交反馈失败:', error)
    } finally {
        getFeedbackList()
    }
}

const rules = {
    content: [{ required: true, message: '请输入反馈内容', trigger: 'blur' }]
}

const feedbackFormRef = ref()

const submitFeedback = () => {
    feedbackFormRef.value?.validate((valid: boolean) => {
        if (valid) {
            const newFeedback = {
                content: feedbackForm.value.content
            }
            AddFeedback(newFeedback)
            ElMessage.success('反馈提交成功！')
            resetForm()
        } else {
            ElMessage.error('请填写完整的反馈信息！')
        }
    })
}

const resetForm = () => {
    feedbackForm.value.content = ''
}

onMounted(() => {
    getFeedbackList()
})
</script>

<style scoped>
.feedback-wrapper {
    padding: 20px;
}

.feedback-form {
    margin-bottom: 20px;
}

.box-card {
    margin-bottom: 20px;
}
</style>
