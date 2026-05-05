<template>
  <div class="space-y-6 animate-fade-in">
    <section class="flex flex-col gap-3 px-1 sm:flex-row sm:items-end sm:justify-between">
      <div>
        <div
          class="inline-flex items-center gap-2 rounded-full bg-primary/8 px-3 py-1 text-xs font-semibold text-primary">
          <div class="i-fa6-solid-fire text-[10px]"></div>
          热门歌单
        </div>
        <h1 class="mt-3 text-3xl font-bold tracking-tight text-slate-800">热门公开歌单</h1>
      </div>
    </section>

    <section class="rounded-[2rem] border border-white bg-white/55 p-5 shadow-sm backdrop-blur-md md:p-6">
      <div v-if="initialLoading"
        class="grid gap-3 grid-cols-[repeat(auto-fill,minmax(132px,1fr))] md:grid-cols-[repeat(auto-fill,minmax(148px,1fr))]">
        <div v-for="item in pageSize" :key="item" class="space-y-3">
          <div class="aspect-square w-full rounded-xl bg-slate-200/80 animate-pulse"></div>
          <div class="h-4 w-4/5 rounded-full bg-slate-200/80 animate-pulse"></div>
          <div class="h-3 w-1/2 rounded-full bg-slate-200/70 animate-pulse"></div>
        </div>
      </div>

      <template v-else-if="playlists.length > 0">
        <div
          class="grid gap-3 grid-cols-[repeat(auto-fill,minmax(132px,1fr))] md:grid-cols-[repeat(auto-fill,minmax(148px,1fr))]">
          <PlaylistCard v-for="playlist in playlists" :key="playlist.id" :playlist="playlist" compact
            :show-approved-badge="false" />
        </div>

        <div class="mt-6">
          <div v-if="loadingMore" class="flex items-center justify-center gap-2 text-sm text-slate-400">
            <div class="i-fa6-solid-spinner animate-spin"></div>
            正在加载更多
          </div>
          <div v-else-if="!hasMore" class="text-center text-sm text-slate-400">
            没有更多歌单了
          </div>
          <div ref="loadMoreTrigger" class="h-2 w-full"></div>
        </div>
      </template>

      <Empty v-else message="暂无热门歌单" sub-message="当前没有可展示的公开歌单">
        <template #icon>
          <div class="i-fa6-solid-compact-disc text-4xl text-slate-300"></div>
        </template>
      </Empty>
    </section>
  </div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { fetchRecommendedPlaylists } from '@/api/playlists'
import PlaylistCard from '@/components/playlist/PlaylistCard.vue'
import Empty from '@/components/ui/Empty.vue'
import type { Playlist } from '@/types/api'

const playlists = ref<Playlist[]>([])
const total = ref(0)
const pageIndex = ref(1)
const pageSize = 14
const hasMore = ref(true)
const initialLoading = ref(true)
const loadingMore = ref(false)
const loadMoreTrigger = ref<HTMLElement | null>(null)

let observer: IntersectionObserver | null = null

function mergePlaylists(nextPage: Playlist[]) {
  const seen = new Set(playlists.value.map((playlist) => playlist.id))
  const merged = [...playlists.value]

  nextPage.forEach((playlist) => {
    if (!seen.has(playlist.id)) {
      merged.push(playlist)
      seen.add(playlist.id)
    }
  })

  playlists.value = merged
}

async function loadPlaylists(reset = false) {
  if (initialLoading.value && !reset) return
  if (loadingMore.value && !reset) return
  if (!hasMore.value && !reset) return

  if (reset) {
    pageIndex.value = 1
    hasMore.value = true
    total.value = 0
    playlists.value = []
    initialLoading.value = true
  } else {
    loadingMore.value = true
  }

  try {
    const res = await fetchRecommendedPlaylists({
      pageIndex: pageIndex.value,
      pageSize,
    })
    const response = res.data
    const nextPage = response?.playlists || []

    if (pageIndex.value === 1) {
      playlists.value = nextPage
    } else {
      mergePlaylists(nextPage)
    }

    total.value = Number(response?.total || 0)
    hasMore.value = Boolean(response?.hasMore)
    pageIndex.value += 1
  } catch (error) {
    console.error('获取热门歌单失败:', error)
    hasMore.value = false
  } finally {
    initialLoading.value = false
    loadingMore.value = false
  }
}

async function initObserver() {
  await nextTick()

  if (!loadMoreTrigger.value) return

  observer?.disconnect()
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting) {
        loadPlaylists(false)
      }
    },
    {
      root: null,
      rootMargin: '200px 0px',
      threshold: 0,
    }
  )

  observer.observe(loadMoreTrigger.value)
}

onMounted(async () => {
  await loadPlaylists(true)
  await initObserver()
})

onBeforeUnmount(() => {
  observer?.disconnect()
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
