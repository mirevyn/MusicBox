<template>
    <div class="px-6 py-4 border-t border-gray-100 flex items-center justify-between flex-wrap gap-4 select-none">
        <!-- 统计摘要 -->
        <span class="text-sm text-gray-500">
            显示 {{ startItem }}-{{ endItem }} 共 {{ total }} 条
        </span>

        <div class="flex items-center gap-4">
            <!-- 页码按钮 -->
            <div class="flex gap-2">
                <button 
                    @click="changePage(current - 1)"
                    :disabled="current <= 1"
                    class="w-8 h-8 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 hover:border-primary hover:text-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:border-gray-200 disabled:hover:text-gray-500">
                    <div class="i-fa6-solid-chevron-left"></div>
                </button>

                <template v-for="(page, index) in displayedPages" :key="index">
                    <div 
                        v-if="page === '...'"
                        class="w-8 h-8 flex items-center justify-center text-gray-400 cursor-default">
                        ...
                    </div>
                    <button 
                        v-else
                        @click="changePage(page as number)"
                        :class="[
                            'w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold transition-all',
                            current === page 
                                ? 'bg-primary text-white shadow-md shadow-primary/20 cursor-default' 
                                : 'border border-gray-200 text-gray-600 hover:border-primary hover:text-primary'
                        ]">
                        {{ page }}
                    </button>
                </template>

                <button 
                    @click="changePage(current + 1)"
                    :disabled="current >= totalPages"
                    class="w-8 h-8 rounded-lg border border-gray-200 flex items-center justify-center text-gray-500 hover:border-primary hover:text-primary transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:border-gray-200 disabled:hover:text-gray-500">
                    <div class="i-fa6-solid-chevron-right"></div>
                </button>
            </div>

            <!-- 跳转页码 -->
            <div class="flex items-center gap-2 text-sm text-gray-500">
                <span>跳至</span>
                <div class="relative">
                    <input 
                        type="number" 
                        min="1" 
                        :max="totalPages"
                        v-model.number="jumpPage"
                        @keyup.enter="handleJump"
                        @blur="handleJump"
                        class="w-14 h-8 rounded-lg border border-gray-200 bg-gray-50 text-center text-gray-700 hover:bg-white hover:border-gray-300 focus:bg-white focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all appearance-none shadow-sm"
                    >
                </div>
                <span>页</span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

const props = withDefaults(defineProps<{
    current: number
    total: number
    pageSize?: number
}>(), {
    pageSize: 10
})

const emit = defineEmits<{
    (e: 'update:current', page: number): void
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))
const jumpPage = ref<number | ''>('')

// 计算当前显示的记录范围（例如 1-10）
const startItem = computed(() => ((props.current - 1) * props.pageSize) + 1)
const endItem = computed(() => Math.min(props.current * props.pageSize, props.total))

// 智能页码显示逻辑（例如 1 ... 4 5 6 ... 10）
const displayedPages = computed(() => {
    const total = totalPages.value
    const current = props.current
    const delta = 2 // 当前页周围显示的页码数量
    const range: (number | string)[] = []

    for (let i = 1; i <= total; i++) {
        // 始终显示第一页、最后一页以及当前页前后的页码
        if (
            i === 1 || 
            i === total || 
            (i >= current - delta && i <= current + delta)
        ) {
            range.push(i)
        } 
        // 如果存在间隙，添加省略号
        else if (
            range[range.length - 1] !== '...' && 
            (i < current - delta || i > current + delta)
        ) {
            range.push('...')
        }
    }
    return range
})

const changePage = (page: number) => {
    if (page === props.current) return
    if (page >= 1 && page <= totalPages.value) {
        emit('update:current', page)
    }
}

const handleJump = () => {
    if (jumpPage.value === '') return
    
    // 确保输入跳转值是有效的整数数字
    let target = Math.floor(Number(jumpPage.value))
    
    if (Number.isNaN(target)) {
        jumpPage.value = ''
        return
    }

    if (target < 1) target = 1
    if (target > totalPages.value) target = totalPages.value
    
    if (target !== props.current) {
        changePage(target)
    }
    
    jumpPage.value = '' // 跳转后重置输入框
}
</script>

<style scoped>
/* Chrome, Safari, Edge, Opera - 隐藏数字输入框箭头 */
input::-webkit-outer-spin-button,
input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  appearance: none; 
  margin: 0;
}

/* Firefox */
input[type=number] {
  -moz-appearance: textfield;
}
</style>
