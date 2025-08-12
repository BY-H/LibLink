<template>
    <el-drawer v-model="localVisible" title="试卷详情" @close="handleClose" size="50%">
        <div class="paper-details">
            <h2>{{ paperData.title }}</h2>
            <p class="description">{{ paperData.description }}</p>
            <el-divider></el-divider>
            <div class="questions">
                <h3>试题内容</h3>
                <el-collapse>
                    <el-collapse-item 
                        v-for="(question, index) in questionList" 
                        :key="index" 
                        :title="'题目 ' + (index + 1)">
                        <!-- 单选题 -->
                        <template v-if="question.question_type === 'single_choice'">
                            <p>【单选题】{{ question.question }}</p>
                            <ul>
                                <li v-for="(option, optIndex) in parseOptions(question.options)" :key="optIndex">
                                    {{ option }}
                                </li>
                            </ul>
                            <p>答案：{{ question.answer }}</p>
                        </template>
                        <!-- 多选题 -->
                        <template v-else-if="question.question_type === 'multiple_choice'">
                            <p>【多选题】{{ question.question }}</p>
                            <ul>
                                <li v-for="(option, optIndex) in parseOptions(question.options)" :key="optIndex">
                                    {{ option }}
                                </li>
                            </ul>
                            <p>答案：{{ question.answer }}</p>
                        </template>
                        <!-- 判断题 -->
                        <template v-else-if="question.question_type === 'true_false'">
                            <p>【判断题】{{ question.question }}</p>
                            <p>答案：{{ question.answer === 'true' ? '正确' : '错误' }}</p>
                        </template>
                        <!-- 填空题 -->
                        <template v-else-if="question.question_type === 'fill_blank'">
                            <p>【填空题】{{ question.question }}</p>
                            <p>答案：{{ question.answer }}</p>
                        </template>
                        <!-- 简答题 -->
                        <template v-else-if="question.question_type === 'short_answer'">
                            <p>【简答题】{{ question.question }}</p>
                            <p>答案：{{ question.answer }}</p>
                        </template>
                        <!-- 未知题型 -->
                        <template v-else>
                            <p>未知题型：{{ question.question }}</p>
                        </template>
                    </el-collapse-item>
                </el-collapse>
            </div>
        </div>
    </el-drawer>
</template>

<script lang="ts" setup>
import { computed, onMounted, watch, ref } from 'vue';
import { getQuestions } from '@/api/question';

const props = defineProps(['visible', 'paperData'])
const emit = defineEmits(['update:visible'])

const localVisible = computed({
    get: () => props.visible,
    set: value => emit('update:visible', value)
})

const handleClose = () => {
    emit('update:visible', false)
}

const questionList = ref([] as any[]);

const getQuestionList = async () => {
    try {
        const response: any = await getQuestions({
            question_ids: props.paperData.questions
        });
        questionList.value = response.list;
        console.log(response)
    } catch (error) {
        console.error('获取试题列表失败:', error);
    }
}

// 解析选项字符串为数组
const parseOptions = (options: string): string[] => {
    return options ? options.split('\t') : [];
}

// 监听 paperData 的变化
watch(
    () => props.paperData,
    (newVal, oldVal) => {
        if (newVal !== oldVal) {
            getQuestionList();
        }
    },
    { deep: true } // 深度监听，确保监听对象内部的变化
);

onMounted(() => {
    getQuestionList()
});
</script>

<style scoped>
.paper-details {
    margin-top: -40px;
    padding: 0px;
}
.description {
    color: #666;
    margin-bottom: 16px;
}
.questions {
    margin-top: 16px;
}
</style>