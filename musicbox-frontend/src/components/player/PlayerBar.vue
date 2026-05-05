<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue'
import { usePlayerStore, PLAY_MODE } from '../../stores/player'
import MusicPlayer from '@/components/player/MusicPlayer.vue'
import SongCover from '@/components/common/SongCover.vue'
import PlayerQueue from '@/components/player/PlayerQueue.vue'
import { recordPlayHistory } from '@/api/songs'
import { resolveUrl, formatDuration, showToast } from '@/utils/common'

// --- 常量定义 ---
const MODE_ICONS: Record<number, string> = {
    [PLAY_MODE.SEQUENCE]: 'i-fa6-solid-arrow-right-arrow-left',
    [PLAY_MODE.LOOP]: 'i-fa6-solid-repeat',
    [PLAY_MODE.RANDOM]: 'i-fa6-solid-shuffle'
}

// --- 状态管理 ---
const player = usePlayerStore()
const audioRef = ref<HTMLAudioElement | null>(null)
const playerQueueRef = ref<InstanceType<typeof PlayerQueue> | null>(null)
const toggleBtnRef = ref<HTMLElement | null>(null)

const currentTime = ref(0)
const duration = ref(0)
const showLyrics = ref(false)
const audioRetryCount = ref(0) // 单首歌曲内的重试次数
const consecutivePlayErrors = ref(0) // 连续切歌失败次数
const lyricContent = ref('')
const isDragging = ref(false)
const retryTimer = ref<ReturnType<typeof setTimeout> | null>(null)
const lastReportedSongId = ref<number | string | null>(null)

// --- 计算属性 ---
const audioUrl = computed(() => resolveUrl(player.currentSong?.fileUrl))
const coverUrl = computed(() => resolveUrl(player.currentSong?.coverUrl))

const seekPercent = computed(() => {
    if (!duration.value) return '0%'
    return `${(currentTime.value / duration.value) * 100}%`
})

const volumePercent = computed(() => `${player.volume * 100}%`)
const playModeIcon = computed(() => MODE_ICONS[player.mode] || MODE_ICONS[PLAY_MODE.SEQUENCE])

const volumeIcon = computed(() => {
    if (player.volume === 0) return 'i-fa6-solid-volume-xmark'
    return player.volume > 0.5 ? 'i-fa6-solid-volume-high' : 'i-fa6-solid-volume-low'
})

// --- 点击外部关闭逻辑 ---
const handleClickOutside = (event: MouseEvent) => {
    if (!player.showPlayList) return

    const target = event.target as HTMLElement
    if (toggleBtnRef.value?.contains(target)) {
        return
    }

    // 如果点击的是播放列表内部，忽略
    const queueEl = playerQueueRef.value?.$el
    if (queueEl && queueEl.nodeType === 1 && queueEl.contains(target)) {
        return
    }

    player.showPlayList = false
}

// 移动端点击左侧展开播放器
const handleMobileClick = () => {
    if (window.innerWidth < 768) {
        showLyrics.value = true
    }
}

// 获取歌词
const fetchLyric = async (url?: string) => {
    if (!url) {
        lyricContent.value = ''
        return
    }
    try {
        const finalUrl = resolveUrl(url)
        if (!finalUrl) return
        const response = await fetch(finalUrl)

        if (response.ok) {
            lyricContent.value = await response.text()
        } else {
            lyricContent.value = ''
        }
    } catch (e) {
        console.warn('歌词加载失败', e)
        lyricContent.value = ''
    }
}

const tryPlayAudio = async () => {
    const audio = audioRef.value
    if (!audio || !audio.src) return

    try {
        await audio.play()
        player.playing = true
    } catch (error: any) {
        if (error.name === 'NotAllowedError') {
            audio.muted = true
            try {
                await audio.play()
                showToast('已静音播放，请与页面交互以开启声音', 'error')
            } catch {
                player.playing = false
            }
        } else if (error.name !== 'AbortError') {
            console.error('播放失败:', error)
            player.playing = false
        }
    }
}

const handleAudioError = () => {
    const audio = audioRef.value
    // 如果没有源或者源是当前页面(空src行为)，忽略
    if (!audio?.src || audio.src === window.location.href) {
        player.playing = false
        return
    }

    // 如果已经连续失败超过3首歌曲，强制停止，防止无限循环
    if (consecutivePlayErrors.value >= 3) {
        console.error("连续播放失败，触发熔断保护")
        showToast('多首歌曲无法播放，已停止', 'error')
        player.playing = false
        consecutivePlayErrors.value = 0 // 重置
        return
    }

    // 清除之前的定时器，防止叠加
    if (retryTimer.value) clearTimeout(retryTimer.value)

    // 单曲重试逻辑
    if (audioRetryCount.value < 3) {
        audioRetryCount.value++
        console.warn(`音频加载失败，正在重试 (${audioRetryCount.value}/3)...`)

        // 延时重试
        retryTimer.value = setTimeout(() => {
            if (audio) audio.load()
        }, 1500) as unknown as number
    } else {
        // 当前歌曲彻底失败，准备切下一首
        consecutivePlayErrors.value++ // 增加连续失败计数
        showToast(`加载失败，自动跳过`, 'error')

        if (player.playList.length > 1) {
            player.playNext()
        } else {
            player.playing = false
        }
    }
}

// --- 音频事件 ---
const onCanPlay = () => {
    if (audioRef.value) {
        duration.value = audioRef.value.duration || 0
        if (player.playing) {
            tryPlayAudio()
        }
    }
    audioRetryCount.value = 0
    consecutivePlayErrors.value = 0
    void reportCurrentSongHistory()
}

const onTimeUpdate = () => {
    if (audioRef.value && !isDragging.value) {
        currentTime.value = audioRef.value.currentTime
        if (!duration.value) duration.value = audioRef.value.duration || 0
    }
}

const onEnded = () => {
    if (player.mode === PLAY_MODE.LOOP && audioRef.value) {
        audioRef.value.currentTime = 0
        tryPlayAudio()
    } else {
        player.playNext()
    }
}

// 拖动逻辑
const onSeekInput = (e: Event) => {
    isDragging.value = true
    currentTime.value = parseFloat((e.target as HTMLInputElement).value)
}

const onSeekChange = (e: Event) => {
    const value = parseFloat((e.target as HTMLInputElement).value)
    if (audioRef.value) {
        audioRef.value.currentTime = value
        currentTime.value = value
    }
    requestAnimationFrame(() => {
        isDragging.value = false
    })
}

const handleFullScreenSeek = (time: number) => {
    if (audioRef.value) {
        audioRef.value.currentTime = time
        currentTime.value = time
    }
}

// --- Media Session API 集成 ---
const updateMediaSession = () => {
    if (!('mediaSession' in navigator) || typeof window.MediaMetadata === 'undefined') return

    try {
        const song = player.currentSong
        if (!song) {
            navigator.mediaSession.metadata = null
            return
        }

        navigator.mediaSession.metadata = new window.MediaMetadata({
            title: song.title || '未知歌曲',
            artist: song.artist || 'Music Box',
            album: 'Music Box',
            artwork: [
                { src: coverUrl.value, sizes: '512x512', type: 'image/jpeg' }
            ]
        })
    } catch (e) {
        console.warn('Media Session metadata update failed:', e)
    }
}

const updateMediaSessionState = () => {
    if (!('mediaSession' in navigator)) return
    try {
        navigator.mediaSession.playbackState = player.playing ? 'playing' : 'paused'
    } catch (e) {
        console.warn('Update playbackState failed:', e)
    }
}

const setupMediaSessionActions = () => {
    if (!('mediaSession' in navigator)) return

    const actionHandlers = [
        ['play', () => player.togglePlay()],
        ['pause', () => player.togglePlay()],
        ['previoustrack', () => player.playPrev()],
        ['nexttrack', () => player.playNext()],
        ['stop', () => {
            player.playing = false
            if (audioRef.value) audioRef.value.currentTime = 0
        }]
    ] as const

    for (const [action, handler] of actionHandlers) {
        try {
            navigator.mediaSession.setActionHandler(action, handler)
        } catch (error) {
            // 某些浏览器可能不支持部分 action，忽略错误
            console.debug(`Media Session action ${action} not supported`)
        }
    }
}

// --- 监听器 ---
watch(audioUrl, (newUrl) => {
    if (retryTimer.value) {
        clearTimeout(retryTimer.value)
        retryTimer.value = null
    }

    if (!newUrl || !audioRef.value) return

    // 重置单曲重试计数
    audioRetryCount.value = 0

    audioRef.value.load()
    if (player.playing) {
        tryPlayAudio()
    }
})

watch(() => player.currentSong, (newSong) => {
    if (newSong?.lyricUrl) {
        fetchLyric(newSong.lyricUrl)
    } else {
        lyricContent.value = ''
        if (!newSong) {
            currentTime.value = 0
            duration.value = 0
            showLyrics.value = false
        }
    }
    if (newSong) {
        player.syncCurrentSongLikeStatus();
        if (newSong.songId !== lastReportedSongId.value) {
            lastReportedSongId.value = null
        }
    }
    updateMediaSession()
}, { immediate: true })

watch(() => player.volume, (val) => {
    if (audioRef.value) audioRef.value.volume = val
})

watch(() => player.playing, (isPlaying) => {
    updateMediaSessionState()
    if (!audioRef.value) return
    isPlaying ? tryPlayAudio() : audioRef.value.pause()
})

watch(() => player.playList, (val) => {
    localStorage.setItem('playList', JSON.stringify(val))
}, { deep: true })

// --- 生命周期 ---
onMounted(() => {
    if (player.playList.length === 0) {
        try {
            const saved = localStorage.getItem('playList')
            if (saved) player.setPlayList(JSON.parse(saved))
        } catch (e) { console.error('解析播放列表失败', e) }
    }

    document.addEventListener('click', handleClickOutside)

    setupMediaSessionActions()
})

onBeforeUnmount(() => {
    if (retryTimer.value) clearTimeout(retryTimer.value)
    document.removeEventListener('click', handleClickOutside)
})

// 同一首歌只在本次切歌成功后上报一次，避免 canplay 重复触发刷历史。
const reportCurrentSongHistory = async () => {
    const token = localStorage.getItem('token')
    const song = player.currentSong
    if (!token || !song?.songId) return
    if (lastReportedSongId.value === song.songId) return

    try {
        await recordPlayHistory({
            songId: song.songId,
            duration: Math.round(duration.value || song.duration || 0),
            source: player.playlistSource || '播放列表'
        })
        lastReportedSongId.value = song.songId
    } catch (error) {
        console.warn('记录播放历史失败', error)
    }
}
</script>

<template>
    <!-- 播放器容器 -->
    <div
        class="fixed bottom-0 left-0 right-0 z-[100] h-[64px] md:h-20 bg-white/95 backdrop-blur-xl shadow-[0_-2px_10px_rgba(0,0,0,0.03)] flex flex-col justify-end transition-all duration-300">

        <!-- 进度条 -->
        <div class="absolute -top-1.5 md:-top-1 left-0 w-full h-3 md:h-2 group cursor-pointer z-20 flex items-center">
            <div class="w-full h-1 md:h-1 group-hover:h-1.5 relative bg-gray-200 transition-all duration-200">
                <div class="absolute top-0 left-0 h-full bg-blue-600 ease-linear"
                    :class="isDragging ? 'transition-none' : 'transition-all duration-100'"
                    :style="{ width: seekPercent }">
                    <div
                        class="absolute right-0 top-1/2 -translate-y-1/2 translate-x-1/2 w-3 h-3 bg-white border border-gray-200 rounded-full shadow-sm scale-0 group-hover:scale-100 transition-transform duration-200">
                    </div>
                </div>
            </div>
            <input type="range" min="0" :max="duration || 100" step="0.01" :value="currentTime"
                class="absolute top-0 left-0 w-full h-full opacity-0 cursor-pointer z-30" @input="onSeekInput"
                @change="onSeekChange">

            <!-- 时间提示 -->
            <div
                class="absolute -top-8 left-0 w-full px-2 hidden group-hover:flex justify-between text-xs font-medium text-gray-500 pointer-events-none">
                <span class="bg-white/90 px-1.5 py-0.5 rounded shadow-sm">{{ formatDuration(currentTime) }}</span>
                <span class="bg-white/90 px-1.5 py-0.5 rounded shadow-sm">{{ formatDuration(duration) }}</span>
            </div>
        </div>

        <!-- 主控制区域 -->
        <div
            class="h-full px-3 md:px-8 flex items-center justify-between gap-2 md:gap-4 max-w-[1400px] mx-auto w-full relative">

            <!-- 左侧：歌曲信息 -->
            <!-- 移动端设为可点击唤出全屏播放器，提升体验 -->
            <div class="flex items-center min-w-0 gap-3 md:gap-4 flex-1 md:flex-none md:w-[30%] cursor-pointer md:cursor-default"
                @click="handleMobileClick">
                <div
                    class="relative flex-shrink-0 w-10 h-10 sm:w-12 sm:h-12 md:w-14 md:h-14 rounded-md md:rounded-lg overflow-hidden shadow-sm bg-gray-100 flex items-center justify-center">
                    <SongCover :src="coverUrl" class="w-full h-full object-cover block"
                        style="width: 100%; height: 100%; object-fit: cover;" icon-class="text-xl" />
                </div>
                <div class="flex flex-col justify-center min-w-0 mr-1">
                    <div class="text-sm md:text-base font-semibold text-gray-800 truncate leading-tight"
                        :title="player.currentSong?.title">
                        {{ player.currentSong?.title || '未播放音乐' }}
                    </div>
                    <div class="text-[11px] md:text-sm text-gray-500 truncate leading-tight mt-0.5"
                        :title="player.currentSong?.artist">
                        {{ player.currentSong?.artist || 'Music Box' }}
                    </div>
                </div>
                <!-- 喜欢按钮：仅在足够宽的屏幕显示 -->
                <button
                    class="hidden md:flex items-center justify-center w-8 h-8 rounded-full text-gray-400 hover:text-rose-500 hover:bg-rose-50 transition-all active:scale-95 flex-shrink-0"
                    :class="{ 'text-rose-500': player.isLiked }" title="喜欢" @click.stop="player.toggleLike">
                    <div :class="player.isLiked ? 'i-fa6-solid-heart text-base' : 'i-fa6-regular-heart text-base'">
                    </div>
                </button>
            </div>

            <!-- 中间：桌面端专属主播放控制 (绝对居中防偏移) -->
            <div
                class="absolute left-1/2 -translate-x-1/2 hidden md:flex items-center justify-center gap-6 pointer-events-none">
                <button
                    class="pointer-events-auto w-10 h-10 flex items-center justify-center text-gray-600 hover:text-blue-600 hover:bg-gray-100 rounded-full transition-all active:scale-95"
                    title="上一曲" @click.stop="player.playPrev()">
                    <div class="i-fa6-solid-backward-step text-xl"></div>
                </button>
                <button
                    class="pointer-events-auto w-14 h-14 bg-blue-600 text-white rounded-full flex items-center justify-center shadow-lg shadow-blue-500/30 hover:shadow-blue-500/50 hover:scale-105 active:scale-95 transition-all duration-300 group"
                    title="播放/暂停" @click.stop="player.togglePlay">
                    <div :class="player.playing ? 'i-fa6-solid-pause' : 'i-fa6-solid-play ml-1'"
                        class="text-2xl transition-transform"></div>
                </button>
                <button
                    class="pointer-events-auto w-10 h-10 flex items-center justify-center text-gray-600 hover:text-blue-600 hover:bg-gray-100 rounded-full transition-all active:scale-95"
                    title="下一曲" @click.stop="player.playNext()">
                    <div class="i-fa6-solid-forward-step text-xl"></div>
                </button>
            </div>

            <!-- 右侧：功能按钮及移动端简化操作 -->
            <div class="flex items-center justify-end gap-2 md:gap-4 flex-shrink-0 md:w-[30%]">

                <!-- 移动端：紧凑的播放按钮 -->
                <button
                    class="md:hidden w-9 h-9 bg-blue-600 text-white rounded-full flex items-center justify-center shadow-md active:scale-95 transition-all"
                    title="播放/暂停" @click.stop="player.togglePlay">
                    <div :class="player.playing ? 'i-fa6-solid-pause' : 'i-fa6-solid-play ml-0.5'" class="text-[15px]">
                    </div>
                </button>

                <button ref="toggleBtnRef"
                    class="w-9 h-9 flex items-center justify-center text-gray-500 hover:text-gray-800 hover:bg-gray-100 rounded-full transition-all relative"
                    :class="{ 'text-blue-600 bg-blue-50 hover:bg-blue-100 hover:text-blue-600': player.showPlayList }"
                    title="播放列表" @click.stop="player.togglePlayList">
                    <div class="i-fa6-solid-list text-lg md:text-base"></div>
                    <span v-if="player.playList.length"
                        class="absolute -top-1 -right-1 bg-gray-200 text-[10px] text-gray-600 px-1 rounded-full min-w-[16px] text-center border border-white">
                        {{ player.playList.length }}
                    </span>
                </button>

                <!-- 其他只在桌面端显示的功能 -->
                <button
                    class="hidden md:flex w-9 h-9 items-center justify-center text-gray-500 hover:text-gray-800 hover:bg-gray-100 rounded-full transition-all"
                    :class="{ 'text-blue-600': player.mode === PLAY_MODE.LOOP }" title="播放模式"
                    @click.stop="player.togglePlayMode">
                    <div :class="playModeIcon" class="text-base"></div>
                </button>

                <div class="hidden lg:flex items-center gap-2 group/volume relative">
                    <button
                        class="w-9 h-9 flex items-center justify-center text-gray-500 hover:text-blue-600 rounded-full transition-colors"
                        @click.stop="player.setVolume(player.volume > 0 ? 0 : 0.8)">
                        <div :class="volumeIcon" class="text-base"></div>
                    </button>
                    <div
                        class="w-20 xl:w-24 h-1 bg-gray-200 rounded-full relative cursor-pointer group-hover/volume:h-1.5 transition-all">
                        <input type="range" min="0" max="1" step="0.01" v-model.number="player.volume"
                            class="absolute top-0 left-0 w-full h-full opacity-0 cursor-pointer z-10" title="音量">
                        <div class="h-full bg-blue-600 rounded-full group-hover/volume:bg-blue-700 transition-colors"
                            :style="{ width: volumePercent }"></div>
                    </div>
                </div>

                <button
                    class="hidden md:flex w-9 h-9 items-center justify-center text-gray-500 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-all"
                    title="展开播放器" @click.stop="showLyrics = !showLyrics">
                    <div class="i-fa6-solid-chevron-up text-base"></div>
                </button>
            </div>
        </div>

        <!-- 播放列表弹窗 -->
        <PlayerQueue ref="playerQueueRef" />

        <Teleport to="body">
            <Transition name="slide-up">
                <MusicPlayer v-show="showLyrics" :visible="showLyrics" :title="player.currentSong?.title"
                    :artist="player.currentSong?.artist" :cover="coverUrl" :current-time="currentTime"
                    :duration="duration" :is-playing="player.playing" :lrc="lyricContent" :liked="player.isLiked"
                    :audio-element="audioRef" :mode="player.mode" :volume="player.volume" @close="showLyrics = false"
                    @togglePlay="player.togglePlay" @prev="player.playPrev" @next="player.playNext"
                    @toggleLike="player.toggleLike" @toggleMode="player.togglePlayMode" @setVolume="player.setVolume"
                    @seek="handleFullScreenSeek" />
            </Transition>
        </Teleport>

        <audio ref="audioRef" :src="audioUrl" preload="auto" crossorigin="anonymous" class="hidden"
            @timeupdate="onTimeUpdate" @ended="onEnded" @error="handleAudioError" @canplay="onCanPlay"></audio>
    </div>
</template>

<style scoped>
/* 进度条动画延迟类 */
.animation-delay-75 {
    animation-delay: 75ms;
}

.animation-delay-150 {
    animation-delay: 150ms;
}

.slide-up-enter-active,
.slide-up-leave-active {
    transition: transform 0.5s cubic-bezier(0.32, 0.72, 0, 1);
    will-change: transform;
}

.slide-up-enter-from,
.slide-up-leave-to {
    transform: translateY(100%);
}
</style>