<template>
    <!-- 歌曲行组件 -->
    <div class="group grid grid-cols-[3rem_1fr_3rem] md:grid-cols-[4rem_4fr_3fr_2fr_6rem] gap-4 px-6 py-3 items-center hover:bg-white/80 transition-all duration-200 cursor-pointer"
        :class="activeClass" @dblclick="$emit('play', index)">

        <!-- 序号/播放状态区域 -->
        <div class="text-center text-sm font-medium transition-colors cursor-pointer" :class="statusClass"
            @click.stop="$emit('play', index)">
            <!-- 正在播放动画 -->
            <div v-if="isActive && isPlaying" class="flex items-end justify-center gap-0.5 h-3">
                <div class="w-0.5 h-3 animate-[bounce_1s_infinite]" :class="bgClass"></div>
                <div class="w-0.5 h-2 animate-[bounce_1.2s_infinite]" :class="bgClass"></div>
                <div class="w-0.5 h-3 animate-[bounce_0.8s_infinite]" :class="bgClass"></div>
            </div>
            <!-- 暂停状态图标 -->
            <div v-else-if="isActive" class="i-fa6-solid-play text-xs mx-auto"></div>
            <!-- 普通状态: 显示序号，Hover时显示播放图标 -->
            <template v-else>
                <span class="group-hover:hidden font-mono text-gray-400">{{ index + 1 }}</span>
                <div class="hidden group-hover:block i-fa6-solid-play mx-auto text-xs text-gray-600"></div>
            </template>
        </div>

        <!-- 标题 + 封面区域 -->
        <div class="flex items-center gap-4 min-w-0">
            <!-- 封面图 -->
            <div
                class="w-10 h-10 rounded-lg shrink-0 overflow-hidden relative shadow-sm group-hover:shadow transition-shadow flex items-center justify-center">
                <SongCover :src="resolveCover(song.coverUrl)" class="w-full h-full bg-gray-100" />
            </div>
            <!-- 标题文本 -->
            <div class="min-w-0 flex-1">
                <div class="text-sm font-bold truncate transition-colors" :class="titleClass">
                    {{ song.title }}
                </div>
                <!-- Mobile端显示的歌手名 -->
                <div class="md:hidden text-xs text-gray-500 truncate mt-0.5">{{ song.artist }}</div>
            </div>
        </div>

        <!-- 专辑 -->
        <div class="hidden md:block text-sm text-gray-500 truncate group-hover:text-gray-700 transition-colors">
            {{ song.album || '未知专辑' }}
        </div>

        <!-- 歌手 -->
        <div class="hidden md:block text-sm text-gray-500 truncate group-hover:text-gray-900 transition-colors">
            {{ song.artist }}
        </div>

        <!-- 操作区 / 时长 -->
        <div class="flex items-center justify-end md:justify-end md:pr-2 gap-2">
            <div class="text-xs text-gray-400 font-medium font-mono group-hover:hidden">
                {{ formatDuration(song.duration || 0) }}
            </div>

            <div class="hidden group-hover:block">
                <slot name="actions">
                    <Dropdown width="w-32" placement="bottom-right">
                        <template #trigger>
                            <button
                                class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-colors"
                                @click.stop>
                                <div class="i-fa6-solid-ellipsis pointer-events-none"></div>
                            </button>
                        </template>

                        <button
                            class="w-full text-left px-3 py-2 rounded-lg hover:bg-gray-50 text-xs font-medium text-gray-700 flex items-center gap-2 transition-colors"
                            @click="$emit('addToQueue', song)">
                            <div class="i-fa6-solid-list-check text-gray-400"></div>
                            <span>下一首播放</span>
                        </button>

                        <button
                            class="w-full text-left px-3 py-2 rounded-lg hover:bg-gray-50 text-xs font-medium text-gray-700 flex items-center gap-2 transition-colors"
                            @click="$emit('addToPlaylist', song)">
                            <div class="i-fa6-solid-folder-plus text-gray-400"></div>
                            <span>收藏到歌单</span>
                        </button>

                        <button v-if="showRemove"
                            class="w-full text-left px-3 py-2 rounded-lg hover:bg-red-50 text-xs font-medium text-gray-700 hover:text-red-500 flex items-center gap-2 transition-colors"
                            @click="$emit('remove', song.songId)">
                            <div class="i-fa6-solid-trash text-gray-400 group-hover:text-red-500"></div>
                            <span>{{ removeLabel }}</span>
                        </button>
                    </Dropdown>
                </slot>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { resolveCover, formatDuration } from '@/utils/common'
import type { Song } from '@/stores/player'
import SongCover from '@/components/common/SongCover.vue'
import Dropdown from '@/components/ui/Dropdown.vue'

const props = withDefaults(defineProps<{
    song: Song
    index: number
    isActive?: boolean
    isPlaying?: boolean
    theme?: 'primary' | 'blue' | 'rose' | 'violet'
    showRemove?: boolean
    removeLabel?: string
}>(), {
    isActive: false,
    isPlaying: false,
    theme: 'primary',
    showRemove: false,
    removeLabel: '从列表移除'
})

defineEmits<{
    (e: 'play', index: number): void
    (e: 'addToQueue', song: Song): void
    (e: 'addToPlaylist', song: Song): void
    (e: 'remove', id: number | string): void
}>()



// 样式映射
const statusClass = computed(() => {
    if (!props.isActive) return 'text-gray-400 group-hover:text-gray-600'
    return props.theme === 'rose' ? 'text-rose-500' : (props.theme === 'blue' ? 'text-blue-500' : (props.theme === 'violet' ? 'text-violet-500' : 'text-primary'))
})

const bgClass = computed(() => {
    return props.theme === 'rose' ? 'bg-rose-500' : (props.theme === 'blue' ? 'bg-blue-500' : (props.theme === 'violet' ? 'bg-violet-500' : 'bg-primary'))
})

const titleClass = computed(() => {
    if (!props.isActive) return 'text-gray-800'
    return props.theme === 'rose' ? 'text-rose-600' : (props.theme === 'blue' ? 'text-blue-600' : (props.theme === 'violet' ? 'text-violet-600' : 'text-primary'))
})

const activeClass = computed(() => {
    if (!props.isActive) return ''
    const color = props.theme === 'rose' ? 'rose' : (props.theme === 'blue' ? 'blue' : (props.theme === 'violet' ? 'violet' : 'primary'))
    return `bg-${color}-50/50 hover:bg-${color}-50`
})
</script>
