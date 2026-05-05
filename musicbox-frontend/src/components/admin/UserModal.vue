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
            <div v-if="modelValue" class="fixed inset-0 z-[110] flex items-center justify-center p-4">
                <!-- 背景遮罩 -->
                <div class="absolute inset-0 bg-black/50 transition-opacity" @click="close"></div>

                <!-- 模态框 -->
                <div class="bg-white rounded-2xl shadow-xl w-full max-w-md relative z-10 flex flex-col transition-all">
                    <!-- 头部 -->
                    <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between rounded-t-2xl">
                        <h3 class="text-lg font-bold text-gray-800">
                            {{ isEdit ? '编辑用户' : '添加用户' }}
                        </h3>
                        <button @click="close"
                            class="w-8 h-8 rounded-full flex items-center justify-center text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors">
                            <div class="i-fa6-solid-xmark text-sm"></div>
                        </button>
                    </div>

                    <!-- 主体 -->
                    <div class="p-6 overflow-visible">
                        <form id="user-form" @submit.prevent="submit" class="space-y-4">
                            <!-- 用户名 -->
                            <div class="space-y-1.5">
                                <label class="text-sm font-medium text-gray-700">用户名</label>
                                <input v-model="form.username" type="text" placeholder="请输入用户名" required
                                    :disabled="isEdit"
                                    :class="[isEdit ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : 'bg-gray-50 focus:bg-white focus:border-primary focus:ring-4 focus:ring-primary/10']"
                                    class="w-full border border-gray-200 rounded-xl px-4 py-2.5 text-sm transition-all outline-none" />
                            </div>

                            <!-- 密码 -->
                            <div class="space-y-1.5">
                                <label class="text-sm font-medium text-gray-700">
                                    {{ isEdit ? '密码 (留空不修改)' : '密码' }}
                                </label>
                                <input v-model="form.password" type="password" :placeholder="isEdit ? '••••••' : '请输入密码 (至少6位)'"
                                    :required="!isEdit" minlength="6"
                                    class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm focus:bg-white focus:border-primary focus:ring-4 focus:ring-primary/10 transition-all outline-none" />
                            </div>

                            <!-- 角色 -->
                            <div class="space-y-1.5">
                                <label class="text-sm font-medium text-gray-700">角色</label>
                                <Dropdown width="w-full">
                                    <template #trigger>
                                        <div class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm flex items-center justify-between cursor-pointer hover:bg-white focus:bg-white transition-all">
                                            <span class="text-gray-700">{{ form.role === 'Admin' ? '管理员' : '普通用户' }}</span>
                                            <div class="i-fa6-solid-chevron-down text-gray-400 text-xs"></div>
                                        </div>
                                    </template>
                                    <div class="w-full py-1">
                                        <button @click="form.role = 'User'" type="button" class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors" :class="{'text-primary font-medium': form.role === 'User'}">
                                            普通用户
                                        </button>
                                        <button @click="form.role = 'Admin'" type="button" class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors" :class="{'text-primary font-medium': form.role === 'Admin'}">
                                            管理员
                                        </button>
                                    </div>
                                </Dropdown>
                            </div>

                            <!-- 状态 (仅编辑模式) -->
                            <div v-if="isEdit" class="space-y-1.5">
                                <label class="text-sm font-medium text-gray-700">状态</label>
                                <div class="flex gap-4">
                                    <label class="flex items-center gap-2 cursor-pointer">
                                        <input type="radio" v-model="form.status" :value="1" class="text-primary focus:ring-primary" />
                                        <span class="text-sm text-gray-600">正常</span>
                                    </label>
                                    <label class="flex items-center gap-2 cursor-pointer">
                                        <input type="radio" v-model="form.status" :value="0" class="text-red-500 focus:ring-red-500" />
                                        <span class="text-sm text-gray-600">禁用</span>
                                    </label>
                                </div>
                            </div>
                            
                            <!-- 底部按钮 (表单内部) -->
                            <div class="pt-4 border-t border-gray-100 flex justify-end gap-3 mt-auto">
                                <button @click="close" type="button"
                                    class="px-4 py-2 text-sm font-medium text-gray-600 bg-white border border-gray-200 rounded-lg hover:bg-gray-50 hover:border-gray-300 transition-colors">
                                    取消
                                </button>
                                <button type="submit" :disabled="loading"
                                    class="px-4 py-2 text-sm font-bold text-white bg-primary rounded-lg shadow-lg shadow-primary/20 hover:bg-primary/90 transition-all transform active:scale-95 disabled:opacity-70 disabled:cursor-not-allowed flex items-center gap-2">
                                    <div v-if="loading" class="i-fa6-solid-spinner animate-spin"></div>
                                    <span>{{ loading ? '保存中...' : '保存' }}</span>
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { reactive, watch, computed } from 'vue'
import Dropdown from '@/components/ui/Dropdown.vue'

interface User {
    id?: number
    username: string
    role: string
    status: number
}

const props = defineProps<{
    modelValue: boolean
    user: User | null
    loading: boolean
}>()

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'submit', data: any): void
}>()

const isEdit = computed(() => !!props.user)

const form = reactive({
    username: '',
    password: '',
    role: 'User',
    status: 1
})

watch(() => props.modelValue, (val) => {
    if (val) {
        if (props.user) {
            form.username = props.user.username
            form.role = props.user.role
            form.status = props.user.status
            form.password = ''
        } else {
            form.username = ''
            form.password = ''
            form.role = 'User'
            form.status = 1
        }
    }
})

const close = () => {
    emit('update:modelValue', false)
}

const submit = () => {
    emit('submit', { ...form })
}
</script>
