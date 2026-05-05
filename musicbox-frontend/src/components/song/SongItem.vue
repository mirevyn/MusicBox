<template>
    <div class="group relative flex items-center p-2.5 md:p-2 rounded-xl bg-white/40 border border-white/20 hover:bg-white/60 hover:shadow-md transition-all duration-300 cursor-pointer hover:z-[60]"
        :class="{ '!bg-white/80 shadow-md z-10': isActive }" @click="$emit('play', song)">

        <!-- 封面图 -->
        <div
            class="relative w-11 h-11 md:w-12 md:h-12 flex-shrink-0 mr-3 overflow-hidden rounded-lg shadow-sm group-hover:shadow transition-shadow">
            <SongCover :src="resolveCover(song.coverUrl)" :alt="song.title" class="w-full h-full" />
            
            <!-- 播放中遮罩 -->
            <div v-if="isActive" class="absolute inset-0 bg-black/20 flex items-center justify-center pointer-events-none">
                <div class="i-fa6-solid-chart-simple text-white animate-pulse text-sm"></div>
            </div>

            <!-- 悬停播放图标 (非播放状态下，仅桌面端显示以避免移动端触摸粘滞) -->
            <div v-else
                class="hidden md:flex absolute inset-0 bg-black/20 items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none">
                <div class="i-fa6-solid-play text-white text-base"></div>
            </div>
        </div>

        <!-- 歌曲信息 -->
        <div class="flex-1 flex flex-col min-w-0 gap-0.5 md:gap-1">
            <div class="font-semibold text-[13px] md:text-sm truncate text-gray-800 transition-colors"
                :class="{ 'text-blue-600': isActive }" :title="song.title">
                {{ song.title }}
            </div>
            <div class="text-[11px] md:text-xs text-gray-500 truncate" :title="song.artist">
                {{ song.artist }}
            </div>
        </div>

        <!-- 操作菜单 -->
        <Dropdown class="ml-2 block" width="w-36" placement="bottom-right">
            <template #trigger>
                <!-- 更多按钮: 移动端常驻可见，桌面端悬停显示 -->
                <button
                    class="w-9 h-9 md:w-8 md:h-8 flex items-center justify-center rounded-full bg-transparent md:bg-gray-50 text-gray-500 hover:bg-gray-200 hover:text-gray-700 transition-all opacity-100 md:opacity-0 group-hover:opacity-100 scale-100 md:scale-90 group-hover:scale-100 active:scale-95"
                    @click.stop>
                    <div class="i-fa6-solid-ellipsis pointer-events-none text-base md:text-sm"></div>
                </button>
            </template>

            <!-- 下拉内容 -->
            <button
                class="w-full text-left px-3 py-2.5 md:py-2 rounded-lg hover:bg-gray-50 text-[13px] md:text-xs font-medium text-gray-700 flex items-center gap-2.5 md:gap-2 transition-colors active:bg-gray-100"
                @click="$emit('add-to-queue', song)">
                <div class="i-fa6-solid-list-check text-gray-400 text-sm md:text-xs"></div>
                <span>下一首播放</span>
            </button>
            <button
                class="w-full text-left px-3 py-2.5 md:py-2 rounded-lg hover:bg-gray-50 text-[13px] md:text-xs font-medium text-gray-700 flex items-center gap-2.5 md:gap-2 transition-colors active:bg-gray-100"
                @click="$emit('add-to-playlist', song)">
                <div class="i-fa6-solid-folder-plus text-gray-400 text-sm md:text-xs"></div>
                <span>收藏到歌单</span>
            </button>
        </Dropdown>
    </div>
</template>

<script setup lang="ts">
import { resolveCover } from '@/utils/common'
import type { Song } from '@/stores/player'
import Dropdown from '@/components/ui/Dropdown.vue'
import SongCover from '@/components/common/SongCover.vue'

withDefaults(defineProps<{
    song: Song
    isActive?: boolean
}>(), {
    isActive: false
})

defineEmits<{
    (e: 'play', song: Song): void
    (e: 'add-to-queue', song: Song): void
    (e: 'add-to-playlist', song: Song): void
}>()
</script>