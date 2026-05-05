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
            <div v-if="modelValue" class="fixed inset-0 z-[100] flex items-center justify-center">
                <!-- 背景遮罩 -->
                <div class="absolute inset-0 bg-black/20" @click="close"></div>

                <!-- 弹窗主体 -->
                <div class="relative bg-white rounded-2xl w-full max-w-sm overflow-hidden shadow-xl mx-8 flex flex-col max-h-[60vh] transition-all duration-300"
                    :class="modelValue ? 'opacity-100 scale-100' : 'opacity-0 scale-95'">

                    <!-- 标题 -->
                    <div class="p-5 pb-2 flex justify-between items-center bg-white z-10">
                        <h3 class="font-bold text-lg text-gray-800">收藏到歌单</h3>
                        <button @click="close" class="w-8 h-8 rounded-full hover:bg-gray-100 flex items-center justify-center text-gray-400 hover:text-gray-600 transition-colors">
                            <div class="i-fa6-solid-xmark text-sm"></div>
                        </button>
                    </div>

                    <!-- 歌单列表 -->
                    <div class="flex-1 overflow-y-auto px-3 py-2 custom-scrollbar">
                        <!-- 新建歌单入口 -->
                        <button @click="handleCreateNew"
                            class="w-full flex items-center gap-3 p-3 rounded-xl hover:bg-blue-50 text-blue-600 transition-colors mb-2 group">
                            <div class="w-12 h-12 rounded-lg bg-blue-100 flex items-center justify-center group-hover:bg-blue-200 transition-colors">
                                <div class="i-fa6-solid-plus text-lg"></div>
                            </div>
                            <span class="font-bold">新建歌单</span>
                        </button>

                        <div v-if="loading" class="py-10 text-center text-gray-400 text-sm">
                            <div class="i-fa6-solid-spinner animate-spin text-xl mb-2 mx-auto"></div>
                            加载中...
                        </div>

                        <div v-else-if="playlists.length === 0" class="py-8 text-center text-gray-400 text-sm">
                            暂无创建的歌单
                        </div>

                        <div v-else class="space-y-1">
                            <button v-for="list in playlists" :key="list.id" @click="selectPlaylist(list.id)"
                                class="w-full flex items-center gap-3 p-2 rounded-xl hover:bg-gray-100 transition-colors text-left group">
                                <!-- 歌单封面 -->
                                <div class="w-12 h-12 rounded-lg bg-gray-100 overflow-hidden shrink-0 border border-gray-100">
                                    <SongCover 
                                        :src="resolveCover(list.coverUrl)" 
                                        class="w-full h-full"
                                    />
                                </div>
                                <!-- 信息 -->
                                <div class="flex-1 min-w-0">
                                    <div class="font-bold text-gray-800 truncate group-hover:text-primary transition-colors">{{ list.title }}</div>
                                    <div class="text-xs text-gray-400 mt-0.5">{{ list.songCount }} 首歌曲</div>
                                </div>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { fetchMyPlaylists, addSongToPlaylist } from '@/api/playlists'
import { resolveCover, showToast } from '@/utils/common'
import SongCover from '@/components/common/SongCover.vue'

const props = defineProps<{
    modelValue: boolean,
    songId: number | string | null // 当前要添加的歌曲ID
}>()
const emit = defineEmits(['update:modelValue', 'create-new'])

const playlists = ref<any[]>([])
const loading = ref(false)
const router = useRouter()

const ensureAuthenticated = () => {
    const token = localStorage.getItem('token')
    if (token) return true

    showToast('登录后才能收藏到歌单', 'info')
    close()
    router.push('/auth')
    return false
}

// 监听打开，加载歌单列表
watch(() => props.modelValue, (val) => {
    if (val) {
        if (!ensureAuthenticated()) return
        loadPlaylists()
    }
})

const loadPlaylists = async () => {
    loading.value = true
    try {
        const res = await fetchMyPlaylists() as any;
        if (res.code === 200) {
            playlists.value = res.data || []
        }
    } catch (error) {
        showToast('获取歌单失败', 'error')
    } finally {
        loading.value = false
    }
}

const selectPlaylist = async (playlistId: number) => {
    if (!props.songId) return
    if (!ensureAuthenticated()) return

    try {
        await addSongToPlaylist(playlistId, props.songId)
        showToast('已添加到歌单')
        close()
    } catch (error: any) {
        close()
    }
}

const handleCreateNew = () => {
    if (!ensureAuthenticated()) return
    close()
    emit('create-new') // 触发父组件打开创建弹窗
}

const close = () => emit('update:modelValue', false)
</script>
