import { createRouter, createMemoryHistory, type RouteRecordRaw } from 'vue-router'
import { isTokenValid } from '@/utils/auth'

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Home',
        component: () => import('@/views/Home.vue'),
        meta: { title: '首页', icon: 'HomeFilled' }
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/Login.vue')
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/Register.vue')
    },
    {
        path: '/question',
        name: 'Question',
        component: () => import('@/views/Question/Question.vue'),
        meta: { title: '题库管理', icon: 'List' }
    },
    {
        path: '/exam',
        name: 'Exam',
        component: () => import('@/views/Exam.vue'),
        meta: { title: '出卷', icon: 'Edit' }
    },
    {
        path: '/paper',
        name: 'Paper',
        component: () => import('@/views/Paper/Paper.vue'),
        meta: { title: '出题历史', icon: 'Paperclip' }
    },
    {
        path: '/feedback',
        name: 'Feedback',
        component: () => import('@/views/Feedback.vue'),
        meta: { title: '用户反馈', icon: 'Message' }
    },
    {
        path: '/user',
        name: 'User',
        component: () => import('@/views/User.vue'),
        meta: { title: '用户管理', icon: 'User' }
    }
]

const router = createRouter({
    routes,
    history: createMemoryHistory()
})

router.beforeEach(async (to, from, next) => {
    if (to.path === '/login' || to.path === '/register') {
        next()
        return
    }

    if (isTokenValid()) {
        next()
    } else {
        next({ name: 'Login' })
    }
})

export default router
