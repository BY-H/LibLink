<template>
    <div class="papers-wrapper">
        <el-card class="box-card">
            <div slot="header" class="clearfix">
                <span>试卷列表</span>
            </div>
            <el-table :data="papers" style="width: 100%">
                <el-table-column prop="ID" label="ID" width="50"></el-table-column>
                <el-table-column prop="title" label="标题" width="200"></el-table-column>
                <el-table-column prop="description" label="描述" width="200"></el-table-column>
                <el-table-column prop="creator" label="创建人" width="150"></el-table-column>
                <el-table-column prop="CreatedAt" label="创建时间" width="180"></el-table-column>
                <el-table-column prop="UpdatedAt" label="更新时间" width="180"></el-table-column>
                <el-table-column label="操作" width="250">
                    <template v-slot="scope">
                        <el-button size="mini" @click="handleView(scope.row)">查看</el-button>
                        <el-button size="mini" @click="handleDownload(scope.row)">下载</el-button>
                        <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </el-card>
        <Pagination :pageObj="pageObj" :total="total" :onUpdate="onUpdate" />
    </div>
    <ViewPaper v-model:visible="viewPaperVisible" :paperData="viewPaper"/>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import Pagination from '@/components/Pagination.vue'
import { getPapers } from '@/api/paper'
import ViewPaper from './ViewPaper.vue'
import { exportPaper } from '@/api/paper'

interface Paper {
    id: number
    title: string
    description: string
    creator: string
    created_at: string
    updated_at: string
    status: string
}

const total = ref(0)
const pageObj = ref({
    page: 1,
    page_size: 20
})
const onUpdate = () => {
    getPaper()
}

const papers = ref<Paper[]>([])

const getPaper = async() => {
    const params = {
        data_start: '',
        data_end: '',
        ...pageObj.value
    }
    const res: any = await getPapers(params)
    papers.value = res.list
    total.value = res.total
}

const viewPaperVisible = ref(false)
const viewPaper = ref({})
const handleView = (row: Paper) => {
    console.log(row)
    viewPaperVisible.value = true
    viewPaper.value = row 
}

const handleDelete = (row: Paper) => {
    console.log('删除试卷', row)
}

const handleDownload = async (row: any) => {
    try {
        // 确保请求返回的是二进制数据
        const response: any = await exportPaper(
            { paper_id: String(row.ID) },
            { responseType: 'blob' } // 确保 responseType 为 'blob'
        );

        // 直接使用 response.data，不需要手动添加 "\ufeff"
        const blob = new Blob([response], {
            type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
        });

        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', `${row.title}.docx`);
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        window.URL.revokeObjectURL(url);

        console.log('试卷下载成功');
    } catch (error) {
        console.error('试卷下载失败:', error);
    }


};


onMounted(() => {
    getPaper()
})
</script>

<style scoped>
.papers-wrapper {
    padding: 20px;
}

.box-card {
    margin-bottom: 20px;
}
</style>
