<template>
    <Teleport to="body">
        <Transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0"
            enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in"
            leave-from-class="opacity-100" leave-to-class="opacity-0">
            <div v-if="modelValue" class="fixed inset-0 z-[100] flex items-center justify-center">
                <!-- 背景遮罩 -->
                <div class="absolute inset-0 bg-black/40" @click="close"></div>

                <!-- 模态框主体 -->
                <div class="relative bg-white rounded-[2rem] w-full max-w-md p-8 shadow-2xl border border-white/50 mx-4 transition-all duration-300 ease-[cubic-bezier(0.16,1,0.3,1)] will-change-transform"
                    :class="modelValue ? 'opacity-100 scale-100' : 'opacity-0 scale-95'">

                    <!-- 装饰背景 -->
                    <div class="absolute -top-20 -right-20 w-60 h-60 rounded-full pointer-events-none"
                        style="background: radial-gradient(circle, rgba(51,91,241,0.2) 0%, rgba(255,255,255,0) 70%);">
                    </div>
                    <div class="absolute -bottom-20 -left-20 w-60 h-60 rounded-full pointer-events-none"
                        style="background: radial-gradient(circle, rgba(51,91,241,0.2) 0%, rgba(255,255,255,0) 70%);">
                    </div>

                    <div class="relative z-10">
                        <h2 class="text-2xl font-bold mb-6 text-gray-800 text-center">
                            {{ mode === 'edit' ? '编辑歌单' : '新建歌单' }}
                        </h2>

                        <div class="space-y-5">
                            <!-- 封面上传 -->
                            <div class="flex justify-center">
                                <div class="relative w-36 h-36 bg-gray-50 rounded-2xl flex items-center justify-center cursor-pointer hover:bg-gray-100 transition-all duration-300 overflow-hidden group shadow-inner border border-gray-100"
                                    @click="triggerUpload">
                                    <img v-if="previewUrl" :src="previewUrl"
                                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                        @error="previewUrl = ''" />
                                    <!-- 编辑模式下如果有初始图片且没有预览图，显示初始图片 -->
                                    <img v-else-if="mode === 'edit' && initialData?.coverUrl && !coverFile && !imageLoadError"
                                        :src="resolveCover(initialData.coverUrl)"
                                        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                        @error="imageLoadError = true" />

                                    <div v-else
                                        class="text-gray-400 flex flex-col items-center group-hover:text-primary transition-colors">
                                        <div class="i-fa6-solid-camera text-3xl mb-2"></div>
                                        <span class="text-xs font-medium">上传封面</span>
                                    </div>
                                    <!-- 遮罩 -->
                                    <div
                                        class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-[1px]">
                                        <div
                                            class="text-white text-xs font-medium px-3 py-1 rounded-full border border-white/50">
                                            更换图片</div>
                                    </div>
                                    <input type="file" ref="fileInput" class="hidden" accept="image/*"
                                        @change="handleFileChange" />
                                </div>
                            </div>

                            <!-- 标题 -->
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">歌单标题</label>
                                <input v-model="form.title" type="text" placeholder="给歌单起个好听的名字"
                                    class="w-full px-5 py-3 bg-white/60 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all placeholder:text-gray-400" />
                            </div>

                            <!-- 描述 -->
                            <div>
                                <label class="block text-sm font-bold text-gray-700 mb-1.5 ml-1">描述</label>
                                <textarea v-model="form.description" rows="3" placeholder="介绍一下这个歌单..."
                                    class="w-full px-5 py-3 bg-white/60 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary transition-all resize-none placeholder:text-gray-400"></textarea>
                            </div>

                            <!-- 公开选项 -->
                            <div class="flex items-center justify-between px-1">
                                <div class="flex flex-col">
                                    <span class="text-sm font-bold text-gray-700">公开歌单</span>
                                    <span class="text-xs text-gray-400">审核通过后，所有人可见</span>
                                </div>
                                <button
                                    class="relative w-12 h-7 rounded-full transition-colors duration-300 focus:outline-none shadow-inner"
                                    :class="form.isPublic ? 'bg-primary' : 'bg-gray-200'"
                                    @click="form.isPublic = !form.isPublic">
                                    <span
                                        class="absolute top-1 left-1 bg-white w-5 h-5 rounded-full shadow-sm transition-transform duration-300"
                                        :class="form.isPublic ? 'translate-x-5' : 'translate-x-0'"></span>
                                </button>
                            </div>
                        </div>

                        <div class="mt-8 flex justify-end gap-3">
                            <button @click="close"
                                class="px-5 py-2.5 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-xl transition-colors font-medium">取消</button>
                            <button @click="submit" :disabled="!form.title || loading"
                                class="px-8 py-2.5 bg-gradient-to-r from-primary to-primary/80 text-white rounded-xl shadow-lg shadow-primary/30 hover:shadow-primary/40 hover:-translate-y-0.5 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed disabled:shadow-none disabled:translate-y-0 flex items-center gap-2 transition-all font-bold">
                                <div v-if="loading" class="i-fa6-solid-spinner animate-spin"></div>
                                <span>{{ mode === 'edit' ? '保存' : '创建' }}</span>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { createPlaylist, updatePlaylist } from '@/api/playlists'
import { showToast, resolveCover } from '@/utils/common'

interface InitialData {
    id?: number
    title: string
    description?: string
    isPublic: boolean
    coverUrl?: string
}

const props = withDefaults(defineProps<{
    modelValue: boolean
    mode?: 'create' | 'edit'
    initialData?: InitialData
}>(), {
    mode: 'create'
})

const emit = defineEmits(['update:modelValue', 'success'])

const loading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const previewUrl = ref('')
const coverFile = ref<File | null>(null)
const imageLoadError = ref(false)

const form = reactive({
    title: '',
    description: '',
    isPublic: true
})

// 监听打开状态，填充表单
watch(() => props.modelValue, (val) => {
    if (val) {
        imageLoadError.value = false
        if (props.mode === 'edit' && props.initialData) {
            form.title = props.initialData.title
            form.description = props.initialData.description || ''
            form.isPublic = props.initialData.isPublic
            previewUrl.value = ''
            coverFile.value = null
        } else {
            form.title = ''
            form.description = ''
            form.isPublic = true
            previewUrl.value = ''
            coverFile.value = null
        }
    }
})

const close = () => {
    emit('update:modelValue', false)
}

const triggerUpload = () => fileInput.value?.click()

const handleFileChange = (e: Event) => {
    const files = (e.target as HTMLInputElement).files
    if (files && files[0]) {
        coverFile.value = files[0]
        previewUrl.value = URL.createObjectURL(files[0])
    }
}

const submit = async () => {
    if (!form.title) return
    loading.value = true
    try {
        const becomesPublicOnCreate = props.mode === 'create' && form.isPublic
        const becomesPublicOnEdit =
            props.mode === 'edit' &&
            props.initialData &&
            !props.initialData.isPublic &&
            form.isPublic

        if (props.mode === 'create') {
            await createPlaylist({
                title: form.title,
                description: form.description,
                isPublic: form.isPublic,
                coverFile: coverFile.value || undefined
            })
            showToast(becomesPublicOnCreate ? '歌单已创建，等待审核' : '创建成功')
        } else {
            if (!props.initialData?.id) return
            await updatePlaylist(props.initialData.id, {
                title: form.title,
                description: form.description,
                isPublic: form.isPublic,
                coverFile: coverFile.value || undefined
            })
            showToast(becomesPublicOnEdit ? '歌单已提交审核' : '更新成功')
        }

        emit('success')
        close()
    } catch (error) {
        showToast(props.mode === 'create' ? '创建失败' : '更新失败', 'error')
    } finally {
        loading.value = false
    }
}
</script>
