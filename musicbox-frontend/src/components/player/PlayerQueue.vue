<template>
    <Transition
        enter-active-class="transition ease-out duration-300 transform"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transition ease-in duration-300 transform"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
    >
        <div v-if="player.showPlayList"
            class="fixed right-2 bottom-24 w-[90vw] sm:w-80 md:w-96 max-h-[60vh] bg-white/95 backdrop-blur-2xl rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-white/50 z-50">
            <!-- 头部 -->
            <div class="px-5 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
                <div class="flex items-center gap-2">
                    <span class="font-bold text-gray-800 text-base">{{ title }}</span>
                    <span class="text-xs bg-gray-200 text-gray-600 px-1.5 py-0.5 rounded-md font-medium">{{ player.playList.length }}</span>
                </div>
                <button
                    class="text-xs flex items-center gap-1.5 text-gray-500 hover:text-rose-500 hover:bg-rose-50 px-2.5 py-1.5 rounded-full transition-colors"
                    @click="player.clearPlayList()">
                    <div class="i-fa6-solid-trash-can"></div> 清空
                </button>
            </div>

            <!-- 列表 -->
            <div ref="queueListRef" class="flex-1 overflow-y-auto custom-scrollbar p-2">
                <div v-for="(song, index) in player.playList" :key="song.songId"
                    class="queue-song-row group flex items-center justify-between p-2 rounded-xl hover:bg-gray-100/80 cursor-pointer transition-colors mb-1"
                    :class="{ 'bg-blue-50 hover:!bg-blue-50/80': isCurrentSong(song) }"
                    @click="handleQueueItemClick(index)">
                    <div class="flex items-center gap-3 min-w-0 flex-1">
                        <button
                            class="queue-drag-handle w-7 h-10 flex flex-shrink-0 items-center justify-center text-gray-300 hover:text-gray-500 cursor-grab active:cursor-grabbing rounded-lg transition-colors"
                            type="button"
                            title="拖动排序"
                            aria-label="拖动排序"
                            @click.stop>
                            <div class="i-fa6-solid-grip-lines text-sm"></div>
                        </button>
                        <div class="relative w-10 h-10 flex-shrink-0">
                            <SongCover 
                                :src="resolveUrl(song.coverUrl)" 
                                class="w-full h-full rounded-lg border border-black/5"
                            >
                                <div v-if="isCurrentSong(song)"
                                    class="absolute inset-0 bg-black/20 rounded-lg flex items-center justify-center">
                                    <div class="w-1 h-3 bg-white mx-0.5 animate-bounce"></div>
                                    <div class="w-1 h-4 bg-white mx-0.5 animate-bounce animation-delay-75"></div>
                                    <div class="w-1 h-2 bg-white mx-0.5 animate-bounce animation-delay-150"></div>
                                </div>
                            </SongCover>
                        </div>
                        <div class="flex flex-col min-w-0">
                            <div class="text-sm font-medium truncate"
                                :class="isCurrentSong(song) ? 'text-blue-600' : 'text-gray-700'">{{ song.title }}
                            </div>
                            <div class="text-xs text-gray-400 truncate">{{ song.artist }}</div>
                        </div>
                    </div>
                    <button
                        class="w-7 h-7 flex items-center justify-center text-gray-300 hover:text-rose-500 hover:bg-rose-50 rounded-full transition-all opacity-0 group-hover:opacity-100"
                        @click.stop="player.removeFromPlayList(song.songId)">
                        <div class="i-fa6-solid-xmark text-sm"></div>
                    </button>
                </div>
            </div>
        </div>
    </Transition>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import Sortable from 'sortablejs'
import { usePlayerStore } from '@/stores/player'
import { resolveUrl } from '@/utils/common'
import SongCover from '@/components/common/SongCover.vue'

const player = usePlayerStore()
const queueListRef = ref<HTMLElement | null>(null)
const sortable = ref<Sortable | null>(null)
const suppressClickUntil = ref(0)

const title = computed(() => player.playlistSource || '播放列表')

const isCurrentSong = (song: any) => player.currentSong?.songId === song.songId

const destroySortable = () => {
    sortable.value?.destroy()
    sortable.value = null
}

const initSortable = async () => {
    if (!player.showPlayList || sortable.value) return

    await nextTick()
    const listEl = queueListRef.value
    if (!listEl) return

    sortable.value = Sortable.create(listEl, {
        animation: 150,
        handle: '.queue-drag-handle',
        draggable: '.queue-song-row',
        ghostClass: 'queue-drag-ghost',
        chosenClass: 'queue-drag-chosen',
        dragClass: 'queue-drag-active',
        delay: 120,
        delayOnTouchOnly: true,
        fallbackOnBody: true,
        fallbackTolerance: 4,
        onStart: () => {
            suppressClickUntil.value = Date.now() + 300
        },
        onEnd: (event: Sortable.SortableEvent) => {
            suppressClickUntil.value = Date.now() + 300
            const { oldIndex, newIndex } = event
            if (typeof oldIndex !== 'number' || typeof newIndex !== 'number') return
            player.reorderPlayList(oldIndex, newIndex)
        },
    })
}

const handleQueueItemClick = (index: number) => {
    if (Date.now() < suppressClickUntil.value) return
    player.playByIndex(index)
}

watch(() => player.showPlayList, (visible) => {
    if (visible) {
        void initSortable()
    } else {
        destroySortable()
    }
}, { flush: 'post' })

onMounted(() => {
    void initSortable()
})

onBeforeUnmount(() => {
    destroySortable()
})
</script>

<style scoped>
:deep(.queue-drag-ghost) {
    opacity: 0.35;
}

:deep(.queue-drag-chosen) {
    background: rgba(239, 246, 255, 0.95);
}

:deep(.queue-drag-active) {
    cursor: grabbing;
}

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

.animation-delay-75 {
    animation-delay: 75ms;
}

.animation-delay-150 {
    animation-delay: 150ms;
}
</style>
