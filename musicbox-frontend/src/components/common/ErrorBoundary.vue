<template>
  <slot v-if="!hasError" />
  <div v-else class="flex items-center justify-center min-h-[200px] p-8">
    <div class="text-center max-w-[400px]">
      <!-- 错误图标 -->
      <div class="mb-4">
        <div class="i-fa6-solid-triangle-exclamation text-4xl text-amber-500"></div>
      </div>
      <!-- 错误标题 -->
      <h2 class="text-xl font-semibold text-gray-800 mb-2">出现了一些问题</h2>
      <!-- 错误信息 -->
      <p class="text-gray-500 text-sm mb-6">{{ errorMessage }}</p>
      <!-- 重试按钮 -->
      <button 
        class="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-lg text-sm font-medium shadow-lg shadow-blue-500/30 hover:shadow-blue-500/40 hover:-translate-y-0.5 active:scale-95 transition-all"
        @click="reset"
      >
        <div class="i-fa6-solid-rotate-right text-xs"></div>
        重试
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onErrorCaptured } from 'vue'

const props = defineProps<{
  fallbackMessage?: string
}>()

const hasError = ref(false)
const errorMessage = ref('')

onErrorCaptured((error: Error) => {
  hasError.value = true
  errorMessage.value = props.fallbackMessage || error.message || '加载组件时发生错误'
  console.error('ErrorBoundary 捕获到错误:', error)
  // 返回 false 阻止错误继续向上传播
  return false
})

function reset() {
  hasError.value = false
  errorMessage.value = ''
}

defineExpose({ reset })
</script>

