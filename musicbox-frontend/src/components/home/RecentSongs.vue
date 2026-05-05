<template>
    <div class="p-8 md:p-10 h-full flex flex-col justify-center">
        <!-- 头部区域 -->
        <div class="flex items-end justify-between mb-6">
            <h2
                class="text-2xl font-bold text-gray-800 relative after:content-[''] after:absolute after:left-0 after:-bottom-2 after:w-10 after:h-1 after:bg-blue-600 after:rounded-full">
                最近添加
            </h2>
            <button
                class="group flex items-center gap-2 px-5 py-2.5 bg-blue-600 text-white rounded-full text-sm font-medium transition-all shadow-lg shadow-blue-600/20 hover:bg-blue-700 hover:-translate-y-0.5 active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed disabled:shadow-none"
                @click="playAllSongs" :disabled="loading || !songs.length">
                <div class="i-fa6-solid-play text-xs group-hover:scale-110 transition-transform"></div>
                <span>播放全部</span>
            </button>
        </div>

        <!-- 骨架屏 Loading -->
        <div v-if="loading && songs.length === 0" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-5.5">
            <div v-for="n in 12" :key="n" class="flex items-center p-2 rounded-xl bg-white/40 border border-transparent">
                <div class="w-12 h-12 rounded-lg bg-gray-200 animate-pulse mr-3 flex-shrink-0"></div>
                <div class="flex-1 space-y-2">
                    <div class="h-4 bg-gray-200 rounded animate-pulse w-3/4"></div>
                    <div class="h-3 bg-gray-200 rounded animate-pulse w-1/2"></div>
                </div>
            </div>
        </div>

        <!-- 歌曲列表 -->
        <div v-else-if="songs.length > 0" class="w-full">
            <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-5.5">
                <SongItem 
                    v-for="song in songs" 
                    :key="song.songId" 
                    :song="song"
                    :is-active="player.currentSong?.songId === song.songId"
                    @play="playSong"
                    @add-to-queue="addToQueue"
                    @add-to-playlist="onAddToPlaylist"
                />
            </div>
        </div>

        <!-- 无歌曲状态 -->
        <Empty v-else message="暂无歌曲">
            <template #icon>
                <div class="i-fa6-solid-music text-4xl text-gray-300"></div>
            </template>
        </Empty>

        <!-- 收藏到歌单弹窗 -->
        <AddToPlaylistModal
            v-model="showAddModal"
            :song-id="currentSongId"
            @create-new="openCreateModal"
        />

        <!-- 新建歌单弹窗 -->
        <PlaylistModal
            v-model="showCreateModal"
            mode="create"
            @success="handleCreateSuccess"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { fetchSongs } from '@/api/songs'
import { usePlayerStore, type Song } from '@/stores/player'
import Empty from '@/components/ui/Empty.vue'
import SongItem from '@/components/song/SongItem.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import PlaylistModal from '@/components/playlist/PlaylistModal.vue'
import { showToast, formatSongs } from '@/utils/common'

const player = usePlayerStore()
const songs = ref<Song[]>([])
const loading = ref(true)
const showAddModal = ref(false)
const showCreateModal = ref(false)
const currentSongId = ref<number | string | null>(null)

// 播放单曲
const playSong = (song: Song) => {
    player.playSong(song)
}

// 添加到下一首播放
const addToQueue = (song: Song) => {
    player.addToQueue(song)
}

// 打开收藏弹窗
const onAddToPlaylist = (song: Song) => {
    currentSongId.value = song.songId || song.id
    showAddModal.value = true
}

const openCreateModal = () => {
    showAddModal.value = false
    showCreateModal.value = true
}

const handleCreateSuccess = () => {
    showToast('歌单创建成功，请重新选择添加到歌单')
    showAddModal.value = true
}

// 播放全部
const playAllSongs = () => {
    if (!songs.value.length) return
    player.setPlayList([...songs.value], '最近添加')
    player.playByIndex(0)
    showToast('开始播放全部歌曲')
}

// 加载数据
async function loadSongs() {
    loading.value = true
    try {
        const res = await fetchSongs({ 
            pageIndex: 1, 
            pageSize: 12, 
            sortBy: 'upload_at', 
            order: 'desc' 
        })
        const songList = res.data?.list || []
        
        if (songList.length > 0) {
            songs.value = formatSongs(songList)
        } else {
            songs.value = []
        }
    } catch (e) {
        console.error('获取歌曲数据失败:', e)
        showToast('加载失败', 'error')
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    nextTick(() => {
        loadSongs()
    })
})
</script>
