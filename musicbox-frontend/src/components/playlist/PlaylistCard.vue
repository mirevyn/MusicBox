<template>
    <div
        class="group cursor-pointer"
        :class="compact ? 'w-full space-y-2' : 'space-y-3'"
        @click="$router.push(`/playlist/${playlist.id}`)"
    >
        <!-- 封面 -->
        <div
            class="relative overflow-hidden bg-gray-100 shadow-sm group-hover:shadow-lg transition-all duration-300 group-hover:-translate-y-1 transform-gpu"
            :class="compact ? 'aspect-square rounded-xl' : 'aspect-square rounded-2xl'"
            style="-webkit-mask-image: -webkit-radial-gradient(white, black); mask-image: radial-gradient(white, black);">
            
            <SongCover 
                :src="resolveCover(playlist.coverUrl)" 
                class="w-full h-full"
                :icon-class="compact ? 'text-3xl opacity-45' : 'text-4xl opacity-50'"
            />

            <!-- 播放中状态 / 播放按钮遮罩 -->
            <div class="absolute inset-0 bg-black/40 flex items-center justify-center transition-all duration-300 z-20 pointer-events-none" 
                 :class="{ 'opacity-100 pointer-events-auto': loading || isPlaying, 'opacity-0 group-hover:opacity-100 group-hover:pointer-events-auto': !loading && !isPlaying }">
                 
                 <!-- Loading -->
                 <div v-if="loading" class="text-white">
                     <div class="i-svg-spinners-90-ring-with-bg text-2xl"></div>
                 </div>
                 
                 <!-- 交互区 -->
                 <template v-else>
                     <!-- 播放中未悬浮：显示律动条 -->
                     <div v-if="isPlaying" class="absolute inset-0 flex items-center justify-center transition-opacity duration-300 group-hover:opacity-0 pointer-events-none">
                         <div class="flex items-end gap-[3px] h-4">
                            <div v-for="(item, i) in bars" :key="i" class="w-[3px] bg-white rounded-t-sm" 
                                 :style="`height:${item.h}%; animation: spectrum 0.6s ease-in-out infinite alternate; animation-delay:${item.delay}ms`"></div>
                         </div>
                     </div>

                     <!-- 播放/暂停按钮：悬浮时显示 -->
                     <div @click.stop="handlePlay" 
                          class="w-10 h-10 rounded-full bg-primary/95 text-white flex items-center justify-center shadow-lg shadow-primary/40 transition-all cursor-pointer hover:shadow-primary/60 hover:scale-105 active:scale-95 z-10"
                          :class="isPlaying ? 'opacity-0 group-hover:opacity-100 scale-90 group-hover:scale-100' : ''">
                         <div :class="isPlaying ? 'i-fa6-solid-pause' : 'i-fa6-solid-play ml-0.5'" class="text-sm"></div>
                     </div>
                 </template>
            </div>

            <div
                class="absolute z-10 flex flex-col items-end"
                :class="compact ? 'top-1.5 right-1.5 gap-1' : 'top-2 right-2 gap-1.5'"
            >
                <div v-if="!playlist.isPublic"
                    class="bg-black/50 backdrop-blur-md text-white/90 rounded-full flex items-center gap-1"
                    :class="compact ? 'px-1.5 py-0.5 text-[9px]' : 'px-2 py-0.5 text-[10px]'"
                >
                    <div class="i-fa6-solid-lock text-[8px]"></div> 私密
                </div>
                <div
                    v-if="playlist.isPublic && playlist.status === 0"
                    class="bg-amber-500/90 backdrop-blur-md text-white rounded-full flex items-center gap-1"
                    :class="compact ? 'px-1.5 py-0.5 text-[9px]' : 'px-2 py-0.5 text-[10px]'"
                >
                    <div class="i-fa6-solid-hourglass-half text-[8px]"></div> 待审核
                </div>
                <div
                    v-else-if="playlist.isPublic && playlist.status === 2"
                    @click.stop
                >
                    <Dropdown placement="bottom-right" width="w-48">
                        <template #trigger>
                            <div class="bg-red-500/90 backdrop-blur-md text-white rounded-full flex items-center gap-1 cursor-help hover:bg-red-600 transition-colors"
                                :class="compact ? 'px-1.5 py-0.5 text-[9px]' : 'px-2 py-0.5 text-[10px]'"
                            >
                                <div class="i-fa6-solid-ban text-[8px]"></div> 已驳回
                            </div>
                        </template>
                        <div class="p-2">
                            <div class="flex items-center gap-1.5 text-red-500 font-bold text-xs mb-1.5">
                                <div class="i-fa6-solid-triangle-exclamation"></div>
                                <span>驳回原因</span>
                            </div>
                            <p class="text-xs text-gray-600 leading-relaxed whitespace-pre-wrap break-words">
                                {{ playlist.rejectReason || '内容不合规，未能通过审核' }}
                            </p>
                        </div>
                    </Dropdown>
                </div>
                <div
                    v-else-if="showApprovedBadge && playlist.isPublic && playlist.status === 1"
                    class="bg-emerald-500/90 backdrop-blur-md text-white rounded-full flex items-center gap-1"
                    :class="compact ? 'px-1.5 py-0.5 text-[9px]' : 'px-2 py-0.5 text-[10px]'"
                >
                    <div class="i-fa6-solid-earth-asia text-[8px]"></div> 已公开
                </div>
            </div>
        </div>

        <!-- 信息 -->
        <div>
            <h3
                class="font-bold text-gray-800 truncate group-hover:text-primary transition-colors"
                :class="compact ? 'text-sm' : ''"
            >
                {{ playlist.title }}
            </h3>
            <p class="text-xs text-gray-500 truncate" :class="compact ? 'mt-0.5' : 'mt-1'">
                {{ playlist.songCount || 0 }} 首歌曲
                <span v-if="playlist.isPublic && playlist.status === 0" class="text-amber-500"> · 审核中</span>
                <span v-else-if="playlist.isPublic && playlist.status === 2" class="text-red-500"> · 审核未通过</span>
            </p>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { resolveCover, formatSongs, showToast } from '@/utils/common'
import type { Playlist, SongDTO } from '@/types/api' 
import { getPlaylistDetails } from '@/api/playlists'
import { usePlayerStore } from '@/stores/player'
import SongCover from '@/components/common/SongCover.vue'
import Dropdown from '@/components/ui/Dropdown.vue'

const player = usePlayerStore()
const loading = ref(false)

const props = withDefaults(defineProps<{
    playlist: Playlist | any
    compact?: boolean
    showApprovedBadge?: boolean
}>(), {
    compact: false,
    showApprovedBadge: true
})

const isPlaying = computed(() => {
    return player.playlistSource === props.playlist.title && player.playing
})

const handlePlay = async () => {
    if (player.playlistSource === props.playlist.title) {
        player.togglePlay()
        return
    }

    if (loading.value) return
    loading.value = true
    try {
        const res = await getPlaylistDetails(props.playlist.id) as unknown as { code: number; data: { songs: SongDTO[] } }
        if (res.code === 200 && res.data?.songs) {
            const songs = formatSongs(res.data.songs)
            if (songs.length > 0) {
                player.setPlayList(songs, props.playlist.title)
                player.playByIndex(0)
                showToast(`开始播放: ${props.playlist.title}`)
            } else {
                showToast('歌单为空', 'warning')
            }
        }
    } catch (e) {
        showToast('获取歌单失败', 'error')
    } finally {
        loading.value = false
    }
}

const bars = [
    { h: 30, delay: 0 },
    { h: 85, delay: 300 },
    { h: 45, delay: 150 },
    { h: 100, delay: 450 },
]
</script>

<style scoped>
@keyframes spectrum {
  0% { height: 25%; opacity: 0.7; }
  100% { height: 100%; opacity: 1; }
}
</style>
