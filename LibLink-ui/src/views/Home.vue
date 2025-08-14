<template>
    <div class="home-wrapper">
        <el-row :gutter="20">
            <el-col :span="24">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <span>欢迎来到 LibLink</span>
                    </div>
                    <div class="text item">欢迎使用 LibLink 管理文献</div>
                </el-card>
            </el-col>
        </el-row>
        <el-row :gutter="20" class="stats-row">
            <el-col :span="8">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <el-icon><i class="el-icon-document"></i></el-icon>
                        <span>文献总数</span>
                    </div>
                    <div class="text item">{{ questions }}</div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <el-icon><i class="el-icon-edit"></i></el-icon>
                        <span>借阅记录总数</span>
                    </div>
                    <div class="text item">{{ papers }}</div>
                </el-card>
            </el-col>
            <el-col :span="8">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <el-icon><i class="el-icon-user"></i></el-icon>
                        <span>用户总数</span>
                    </div>
                    <div class="text item">{{ users }}</div>
                </el-card>
            </el-col>
        </el-row>
        <el-row :gutter="20" class="announcement-row">
            <el-col :span="24">
                <el-card class="box-card">
                    <div slot="header" class="clearfix">
                        <el-icon><i class="el-icon-bell"></i></el-icon>
                        <span>系统公告</span>
                    </div>
                    <el-alert v-for="(announcement, index) in announcements" :key="index" :title="announcement.title" type="info" show-icon :description="announcement.content" class="announcement-item" />
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { getNotifications } from '@/api/system'
import { summary as userSummary } from '@/api/user'

const papers = ref(0)
const users = ref(0)
const questions = ref(0)

interface Announcement {
    title: string
    content: string
}

const announcements = ref<Announcement[]>([])

const getNotificationsData = async () => {
    try {
        const response: any = await getNotifications()
        announcements.value = response.list
    } catch (error) {
        console.error('获取系统公告失败:', error)
    }
}

// 在组件挂载时调用
onMounted(() => {
    getNotificationsData()
})
</script>

<style scoped>
.home-wrapper {
    padding: 20px;
    background-color: #f5f5f5;
}

.box-card {
    margin-bottom: 20px;
    border-radius: 10px;
    overflow: hidden;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.box-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.stats-row {
    margin-top: 5px;
}

.nav-row {
    margin-top: 40px;
    text-align: center;
}

.announcement-row {
    margin-top: 20px;
}

.text {
    font-size: 24px;
    text-align: center;
    margin-top: 20px;
    color: #333;
}

.el-icon {
    margin-right: 10px;
    color: #409eff;
}

.el-card .header {
    display: flex;
    align-items: center;
    font-size: 18px;
    font-weight: bold;
    color: #333;
}

.el-card .item {
    font-size: 20px;
    font-weight: bold;
    color: #409eff;
}

.announcement-item {
    margin-top: 10px;
}
</style>
