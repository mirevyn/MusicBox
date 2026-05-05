<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-1">歌曲管理</h2>
                <p class="text-gray-500 text-sm">管理系统中的所有音乐资源</p>
            </div>
        </div>

        <!-- 歌曲表格卡片 -->
        <div class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden min-h-[620px] flex flex-col">
            <!-- 工具栏 -->
            <div class="p-4 border-b border-gray-100 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <!-- 搜索 -->
                <div class="flex gap-2 flex-1 max-w-lg">
                    <div class="relative flex-1">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                            <div class="i-fa6-solid-magnifying-glass"></div>
                        </div>
                        <input v-model="searchKeyword" @keydown.enter="executeSearch" type="text"
                            placeholder="搜索歌曲名、歌手、专辑..."
                            class="w-full bg-gray-50 border-none rounded-lg py-2 pl-10 pr-4 text-sm focus:ring-2 focus:ring-primary/20 focus:bg-white transition-all">
                    </div>
                    <button @click="executeSearch"
                        class="px-4 py-2 bg-gray-50 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors flex items-center gap-2 whitespace-nowrap">
                        搜索
                    </button>
                </div>

                <!-- 操作按钮 -->
                <div class="flex gap-2">
                    <button @click="handleExport" :disabled="exportLoading"
                        class="px-4 py-2 bg-gray-50 border border-gray-100 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors flex items-center gap-2 whitespace-nowrap disabled:opacity-50">
                        <div v-if="exportLoading" class="i-fa6-solid-spinner animate-spin text-primary"></div>
                        <div v-else class="i-fa6-solid-download opacity-70"></div>
                        批量导出
                    </button>
                    <button @click="handleCreate"
                        class="bg-primary text-white hover:bg-primary/90 transition-colors rounded-lg px-4 py-2 text-sm font-medium shadow-lg shadow-primary/20 flex items-center gap-2 whitespace-nowrap">
                        <div class="i-fa6-solid-plus"></div> 上传歌曲
                    </button>
                </div>
            </div>

            <!-- 内容区域 -->
            <div class="relative flex-1 flex flex-col">
                <!-- 加载状态 -->
                <div v-if="loading && songs.length === 0"
                    class="absolute inset-0 flex items-center justify-center bg-white z-10">
                    <div class="flex flex-col items-center gap-3">
                        <div class="i-fa6-solid-spinner animate-spin text-2xl text-primary"></div>
                        <div class="text-sm text-gray-400">加载中...</div>
                    </div>
                </div>

                <!-- 空状态 -->
                <div v-if="!loading && songs.length === 0" class="absolute inset-0 flex items-center justify-center">
                    <div class="flex flex-col items-center gap-3 text-gray-400">
                        <div class="i-fa6-solid-music text-4xl opacity-20"></div>
                        <div>暂无歌曲数据</div>
                    </div>
                </div>

                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr
                            class="bg-gray-50/50 border-b border-gray-100 text-xs uppercase text-gray-500 font-semibold tracking-wider">
                            <th class="px-6 py-4 w-16">#</th>
                            <th class="px-6 py-4 w-24">封面</th>
                            <th class="px-6 py-4 w-1/4">标题</th>
                            <th class="px-6 py-4 w-1/6">歌手</th>
                            <th class="px-6 py-4 w-1/6">专辑</th>
                            <th class="px-6 py-4 w-24">时长</th>
                            <th class="px-6 py-4 w-24 text-right">操作</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-50">
                        <tr v-for="(song, index) in songs" :key="song.songId || song.id"
                            class="hover:bg-gray-50/80 transition-colors group">
                            <td class="px-6 py-4 text-gray-400 font-medium text-sm">{{ (currentPage - 1) * pageSize +
                                index
                                + 1 }}</td>
                            <td class="px-6 py-4">
                                <div class="w-10 h-10 rounded-lg overflow-hidden shadow-sm border border-gray-100">
                                    <SongCover :src="song.coverUrl" :alt="song.title" />
                                </div>
                            </td>
                            <td class="px-6 py-4">
                                <div class="text-sm font-bold text-gray-800 line-clamp-1">{{ song.title }}</div>
                                <div class="text-xs text-gray-400 md:hidden">{{ song.artist }}</div>
                            </td>
                            <td class="px-6 py-4 text-sm text-gray-600">{{ song.artist }}</td>
                            <td class="px-6 py-4 text-sm text-gray-500">{{ song.album || '-' }}</td>
                            <td class="px-6 py-4 text-sm text-gray-500 font-mono">{{ formatDuration(song.duration) }}
                            </td>
                            <td class="px-6 py-4 text-right">
                                <div class="flex items-center justify-end gap-2">
                                    <button @click="handleEdit(song)"
                                        class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-500 hover:bg-blue-50 hover:text-blue-600 transition-colors"
                                        title="编辑">
                                        <div class="i-fa6-solid-pen"></div>
                                    </button>
                                    <button @click="handleDelete(song)"
                                        class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-500 hover:bg-red-50 hover:text-red-500 transition-colors"
                                        title="删除">
                                        <div class="i-fa6-solid-trash"></div>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>

                <!-- 分页 -->
                <div class="mt-auto">
                     <Pagination v-if="total > 0" :current="currentPage" :total="total" :page-size="pageSize"
                    @update:current="handlePageChange" />
                </div>
               
            </div>
        </div>

        <!-- 弹窗 -->
        <SongModal v-model="showModal" :song="currentSong" :loading="modalLoading" :upload-progress="uploadProgress"
            @submit="handleSubmit" />

        <ConfirmModal v-model="showDeleteModal" type="danger" title="确认删除"
            :content="`确定要删除歌曲 “${deleteTarget?.title}” 吗？此操作无法撤销。`" :loading="confirmLoading"
            @confirm="confirmDelete" />

        <ConfirmModal v-model="showExportConfirmModal" type="info" title="确认导出"
            content="确定要导出数据库中所有的歌曲吗？系统将打包音频、歌词和封面文件，这可能需要一点时间。" :loading="exportLoading"
            @confirm="confirmExport" />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { fetchSongs, deleteSong, uploadSong, updateSong, exportSongs } from '@/api/songs'
import { formatSongs, showToast } from '@/utils/common'
import { getApiBlobErrorMessage, getApiErrorMessage } from '@/utils/request'
import type { Song } from '@/stores/player'
import Pagination from '@/components/ui/Pagination.vue'
import SongModal from './SongModal.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'
import SongCover from '@/components/common/SongCover.vue'


// 状态变量
const loading = ref(false)
const exportLoading = ref(false)
const songs = ref<Song[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(6)
const searchKeyword = ref('')

// 弹窗状态变量
const showModal = ref(false)
const currentSong = ref<Song | null>(null)
const modalLoading = ref(false)
const uploadProgress = ref(0)

const showDeleteModal = ref(false)
const deleteTarget = ref<Song | null>(null)
const confirmLoading = ref(false)

const showExportConfirmModal = ref(false)

// 数据获取函数
const fetchData = async () => {
    loading.value = true
    try {
        const res = await fetchSongs({
            pageIndex: currentPage.value,
            pageSize: pageSize.value,
            keyword: searchKeyword.value
        })
        if (res.code === 200 && res.data) {
            songs.value = formatSongs(res.data.list)
            total.value = res.data.total || 0
        }
    } catch (error) {
        console.error('Failed to fetch songs:', error)
        showToast(getApiErrorMessage(error, '获取歌曲失败'), 'error')
    } finally {
        loading.value = false
    }
}

const handlePageChange = (page: number) => {
    currentPage.value = page
    fetchData()
}



const executeSearch = () => {
    currentPage.value = 1
    fetchData()
}

// 创建与编辑处理
const handleCreate = () => {
    currentSong.value = null
    uploadProgress.value = 0
    showModal.value = true
}

const handleEdit = (song: Song) => {
    currentSong.value = song
    uploadProgress.value = 0
    showModal.value = true
}

const handleSubmit = async (formData: FormData) => {
    modalLoading.value = true
    uploadProgress.value = 0

    try {
        let res;
        const width = (percent: number) => {
            uploadProgress.value = percent
        }

        if (currentSong.value) {
            // 更新歌曲
            const id = currentSong.value.songId
            res = await updateSong(id, formData, width)
        } else {
            // 创建新歌曲
            res = await uploadSong(formData, width)
        }

        const successCode = currentSong.value ? 200 : 201
        if (res.code === successCode) {
            showModal.value = false
            fetchData() // 刷新列表数据
            showToast(currentSong.value ? '保存成功' : '上传成功')
        } else {
            showToast(res.msg || '操作失败', 'error')
        }
    } catch (error) {
        console.error('Submit failed:', error)
        showToast(getApiErrorMessage(error, '操作失败，请重试'), 'error')
    } finally {
        modalLoading.value = false
    }
}

// 删除操作处理函数
const handleDelete = (song: Song) => {
    deleteTarget.value = song
    showDeleteModal.value = true
}

const confirmDelete = async () => {
    if (!deleteTarget.value) return

    confirmLoading.value = true
    try {
        const id = deleteTarget.value.songId
        const res = await deleteSong(id)
        if (res.code === 200) {
            showDeleteModal.value = false
            // 检查是否需要返回上一页（如果当前页被删光了）
            if (songs.value.length === 1 && currentPage.value > 1) {
                currentPage.value--
            }
            fetchData()
            showToast('删除成功')
        } else {
            showToast(res.msg || '删除失败', 'error')
        }
    } catch (error) {
        console.error('Delete failed:', error)
        showToast(getApiErrorMessage(error, '删除失败，请重试'), 'error')
    } finally {
        confirmLoading.value = false
    }
}

// 导出全部功能逻辑
const handleExport = () => {
    showExportConfirmModal.value = true
}

const confirmExport = async () => {
    exportLoading.value = true
    showToast('正在准备导出文件，请稍候...', 'info')
    
    try {
        const blob = await exportSongs()
        if (!(blob instanceof Blob)) {
            throw new Error('导出返回数据格式不正确')
        }
        
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `MusicBox_Export_${new Date().getTime()}.zip`)
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)
        
        showExportConfirmModal.value = false
        showToast('导出成功，请查看浏览器下载栏', 'success')
    } catch (error) {
        console.error('Export failed:', error)
        showToast(await getApiBlobErrorMessage(error, '导出失败，请重试'), 'error')
    } finally {
        exportLoading.value = false
    }
}

// 工具辅助函数
const formatDuration = (seconds?: number) => {
    if (!seconds) return '--:--'
    const min = Math.floor(seconds / 60)
    const sec = Math.floor(seconds % 60)
    return `${min}:${sec.toString().padStart(2, '0')}`
}

onMounted(() => {
    fetchData()
})
</script>
