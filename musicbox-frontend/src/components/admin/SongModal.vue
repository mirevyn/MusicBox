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
                <div class="absolute inset-0 bg-black/60 transition-opacity" @click="onCancel"></div>

                <!-- 模态框内容 -->
                <div class="relative bg-white rounded-3xl w-full max-w-2xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
                    
                    <!-- 头部 -->
                    <div class="px-8 py-6 border-b border-gray-100 flex items-center justify-between shrink-0 bg-white z-10">
                        <div>
                            <h3 class="text-xl font-bold text-gray-900">{{ isEdit ? '编辑歌曲信息' : '上传新歌曲' }}</h3>
                            <p class="text-xs text-gray-500 mt-1">完善歌曲详情与媒体文件</p>
                        </div>
                        <button @click="onCancel" class="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center text-gray-400 hover:bg-gray-100 hover:text-gray-900 transition-colors">
                            <div class="i-fa6-solid-xmark text-lg"></div>
                        </button>
                    </div>

                    <!-- 可滚动主体 -->
                    <div class="p-8 overflow-y-auto custom-scrollbar flex-1">
                        <form @submit.prevent="onSubmit" class="space-y-8">
                            
                            <!-- 顶部区域：封面与元数据 -->
                            <div class="flex flex-col md:flex-row gap-8">
                                <!-- 左侧：封面图片 -->
                                <div class="shrink-0 flex flex-col gap-2 group">
                                    <div class="relative w-40 h-40 bg-gray-50 rounded-2xl border border-gray-100 overflow-hidden shadow-sm group-hover:shadow-md transition-all">
                                        <input type="file" accept="image/*" @change="handleCoverChange" 
                                            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-20">
                                        
                                        <img v-if="previewCover" :src="previewCover" class="w-full h-full object-cover" @error="imageLoadError = true">
                                        <img v-else-if="song?.coverUrl && !imageLoadError" :src="resolveCover(song.coverUrl)" class="w-full h-full object-cover" @error="imageLoadError = true">
                                        
                                        <div v-else class="w-full h-full flex flex-col items-center justify-center text-gray-400 group-hover:text-primary transition-colors">
                                            <div class="i-fa6-solid-image text-3xl mb-2 opacity-50"></div>
                                            <span class="text-xs font-medium">点击上传</span>
                                        </div>

                                        <div class="absolute inset-0 bg-black/50 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity z-10 pointer-events-none">
                                            <div class="i-fa6-solid-pen text-white text-xl"></div>
                                        </div>
                                    </div>
                                    <div class="text-center">
                                        <div class="text-xs text-gray-500">建议 1:1 尺寸</div>
                                    </div>
                                </div>

                                <!-- 右侧：文本字段 -->
                                <div class="flex-1 space-y-5">
                                    <div class="space-y-1.5">
                                        <label class="text-sm font-semibold text-gray-700">歌曲标题 <span class="text-red-500">*</span></label>
                                        <input v-model="form.title" type="text" required
                                            class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-gray-400"
                                            placeholder="输入歌曲名称">
                                    </div>
                                    
                                    <div class="grid grid-cols-2 gap-4">
                                        <div class="space-y-1.5">
                                            <label class="text-sm font-semibold text-gray-700">歌手 <span class="text-red-500">*</span></label>
                                            <input v-model="form.artist" type="text" required
                                                class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-gray-400"
                                                placeholder="输入歌手名">
                                        </div>
                                        <div class="space-y-1.5">
                                            <label class="text-sm font-semibold text-gray-700">专辑</label>
                                            <input v-model="form.album" type="text"
                                                class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 focus:border-primary outline-none transition-all placeholder:text-gray-400"
                                                placeholder="输入专辑名">
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <!-- 分割线 -->
                            <div class="h-px bg-gray-100"></div>

                            <!-- 底部区域：文件 -->
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <!-- 音频文件 -->
                                <div class="space-y-2">
                                    <label class="text-sm font-semibold text-gray-700 flex justify-between">
                                        <span>音频文件 <span v-if="!isEdit" class="text-red-500">*</span></span>
                                        <span v-if="audioFile" class="text-xs text-primary truncate max-w-[150px]">{{ audioFile.name }}</span>
                                    </label>
                                    <div class="relative">
                                        <input type="file" ref="audioInput" accept="audio/*" @change="handleAudioChange"
                                            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10" :required="!isEdit">
                                        <div class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-3 flex items-center gap-3 hover:bg-gray-100 hover:border-gray-300 transition-all">
                                            <div class="w-8 h-8 rounded-lg bg-blue-50 text-blue-500 flex items-center justify-center shrink-0">
                                                <div class="i-fa6-solid-music"></div>
                                            </div>
                                            <span class="text-sm text-gray-500 flex-1 truncate">
                                                {{ audioFile ? '已选择文件' : (isEdit ? '点击更换音频' : '选择音频文件') }}
                                            </span>
                                        </div>
                                    </div>
                                    <audio v-if="previewAudio" controls :src="previewAudio" class="w-full h-8 mt-1 rounded"></audio>
                                </div>

                                <!-- 歌词文件 -->
                                <div class="space-y-2">
                                    <label class="text-sm font-semibold text-gray-700 flex justify-between">
                                        <span>歌词文件</span>
                                        <span v-if="lyricFile" class="text-xs text-teal-600 truncate max-w-[150px]">{{ lyricFile.name }}</span>
                                    </label>
                                    <div class="relative">
                                        <input type="file" accept=".lrc" @change="handleLyricChange"
                                            class="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10">
                                        <div class="w-full bg-gray-50 border border-gray-200 rounded-xl px-4 py-3 flex items-center gap-3 hover:bg-gray-100 hover:border-gray-300 transition-all">
                                            <div class="w-8 h-8 rounded-lg bg-teal-50 text-teal-500 flex items-center justify-center shrink-0">
                                                <div class="i-fa6-solid-file-lines"></div>
                                            </div>
                                            <span class="text-sm text-gray-500 flex-1 truncate">
                                                {{ lyricFile ? '已选择歌词' : (song?.lyricUrl ? '已上传 (点击更换)' : '选择歌词 (.lrc)') }}
                                            </span>
                                        </div>
                                    </div>
                                </div>
                            </div>

                        </form>
                    </div>

                    <!-- 底部按钮 -->
                    <div class="px-8 py-5 bg-gray-50 border-t border-gray-100 flex items-center gap-4 shrink-0">
                         <!-- 上传进度 -->
                        <div v-if="uploadProgress > 0 && uploadProgress < 100" class="flex-1 mr-4">
                            <div class="flex justify-between text-xs font-medium text-gray-500 mb-1">
                                <span>{{ isEdit ? '更新中' : '上传中' }}</span>
                                <span>{{ uploadProgress }}%</span>
                            </div>
                            <div class="h-1.5 w-full bg-gray-200 rounded-full overflow-hidden">
                                <div class="h-full bg-primary transition-all duration-300" :style="{ width: `${uploadProgress}%` }"></div>
                            </div>
                        </div>
                        <div v-else class="flex-1"></div>

                        <button @click="onCancel" :disabled="loading"
                            class="px-6 py-2.5 rounded-xl text-sm font-semibold text-gray-600 hover:bg-gray-200 hover:text-gray-900 transition-colors disabled:opacity-50">
                            取消
                        </button>
                        <button @click="onSubmit" :disabled="loading"
                            class="px-6 py-2.5 rounded-xl text-sm font-semibold bg-primary text-white shadow-lg shadow-primary/25 hover:bg-primary/90 hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-2 disabled:opacity-70 disabled:cursor-not-allowed">
                            <div v-if="loading" class="i-fa6-solid-spinner animate-spin"></div>
                            <span>{{ isEdit ? '保存更改' : '立即上传' }}</span>
                        </button>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import type { Song as PlayerSong } from '@/stores/player';
import { resolveUrl, resolveCover } from '@/utils/common';

const props = withDefaults(defineProps<{
    modelValue: boolean
    song?: PlayerSong | null
    loading?: boolean
    uploadProgress?: number
}>(), {
    song: null,
    loading: false,
    uploadProgress: 0
})

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
    (e: 'submit', formData: FormData): void
}>()

const isEdit = computed(() => !!props.song)

const form = ref({
    title: '',
    artist: '',
    album: ''
})

const audioFile = ref<File | null>(null)
const coverFile = ref<File | null>(null)
const lyricFile = ref<File | null>(null)
const previewCover = ref<string>('')
const previewAudio = ref<string>('')
const audioInput = ref<HTMLInputElement | null>(null)
const imageLoadError = ref(false)

// 当模态框打开或歌曲改变时初始化表单
watch(() => props.song, (newSong) => {
    if (newSong) {
        form.value = {
            title: newSong.title || '',
            artist: newSong.artist || '',
            album: newSong.album || ''
        }


        imageLoadError.value = false
        previewCover.value = '' // 重置预览，依赖于下面的逻辑
        // 如果有封面URL，我们将通过模板中的 resolveCover 使用它，除非选择了新文件
        previewAudio.value = resolveUrl(newSong.fileUrl)
    } else {
        resetForm()
    }
}, { immediate: true })

watch(() => props.modelValue, (isOpen) => {
    if (!isOpen) {
        // 可选：清除验证错误
    } else if (!props.song) {
        resetForm()
    }
})

function resetForm() {
    form.value = { title: '', artist: '', album: '' }
    audioFile.value = null
    coverFile.value = null
    lyricFile.value = null
    previewCover.value = ''
    imageLoadError.value = false
    previewAudio.value = ''
    if (audioInput.value) audioInput.value.value = ''
}

const handleAudioChange = (event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) {
        audioFile.value = input.files[0]
        
        // 如果为空，则自动填充元数据
        const name = audioFile.value.name.replace(/\.[^/.]+$/, "")
        if (name.includes('-')) {
            const parts = name.split('-').map(s => s.trim())
            if (!form.value.artist) form.value.artist = parts[0] || ''
            if (!form.value.title) form.value.title = parts[1] || parts[0] || ''
        } else {
            if (!form.value.title) form.value.title = name
        }
        previewAudio.value = URL.createObjectURL(audioFile.value)
    }
}

const handleCoverChange = (event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) {
        const file = input.files[0]
        coverFile.value = file
        previewCover.value = URL.createObjectURL(file)
    }
}

const handleLyricChange = (event: Event) => {
    const input = event.target as HTMLInputElement
    if (input.files && input.files[0]) {
        lyricFile.value = input.files[0]
    }
}

const onCancel = () => {
    if (!props.loading) {
        emit('update:modelValue', false)
    }
}

const onSubmit = () => {
    const formData = new FormData()
    formData.append('title', form.value.title)
    formData.append('artist', form.value.artist)
    formData.append('album', form.value.album)
    
    if (audioFile.value) {
        formData.append('file', audioFile.value)
    }
    
    if (coverFile.value) {
        formData.append('cover', coverFile.value)
    }
    
    if (lyricFile.value) {
        formData.append('lyric', lyricFile.value)
    }

    emit('submit', formData)
}
</script>
