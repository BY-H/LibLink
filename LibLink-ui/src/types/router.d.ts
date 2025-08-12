import 'vue-router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

declare module 'vue-router' {
    interface RouteMeta {
        title?: string // 路由标题
        icon?: keyof typeof ElementPlusIconsVue
    }
}

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
    const component: DefineComponent<{}, {}, any>
    export default component
}
