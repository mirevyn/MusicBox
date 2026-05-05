<template>
  <BaseCard @click="$router.push('/daily')" tag="Personal Mix" tag-icon="i-fa6-solid-music" title="私享推荐"
    description="基于你的听歌偏好，每日更新专属旋律" glow-class="from-blue-500/30 to-transparent" class="min-h-[200px]">
    <!-- 背景 -->
    <template #background>
      <div class="absolute inset-0 transition-transform duration-700 group-hover:scale-[1.03]"
        style="background: linear-gradient(135deg, #0d1117 0%, #0f1c3f 55%, #162350 100%);">
      </div>
      <!-- 网格纹理 -->
      <div class="absolute inset-0"
        style="background-image: linear-gradient(rgba(51,91,241,0.06) 1px, transparent 1px), linear-gradient(90deg, rgba(51,91,241,0.06) 1px, transparent 1px); background-size: 28px 28px;">
      </div>
      <!-- 光晕 -->
      <div class="absolute -top-8 -right-8 w-44 h-44 rounded-full blur-[70px]"
        style="background: rgba(51,91,241,0.22);">
      </div>
    </template>

    <!-- 底部 -->
    <template #footer>
      <div class="flex items-end justify-between">

        <!-- 频谱律动条 -->
        <div class="flex items-end gap-1 h-9 pb-1 px-1">
          <div v-for="(item, i) in bars" :key="i" class="w-[3px] rounded-t-sm" :class="[item.cls, { 'playing-bar': isPlaying }]"
            :style="`height:${item.h}%; animation-delay:${item.delay}ms`">
          </div>
        </div>

        <!-- 播放按钮 -->
        <button @click.stop="handlePlay" :disabled="loading"
          class="w-11 h-11 rounded-full flex items-center justify-center shadow-lg transition-all duration-300 group-hover:scale-105 active:scale-95 bg-primary shadow-primary/45 disabled:opacity-50">
          <div v-if="loading" class="i-svg-spinners-90-ring-with-bg text-sm text-white"></div>
          <div v-else :class="isPlaying ? 'i-fa6-solid-pause text-sm' : 'i-fa6-solid-play ml-0.5 text-sm'" class="text-white"></div>
        </button>

      </div>
    </template>
  </BaseCard>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseCard from '@/components/common/BaseCard.vue'
import { getDailyRecommendations } from '@/api/songs'
import { usePlayerStore } from '@/stores/player'
import { formatSongs, showToast } from '@/utils/common'

const player = usePlayerStore()
const loading = ref(false)

const isPlaying = computed(() => player.playlistSource === '私享推荐' && player.playing)

const handlePlay = async () => {
    if (player.playlistSource === '私享推荐') {
        player.togglePlay()
        return
    }

    if (loading.value) return
    loading.value = true
    try {
        const res = await getDailyRecommendations(20)
        const rawSongs = res.data?.songs || []
        const songs = formatSongs(rawSongs)

        if (songs.length > 0) {
            player.setPlayList(songs, '私享推荐')
            player.playByIndex(0)
            showToast('开始播放私享推荐', 'success')
        } else {
            showToast('暂无推荐歌曲', 'warning')
        }
    } catch (e) {
        showToast('获取推荐失败', 'error')
    } finally {
        loading.value = false
    }
}

// 频谱条数据：静态高度 + hover 动画配置
const bars = [
  { h: 30, delay: 0, cls: 'bar-static group-hover:bar-ani' },
  { h: 65, delay: 300, cls: 'bar-static group-hover:bar-ani' },
  { h: 45, delay: 150, cls: 'bar-static group-hover:bar-ani' },
  { h: 80, delay: 450, cls: 'bar-static group-hover:bar-ani' },
  { h: 55, delay: 100, cls: 'bar-static group-hover:bar-ani' },
  { h: 35, delay: 250, cls: 'bar-static group-hover:bar-ani' },
]
</script>

<style scoped>
/* 频谱条默认样式 */
.bar-static {
  background: rgba(255, 255, 255, 0.3);
  transition: background 0.3s;
}

/* hover 时变蓝并动起来 */
.group:hover .bar-static {
  background: theme('colors.primary');
  animation: spectrum 1s ease-in-out infinite alternate;
}

/* 各条不同延迟通过 animation-delay inline style 控制 */
.group:hover .bar-static:nth-child(1) {
  animation-duration: 0.8s;
}

.group:hover .bar-static:nth-child(2) {
  animation-duration: 1.2s;
}

.group:hover .bar-static:nth-child(3) {
  animation-duration: 0.9s;
}

.group:hover .bar-static:nth-child(4) {
  animation-duration: 1.5s;
}

.group:hover .bar-static:nth-child(5) {
  animation-duration: 1.1s;
}

.group:hover .bar-static:nth-child(6) {
  animation-duration: 0.75s;
}

/* 播放中动画 */
.playing-bar {
  background: theme('colors.primary') !important;
  animation: spectrum 1s ease-in-out infinite alternate !important;
}

.playing-bar:nth-child(1) { animation-duration: 0.8s !important; }
.playing-bar:nth-child(2) { animation-duration: 1.2s !important; }
.playing-bar:nth-child(3) { animation-duration: 0.9s !important; }
.playing-bar:nth-child(4) { animation-duration: 1.5s !important; }
.playing-bar:nth-child(5) { animation-duration: 1.1s !important; }
.playing-bar:nth-child(6) { animation-duration: 0.75s !important; }

@keyframes spectrum {
  0% {
    height: 25%;
    opacity: 0.7;
  }

  100% {
    height: 100%;
    opacity: 1;
  }
}
</style>
