<template>
  <div class="relative w-full max-w-full md:max-w-md" ref="containerRef">
    <!-- 搜索框容器 -->
    <div
      class="relative flex items-center bg-gray-100/60 hover:bg-gray-200/40 rounded-full transition-all duration-300 border border-transparent focus-within:border-primary/30 focus-within:bg-white focus-within:shadow-lg focus-within:shadow-primary/5 px-4 h-10 group"
    >
      <div class="i-fa6-solid-magnifying-glass text-gray-400 mr-2 group-focus-within:text-primary transition-colors text-sm"></div>
      <input
        v-model="keyword"
        type="text"
        class="bg-transparent border-none outline-none text-sm text-gray-700 w-full placeholder-gray-400 font-medium"
        placeholder="搜索歌曲、歌手..."
        @focus="onFocus"
        @keydown.esc="closeResults"
        @keyup.enter="handleSearch"
      />
      <button
        v-if="keyword"
        class="ml-2 w-5 h-5 flex items-center justify-center rounded-full hover:bg-gray-100 text-gray-400 hover:text-gray-600 transition-all font-bold"
        @click="clearSearch"
      >
        <div class="i-fa6-solid-xmark text-[10px]"></div>
      </button>
    </div>

    <!-- 搜索结果列表 -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 translate-y-2 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 translate-y-2 scale-95"
    >
      <div
        v-if="showResults && (loading || searchResults.length > 0 || lastKeyword)"
        class="absolute top-full right-0 md:left-0 mt-3 w-full bg-white rounded-2xl shadow-2xl border border-gray-100 overflow-hidden z-[110] transform-gpu origin-top"
      >
        <!-- 加载状态 -->
        <div v-if="loading" class="p-8 text-center text-gray-400">
          <div class="i-fa6-solid-spinner animate-spin text-primary text-2xl mb-3 mx-auto"></div>
          <p class="text-xs font-medium tracking-wide">正在寻找...</p>
        </div>

        <!-- 结果列表 -->
        <div v-else-if="searchResults.length > 0 || playlistResults.length > 0" class="py-2">
          
          <!-- 歌曲结果块 -->
          <div v-if="searchResults.length > 0">
            <div class="px-4 py-2 text-[10px] font-bold text-primary/60 uppercase tracking-widest bg-blue-50/30 mb-1 flex items-center justify-between">
              <div class="flex items-center gap-2">
                  <div class="i-fa6-solid-music text-[8px]"></div>
                  匹配歌曲
              </div>
              <button @click="handleSearch" class="hover:underline text-primary/80">查看全部</button>
            </div>
            <div class="px-2 mb-2">
              <button
                v-for="song in searchResults"
                :key="song.songId"
                class="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-50 transition-all group/item text-left"
                @click="handlePlay(song)"
              >
                <div class="relative w-11 h-11 shrink-0 shadow-sm border border-gray-100 flex items-center justify-center rounded-lg overflow-hidden">
                  <SongCover 
                      :src="resolveCover(song.coverUrl)" 
                      class="w-full h-full bg-gray-100 object-cover"
                  />
                  <div class="absolute inset-0 bg-black/30 w-full h-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
                      <div class="i-fa6-solid-play text-white text-xs"></div>
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-bold text-gray-800 truncate group-hover/item:text-primary transition-colors">
                    {{ song.title }}
                  </div>
                  <div class="text-[11px] text-gray-500 truncate mt-0.5">{{ song.artist }}</div>
                </div>
              </button>
            </div>
          </div>

          <!-- 歌单结果块 -->
          <div v-if="playlistResults.length > 0">
            <div class="px-4 py-2 text-[10px] font-bold text-orange-500/60 uppercase tracking-widest bg-orange-50/30 mb-1 flex items-center justify-between mt-2">
              <div class="flex items-center gap-2">
                  <div class="i-fa6-solid-compact-disc text-[8px]"></div>
                  匹配歌单
              </div>
              <button @click="handleSearchPlaylist" class="hover:underline text-orange-500/80">查看全部</button>
            </div>
            <div class="px-2 mb-2">
              <button
                v-for="pl in playlistResults"
                :key="pl.id"
                class="w-full flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-50 transition-all group/item text-left"
                @click="goToPlaylist(pl.id)"
              >
                <div class="relative w-11 h-11 shrink-0 shadow-sm border border-gray-100 flex items-center justify-center rounded-lg overflow-hidden">
                  <SongCover 
                      :src="resolveCover(pl.coverUrl)" 
                      class="w-full h-full bg-gray-100 object-cover"
                  />
                  <!-- 歌单角标指示 -->
                  <div class="absolute bottom-0 right-0 bg-orange-500/90 text-white p-0.5 rounded-tl shadow-sm">
                    <div class="i-fa6-solid-compact-disc text-[8px]"></div>
                  </div>
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-bold text-gray-800 truncate group-hover/item:text-orange-500 transition-colors">
                    {{ pl.title }}
                  </div>
                  <div class="text-[11px] text-gray-500 truncate mt-0.5 flex items-center gap-1">
                    <span>{{ pl.songCount || 0 }} 首</span>
                    <span class="text-gray-300">·</span>
                    <span>by {{ pl.user?.nickname || pl.user?.username || '用户' }}</span>
                  </div>
                </div>
              </button>
            </div>
          </div>

        </div>

        <!-- 未找到结果 -->
        <div v-else-if="lastKeyword" class="p-10 text-center">
          <div class="w-16 h-16 bg-gray-50 rounded-full flex items-center justify-center mb-4 mx-auto text-gray-200">
            <div class="i-fa6-solid-magnifying-glass text-3xl"></div>
          </div>
          <p class="text-sm text-gray-600 font-bold mb-1">未找到结果</p>
          <p class="text-xs text-gray-400 px-4 leading-relaxed">未找到与 "{{ lastKeyword }}" 相关的歌曲记录</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { onClickOutside, watchDebounced } from '@vueuse/core'
import { fetchSongs } from '@/api/songs'
import { searchPlaylists } from '@/api/playlists'
import { usePlayerStore, type Song } from '@/stores/player'
import type { Playlist } from '@/types/api'
import { resolveCover, formatSongs } from '@/utils/common'
import SongCover from '@/components/common/SongCover.vue'

const router = useRouter()
const keyword = ref('')
const lastKeyword = ref('')
const searchResults = ref<Song[]>([])
const playlistResults = ref<Playlist[]>([])
const loading = ref(false)
const showResults = ref(false)
const containerRef = ref(null)
const player = usePlayerStore()

// 使用 VueUse 的 onClickOutside 替代自定义指令
onClickOutside(containerRef, () => {
  closeResults()
})

const performSearch = async () => {
  if (!keyword.value.trim()) {
    searchResults.value = []
    lastKeyword.value = ''
    return
  }
  
  loading.value = true
  lastKeyword.value = keyword.value
  try {
    const [songRes, plRes] = await Promise.all([
      fetchSongs({ keyword: keyword.value, pageSize: 5 }),
      searchPlaylists({ keyword: keyword.value, pageSize: 3 })
    ]);

    const rawList = songRes.data?.list || []
    searchResults.value = formatSongs(rawList)
    
    playlistResults.value = (plRes as any).data?.playlists || []
  } catch (error) {
    console.error('Search error:', error)
  } finally {
    loading.value = false
  }
}

// 使用 watchDebounced 实现防抖
watchDebounced(
  keyword,
  (newVal) => {
    if (!newVal.trim()) {
      searchResults.value = []
      playlistResults.value = []
      showResults.value = false
      return
    }

    showResults.value = true
    performSearch()
  },
  { debounce: 400 }
)

const onFocus = () => {
  if (keyword.value.trim()) {
    showResults.value = true
  }
}

const closeResults = () => {
  showResults.value = false
}

const clearSearch = () => {
  keyword.value = ''
  searchResults.value = []
  playlistResults.value = []
  showResults.value = false
}

const handlePlay = (song: Song) => {
  player.playSong(song)
  closeResults()
}

const goToPlaylist = (id: number) => {
  router.push(`/playlist/${id}`)
  closeResults()
}

const handleSearch = () => {
  if (keyword.value.trim()) {
    router.push({ path: '/search', query: { q: keyword.value.trim() } })
    closeResults()
  }
}

const handleSearchPlaylist = () => {
  if (keyword.value.trim()) {
    router.push({ path: '/search', query: { q: keyword.value.trim(), tab: 'playlists' } })
    closeResults()
  }
}
</script>
