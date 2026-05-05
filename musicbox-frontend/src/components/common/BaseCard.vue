<template>
  <div
    class="group relative overflow-hidden rounded-[2rem] transition-all duration-500 hover:-translate-y-2 cursor-pointer shadow-xl hover:shadow-2xl"
    :class="containerClass">
    <!-- 背景层 -->
    <div class="absolute inset-0 z-0">
      <slot name="background"></slot>
      <!-- 默认噪点纹理 -->
      <div class="absolute inset-0 opacity-[0.03] pointer-events-none mix-blend-overlay bg-noise"></div>
    </div>

    <!-- 内容层 -->
    <div class="relative z-10 flex flex-col h-full p-6 text-white">
      <!-- 顶部标签区域 -->
      <div v-if="$slots.tag || tag" class="flex items-center justify-between mb-4">
        <slot name="tag">
          <div
            class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/15 backdrop-blur-xl border border-white/20 text-[11px] font-semibold tracking-wide uppercase">
            <span v-if="tagIcon" :class="tagIcon" class="text-xs"></span>
            <span>{{ tag }}</span>
          </div>
        </slot>
        <slot name="header-right"></slot>
      </div>

      <!-- 主体内容 -->
      <div class="flex-1">
        <h3 v-if="title" class="text-2xl font-black mb-2 tracking-tight leading-tight">
          <slot name="title">{{ title }}</slot>
        </h3>
        <p v-if="description" class="text-white/70 text-sm leading-relaxed max-w-[90%]">
          <slot name="description">{{ description }}</slot>
        </p>
        <slot></slot>
      </div>

      <!-- 底部区域 -->
      <div class="mt-4">
        <slot name="footer"></slot>
      </div>
    </div>

    <!-- 悬浮时的外发光效果 -->
    <div
      class="absolute inset-x-0 bottom-0 h-1/2 bg-gradient-to-t opacity-0 group-hover:opacity-40 transition-opacity duration-500 -z-1"
      :class="glowClass"></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  containerClass?: string
  glowClass?: string
  tag?: string
  tagIcon?: string
  title?: string
  description?: string
}>()
</script>

<style scoped>
/* 可以在这里添加一些特殊的微动画 */
</style>
