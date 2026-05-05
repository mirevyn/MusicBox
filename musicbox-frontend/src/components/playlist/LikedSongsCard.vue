<template>
  <div 
    class="group relative overflow-hidden rounded-[1.5rem] transition-all duration-300 cursor-pointer bg-white border h-full flex flex-col hover:-translate-y-1"
    :class="isPlaying ? 'border-primary/30 shadow-xl shadow-primary/10 ring-1 ring-primary/20' : 'border-gray-100 hover:border-primary/20 hover:shadow-xl hover:shadow-primary/5'"
    @click="router.push('/favorites')"
  >
    <!-- 背景层 -->
    <div class="absolute inset-0 z-0 bg-gradient-to-br from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
    
    <!-- 装饰性元素 -->
    <div class="absolute -right-8 -bottom-8 w-48 h-48 bg-primary/10 rounded-full blur-3xl pointer-events-none transition-transform duration-700 group-hover:scale-110"></div>
    <div class="absolute top-0 right-0 p-6 text-primary/5 transition-transform duration-700 group-hover:scale-110 group-hover:-rotate-12">
      <div class="i-fa6-solid-heart text-8xl md:text-9xl"></div>
    </div>

    <!-- 内容层 -->
    <div class="relative z-10 flex flex-col h-full p-6 md:p-8">
      
      <!-- 顶部信息 -->
      <div class="flex items-center gap-4 mb-auto">
         <div class="w-12 h-12 rounded-[14px] flex items-center justify-center transition-all duration-300 shadow-sm"
              :class="isPlaying ? 'bg-primary text-white shadow-md shadow-primary/30 scale-105' : 'bg-primary/10 text-primary group-hover:bg-primary group-hover:text-white'">
            <div class="i-fa6-solid-heart text-xl"></div>
         </div>
         <div>
            <h3 class="text-xl md:text-2xl font-bold tracking-tight transition-colors flex items-center gap-2"
                :class="isPlaying ? 'text-primary' : 'text-gray-800 group-hover:text-primary'">
              我喜欢的音乐
              <!-- 播放指示器 -->
              <div v-if="isPlaying" class="flex items-end gap-[2px] h-3.5 mb-0.5 ml-1">
                <div class="w-[3px] bg-primary rounded-t-sm animate-[equalizer_0.8s_ease-in-out_infinite_alternate]" style="height: 100%;"></div>
                <div class="w-[3px] bg-primary rounded-t-sm animate-[equalizer_1.2s_ease-in-out_infinite_alternate]" style="height: 60%;"></div>
                <div class="w-[3px] bg-primary rounded-t-sm animate-[equalizer_0.9s_ease-in-out_infinite_alternate]" style="height: 80%;"></div>
                <div class="w-[3px] bg-primary rounded-t-sm animate-[equalizer_1.0s_ease-in-out_infinite_alternate]" style="height: 40%;"></div>
              </div>
            </h3>
            <p class="text-sm text-gray-400 font-medium mt-0.5">{{ count }} 首歌曲</p>
         </div>
      </div>

      <!-- 底部操作与展示区 -->
      <div class="flex items-end justify-between mt-8">
          <!-- 左侧：封面堆叠 -->
          <div class="flex flex-col gap-3">
              <!-- 堆叠歌曲封面 -->
              <div class="flex items-center -space-x-3 transform-gpu group-hover:translate-x-2 transition-transform duration-500">
                <template v-if="songs && songs.length > 0">
                  <div v-for="(song, index) in topSongs" :key="song.id || (song as any).songId"
                       class="group/cover w-10 h-10 md:w-12 md:h-12 rounded-full border-2 overflow-hidden shadow-sm relative transform-gpu transition-all duration-300 hover:-translate-y-1 hover:z-40 hover:shadow-md hover:border-primary/20"
                       :class="[
                         isPlaying && index === 0 ? 'border-primary/50 shadow-md shadow-primary/20' : 'border-white',
                         index === 0 ? 'z-30' : index === 1 ? 'z-20' : 'z-10'
                       ]">
                    <img :src="resolveUrl(song.coverUrl || (song as any).cover) || '/default-cover.png'" class="w-full h-full object-cover transition-transform duration-500 group-hover/cover:scale-110" :alt="song.title || (song as any).name" />
                  </div>
                </template>
                <template v-else>
                   <div v-for="_i in 3" :key="_i"
                       class="w-10 h-10 md:w-12 md:h-12 rounded-full border-2 border-white overflow-hidden bg-gray-50 flex items-center justify-center shadow-sm relative"
                       :class="_i === 1 ? 'z-30' : _i === 2 ? 'z-20' : 'z-10'">
                    <div class="i-fa6-solid-music text-gray-300 text-xs"></div>
                  </div>
                </template>
              </div>

               <!-- 最新歌曲摘要 -->
               <div class="space-y-0.5 max-w-[12rem] md:max-w-xs transition-opacity duration-300 opacity-70 group-hover:opacity-100 min-h-[1.25rem]">
                  <template v-if="songs && songs.length > 0">
                      <div v-for="(song, _i) in topSongs.slice(0, 1)" :key="song.id" class="text-xs md:text-sm line-clamp-1 text-gray-500">
                          <span class="font-medium text-gray-700">{{ song.title || (song as any).name }}</span>
                          <span class="text-gray-300 mx-1">-</span>
                          <span>{{ song.artist || (song as any).singer }}</span>
                      </div>
                  </template>
                  <template v-else>
                      <div class="text-xs md:text-sm text-gray-400">去添加你喜欢的音乐吧</div>
                  </template>
              </div>
          </div>

          <!-- 右侧：播放按钮 (主题色) -->
          <button
            class="w-14 h-14 md:w-16 md:h-16 bg-primary text-white rounded-full flex items-center justify-center transform active:scale-95 transition-all duration-300 shrink-0"
            :class="[
              isPlaying ? 'shadow-xl shadow-primary/40 scale-105' : 'shadow-lg shadow-primary/20 hover:scale-105 hover:shadow-xl hover:shadow-primary/30 group-hover:-translate-y-1',
              { 'opacity-50 !cursor-not-allowed hover:!scale-100 group-hover:!-translate-y-0 active:!scale-100': !songs || songs.length === 0 }
            ]"
            :disabled="!songs || songs.length === 0"
            @click.stop="playSongs"
            :title="!songs || songs.length === 0 ? '暂无歌曲' : (isPlaying ? '暂停' : '播放全部')"
          >
            <div :class="isPlaying ? 'i-fa6-solid-pause' : 'i-fa6-solid-play ml-1 md:ml-1.5'" class="text-xl md:text-2xl transition-all"></div>
          </button>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore, type Song } from '@/stores/player'
import { showToast, resolveUrl } from '@/utils/common'

const router = useRouter()

const props = defineProps<{
    count: number;
    songs: Song[];
}>()

const player = usePlayerStore()

const isPlaying = computed(() => player.playlistSource === '我喜欢的音乐' && player.playing)

const topSongs = computed(() => {
    return (props.songs || []).slice(0, 3)
})

const playSongs = () => {
    if (props.count === 0 || !props.songs || props.songs.length === 0) return
    
    if (player.playlistSource === '我喜欢的音乐') {
        player.togglePlay()
        return
    }

    player.setPlayList(props.songs)
    player.playByIndex(0)
    player.playlistSource = '我喜欢的音乐'
    showToast('开始播放喜欢的音乐')
}


</script>

<style scoped>
@keyframes equalizer {
  0% { transform: scaleY(0.3); transform-origin: bottom; }
  100% { transform: scaleY(1); transform-origin: bottom; }
}
</style>
