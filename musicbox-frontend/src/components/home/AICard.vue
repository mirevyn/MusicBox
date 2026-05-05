<template>
  <BaseCard tag="AI 音乐实验室" tag-icon="i-fa6-solid-wand-magic-sparkles" title="AI 音乐向导"
    description="接入 Ollama 大模型，为您提供基于情绪与场景的智能音乐推荐。" glow-class="from-blue-500/20 to-transparent" class="min-h-[200px]">
    <!-- 背景 -->
    <template #background>
      <div class="absolute inset-0 bg-[#0d1117] transition-transform duration-700 group-hover:scale-[1.03]"></div>
      <!-- 网格纹理 -->
      <div class="absolute inset-0 opacity-100"
        style="background-image: linear-gradient(rgba(51,91,241,0.07) 1px, transparent 1px), linear-gradient(90deg, rgba(51,91,241,0.07) 1px, transparent 1px); background-size: 28px 28px;">
      </div>
      <!-- 蓝色光晕 -->
      <div class="absolute -top-16 -right-16 w-48 h-48 rounded-full blur-[70px] transition-all duration-700 bg-primary/20"
        :class="{ '!opacity-60': false }"></div>
      <div class="absolute -bottom-12 -left-12 w-40 h-40 rounded-full blur-[60px] transition-all duration-700 bg-primary/10">
      </div>
    </template>

    <!-- 顶部右侧 -->
    <template #header-right>
      <div class="flex items-end gap-[3px] h-5 opacity-40 group-hover:opacity-90 transition-opacity duration-300">
        <span v-for="(h, i) in [55, 85, 65, 100, 70, 45]" :key="i"
          class="w-[3px] rounded-full bg-primary group-hover:animate-[bar-dance_1.2s_ease-in-out_infinite]"
          :style="`height:${h}%; animation-delay:${i * 80}ms`">
        </span>
      </div>
    </template>

    <!-- 底部按钮 -->
    <template #footer>
      <button @click.stop="openAIChat"
        class="w-full py-3 px-4 rounded-xl text-sm font-semibold flex items-center justify-center gap-2 transition-all active:scale-95 group/btn relative overflow-hidden bg-primary text-white shadow-lg shadow-primary/35 hover:bg-primary/90">
        <!-- 高光 -->
        <span class="absolute inset-0 pointer-events-none"
          style="background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, transparent 60%);">
        </span>
        <span>开启智能对话</span>
        <div class="i-fa6-solid-arrow-right text-xs transition-transform group-hover/btn:translate-x-1"></div>
      </button>
    </template>
  </BaseCard>
</template>

<script setup lang="ts">
import BaseCard from '@/components/common/BaseCard.vue'

function openAIChat() {
  window.dispatchEvent(new Event('toggle-ai-chat'))
}
</script>

<style scoped>
@keyframes bar-dance {

  0%,
  100% {
    transform: scaleY(0.4);
    opacity: 0.5;
  }

  50% {
    transform: scaleY(1);
    opacity: 1;
  }
}
</style>