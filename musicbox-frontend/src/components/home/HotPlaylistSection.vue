<template>
  <section class="space-y-5">
    <div class="flex items-end justify-between gap-4 px-1">
      <div class="space-y-2">
        <div
          class="inline-flex items-center gap-2 rounded-full bg-primary/8 px-3 py-1 text-xs font-semibold text-primary">
          <div class="i-fa6-solid-fire text-[10px]"></div>
          热门推荐
        </div>
        <div>
          <h2 class="text-2xl font-bold tracking-tight text-slate-800">热门歌单推荐</h2>
        </div>
      </div>

      <RouterLink :to="{ name: 'HotPlaylists' }"
        class="inline-flex items-center gap-2 rounded-full border border-slate-200 bg-white/80 px-4 py-2 text-sm font-medium text-slate-600 transition-colors hover:border-primary/30 hover:text-primary">
        查看更多
        <div class="i-fa6-solid-chevron-right text-[10px]"></div>
      </RouterLink>
    </div>

    <div class="rounded-[2rem] border border-white bg-white/55 p-6 shadow-sm backdrop-blur-md md:p-7">
      <div v-if="loading && playlists.length === 0"
        class="grid gap-3 grid-cols-[repeat(auto-fill,minmax(132px,1fr))] md:grid-cols-[repeat(auto-fill,minmax(148px,1fr))]">
        <div v-for="item in 7" :key="item" class="space-y-3">
          <div class="aspect-square w-full rounded-xl bg-slate-200/80 animate-pulse"></div>
          <div class="h-4 w-4/5 rounded-full bg-slate-200/80 animate-pulse"></div>
          <div class="h-3 w-1/2 rounded-full bg-slate-200/70 animate-pulse"></div>
        </div>
      </div>

      <div v-else-if="playlists.length > 0"
        class="grid gap-3 grid-cols-[repeat(auto-fill,minmax(132px,1fr))] md:grid-cols-[repeat(auto-fill,minmax(148px,1fr))]">
        <PlaylistCard v-for="playlist in playlists" :key="playlist.id" :playlist="playlist" compact
          :show-approved-badge="false" />
      </div>

      <Empty v-else message="暂无可推荐的公开歌单">
        <template #icon>
          <div class="i-fa6-solid-compact-disc text-4xl text-slate-300"></div>
        </template>
      </Empty>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { fetchRecommendedPlaylists } from '@/api/playlists'
import PlaylistCard from '@/components/playlist/PlaylistCard.vue'
import Empty from '@/components/ui/Empty.vue'
import type { Playlist } from '@/types/api'

const loading = ref(true)
const playlists = ref<Playlist[]>([])

async function loadRecommendedPlaylists() {
  loading.value = true
  try {
    const res = await fetchRecommendedPlaylists({ limit: 7 })
    playlists.value = res.data?.playlists || []
  } catch (error) {
    console.error('获取热门歌单推荐失败:', error)
    playlists.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecommendedPlaylists()
})
</script>
