<template>
    <div class="space-y-6 animate-fade-in h-full flex flex-col">
        <!-- 头部信息卡片 -->
        <div v-if="loading" class="h-64 bg-gray-200/50 rounded-[2rem] animate-pulse"></div>
        <div v-else
            class="relative flex flex-col md:flex-row gap-6 items-center md:items-end bg-white/60 backdrop-blur-md p-6 md:p-8 rounded-[2rem] border border-white shadow-sm overflow-hidden shrink-0">

            <!-- 背景装饰 -->
            <div
                class="absolute top-0 right-0 w-64 h-64 bg-rose-200/20 rounded-full blur-3xl -translate-y-1/2 translate-x-1/3 pointer-events-none">
            </div>

            <!-- 封面区 -->
            <div class="relative group shrink-0 z-10">
                <div
                    class="w-40 h-40 md:w-48 md:h-48 bg-gradient-to-br from-rose-400 to-pink-500 rounded-3xl flex items-center justify-center text-white shadow-xl shadow-rose-500/20 transform group-hover:scale-105 transition-transform duration-500">
                    <div class="i-fa6-solid-heart text-5xl md:text-6xl drop-shadow-md animate-pulse"></div>
                </div>
            </div>

            <!-- 信息区 -->
            <div class="flex-1 text-center md:text-left space-y-4 w-full z-10">
                <div>
                    <h5
                        class="text-xs font-bold text-rose-500 uppercase tracking-widest mb-2 flex items-center justify-center md:justify-start gap-2">
                        <div class="i-fa6-solid-star"></div>
                        My Favorites
                    </h5>
                    <h1 class="text-3xl md:text-4xl font-extrabold text-gray-800 tracking-tight leading-tight">
                        我喜欢的音乐
                    </h1>
                </div>

                <div class="flex items-center justify-center md:justify-start gap-3 text-sm text-gray-500">
                    <div
                        class="flex items-center gap-1.5 p-1 pr-3 bg-white/50 rounded-full border border-white/50 shadow-sm">
                        <!-- 用户头像 -->
                        <div v-if="userInfo.avatarUrl && !avatarImgError" class="w-6 h-6 rounded-full overflow-hidden">
                            <img :src="resolveCover(userInfo?.avatarUrl)"
                    class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-105" @error="avatarImgError = true" />
                        </div>
                        <div v-else
                            class="w-6 h-6 rounded-full bg-gradient-to-br from-indigo-400 to-purple-500 text-white flex items-center justify-center text-xs font-bold">
                            {{ (userInfo.username || 'U').charAt(0).toUpperCase() }}
                        </div>
                        <span class="font-bold text-gray-700 text-xs">{{ userInfo.nickname || userInfo.username ||
                            'User' }}</span>
                    </div>
                    <span class="text-gray-300">|</span>
                    <span>{{ songs.length }} 首歌曲</span>
                </div>

                <!-- 操作按钮组 -->
                <div class="flex gap-3 justify-center md:justify-start pt-2">
                    <button @click="playAll" :disabled="songs.length === 0"
                        class="px-8 py-2.5 bg-rose-500 text-white rounded-full font-bold shadow-lg shadow-rose-500/30 hover:bg-rose-600 hover:shadow-rose-500/40 hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                        <div class="i-fa6-solid-play"></div> 播放全部
                    </button>
                </div>
            </div>
        </div>

        <!-- 歌曲列表区域 -->
        <div v-if="loading" class="p-4 bg-white/60 backdrop-blur-md rounded-[2rem] border border-white shadow-sm space-y-2">
            <div v-for="n in 5" :key="n" class="h-14 bg-white/20 rounded-xl animate-pulse"></div>
        </div>
        <SongList v-else :is-empty="songs.length === 0" :count="songs.length" theme="rose">
            <SongRow v-for="(song, index) in songs" :key="song.songId" :song="song" :index="index"
                :is-active="isCurrentSong(song)" :is-playing="!!player.playing" theme="rose"
                :show-remove="true" remove-label="取消喜欢"
                @play="playSong" @add-to-queue="player.addToQueue" @add-to-playlist="onAddToPlaylist"
                @remove="onRemoveLike" />

            <template #empty>
                <div class="w-20 h-20 bg-rose-50 rounded-full flex items-center justify-center mb-4">
                    <div class="i-fa6-regular-heart text-4xl text-rose-300"></div>
                </div>
                <p class="text-gray-500 font-medium">还没有喜欢的歌曲</p>
                <p class="text-sm text-gray-400 mt-1">点击爱心图标，将你喜爱的音乐收藏到这里</p>
                <router-link to="/playlists/hot"
                    class="mt-6 px-6 py-2 bg-white border border-gray-200 rounded-full text-sm text-gray-600 hover:border-rose-300 hover:text-rose-500 transition-colors shadow-sm">
                    去看看歌单
                </router-link>
            </template>
        </SongList>

        <!-- 弹窗对话框组 -->
        <AddToPlaylistModal
            v-model="showAddModal"
            :song-id="currentSongId"
            @create-new="openCreateModal"
        />

        <PlaylistModal
            v-model="showCreateModal"
            mode="create"
            @success="handleCreateSuccess"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getUserLikedSongs, toggleSongLike } from '@/api/songs'
import { getUserInfo } from '@/api/user'
import { usePlayerStore, type Song } from '@/stores/player'
import { resolveCover, showToast, formatSongs } from '@/utils/common'
import SongList from '@/components/song/SongList.vue'
import SongRow from '@/components/song/SongRow.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import PlaylistModal from '@/components/playlist/PlaylistModal.vue'

const player = usePlayerStore()
const loading = ref(true)
const songs = ref<Song[]>([])
const userInfo = ref<any>({})
const avatarImgError = ref(false)

watch(() => userInfo.value.avatarUrl, () => {
    avatarImgError.value = false
})

// 判断当前歌曲是否正在播放
const isCurrentSong = (song: Song) => player.currentSong?.songId === song.songId

// 加载数据
const initData = async () => {
    loading.value = true
    try {
        // 并行加载用户信息和喜欢的歌曲
        const [userRes, likesRes] = await Promise.all([
            getUserInfo().catch(() => ({ data: {} })), // 容错处理
            getUserLikedSongs()
        ])

        // 设置用户信息
        if (userRes && userRes.data) {
            userInfo.value = userRes.data
        }

        // 处理歌曲数据
        if (likesRes && likesRes.data) {
            const list = Array.isArray(likesRes.data) ? likesRes.data : []
            songs.value = formatSongs(list)
        }
    } catch (error) {
        console.error('加载喜欢列表失败:', error)
        showToast('加载失败，请重试', 'error')
    } finally {
        loading.value = false
    }
}

// 播放全部
const playAll = () => {
    if (songs.value.length === 0) return
    player.setPlayList([...songs.value])
    player.playByIndex(0)
    player.playlistSource = '我喜欢的音乐'
    showToast('开始播放收藏列表')
}

// 播放单曲
const playSong = (index: number) => {
    if (player.playlistSource === '我喜欢的音乐' && player.playList.length === songs.value.length) {
        player.playByIndex(index)
    } else {
        player.setPlayList([...songs.value])
        player.playByIndex(index)
        player.playlistSource = '我喜欢的音乐'
    }
}

// 添加到歌单业务逻辑控制
const showAddModal = ref(false)
const showCreateModal = ref(false)
const currentSongId = ref<number | string | null>(null)

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

const onRemoveLike = async (songId: number | string) => {
    if (!songId) return
    try {
        await toggleSongLike(songId)
        songs.value = songs.value.filter((song) => String(song.songId) !== String(songId))

        const newSet = new Set(player.likedSongs as Set<string>)
        newSet.delete(String(songId))
        player.likedSongs = newSet

        if (player.playlistSource === '我喜欢的音乐') {
            player.removeFromPlayList(songId)
        }

        showToast('已取消喜欢')
    } catch (error) {
        showToast('操作失败，请重试', 'error')
    }
}

onMounted(() => {
    initData()
})
</script>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(15px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
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
</style>
