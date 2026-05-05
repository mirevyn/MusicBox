<template>
    <Teleport to="body">
        <Transition
            enter-active-class="transition duration-300 ease-out"
            enter-from-class="opacity-0"
            enter-to-class="opacity-100"
            leave-active-class="transition duration-200 ease-in"
            leave-from-class="opacity-100"
            leave-to-class="opacity-0"
        >
            <div v-if="modelValue" class="fixed inset-0 z-[110] flex items-center justify-center">
                <!-- 背景遮罩 -->
                <div class="absolute inset-0 bg-black/20" @click="onCancel"></div>

                <!-- 弹窗主体 -->
                <div
                    class="relative bg-white rounded-2xl w-full max-w-sm p-6 shadow-xl mx-8 transition-all duration-300"
                    :class="modelValue ? 'opacity-100 scale-100' : 'opacity-0 scale-95'"
                >
                    <div class="text-center">
                        <!-- 图标 -->
                        <div class="mx-auto w-12 h-12 rounded-full flex items-center justify-center mb-4 text-2xl"
                             :class="type === 'danger' ? 'bg-red-50 text-red-500' : 'bg-blue-50 text-primary'">
                            <div :class="iconClass"></div>
                        </div>

                        <h3 class="text-xl font-bold text-gray-800 mb-2">{{ title }}</h3>
                        <p class="text-sm text-gray-500 mb-6 leading-relaxed px-2">{{ content }}</p>

                        <!-- 按钮组 -->
                        <div class="flex gap-3 justify-center">
                            <button @click="onCancel"
                                class="flex-1 px-4 py-2.5 bg-gray-50 text-gray-600 rounded-xl font-medium hover:bg-gray-100 transition-colors">
                                {{ cancelText }}
                            </button>
                            <button @click="onConfirm"
                                class="flex-1 px-4 py-2.5 text-white rounded-xl font-bold shadow-lg transition-all hover:-translate-y-0.5 active:scale-95 flex items-center justify-center gap-2"
                                :class="type === 'danger' 
                                    ? 'bg-red-500 shadow-red-500/30 hover:bg-red-600' 
                                    : 'bg-primary shadow-primary/30 hover:bg-primary/90'">
                                <div v-if="loading" class="i-fa6-solid-spinner animate-spin text-xs"></div>
                                <span>{{ confirmText }}</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
    modelValue: boolean
    title?: string
    content?: string
    confirmText?: string
    cancelText?: string
    type?: 'info' | 'danger'
    loading?: boolean
}>(), {
    title: '确认操作',
    content: '确定要继续此操作吗？',
    confirmText: '确定',
    cancelText: '取消',
    type: 'info',
    loading: false
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

const iconClass = computed(() => {
    return props.type === 'danger' 
        ? 'i-fa6-solid-triangle-exclamation' 
        : 'i-fa6-solid-circle-info'
})

const onCancel = () => {
    emit('update:modelValue', false)
    emit('cancel')
}

const onConfirm = () => {
    emit('confirm')
}
</script>
