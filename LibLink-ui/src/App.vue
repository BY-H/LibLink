<template>
    <el-container class="full-height">
        <el-aside width="200px" height="100%">
            <AsideNav />
        </el-aside>

        <el-main class="main-content">
            <RouterView />
        </el-main>
    </el-container>
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import AsideNav from './components/AsideNav.vue'
</script>

<style scoped>
html,
body,
#app {
    height: 100%;
    margin: 0;
    padding: 0;
}

.full-height {
    height: 100%;
}

.el-main {
    display: flex;
    flex-direction: column;
    align-items: stretch;
    padding: 0;
    background-color: #f5f5f5;
}

.main-content {
    flex: 1;
    display: flex;
    flex-direction: column;
}
</style>
