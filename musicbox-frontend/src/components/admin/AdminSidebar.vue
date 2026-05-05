<template>
    <aside
        class="w-72 bg-white border-r border-gray-100 flex flex-col fixed h-full inset-y-0 left-0 z-50 transition-transform duration-300 lg:translate-x-0 dark:bg-slate-900 dark:border-slate-800 shadow-xl lg:shadow-none"
        :class="isOpen ? 'translate-x-0' : '-translate-x-full'">
        <div class="h-20 flex items-center px-8 border-b border-gray-100 dark:border-slate-800">
            <div class="flex items-center gap-3">
                <img src="/logo.svg" alt="MusicBox Logo" class="w-8 h-8 rounded-lg shadow-sm object-cover" />
                <h1 class="text-xl font-bold tracking-tight text-gray-900 dark:text-white">MusicBox</h1>
            </div>
        </div>

        <nav class="flex-1 py-6 px-4 space-y-1.5 overflow-y-auto custom-scrollbar">
            <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider px-4 mb-2">Menu</div>

            <RouterLink to="/admin"
                exact-active-class="bg-primary/5 text-primary font-medium dark:bg-primary/20 dark:text-white"
                class="flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-xl transition-all duration-300 group dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-gauge transition-transform group-hover:scale-110"></div>
                <span>仪表盘</span>
            </RouterLink>

            <RouterLink to="/admin/songs"
                active-class="bg-primary/5 text-primary font-medium dark:bg-primary/20 dark:text-white"
                class="flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-xl transition-all duration-300 group dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-music transition-transform group-hover:scale-110"></div>
                <span>歌曲管理</span>
            </RouterLink>

            <RouterLink to="/admin/users"
                active-class="bg-primary/5 text-primary font-medium dark:bg-primary/20 dark:text-white"
                class="flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-xl transition-all duration-300 group dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-users transition-transform group-hover:scale-110"></div>
                <span>用户管理</span>
            </RouterLink>

            <RouterLink to="/admin/playlists"
                active-class="bg-primary/5 text-primary font-medium dark:bg-primary/20 dark:text-white"
                class="flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-xl transition-all duration-300 group dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-list transition-transform group-hover:scale-110"></div>
                <span>歌单管理</span>
            </RouterLink>

            <RouterLink to="/admin/settings"
                active-class="bg-primary/5 text-primary font-medium dark:bg-primary/20 dark:text-white"
                class="flex items-center gap-3 px-4 py-3 text-gray-600 hover:bg-gray-50 hover:text-gray-900 rounded-xl transition-all duration-300 group dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-gear transition-transform group-hover:scale-110"></div>
                <span>系统设置</span>
            </RouterLink>
        </nav>

        <div class="p-4 border-t border-gray-100 dark:border-slate-800">
            <button @click="$router.push('/')"
                class="w-full flex items-center justify-center gap-2 text-gray-500 hover:text-primary hover:bg-primary/5 py-3 rounded-xl transition-all duration-300 text-sm font-medium dark:text-slate-400 dark:hover:bg-slate-800 dark:hover:text-white">
                <div class="i-fa6-solid-arrow-right-from-bracket"></div>
                <span>返回客户端</span>
            </button>
        </div>
    </aside>

    <!-- 移动端遮罩层 -->
    <div v-if="isOpen" @click="$emit('close')"
        class="fixed inset-0 bg-gray-900/20 backdrop-blur-sm z-40 lg:hidden transition-opacity"></div>
</template>

<script setup lang="ts">
defineProps<{
    isOpen: boolean
}>()

defineEmits<{
    (e: 'close'): void
}>()
</script>

<style scoped>
.glass-effect {
    backdrop-filter: blur(4px);
}
</style>
