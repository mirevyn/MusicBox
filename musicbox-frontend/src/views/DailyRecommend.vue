<template>
    <div class="flex flex-col space-y-6 animate-fade-in pb-20">
        <!-- 头部卡片 -->
        <div
            class="relative flex flex-col md:flex-row gap-6 items-center md:items-end bg-white/60 backdrop-blur-md p-6 md:p-8 rounded-[2rem] border border-white shadow-sm overflow-hidden group">

            <!-- 背景装饰 -->
            <div
                class="absolute top-0 right-0 w-96 h-96 bg-violet-200/30 rounded-full blur-3xl -translate-y-1/2 translate-x-1/3 pointer-events-none group-hover:bg-violet-300/30 transition-colors duration-700">
            </div>
            <div
                class="absolute bottom-0 left-0 w-64 h-64 bg-fuchsia-100/30 rounded-full blur-3xl translate-y-1/2 -translate-x-1/3 pointer-events-none">
            </div>

            <!-- 封面区域 -->
            <div class="relative shrink-0 z-10">
                <div
                    class="w-40 h-40 md:w-48 md:h-48 rounded-[2rem] flex flex-col items-center justify-center bg-gradient-to-br from-violet-500 to-fuchsia-600 shadow-xl shadow-violet-500/20 transform group-hover:scale-105 transition-transform duration-500 overflow-hidden text-white relative">
                    <!-- 噪点纹理 -->
                    <div class="absolute inset-0 opacity-20 bg-noise mix-blend-overlay"></div>

                    <div class="relative z-10 flex flex-col items-center">
                        <div
                            class="mb-2 inline-flex items-center gap-2 rounded-full bg-white/15 px-3 py-1 text-xs font-semibold tracking-[0.2em] uppercase">
                            <div class="i-fa6-solid-sparkles text-[10px]"></div>
                            For You
                        </div>
                        <div class="text-5xl font-black tracking-tight">私享</div>
                        <div class="mt-1 text-lg font-medium opacity-90">推荐歌单</div>
                    </div>

                    <!-- 装饰图标 -->
                    <div class="absolute -bottom-4 -right-4 opacity-10 rotate-12">
                        <div class="i-fa6-solid-headphones text-8xl"></div>
                    </div>
                </div>
            </div>

            <!-- 信息区域 -->
            <div class="flex-1 text-center md:text-left space-y-4 w-full z-10">
                <div>
                    <div class="flex items-center justify-center md:justify-start gap-2 mb-3">
                        <span
                            class="px-2.5 py-1 rounded-lg text-[10px] font-bold uppercase tracking-wider border border-violet-200 text-violet-600 bg-violet-50">
                            Personal Mix
                        </span>
                        <span class="text-xs font-bold text-gray-400 uppercase tracking-widest flex items-center gap-1">
                            <div class="i-fa6-solid-bolt text-yellow-400"></div>
                            Realtime
                        </span>
                    </div>
                    <h1 class="text-3xl md:text-5xl font-black text-gray-800 tracking-tight leading-tight">
                        私享推荐
                    </h1>
                    <p class="text-sm text-gray-500 mt-2 max-w-2xl font-medium leading-relaxed">
                        根据你的播放历史和喜欢偏好，实时生成更贴近你当前口味的推荐列表。
                    </p>
                </div>

                <!-- 元数据 -->
                <div class="flex items-center justify-center md:justify-start gap-4 text-sm text-gray-500 font-medium">
                    <div class="flex items-center gap-2">
                        <div class="w-6 h-6 rounded-full bg-gray-100 flex items-center justify-center">
                            <div class="i-fa6-solid-robot text-violet-500 text-xs"></div>
                        </div>
                        <span>MusicBox 推荐引擎</span>
                    </div>
                    <span class="text-gray-300">|</span>
                    <span>{{ songs.length }} 首歌曲</span>
                </div>

                <!-- 操作按钮 -->
                <div class="flex gap-3 justify-center md:justify-start pt-2">
                    <button @click="handlePlayAll" :disabled="loading || songs.length === 0"
                        class="px-8 py-3 bg-violet-600 text-white rounded-full font-bold shadow-lg shadow-violet-500/30 hover:bg-violet-700 hover:shadow-violet-500/40 hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed group-btn">
                        <div class="i-fa6-solid-play group-btn-hover:scale-110 transition-transform"></div>
                        播放全部
                    </button>
                </div>
            </div>
        </div>

        <!-- 歌曲列表 -->
        <SongList :is-empty="loading || songs.length === 0" :count="songs.length" theme="violet" action-label="操作">
            <template v-if="loading">
                <!-- 骨架屏加载 -->
                <div class="p-8 space-y-4">
                    <div v-for="n in 5" :key="n" class="flex items-center gap-4 animate-pulse">
                        <div class="w-8 h-4 bg-gray-100 rounded"></div>
                        <div class="w-10 h-10 bg-gray-100 rounded-lg"></div>
                        <div class="flex-1 space-y-2">
                            <div class="h-4 bg-gray-100 rounded w-1/3"></div>
                            <div class="h-3 bg-gray-100 rounded w-1/4"></div>
                        </div>
                    </div>
                </div>
            </template>

            <template v-else>
                <SongRow v-for="(song, index) in songs" :key="song.songId" :song="song" :index="index"
                    :is-active="isCurrentSong(song)" :is-playing="!!player.playing" theme="violet" @play="playSong"
                    @add-to-queue="player.addToQueue" @add-to-playlist="onAddToPlaylist" />
            </template>

            <template #empty>
                <div class="py-20 text-center">
                    <div class="inline-block p-6 rounded-full bg-violet-50 mb-4">
                        <div class="i-fa6-solid-music text-4xl text-violet-300"></div>
                    </div>
                    <p class="text-gray-500 font-medium">暂无推荐歌曲</p>
                    <p class="text-sm text-gray-400 mt-1">快去多听几首歌，让算法更懂你</p>
                </div>
            </template>
        </SongList>

        <!-- 添加到歌单对话框 -->
        <AddToPlaylistModal v-model="showAddModal" :song-id="currentSongId" @create-new="openCreateModal" />

        <PlaylistModal v-model="showCreateModal" mode="create" @success="handleCreateSuccess" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDailyRecommendations } from '@/api/songs'
import { usePlayerStore, type Song } from '@/stores/player'
import { formatSongs, showToast } from '@/utils/common'
import SongList from '@/components/song/SongList.vue'
import SongRow from '@/components/song/SongRow.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import PlaylistModal from '@/components/playlist/PlaylistModal.vue'

// --- 状态变量 ---
const player = usePlayerStore()
const loading = ref(true)
const songs = ref<Song[]>([])

// 弹窗状态变量管理
const showAddModal = ref(false)
const showCreateModal = ref(false)
const currentSongId = ref<number | string | null>(null)

// --- 业务方法 ---
const isCurrentSong = (song: Song) => {
    return player.currentSong?.songId === song.songId
}

const playSong = (index: number) => {
    const isSameSource = player.playlistSource === '私享推荐'
    const isSameLength = player.playList.length === songs.value.length

    // 如果已经在播放推荐列表，直接切歌
    if (isSameSource && isSameLength) {
        player.playByIndex(index)
    } else {
        // 否则替换播放列表
        player.setPlayList([...songs.value], '私享推荐')
        player.playByIndex(index)
        showToast('开始播放私享推荐', 'success')
    }
}

const handlePlayAll = () => {
    if (songs.value.length === 0) return
    playSong(0)
}

const onAddToPlaylist = (song: Song) => {
    currentSongId.value = song.songId || song.id
    showAddModal.value = true
}

const openCreateModal = () => {
    showAddModal.value = false
    showCreateModal.value = true
}

const handleCreateSuccess = () => {
    showToast('歌单创建成功')
    showAddModal.value = true
}

// --- 异步数据获取 ---
onMounted(async () => {
    loading.value = true
    try {
        // 页面直接消费后端实时推荐结果
        const res = await getDailyRecommendations(20)
        const rawSongs = res.data?.songs || []
        songs.value = formatSongs(rawSongs)
    } catch (e) {
        showToast('获取推荐失败', 'error')
    } finally {
        loading.value = false
    }
})
</script>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(20px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
