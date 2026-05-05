<template>
  <div class="px-4 md:px-8 py-6 max-w-[1400px] mx-auto h-full flex flex-col animate-fade-in">
    <!-- 顶部标题区 -->
    <div class="mb-8 flex flex-col md:flex-row md:items-end justify-between gap-4">
      <div>
        <h2 class="text-3xl font-black text-gray-800 tracking-tight">搜索</h2>
        <p v-if="query" class="text-sm text-gray-500 font-medium mt-2">
          关于 <span class="text-primary font-bold">"{{ query }}"</span> 的搜索结果
        </p>
        <p v-else class="text-sm text-gray-500 font-medium mt-2">
          输入关键字探索更多音乐
        </p>
      </div>
    </div>

    <!-- 加载中骨架屏 -->
    <div v-if="loading" class="flex-1">
       <!-- Tab 骨架 -->
       <div class="flex gap-8 mb-6 border-b border-gray-100 pb-2">
          <div class="h-8 w-16 bg-gray-100 rounded-md animate-pulse"></div>
          <div class="h-8 w-16 bg-gray-100 rounded-md animate-pulse"></div>
       </div>
       <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 md:gap-4">
          <div v-for="n in 6" :key="n" class="h-16 bg-gray-50 rounded-xl animate-pulse"></div>
       </div>
    </div>

    <!-- 结果展示区 -->
    <div v-else class="flex-1 flex flex-col">
      <!-- 极简分类 Tabs -->
      <div class="flex items-center gap-8 mb-6 border-b border-gray-100/80">
        <button 
          v-for="(currentTab, key) in tabs" :key="key"
          class="pb-3 text-base md:text-lg font-bold transition-all relative group flex items-center"
          :class="activeTab === key ? 'text-gray-800' : 'text-gray-400 hover:text-gray-600'"
          @click="activeTab = key"
        >
          {{ currentTab.name }}
          <span class="ml-2 text-[10px] md:text-xs px-2 py-0.5 rounded-full font-medium transition-colors"
                :class="activeTab === key ? 'bg-primary/10 text-primary' : 'bg-gray-100 text-gray-400 group-hover:bg-gray-200'">
            {{ currentTab.count }}
          </span>
          <span
            class="absolute bottom-0 left-0 w-0 h-1 bg-primary rounded-t-full transition-all duration-300 group-hover:w-full"
            :class="{ '!w-full': activeTab === key }"></span>
        </button>
      </div>

      <!-- 歌曲结果 -->
      <div v-if="activeTab === 'songs'" class="flex-1 pb-8">
        <div v-if="songs.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3 md:gap-4 lg:pr-8">
          <SongItem
            v-for="song in songs"
            :key="song.songId"
            :song="song"
            :is-active="player.currentSong?.songId === song.songId"
            @play="handlePlay"
            @add-to-queue="player.addToQueue"
            @add-to-playlist="openAddToPlaylist"
            class="bg-white border border-transparent hover:border-primary/10 hover:shadow-md hover:shadow-primary/5 transition-all duration-300 rounded-xl"
          />
        </div>
        <div v-else class="flex flex-col items-center justify-center py-20 text-gray-400 bg-gray-50/50 rounded-[2rem] border border-dashed border-gray-100">
           <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center shadow-sm mb-4 text-gray-300">
              <div class="i-fa6-solid-music text-2xl"></div>
           </div>
           <p class="font-medium text-sm">未找到相关歌曲</p>
        </div>
      </div>

      <!-- 歌单结果 -->
      <div v-if="activeTab === 'playlists'" class="flex-1 pb-8">
        <div v-if="playlists.length > 0" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 md:gap-5">
          <PlaylistCard
            v-for="pl in playlists"
            :key="pl.id"
            :playlist="pl"
            :show-approved-badge="false"
            @click="goToPlaylist(pl.id)"
          />
        </div>
        <div v-else class="flex flex-col items-center justify-center py-20 text-gray-400 bg-gray-50/50 rounded-[2rem] border border-dashed border-gray-100">
           <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center shadow-sm mb-4 text-gray-300">
              <div class="i-fa6-solid-compact-disc text-2xl"></div>
           </div>
           <p class="font-medium text-sm">未找到相关歌单</p>
        </div>
      </div>
    </div>

    <!-- 收藏弹窗 -->
    <AddToPlaylistModal
        v-model="showAddModal"
        :song-id="selectedSongId"
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
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { fetchSongs } from '@/api/songs'
import { searchPlaylists } from '@/api/playlists'
import { usePlayerStore, type Song } from '@/stores/player'
import type { Playlist } from '@/types/api'
import { formatSongs, showToast } from '@/utils/common'
import SongItem from '@/components/song/SongItem.vue'
import PlaylistCard from '@/components/playlist/PlaylistCard.vue'
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue'
import PlaylistModal from '@/components/playlist/PlaylistModal.vue'

const route = useRoute()
const router = useRouter()
const player = usePlayerStore()
const query = computed(() => route.query.q as string || '')

const activeTab = ref<'songs' | 'playlists'>((route.query.tab as 'songs' | 'playlists') || 'songs')
const songs = ref<Song[]>([])
const playlists = ref<Playlist[]>([])
const loading = ref(false)
const showAddModal = ref(false)
const showCreateModal = ref(false)
const selectedSongId = ref<number | string>('')

const tabs = computed(() => ({
  songs: { name: '单曲', count: songs.value.length },
  playlists: { name: '歌单', count: playlists.value.length }
}))

const search = async () => {
  if (!query.value) {
    songs.value = []
    playlists.value = []
    return
  }

  loading.value = true
  try {
    const [songRes, plRes] = await Promise.all([
      fetchSongs({ keyword: query.value, pageSize: 30 }),
      searchPlaylists({ keyword: query.value, pageSize: 30 })
    ])

    const rawList = songRes.data?.list || []
    songs.value = formatSongs(rawList)

    playlists.value = (plRes as any).data?.playlists || []

    // 智能切换默认展示有结果的Tab (仅当URL没有明确指定tab时)
    const urlTab = route.query.tab as string
    if (urlTab === 'songs' || urlTab === 'playlists') {
      activeTab.value = urlTab
    } else if (songs.value.length === 0 && playlists.value.length > 0) {
      activeTab.value = 'playlists'
    } else {
      activeTab.value = 'songs'
    }

  } catch (error) {
    console.error('Search page error:', error)
  } finally {
    loading.value = false
  }
}

const goToPlaylist = (id: number) => {
  router.push(`/playlist/${id}`)
}

const handlePlay = (song: Song) => {
  player.setPlayList(songs.value, '搜索结果')
  player.playSong(song)
}

const openAddToPlaylist = (song: Song) => {
  selectedSongId.value = song.songId
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

watch(query, () => {
  search()
}, { immediate: true })

watch(() => route.query.tab, (newTab) => {
  if (newTab === 'songs' || newTab === 'playlists') {
    activeTab.value = newTab
  }
})

onMounted(() => {
  search()
})
</script>
