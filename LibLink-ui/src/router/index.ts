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
        path: '/user',
        name: 'User',
        component: () => import('@/views/User.vue'),
        meta: { title: '用户管理', icon: 'User' }
    },
    {
        path: '/borrow',
        name: 'Borrow',
        component: () => import('@/views/borrow/Borrow.vue'),
        meta: { title: '文献借阅', icon: 'Reading' }
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
