<template>
    <!-- 歌曲列表容器 -->
    <div class="bg-white/60 backdrop-blur-md rounded-[2rem] border border-white shadow-sm"
        :class="{ 'min-h-[400px]': isEmpty }">
        <!-- 列表头: 定义网格布局 -->
        <div
            class="grid grid-cols-[3rem_1fr_3rem] md:grid-cols-[4rem_4fr_3fr_2fr_6rem] gap-4 px-6 py-4 border-b border-gray-100 text-xs font-bold text-gray-400 uppercase tracking-wider"
            :class="headerBg">
            <div class="text-center">#</div>
            <div>标题</div>
            <div class="hidden md:block">专辑</div>
            <div class="hidden md:block">歌手</div>
            <div class="text-center md:text-right md:pr-4">{{ actionLabel }}</div>
        </div>

        <!-- 列表内容区域 -->
        <div class="divide-y divide-gray-50/80">
            <slot></slot>
            
            <!-- 空状态插槽 -->
            <div v-if="isEmpty" class="flex flex-col items-center justify-center py-20 text-gray-400">
                <slot name="empty">
                    <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                        <div class="i-fa6-solid-record-vinyl text-3xl opacity-30"></div>
                    </div>
                    <p class="text-sm">暂无歌曲</p>
                </slot>
            </div>
        </div>
        
        <!-- 底部计数信息 -->
        <div v-if="!isEmpty" class="py-4 text-center text-xs text-gray-400 bg-gray-50/50 border-t border-gray-100 rounded-b-[2rem]">
            共 {{ count }} 首歌曲
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
    isEmpty: boolean
    count: number
    theme?: 'primary' | 'blue' | 'rose' | 'violet'
    actionLabel?: string
}>(), {
    theme: 'primary',
    actionLabel: '时长'
})

const headerBg = computed(() => {
    return props.theme === 'rose' ? 'bg-gray-50/50' : ''
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
    width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
    background: #e5e7eb;
    border-radius: 4px;
}
.custom-scrollbar:hover::-webkit-scrollbar-thumb {
    background: #d1d5db;
}
</style>
