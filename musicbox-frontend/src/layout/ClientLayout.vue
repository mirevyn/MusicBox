<template>
    <!-- 外层容器：占满屏幕高度，禁止 body 滚动 -->
    <div
        class="h-screen w-full bg-[#f8fafc] overflow-hidden flex flex-col font-sans text-gray-900 selection:bg-blue-100 selection:text-blue-700">

        <!-- 顶部导航栏 -->
        <Navbar />

        <!-- 中间主要内容区域 -->
        <main ref="mainRef" class="flex-1 h-full w-full overflow-y-auto overflow-x-hidden relative custom-scrollbar scroll-smooth">

            <!-- 内容容器 -->
            <div class="w-full max-w-7xl mx-auto pt-20 pb-32 px-4 sm:px-6 md:px-8">

                <!-- 路由出口 -->
                <router-view v-slot="{ Component }">
                    <transition name="fade" mode="out-in">
                        <component v-if="Component" :is="Component" :key="route.fullPath" />
                    </transition>
                </router-view>

            </div>
        </main>

        <!-- 底部播放条 -->
        <PlayerBar />

        <!-- AI 助手抽屉 -->
        <AIChatDrawer />
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PlayerBar from '@/components/player/PlayerBar.vue'
import AIChatDrawer from '@/components/ai/AIChatDrawer.vue'

const route = useRoute()
const mainRef = ref<HTMLElement | null>(null)

// 路由切换时手动复位滚动条位置
watch(() => route.path, () => {
    if (mainRef.value) {
        mainRef.value.scrollTop = 0
    }
})
</script>

<style scoped>
/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

/* --- 隐藏滚动条 (保留滚动功能) --- */
.custom-scrollbar::-webkit-scrollbar {
    display: none; /* Chrome, Safari, Edge */
}

.custom-scrollbar {
    -ms-overflow-style: none; /* IE and Edge */
    scrollbar-width: none; /* Firefox */
}
</style>