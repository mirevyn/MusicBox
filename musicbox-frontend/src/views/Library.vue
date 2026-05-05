<template>
    <div class="space-y-6 animate-fade-in h-full flex flex-col px-4 md:px-0">
        <!-- 顶部标题区 -->
        <div class="flex items-end justify-between pb-4 border-b border-gray-100">
            <div>
                <h1 class="text-3xl font-black text-gray-800 tracking-tight">我的音乐库</h1>
                <p class="text-sm text-gray-500 mt-2 font-medium">管理与发现你的音乐收藏</p>
            </div>

            <!-- 大圆角新建按钮 (移动至标题区右侧，可选) -->
            <button @click="openCreateModal"
                class="hidden md:flex items-center gap-2 px-5 py-2.5 bg-primary/10 text-primary hover:bg-primary hover:text-white rounded-full text-sm font-bold transition-colors duration-300">
                <div class="i-fa6-solid-plus"></div>
                <span>新建歌单</span>
            </button>
        </div>

        <!-- 主内容区 -->
        <div class="flex-1 pb-8">
            <!-- 加载中骨架屏 -->
            <div v-if="loading" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 md:gap-5">
                <div class="col-span-2 bg-gray-100 rounded-[1.5rem] animate-pulse h-48 md:h-full"></div>
                <div v-for="n in 10" :key="n" class="space-y-3">
                    <div class="aspect-square bg-gray-100 rounded-[1.5rem] animate-pulse"></div>
                    <div class="h-4 bg-gray-100 rounded w-3/4 animate-pulse"></div>
                    <div class="h-3 bg-gray-100 rounded w-1/2 animate-pulse"></div>
                </div>
            </div>

            <!-- 数据展示网格 -->
            <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 md:gap-5 relative">
                <!-- 1. 我喜欢的音乐 (固定占据前两列) -->
                <LikedSongsCard class="col-span-2 row-span-1 h-full" :count="likedCount" :songs="likedSongsList" />

                <!-- 2. 普通歌单卡片 -->
                <PlaylistCard v-for="list in playlists" :key="list.id" :playlist="list" />

                <!-- 3. 新建歌单入口卡片 (固定在网格末尾) -->
                <div @click="openCreateModal" class="group cursor-pointer space-y-3 h-full">
                    <div
                        class="aspect-square bg-gray-50/50 rounded-[1.5rem] border-2 border-dashed border-gray-200 flex flex-col items-center justify-center group-hover:bg-primary/5 group-hover:border-primary/30 transition-all duration-300 transform-gpu group-hover:-translate-y-1">
                        <div
                            class="w-12 h-12 rounded-full bg-white text-gray-400 shadow-sm flex items-center justify-center group-hover:bg-primary group-hover:text-white transition-colors duration-300 shadow-gray-200/50">
                            <div class="i-fa6-solid-plus text-xl"></div>
                        </div>
                    </div>
                    <div>
                        <h3 class="font-bold text-gray-700 group-hover:text-primary transition-colors text-sm">创建新歌单
                        </h3>
                    </div>
                </div>
            </div>
        </div>

        <!-- 弹窗组件 -->
        <PlaylistModal v-model="showCreateModal" @success="loadData" mode="create" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchMyPlaylists } from '@/api/playlists'
import { getUserLikedSongs } from '@/api/songs'
import { type Song } from '@/stores/player'
import { showToast, formatSongs } from '@/utils/common'
import PlaylistModal from '@/components/playlist/PlaylistModal.vue'
import PlaylistCard from '@/components/playlist/PlaylistCard.vue'
import LikedSongsCard from '@/components/playlist/LikedSongsCard.vue'

const loading = ref(true)
const playlists = ref<any[]>([])
const likedCount = ref(0)
const showCreateModal = ref(false)
const likedSongsList = ref<Song[]>([])

const loadData = async () => {
    loading.value = true
    try {
        const [playlistRes, likedRes] = await Promise.all([
            fetchMyPlaylists().catch(() => ({ code: 500, data: [] })),
            getUserLikedSongs().catch(() => ({ code: 500, data: [] }))
        ])

        // 这里的返回值可能是完整响应对象，也可能已经被拦截器解包。
        const pRes = playlistRes as any
        const lRes = likedRes as any

        if (pRes.code === 200) {
            playlists.value = pRes.data || []
        }

        if (lRes.code === 200) {
            const list = Array.isArray(lRes.data) ? lRes.data : []
            likedSongsList.value = formatSongs(list)
            likedCount.value = list.length
        }

    } catch (error) {
        console.error('加载媒体库失败', error)
        showToast('加载失败', 'error')
    } finally {
        loading.value = false
    }
}

const openCreateModal = () => {
    showCreateModal.value = true
}

onMounted(() => {
    loadData()
})
</script>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
