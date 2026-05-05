<template>
  <div class="relative overflow-hidden w-full h-full bg-gray-100 flex items-center justify-center shrink-0">
    <img 
      v-if="src && !isError" 
      :src="src" 
      :alt="alt || 'Cover'"
      class="w-full h-full object-cover transition-opacity duration-300" 
      loading="lazy"
      @error="handleError"
    />
    <div v-else class="flex items-center justify-center text-gray-300 w-full h-full">
      <div :class="[iconClass || 'text-xl', 'i-fa6-solid-music']"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  src?: string
  alt?: string
  iconClass?: string
}>()

const isError = ref(false)

const handleError = () => {
  isError.value = true
}

// 监听 src 变化，重置错误状态
watch(() => props.src, (newVal) => {
  isError.value = !newVal
})
</script>
