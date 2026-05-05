/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_BACKEND_URL: string
    // 其他环境变量
}

interface ImportMeta {
    readonly env: ImportMetaEnv
}

// Vue 组件类型声明
declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<object, object, unknown>
    export default component
}
