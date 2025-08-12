<template>
    <el-drawer v-model="localVisible" title="新增试题" @close="handleClose" size="50%">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
            <el-form-item label="题目正文" prop="question">
                <el-input v-model="form.question" placeholder="请输入题目正文"></el-input>
            </el-form-item>
            <el-form-item label="题目类型" prop="questionType">
                <el-select v-model="form.question_type" placeholder="请选择题目类型">
                    <el-option label="单选题" value="single_choice"></el-option>
                    <el-option label="多选题" value="multiple_choice"></el-option>
                    <el-option label="判断题" value="true_false"></el-option>
                    <el-option label="填空题" value="fill_blank"></el-option>
                    <el-option label="简答题" value="short_answer"></el-option>
                </el-select>
            </el-form-item>
            <!-- 根据题目类型动态显示答案输入方式 -->
            <template v-if="form.question_type === 'single_choice' || form.question_type === 'multiple_choice'">
                <el-form-item label="选项" prop="options">
                    <el-input v-model="form.options" placeholder="请输入选项，用逗号分隔"></el-input>
                </el-form-item>
                <el-form-item label="答案" prop="answer">
                    <el-input v-model="form.answer" placeholder="请输入答案（单选题为一个答案，多选题用逗号分隔）"></el-input>
                </el-form-item>
            </template>
            <template v-else-if="form.question_type === 'true_false'">
                <el-form-item label="答案" prop="answer">
                    <el-radio-group v-model="form.answer">
                        <el-radio label="true">正确</el-radio>
                        <el-radio label="false">错误</el-radio>
                    </el-radio-group>
                </el-form-item>
            </template>
            <template v-else-if="form.question_type === 'fill_blank'">
                <el-form-item label="答案" prop="answer">
                    <el-input v-model="form.answer" placeholder="请输入填空题答案"></el-input>
                </el-form-item>
            </template>
            <template v-else-if="form.question_type === 'short_answer'">
                <el-form-item label="答案" prop="answer">
                    <el-input type="textarea" v-model="form.answer" placeholder="请输入简答题答案" rows="4"></el-input>
                </el-form-item>
            </template>
            <el-form-item label="难度" prop="hardLevel">
                <el-input-number v-model="form.hard_level" :min="1" :max="5" placeholder="难度等级"></el-input-number>
            </el-form-item>
            <el-form-item label="分值" prop="score">
                <el-input-number v-model="form.score" :min="1" placeholder="请输入分值"></el-input-number>
            </el-form-item>
            <el-form-item label="标签" prop="tag">
                <el-input v-model="form.tag" placeholder="请输入标签"></el-input>
            </el-form-item>
            <el-form-item>
                <el-button type="primary" @click="handleSubmit">提交</el-button>
                <el-button @click="handleClose">取消</el-button>
            </el-form-item>
        </el-form>
    </el-drawer>
</template>

<script lang="ts" setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { addQuestion } from '@/api/question'

const props = defineProps(['visible'])
const emit = defineEmits(['update:visible', 'submit'])

const localVisible = computed({
    get: () => props.visible,
    set: value => emit('update:visible', value)
})

const form = ref({
    question: '',
    question_type: '',
    options: '', // 用于单选题和多选题的选项
    answer: '',
    hard_level: 1,
    score: 1,
    tag: '',
    creator: ''
})

const rules = {
    question: [{ required: true, message: '题目正文不能为空', trigger: 'blur' }],
    question_type: [{ required: true, message: '请选择题目类型', trigger: 'change' }],
    answer: [{ required: true, message: '答案不能为空', trigger: 'blur' }],
    hard_level: [{ required: true, message: '难度不能为空', trigger: 'change' }],
    score: [{ required: true, message: '分值不能为空', trigger: 'change' }]
}

const formRef = ref()

function handleClose() {
    emit('update:visible', false)
}

function handleSubmit() {
    formRef.value.validate(async (valid: boolean) => {
        if (valid) {
            const questionData = {
                ...form.value
            }
            const response: any = await addQuestion(questionData)
            if (response.status == "ok") {
                ElMessage.success('试题添加成功')
                emit('submit', questionData)
                handleClose()
            } else {
                ElMessage.error('试题添加失败')
            }
        } else {
            ElMessage.error('请填写完整的试题信息')
            return false
        }
    })
}
</script>

<style scoped>
.el-drawer {
    padding: 20px;
}
</style>
