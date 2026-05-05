<template>
  <div class="flex flex-col gap-8">

    <!-- 顶部问候区 -->
    <div class="flex flex-col gap-1 px-1">
      <div class="flex items-center gap-2 text-blue-600 mb-1">
        <div class="i-fa6-solid-sun text-sm" v-if="isDaytime"></div>
        <div class="i-fa6-solid-moon text-sm" v-else></div>
        <span class="text-xs font-bold tracking-wider uppercase">{{ dateStr }}</span>
      </div>
      <h1 class="text-3xl font-bold text-slate-800 tracking-tight">
        {{ greeting }}
      </h1>
      <p class="text-slate-500 text-sm">准备好聆听今天的旋律了吗？</p>
    </div>

    <!-- 主布局 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 lg:gap-8">

      <!-- 左侧栏：私享推荐 + AI向导 -->
      <div class="lg:col-span-1 lg:sticky lg:top-24 z-10 flex flex-col gap-6 h-full">
        <!-- 私享推荐 -->
        <DiscoverCard class="flex-1 min-h-[200px]" />
        <!-- AI 灵感助手 -->
        <AICard class="flex-1 min-h-[200px]" />
      </div>

      <!-- 右侧：歌曲探索区 -->
      <div class="lg:col-span-2">
        <div class="bg-white/40 backdrop-blur-md rounded-[2rem] border border-white shadow-sm min-h-[360px] h-full relative">
          <RecentSongs />
        </div>
      </div>

    </div>

    <HotPlaylistSection />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, defineAsyncComponent } from 'vue'
import RecentSongs from '@/components/home/RecentSongs.vue'
import AICard from '@/components/home/AICard.vue'
import DiscoverCard from '@/components/home/DiscoverCard.vue'

// 异步加载下方组件
const HotPlaylistSection = defineAsyncComponent(() => 
  import('@/components/home/HotPlaylistSection.vue')
)

// 时间逻辑
const hour = new Date().getHours()
const isDaytime = hour >= 6 && hour < 18

// 动态问候语
const greeting = computed(() => {
  if (hour < 5) return '夜深了，注意休息'
  if (hour < 11) return '早上好，开启活力一天'
  if (hour < 13) return '中午好，休息一下'
  if (hour < 18) return '下午好，来点音乐提神'
  return '晚上好，享受放松时刻'
})

// 日期显示
const dateStr = ref('')

onMounted(() => {
  dateStr.value = new Date().toLocaleDateString('zh-CN', {
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
})
</script>
