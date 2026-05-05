<script setup lang="ts">
import { ref, onUnmounted } from 'vue'

export type ToastType = 'success' | 'error' | 'info' | 'warning'

const visible = ref(false)
const message = ref('')
const type = ref<ToastType>('success')
const progress = ref(100)

let timer: ReturnType<typeof setTimeout> | null = null
let rafId: number | null = null
let startTime = 0
let duration = 3000

const cancelPending = () => {
    if (timer) { clearTimeout(timer); timer = null }
    if (rafId) { cancelAnimationFrame(rafId); rafId = null }
}

// rAF 驱动进度，与浏览器帧对齐，不额外触发 Vue diff
const tickProgress = () => {
    const elapsed = Date.now() - startTime
    const next = Math.max(0, 100 - (elapsed / duration) * 100)
    // 只在值变化超过 0.5% 时才写入 ref，减少响应式触发频率
    if (Math.abs(progress.value - next) > 0.5) {
        progress.value = next
    }
    if (elapsed < duration) {
        rafId = requestAnimationFrame(tickProgress)
    }
}

const show = (msg: string, toastType: ToastType = 'success', dur = 3000) => {
    cancelPending()

    message.value = msg
    type.value = toastType
    progress.value = 100
    visible.value = true
    duration = dur
    startTime = Date.now()

    rafId = requestAnimationFrame(tickProgress)
    timer = setTimeout(() => {
        visible.value = false
        cancelPending()
    }, dur)
}

onUnmounted(cancelPending)

defineExpose({ show })
</script>

<template>
    <Transition enter-active-class="transition duration-400 ease-[cubic-bezier(0.23,1,0.32,1)]"
        enter-from-class="opacity-0 translate-y-[-10px] scale-98" enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition duration-300 ease-in" leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95 translate-y-[-5px]">
        <div v-if="visible" class="fixed top-12 left-1/2 -translate-x-1/2 z-[2000] pointer-events-none">

            <div
                class="relative flex items-center gap-3.5 px-6 py-3.5 rounded-[1.25rem] bg-white/95 backdrop-blur-xl border border-gray-100 min-w-[300px] max-w-md shadow-[0_15px_30px_-10px_rgba(0,0,0,0.08),0_0_1px_rgba(0,0,0,0.1)] overflow-hidden pointer-events-auto">

                <!-- 图标 -->
                <div class="flex-shrink-0"
                    :class="type === 'error' ? 'text-rose-500' : type === 'success' ? 'text-primary' : 'text-gray-500'">
                    <div class="text-base"
                        :class="type === 'error' ? 'i-fa6-solid-circle-xmark' : type === 'success' ? 'i-fa6-solid-circle-check' : 'i-fa6-solid-circle-info'" />
                </div>

                <!-- 消息 -->
                <p
                    class="flex-1 min-w-0 mr-2 text-[13px] font-bold tracking-tight text-gray-700 leading-snug break-words">
                    {{ message }}
                </p>

                <!-- 进度条 -->
                <div class="absolute bottom-0 left-0 right-0 h-[2px] bg-gray-50">
                    <div class="h-full opacity-60 transition-none" :style="{ width: `${progress}%` }"
                        :class="type === 'error' ? 'bg-rose-500' : 'bg-primary'" />
                </div>
            </div>
        </div>
    </Transition>
</template>