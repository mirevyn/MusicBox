<template>
    <div class="relative group/dropdown" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
        <!-- 触发器 -->
        <div class="cursor-pointer">
            <slot name="trigger"></slot>
        </div>

        <!-- 下拉菜单 -->
        <div class="absolute z-[100] bg-white rounded-xl shadow-2xl border border-gray-100 p-1.5 transition-all duration-200 transform origin-top-right scale-95"
            :class="[
                placementClasses,
                widthClass,
                isHovered && !isClicked ? 'opacity-100 visible scale-100' : 'opacity-0 invisible scale-95'
            ]"
            @click.stop="handleClick">
            <slot></slot>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, onBeforeUnmount } from 'vue'

const props = withDefaults(defineProps<{
    placement?: 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right'
    width?: string
}>(), {
    placement: 'bottom-right',
    width: 'w-48'
})

const isHovered = ref(false)
const isClicked = ref(false)
let hoverTimeout: number | null = null

const handleMouseEnter = () => {
    if (hoverTimeout) clearTimeout(hoverTimeout)
    isHovered.value = true
    isClicked.value = false
}

const handleMouseLeave = () => {
    hoverTimeout = window.setTimeout(() => {
        isHovered.value = false
    }, 150) // 添加一个小延迟，防止鼠标在触发元素和菜单之间移动时出现闪烁
}

const handleClick = () => {
    isClicked.value = true
    isHovered.value = false
}

onBeforeUnmount(() => {
    if (hoverTimeout) clearTimeout(hoverTimeout)
})

const placementClasses = computed(() => {
    switch (props.placement) {
        case 'bottom-left':
            return 'top-full left-0 mt-2'
        case 'bottom-right':
            return 'top-full right-0 mt-2'
        case 'top-left':
            return 'bottom-full left-0 mb-2'
        case 'top-right':
            return 'bottom-full right-0 mb-2'
        default:
            return 'top-full right-0 mt-2'
    }
})

const widthClass = computed(() => props.width)
</script>
