<template>
    <div class="h-full w-full flex flex-col items-center justify-center py-10 px-4 text-center">
        <div class="mb-4 transition-transform hover:scale-110 duration-300">
            <slot name="icon">
                <div class="i-fa6-solid-box-open text-6xl text-gray-200/80"></div>
            </slot>
        </div>

        <p class="text-base text-gray-500 font-medium tracking-wide">{{ message }}</p>

        <p v-if="subMessage" class="mt-1 text-sm text-gray-400">{{ subMessage }}</p>

        <div v-if="$slots.action" class="mt-6">
            <slot name="action"></slot>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps({
    message: {
        type: String,
        default: '暂无数据'
    },
    subMessage: {
        type: String,
        default: ''
    }
});
</script>