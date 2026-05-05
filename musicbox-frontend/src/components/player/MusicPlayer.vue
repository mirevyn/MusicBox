<template>
    <div
        class="fixed inset-0 z-[999] bg-[#050505] text-white select-none overflow-hidden flex flex-col pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] antialiased font-sans">

        <!-- 动态背景 -->
        <div class="absolute inset-0 z-0 pointer-events-none overflow-hidden bg-[#050505]">
            <Transition name="bg-fade">
                <div class="absolute inset-0 will-change-transform" :key="props.cover">
                    <img v-if="props.cover && !imgError" :src="props.cover"
                        class="w-full h-full object-cover opacity-40 blur-[50px] saturate-150 scale-[1.2] transform-gpu"
                        @error="imgError = true" />
                    <div v-else class="w-full h-full bg-gradient-to-br from-zinc-800 to-black opacity-40" />
                </div>
            </Transition>

            <div class="absolute inset-0 bg-gradient-to-b from-black/20 via-black/40 to-black/80 pointer-events-none" />
            <div
                class="absolute inset-0 bg-[radial-gradient(ellipse_at_center,transparent_20%,rgba(0,0,0,0.7)_100%)] pointer-events-none" />
        </div>

        <!-- 关闭按钮 -->
        <button @click="emit('close')" aria-label="Close Player"
            class="absolute top-[max(1rem,env(safe-area-inset-top))] right-4 lg:top-6 lg:right-6 z-50 w-9 h-9 lg:w-10 lg:h-10 flex items-center justify-center rounded-full bg-white/5 backdrop-blur-xl hover:bg-white/15 active:scale-90 active:bg-white/20 transition-all duration-300 ease-[cubic-bezier(0.25,1,0.3,1)] shadow-sm border border-white/5 group">
            <div
                class="i-fa6-solid-chevron-down text-base lg:text-lg text-white/60 group-hover:text-white transition-colors" />
        </button>

        <!-- 主体内容区 (横向滚动抽屉) -->
        <div ref="mainScrollRef"
            class="relative z-10 flex-1 flex flex-row items-center lg:justify-center gap-0 lg:gap-16 xl:gap-24 w-full max-w-7xl mx-auto h-full overflow-x-auto overflow-y-hidden lg:overflow-visible snap-x snap-mandatory scrollbar-hide scroll-smooth"
            @scroll.passive="handleMainScroll">

            <!-- 左侧：封面 + 控制面板 -->
            <div class="w-full h-full lg:h-auto flex-none snap-center snap-always flex flex-col items-center justify-center gap-6 lg:gap-8 px-5 sm:px-6 lg:px-0 lg:w-full lg:max-w-[420px] lg:flex-shrink-0 transition-all duration-[800ms] ease-[cubic-bezier(0.25,1,0.3,1)]"
                :class="props.isPlaying ? 'scale-100 opacity-100' : 'scale-[0.97] opacity-90'">

                <!-- 封面容器 -->
                <div
                    class="relative w-[70vw] sm:w-[50vw] max-w-[340px] aspect-square lg:w-[400px] lg:h-[400px] lg:max-w-none transform-gpu transition-transform duration-1000">
                    <!-- 封面主体 -->
                    <div class="relative w-full h-full rounded-2xl lg:rounded-[2rem] overflow-hidden border border-white/[0.08] z-10 bg-[#111] transform-gpu transition-shadow duration-700"
                        :class="props.isPlaying ? 'shadow-[0_20px_50px_-10px_rgba(0,0,0,0.6)]' : 'shadow-[0_10px_20px_-5px_rgba(0,0,0,0.4)]'">
                        <Transition name="fade">
                            <img v-if="props.cover && !imgError" :key="props.cover" :src="props.cover"
                                class="absolute inset-0 w-full h-full object-cover" @error="imgError = true" />
                            <div v-else
                                class="absolute inset-0 w-full h-full flex flex-col items-center justify-center text-white/10 bg-zinc-900/50">
                                <div class="i-fa6-solid-record-vinyl text-6xl lg:text-8xl opacity-50" />
                            </div>
                        </Transition>
                    </div>
                </div>

                <!-- 标题 + 喜欢按钮 -->
                <div class="w-full flex items-center justify-between px-1 lg:px-2">
                    <div class="flex flex-col items-start overflow-hidden mr-2 lg:mr-4 space-y-1">
                        <h1
                            class="text-2xl sm:text-3xl lg:text-3xl font-bold text-white tracking-[-0.02em] truncate leading-tight drop-shadow-md max-w-[11em] lg:max-w-[13em] py-0.5">
                            {{ props.title }}
                        </h1>
                        <p
                            class="text-white/55 font-medium text-sm sm:text-base tracking-[0.05em] uppercase truncate max-w-[14em] lg:max-w-[18em]">
                            {{ props.artist }}
                        </p>
                    </div>
                    <button aria-label="Toggle Like"
                        class="group relative flex items-center justify-center p-3 rounded-full hover:bg-white/10 transition-colors duration-300 active:scale-75 flex-shrink-0"
                        @click="emit('toggleLike')">
                        <div class="text-[1.35rem] transition-all duration-[400ms] ease-[cubic-bezier(0.175,0.885,0.32,1.275)]"
                            :class="props.liked
                                ? 'i-fa6-solid-heart text-[#ff2d55] scale-110 drop-shadow-[0_0_12px_rgba(255,45,85,0.5)]'
                                : 'i-fa6-regular-heart text-white/40 group-hover:text-white'" />
                    </button>
                </div>

                <!-- 进度 + 控制面板玻璃态容器 -->
                <div class="w-full relative mt-2">
                    <div
                        class="absolute inset-0 bg-white/[0.03] backdrop-blur-[30px] rounded-2xl lg:rounded-[2rem] border border-white/[0.06] shadow-[inset_0_1px_1px_rgba(255,255,255,0.08),0_8px_32px_rgba(0,0,0,0.3)] overflow-hidden">
                        <!-- 波形可视化 Canvas 容器 -->
                        <div
                            class="absolute bottom-0 left-0 right-0 h-20 lg:h-32 opacity-[0.25] pointer-events-none mix-blend-screen">
                            <canvas ref="canvasRef" class="w-full h-full" />
                        </div>
                    </div>

                    <div class="relative z-10 p-5 lg:p-6 flex flex-col gap-4 lg:gap-5">
                        <!-- 进度条 -->
                        <div class="w-full h-6 relative flex items-center group cursor-pointer z-20"
                            @mousemove="onProgressMouseMove" @mouseenter="isHoveringProgress = true"
                            @mouseleave="isHoveringProgress = false">

                            <input type="range" min="0" :max="props.duration || 100" step="0.01"
                                aria-label="Seek Progress" :value="isDragging ? dragValue : props.currentTime"
                                @pointerdown="onProgressPointerDown" @input="onDragInput" @pointerup="commitSeek"
                                style="touch-action: none;"
                                class="absolute top-0 left-0 w-full h-full opacity-0 cursor-pointer z-30" />

                            <div
                                class="absolute top-1/2 -translate-y-1/2 w-full h-[4px] bg-white/10 rounded-full overflow-hidden pointer-events-none transition-[height] duration-300 group-hover:h-[5px]">
                                <div class="h-full bg-white/95 rounded-full relative"
                                    :style="{ width: progressPercent, transition: isDragging ? 'none' : 'width 0.1s linear' }">
                                    <div
                                        class="absolute right-0 top-1/2 -translate-y-1/2 w-[4px] h-full bg-white shadow-[0_0_8px_2px_rgba(255,255,255,0.5)]" />
                                </div>
                            </div>

                            <!-- 拖动圆点 -->
                            <div class="absolute top-1/2 -translate-y-1/2 -translate-x-1/2 w-3.5 h-3.5 rounded-full bg-white pointer-events-none transition-[opacity,transform] duration-300 ease-[cubic-bezier(0.25,1,0.3,1)] z-10 shadow-[0_2px_8px_rgba(0,0,0,0.4)]"
                                :class="(isHoveringProgress || isDragging) ? 'opacity-100 scale-125' : 'opacity-0 scale-50'"
                                :style="{ left: progressPercent }" />

                            <!-- hover 预览时间气泡 -->
                            <div class="absolute -top-8 -translate-x-1/2 bg-[#222]/90 backdrop-blur-md px-2.5 py-1 rounded-md text-[11px] font-mono font-bold text-white pointer-events-none transition-[opacity,transform] duration-200 z-20 border border-white/10 shadow-lg"
                                :style="{ left: `${hoverPercent * 100}%` }"
                                :class="isHoveringProgress ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'">
                                {{ hoverTime }}
                            </div>

                            <!-- 时间文本 -->
                            <div
                                class="absolute -bottom-4 lg:-bottom-5 w-full flex justify-between text-[11px] text-white/40 font-semibold tracking-wider font-mono pointer-events-none">
                                <span>{{ formatTime(isDragging ? dragValue : props.currentTime) }}</span>
                                <span>{{ formatTime(props.duration) }}</span>
                            </div>
                        </div>

                        <!-- 控制按钮 -->
                        <div
                            class="flex items-center justify-center gap-5 sm:gap-7 lg:gap-9 relative z-10 mt-3 lg:mt-4">
                            <button @click="emit('toggleMode')" title="切换播放模式" aria-label="Toggle Mode"
                                class="w-10 h-10 flex items-center justify-center text-white/40 hover:text-white transition-colors duration-300 active:scale-90 rounded-full hover:bg-white/10">
                                <div :class="playModeIcon" class="text-lg" />
                            </button>

                            <button @click="emit('prev')" aria-label="Previous"
                                class="w-10 h-10 flex items-center justify-center text-white/60 hover:text-white transition-all duration-300 active:scale-75">
                                <div class="i-fa6-solid-backward-step text-[1.35rem]" />
                            </button>

                            <!-- 播放/暂停 -->
                            <button @click="emit('togglePlay')" aria-label="Play/Pause"
                                class="w-16 h-16 lg:w-[4.5rem] lg:h-[4.5rem] rounded-full flex items-center justify-center transition-all duration-[400ms] ease-[cubic-bezier(0.25,1,0.3,1)] active:scale-90 transform-gpu"
                                :class="props.isPlaying
                                    ? 'bg-white text-black shadow-[0_8px_20px_rgba(255,255,255,0.2)] hover:shadow-[0_12px_28px_rgba(255,255,255,0.3)] hover:scale-105'
                                    : 'bg-white/90 text-black shadow-[0_4px_12px_rgba(255,255,255,0.1)] hover:bg-white hover:scale-105'">
                                <div :class="props.isPlaying ? 'i-fa6-solid-pause' : 'i-fa6-solid-play ml-1'"
                                    class="text-2xl lg:text-[1.7rem]" />
                            </button>

                            <button @click="emit('next')" aria-label="Next"
                                class="w-10 h-10 flex items-center justify-center text-white/60 hover:text-white transition-all duration-300 active:scale-75">
                                <div class="i-fa6-solid-forward-step text-[1.35rem]" />
                            </button>

                            <!-- 音量调节 -->
                            <div class="relative flex items-center justify-center w-10 h-10 group/volume"
                                @wheel.prevent="handleVolumeWheel">
                                <button aria-label="Mute/Unmute"
                                    class="w-full h-full flex items-center justify-center text-white/40 group-hover/volume:text-white transition-colors duration-300 active:scale-90 rounded-full group-hover/volume:bg-white/10"
                                    title="滚轮调节音量 / 点击静音" @click="emit('setVolume', props.volume > 0 ? 0 : 0.8)">
                                    <div :class="volumeIcon" class="text-lg" />
                                </button>
                                <div class="absolute -top-9 left-1/2 -translate-x-1/2 bg-[#222]/90 backdrop-blur-md px-2 py-1 rounded-md text-[11px] font-mono font-bold text-white opacity-0 group-hover/volume:opacity-100 transition-all duration-200 pointer-events-none whitespace-nowrap border border-white/10 shadow-lg"
                                    :class="isHoveringProgress ? 'translate-y-0' : 'translate-y-1'">
                                    {{ Math.round(props.volume * 100) }}%
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 右侧：歌词显示区 -->
            <div
                class="w-full h-full lg:h-[85vh] flex-none snap-center snap-always lg:w-full lg:max-w-[48rem] lg:flex-1 relative flex flex-col items-center justify-center mask-linear-enhanced">

                <!-- 回到当前行按钮 -->
                <Transition name="fade">
                    <button v-if="isUserScrolling && currentLineIndex !== -1"
                        class="absolute bottom-10 left-1/2 -translate-x-1/2 z-30 flex items-center gap-2 px-4 py-2 rounded-full bg-white/15 backdrop-blur-md border border-white/20 text-white/90 text-xs font-bold tracking-wide hover:bg-white/25 active:scale-90 transition-all duration-300 shadow-[0_8px_20px_rgba(0,0,0,0.4)]"
                        @click="jumpToCurrentLine">
                        <div class="i-fa6-solid-circle-play text-sm text-white" />
                        回到当前
                    </button>
                </Transition>

                <!-- 歌词容器 -->
                <div class="w-full h-full relative overflow-hidden text-center cursor-grab active:cursor-grabbing transform-gpu"
                    ref="lyricWrapperRef" @wheel.passive="handleWheel" @touchstart.passive="handleTouchStart"
                    @touchmove.passive="handleTouchMove" @touchend.passive="handleTouchEnd">

                    <div class="absolute w-full px-6 sm:px-8 lg:px-16 transition-transform duration-[600ms] ease-[cubic-bezier(0.25,1,0.3,1)] will-change-transform"
                        :style="{ transform: `translateY(${lyricTranslateY}px)` }">

                        <div v-if="lyrics.length > 0"
                            class="flex flex-col items-center gap-6 lg:gap-8 pb-[50vh] font-bold">
                            <div v-for="(line, index) in lyrics" :key="line.startTime"
                                :ref="(el) => setLineRef(el, index)"
                                class="transition-all duration-[600ms] ease-[cubic-bezier(0.25,1,0.3,1)] origin-center cursor-pointer py-2 lg:py-3 max-w-full w-full flex flex-col items-center"
                                :class="index === currentLineIndex
                                    ? 'text-white opacity-100 drop-shadow-md z-10 scale-[1.02]'
                                    : 'text-white/40 opacity-30 hover:text-white/70 hover:opacity-70 scale-[0.95] z-0'"
                                @click.stop="handleLyricClick(line.startTime)">

                                <div class="w-full break-words"
                                    style="text-wrap: balance; word-break: break-word; white-space: normal;"
                                    :class="line.sizeClass">{{ line.text }}</div>
                                <div v-if="line.translation"
                                    class="w-full mt-2.5 text-base sm:text-lg lg:text-xl font-medium opacity-60 tracking-wider leading-relaxed break-words">
                                    {{ line.translation }}
                                </div>
                            </div>
                        </div>

                        <!-- 加载中占位 -->
                        <div v-else
                            class="flex items-center justify-center h-[50vh] text-white/30 font-semibold tracking-[0.4em] animate-pulse text-xs lg:text-sm">
                            LYRICS LOADING
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 移动端分页指示点 -->
        <div
            class="absolute bottom-5 left-0 w-full flex justify-center items-center gap-2 pointer-events-none lg:hidden z-20">
            <div v-for="i in 2" :key="i"
                class="rounded-full transition-all duration-500 ease-[cubic-bezier(0.25,1,0.3,1)]"
                :class="currentPage === i - 1 ? 'w-5 h-1.5 bg-white/80' : 'w-1.5 h-1.5 bg-white/20'" />
        </div>

        <div class="absolute inset-0 pointer-events-none noise-bg opacity-[0.04] z-[9999]" />
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, shallowRef, onUnmounted, nextTick } from 'vue';

interface Props {
    title?: string;
    artist?: string;
    cover?: string;
    liked?: boolean;
    lrc?: string;
    currentTime: number;
    duration: number;
    isPlaying: boolean;
    audioElement?: HTMLAudioElement | null;
    visible?: boolean;
    mode: number;
    volume: number;
}

const props = withDefaults(defineProps<Props>(), {
    title: 'Unknown',
    artist: 'Unknown',
    cover: '',
    currentTime: 0,
    duration: 0,
    isPlaying: false,
    lrc: '',
    visible: false,
    mode: 0,
    volume: 1.0,
});

const emit = defineEmits<{
    (e: 'close'): void;
    (e: 'togglePlay'): void;
    (e: 'prev'): void;
    (e: 'next'): void;
    (e: 'toggleLike'): void;
    (e: 'seek', time: number): void;
    (e: 'toggleMode'): void;
    (e: 'setVolume', vol: number): void;
}>();

// --- 状态与逻辑 ---
const imgError = ref(false);
watch([() => props.cover, () => props.visible], () => { imgError.value = false; });

const mainScrollRef = ref<HTMLElement | null>(null);
const currentPage = ref(0);
const handleMainScroll = () => {
    const el = mainScrollRef.value;
    if (!el) return;
    currentPage.value = el.scrollLeft > el.scrollWidth / 2 - 50 ? 1 : 0;
};

const formatTime = (s: number): string => {
    const total = Math.floor(s || 0);
    return `${Math.floor(total / 60)}:${String(total % 60).padStart(2, '0')}`;
};

// --- 进度条 ---
const isDragging = ref(false);
const dragValue = ref(0);

const onProgressPointerDown = (e: PointerEvent) => {
    const input = e.currentTarget as HTMLInputElement;
    dragValue.value = props.currentTime;
    isDragging.value = true;
    input.setPointerCapture(e.pointerId);
};

const onDragInput = (e: Event) => {
    dragValue.value = parseFloat((e.target as HTMLInputElement).value);
};

const commitSeek = () => {
    if (!isDragging.value) return;
    emit('seek', dragValue.value);
    requestAnimationFrame(() => { isDragging.value = false; });
};

const onWindowPointerUp = (e: PointerEvent) => {
    if (!isDragging.value) return;
    if (e.button !== 0) { isDragging.value = false; return; }
    commitSeek();
};
window.addEventListener('pointerup', onWindowPointerUp);

const hoverPercent = ref(0);
const isHoveringProgress = ref(false);
const onProgressMouseMove = (e: MouseEvent) => {
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
    hoverPercent.value = Math.min(Math.max((e.clientX - rect.left) / rect.width, 0), 1);
};
const hoverTime = computed(() => formatTime(hoverPercent.value * (props.duration || 0)));

const progressPercent = computed(() => {
    if (!props.duration) return '0%';
    const t = isDragging.value ? dragValue.value : props.currentTime;
    return `${Math.min((t / props.duration) * 100, 100)}%`;
});

const MODE_ICONS: Record<number, string> = {
    0: 'i-fa6-solid-arrow-right-arrow-left',
    1: 'i-fa6-solid-repeat',
    2: 'i-fa6-solid-shuffle',
};
const playModeIcon = computed(() => MODE_ICONS[props.mode] ?? MODE_ICONS[0]);

const volumeIcon = computed(() =>
    props.volume === 0 ? 'i-fa6-solid-volume-xmark' : props.volume > 0.5 ? 'i-fa6-solid-volume-high' : 'i-fa6-solid-volume-low'
);

let volumeThrottle = 0;
const handleVolumeWheel = (e: WheelEvent) => {
    const now = Date.now();
    if (now - volumeThrottle < 20) return;
    volumeThrottle = now;
    const vol = Math.min(1, Math.max(0, props.volume + (e.deltaY < 0 ? 0.05 : -0.05)));
    emit('setVolume', Number(vol.toFixed(2)));
};

// --- 歌词 ---
interface LyricLine {
    startTime: number;
    text: string;
    translation?: string;
    sizeClass: string;
}

const lyrics = shallowRef<LyricLine[]>([]);
const currentLineIndex = ref(-1);
const lineRefs = new Map<number, HTMLElement>();
const lyricTranslateY = ref(200);
const lyricWrapperRef = ref<HTMLElement | null>(null);
const isUserScrolling = ref(false);
let userScrollTimeout: ReturnType<typeof setTimeout> | null = null;

const setLineRef = (el: unknown, index: number) => {
    if (el instanceof HTMLElement) lineRefs.set(index, el);
    else lineRefs.delete(index);
};

const getBaseLyricSize = (len: number): string => {
    if (len < 10) return 'text-[2.2rem] sm:text-5xl lg:text-[3.2rem] leading-[1.2] tracking-tight';
    if (len < 20) return 'text-3xl sm:text-4xl lg:text-[2.6rem] leading-[1.25] tracking-tight';
    return 'text-2xl sm:text-3xl lg:text-3xl leading-[1.3] tracking-normal';
};

const parseLrc = (lrcText: string): LyricLine[] => {
    if (!lrcText) return [];
    const lyricMap = new Map<number, LyricLine>();
    const timeExp = /\[(\d{2}):(\d{2})(?:\.(\d+))?\]/g;

    for (const line of lrcText.split('\n')) {
        const trimmed = line.trim();
        if (!trimmed || /\[[a-zA-Z]+:/.test(trimmed)) continue;
        const matches = [...trimmed.matchAll(timeExp)];
        if (!matches.length) continue;
        const text = trimmed.replace(timeExp, '').trim();
        if (!text) continue;

        for (const match of matches) {
            const ms = parseInt(match[1]!) * 60_000 + parseInt(match[2]!) * 1_000 + parseFloat('0.' + (match[3] ?? '0')) * 1_000;
            const existing = lyricMap.get(ms);
            if (existing) {
                if (!existing.translation) existing.translation = text;
                else if (existing.text !== text && existing.translation !== text) existing.translation += `\n${text}`;
            } else {
                lyricMap.set(ms, { startTime: ms, text, translation: '', sizeClass: getBaseLyricSize(text.length) });
            }
        }
    }
    return Array.from(lyricMap.values()).sort((a, b) => a.startTime - b.startTime);
};

watch(() => props.lrc, (val) => {
    lyrics.value = parseLrc(val || '');
    lineRefs.clear();
    currentLineIndex.value = -1;
    lyricTranslateY.value = lyricWrapperRef.value?.clientHeight ?? window.innerHeight * 0.4;
}, { immediate: true });

const updateLyricPosition = (index: number) => {
    const activeEl = lineRefs.get(index);
    if (!activeEl) return;
    const containerH = lyricWrapperRef.value?.clientHeight ?? window.innerHeight * 0.8;
    lyricTranslateY.value = containerH / 2 - activeEl.offsetTop - activeEl.offsetHeight / 2;
};

const handleLyricClick = (startTime: number) => {
    emit('seek', startTime / 1000);
    if (userScrollTimeout) clearTimeout(userScrollTimeout);
    isUserScrolling.value = false;
    const index = lyrics.value.findIndex(l => l.startTime === startTime);
    if (index !== -1) {
        currentLineIndex.value = index;
        nextTick(() => updateLyricPosition(index));
    }
};

const findActiveLyricIndex = (ms: number): number => {
    const list = lyrics.value;
    let lo = 0, hi = list.length - 1, result = -1;
    while (lo <= hi) {
        const mid = (lo + hi) >> 1;
        if (ms >= list[mid]!.startTime) { result = mid; lo = mid + 1; }
        else hi = mid - 1;
    }
    return result;
};

watch(() => props.currentTime, (time) => {
    if (!props.visible || isDragging.value || isUserScrolling.value || !lyrics.value.length) return;
    const ms = time * 1000;
    const curr = currentLineIndex.value;
    const list = lyrics.value;

    if (curr >= 0 && ms >= list[curr]!.startTime && ms < (list[curr + 1]?.startTime ?? Infinity)) return;

    const activeIndex = findActiveLyricIndex(ms);
    if (activeIndex !== -1 && activeIndex !== curr) {
        currentLineIndex.value = activeIndex;
        updateLyricPosition(activeIndex);
    }
});

let touchStartY = 0, touchStartX = 0, initialTranslateY = 0, isHorizontalSwipe = false;

const startUserScroll = () => {
    isUserScrolling.value = true;
    if (userScrollTimeout) clearTimeout(userScrollTimeout);
};

const resetUserScroll = () => {
    if (userScrollTimeout) clearTimeout(userScrollTimeout);
    userScrollTimeout = setTimeout(() => {
        isUserScrolling.value = false;
        if (currentLineIndex.value !== -1) updateLyricPosition(currentLineIndex.value);
    }, 2500);
};

const handleWheel = (e: WheelEvent) => {
    startUserScroll();
    lyricTranslateY.value -= e.deltaY;
    resetUserScroll();
};

const handleTouchStart = (e: TouchEvent) => {
    const t = e.touches[0];
    if (!t) return;
    startUserScroll();
    touchStartY = t.clientY;
    touchStartX = t.clientX;
    initialTranslateY = lyricTranslateY.value;
    isHorizontalSwipe = false;
};

const handleTouchMove = (e: TouchEvent) => {
    const t = e.touches[0];
    if (!t) return;
    const dx = t.clientX - touchStartX;
    const dy = t.clientY - touchStartY;
    if (!isHorizontalSwipe && Math.abs(dx) > Math.abs(dy) && Math.abs(dx) > 10) {
        isHorizontalSwipe = true;
    }
    if (isHorizontalSwipe) return;
    lyricTranslateY.value = initialTranslateY + dy;
};

const handleTouchEnd = () => resetUserScroll();

const jumpToCurrentLine = () => {
    if (userScrollTimeout) clearTimeout(userScrollTimeout);
    isUserScrolling.value = false;
    if (currentLineIndex.value !== -1) updateLyricPosition(currentLineIndex.value);
};

// --- 波形可视化 ---
const canvasRef = ref<HTMLCanvasElement | null>(null);
let audioCtx: AudioContext | null = null;
let analyser: AnalyserNode | null = null;
let dataArray: Uint8Array<ArrayBuffer> | null = null;
let animationId: number | null = null;
let source: MediaElementAudioSourceNode | null = null;
let isVisualizerInit = false;

let resizeObserver: ResizeObserver | null = null;
let cachedCanvasW = 0;
let cachedCanvasH = 0;

const initVisualizer = () => {
    if (!props.audioElement || isVisualizerInit) return;
    try {
        const AudioCtxCtor: typeof AudioContext = window.AudioContext ?? (window as any).webkitAudioContext;
        audioCtx ??= new AudioCtxCtor();
        if (!analyser) {
            analyser = audioCtx.createAnalyser();
            analyser.fftSize = 1024;
            analyser.smoothingTimeConstant = 0.85;
        }
        if (!source) {
            source = audioCtx.createMediaElementSource(props.audioElement);
            source.connect(analyser);
            analyser.connect(audioCtx.destination);
        }
        dataArray = new Uint8Array(analyser.frequencyBinCount) as Uint8Array<ArrayBuffer>;
        isVisualizerInit = true;

        if (canvasRef.value && !resizeObserver) {
            resizeObserver = new ResizeObserver((entries) => {
                if (entries[0]) {
                    cachedCanvasW = entries[0].contentRect.width;
                    cachedCanvasH = entries[0].contentRect.height;
                }
            });
            resizeObserver.observe(canvasRef.value);
        }
    } catch (err) {
        console.warn('Visualizer initialization failed:', err);
        isVisualizerInit = true;
    }
};

const draw = () => {
    if (!props.visible || !props.isPlaying) {
        if (animationId) { cancelAnimationFrame(animationId); animationId = null; }
        return;
    }
    animationId = requestAnimationFrame(draw);

    const canvas = canvasRef.value;
    if (!analyser || !dataArray || !canvas || cachedCanvasW === 0 || cachedCanvasH === 0) return;

    analyser.getByteFrequencyData(dataArray);

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const dpr = window.devicePixelRatio || 1;
    const targetW = cachedCanvasW * dpr;
    const targetH = cachedCanvasH * dpr;

    if (canvas.width !== targetW || canvas.height !== targetH) {
        canvas.width = targetW;
        canvas.height = targetH;
        ctx.scale(dpr, dpr);
    }

    const w = cachedCanvasW;
    const h = cachedCanvasH;
    ctx.clearRect(0, 0, w, h);

    const dataSize = Math.floor(analyser.frequencyBinCount * 0.35);
    const sliceWidth = w / (dataSize - 1);

    const rawPoints: { x: number; y: number }[] = [{ x: 0, y: h }];
    for (let i = 0; i < dataSize; i++) {
        const pct = Math.pow((dataArray[i] ?? 0) / 255, 1.2);
        rawPoints.push({ x: i * sliceWidth, y: h - pct * h * 0.55 });
    }
    rawPoints.push({ x: w, y: h });

    const pts: { x: number; y: number }[] = [rawPoints[0]!];
    for (let i = 1; i < rawPoints.length - 1; i++) {
        const p = rawPoints[i - 1]!, c = rawPoints[i]!, n = rawPoints[i + 1]!;
        pts.push({ x: c.x, y: (p.y + c.y + n.y) / 3 });
    }
    pts.push(rawPoints[rawPoints.length - 1]!);

    ctx.beginPath();
    ctx.moveTo(pts[0]!.x, pts[0]!.y);
    for (let i = 0; i < pts.length - 1; i++) {
        const a = pts[i]!, b = pts[i + 1]!;
        ctx.quadraticCurveTo(a.x, a.y, (a.x + b.x) / 2, (a.y + b.y) / 2);
    }
    ctx.lineTo(pts[pts.length - 1]!.x, pts[pts.length - 1]!.y);

    // 第一层粗光晕
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.2)';
    ctx.lineWidth = 6;
    ctx.lineJoin = 'round';
    ctx.lineCap = 'round';
    ctx.stroke();

    // 第二层细主线
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.8)';
    ctx.lineWidth = 1.5;
    ctx.stroke();

    // 闭合底部进行填充
    ctx.lineTo(w, h);
    ctx.lineTo(0, h);
    ctx.closePath();

    const gradient = ctx.createLinearGradient(0, h, 0, h * 0.2);
    gradient.addColorStop(0, 'rgba(255,255,255,0.01)');
    gradient.addColorStop(1, 'rgba(255,255,255,0.3)');

    ctx.fillStyle = gradient;
    ctx.fill();
};

watch(
    [() => props.isPlaying, () => props.visible],
    ([playing, visible]) => {
        if (playing && visible) {
            nextTick(() => setTimeout(() => {
                if (!props.isPlaying || !props.visible) return;
                initVisualizer();
                if (audioCtx?.state === 'suspended') audioCtx.resume();
                if (!animationId) draw();
            }, 800));
        } else {
            if (animationId) { cancelAnimationFrame(animationId); animationId = null; }
        }
    },
    { immediate: true }
);

onUnmounted(() => {
    if (animationId) cancelAnimationFrame(animationId);
    if (userScrollTimeout) clearTimeout(userScrollTimeout);
    window.removeEventListener('pointerup', onWindowPointerUp);

    try {
        if (resizeObserver) {
            resizeObserver.disconnect();
            resizeObserver = null;
        }
        if (source) {
            source.disconnect();
            source = null;
        }
        if (analyser) {
            analyser.disconnect();
            analyser = null;
        }
        if (audioCtx && audioCtx.state !== 'closed') {
            audioCtx.close().catch(() => { });
            audioCtx = null;
        }
        isVisualizerInit = false;
    } catch { /* 忽略 */ }
});
</script>

<style scoped>
.noise-bg {
    background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.85' numOctaves='3' stitchTiles='stitch'/%3E%3CfeColorMatrix type='saturate' values='0'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)'/%3E%3C/svg%3E");
    background-repeat: repeat;
    background-size: 100px 100px;
}

.mask-linear-enhanced {
    mask-image: linear-gradient(to bottom, transparent 0%, black 15%, black 85%, transparent 100%);
    -webkit-mask-image: linear-gradient(to bottom, transparent 0%, black 15%, black 85%, transparent 100%);
}

.bg-fade-enter-active,
.bg-fade-leave-active {
    transition: opacity 1.5s ease-in-out;
}

.bg-fade-enter-from,
.bg-fade-leave-to {
    opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

.scrollbar-hide::-webkit-scrollbar {
    display: none;
}

.scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>