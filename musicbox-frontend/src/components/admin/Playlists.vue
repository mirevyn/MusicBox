<template>
    <div class="space-y-6">
        <div class="flex items-center justify-between">
            <div>
                <h2 class="text-2xl font-bold text-gray-800 mb-1">歌单审核</h2>
                <p class="text-sm text-gray-500">检查标题、简介、歌曲内容和可见性设置是否存在违规风险</p>
            </div>
        </div>

        <div
            class="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden min-h-[620px] h-[620px] flex flex-col">
            <!-- 工具栏 -->
            <div class="p-4 border-b border-gray-100 flex flex-col lg:flex-row lg:items-center justify-between gap-4 shrink-0">
                <!-- 左侧：过滤标签 -->
                <div class="flex flex-wrap gap-2">
                    <button v-for="item in filterTabs" :key="item.value" @click="currentFilter = item.value"
                        :class="currentFilter === item.value ? item.activeClass : 'bg-gray-50 text-gray-500 hover:bg-gray-100'"
                        class="px-3.5 py-2 rounded-full text-xs font-semibold transition-all">
                        {{ item.label }}
                    </button>
                </div>
                
                <!-- 右侧：搜索 -->
                <div class="flex gap-2 w-full lg:w-auto lg:max-w-sm flex-1">
                    <div class="relative flex-1">
                        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400">
                            <div class="i-fa6-solid-magnifying-glass"></div>
                        </div>
                        <input v-model="searchQuery" @keydown.enter="handleSearch" type="text"
                            placeholder="搜索歌单标题"
                            class="w-full bg-gray-50 border-none rounded-lg py-2 pl-10 pr-4 text-sm focus:ring-2 focus:ring-primary/20 focus:bg-white transition-all">
                    </div>
                    <button @click="handleSearch"
                        class="px-4 py-2 bg-gray-50 rounded-lg text-sm font-medium text-gray-600 hover:bg-gray-100 transition-colors flex items-center gap-2 whitespace-nowrap">
                        搜索
                    </button>
                </div>
            </div>

            <div class="grid grid-cols-1 xl:grid-cols-[minmax(0,1.15fr)_minmax(360px,0.85fr)] flex-1 min-h-0">
                <section class="border-r border-gray-100 min-w-0 min-h-0 flex flex-col relative">
                    <div
                        class="grid grid-cols-[72px_minmax(0,1.35fr)_88px_72px_88px_96px] gap-3 border-b border-gray-100 bg-gray-50/80 px-4 py-3 text-[11px] font-bold uppercase tracking-widest text-gray-400 shrink-0">
                        <div>封面</div>
                        <div>歌单信息</div>
                        <div>创建者</div>
                        <div>歌曲数</div>
                        <div>提交状态</div>
                        <div>操作</div>
                    </div>

                    <div v-if="loading && playlists.length === 0" class="divide-y divide-gray-100 flex-1 overflow-hidden">
                        <div v-for="n in pageSize" :key="n"
                            class="grid grid-cols-[72px_minmax(0,1.35fr)_88px_72px_88px_96px] gap-3 px-4 py-3 animate-pulse">
                            <div class="h-12 rounded-xl bg-gray-200"></div>
                            <div class="space-y-2 py-1">
                                <div class="h-4 w-2/3 rounded bg-gray-200"></div>
                                <div class="h-3 w-full rounded bg-gray-100"></div>
                                <div class="h-3 w-1/2 rounded bg-gray-100"></div>
                            </div>
                            <div class="h-4 w-16 self-center rounded bg-gray-100"></div>
                            <div class="h-4 w-12 self-center rounded bg-gray-100"></div>
                            <div class="h-6 w-18 self-center rounded-full bg-gray-100"></div>
                            <div class="h-8 w-full self-center rounded-xl bg-gray-100"></div>
                        </div>
                    </div>

                    <div v-else-if="playlists.length === 0"
                        class="flex flex-1 items-center justify-center px-6 text-center">
                        <div>
                            <div
                                class="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-gray-50 text-gray-300">
                                <div class="i-fa6-solid-folder-open text-xl"></div>
                            </div>
                            <div class="text-sm font-semibold text-gray-600">当前筛选下没有待处理歌单</div>
                            <p class="mt-2 text-xs text-gray-400">切换筛选条件，或等待新的歌单提交审核</p>
                        </div>
                    </div>

                    <div v-else class="divide-y divide-gray-100 flex-1 overflow-hidden">
                        <button v-for="playlist in playlists" :key="playlist.id" @click="openDetailsPanel(playlist)"
                            class="grid h-[72px] w-full grid-cols-[72px_minmax(0,1.35fr)_88px_72px_88px_96px] gap-3 px-4 py-3 text-left transition-colors hover:bg-[rgba(51,91,241,0.03)]"
                            :class="selectedPlaylistId === playlist.id ? 'bg-[rgba(51,91,241,0.05)]' : 'bg-white'">
                            <div class="h-12 overflow-hidden rounded-xl bg-gray-100">
                                <img v-if="playlist.coverUrl" :src="resolveUrl(playlist.coverUrl)"
                                    class="h-full w-full object-cover" @error="playlist.coverUrl = ''">
                                <div v-else
                                    class="flex h-full w-full items-center justify-center bg-gradient-to-br from-slate-100 to-slate-200 text-slate-300">
                                    <div class="i-fa6-solid-music"></div>
                                </div>
                            </div>

                            <div class="min-w-0 flex flex-col justify-center">
                                <div class="flex items-center gap-2">
                                    <div class="truncate text-sm font-bold text-gray-800" :title="playlist.title">{{
                                        playlist.title }}</div>
                                    <div v-if="playlist.status === 2 && playlist.rejectReason"
                                        class="shrink-0 max-w-[120px] truncate rounded border border-red-100 bg-red-50/50 px-1.5 py-0.5 text-[10px] font-medium text-red-500" :title="playlist.rejectReason">
                                        驳回: {{ playlist.rejectReason }}
                                    </div>
                                </div>
                                <div class="mt-0.5 line-clamp-1 text-xs leading-5 text-gray-500">
                                    {{ playlist.description || '暂无简介，建议重点检查标题与歌单歌曲组合' }}
                                </div>
                            </div>

                            <div class="self-center min-w-0">
                                <div class="truncate text-xs font-medium text-gray-700">{{ playlist.user?.username ||
                                    'Unknown' }}</div>
                                <div class="text-[11px] text-gray-400">创建者</div>
                            </div>

                            <div class="self-center">
                                <div class="text-sm font-semibold text-gray-800">{{ playlist.songCount || 0 }}</div>
                                <div class="text-[11px] text-gray-400">歌曲</div>
                            </div>

                            <div class="self-center">
                                <span
                                    class="inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-bold text-white"
                                    :class="statusColor(playlist.status)">
                                    {{ statusText(playlist.status) }}
                                </span>
                            </div>

                            <div class="self-center">
                                <span
                                    class="inline-flex items-center justify-center rounded-xl border border-gray-200 px-2.5 py-1.5 text-xs font-medium text-gray-700">
                                    查看审核
                                </span>
                            </div>
                        </button>
                    </div>

                    <Pagination v-if="total > 0" :current="pageIndex" :total="total" :page-size="pageSize"
                        @update:current="changePage" />

                    <Transition name="fade">
                        <div
                            v-if="loading && playlists.length > 0"
                            class="absolute inset-0 z-10 bg-white/40 pointer-events-none"
                        >
                            <div class="absolute top-3 right-3 inline-flex items-center gap-2 rounded-full border border-gray-100 bg-white px-3 py-1.5 text-xs text-gray-500 shadow-sm">
                                <div class="i-fa6-solid-spinner animate-spin text-primary"></div>
                                正在刷新列表
                            </div>
                        </div>
                    </Transition>
                </section>

                <aside class="min-w-0 bg-gray-50/50 min-h-0 flex flex-col relative">
                    <div v-if="currentPlaylistDetails" class="flex flex-1 min-h-0 flex-col">
                        <div class="border-b border-gray-100 bg-white px-5 py-4 shrink-0">
                            <div class="flex items-start gap-4">
                                <div class="h-16 w-16 shrink-0 overflow-hidden rounded-2xl bg-gray-100">
                                    <img v-if="currentPlaylistDetails.coverUrl"
                                        :src="resolveUrl(currentPlaylistDetails.coverUrl)"
                                        class="h-full w-full object-cover"
                                        @error="currentPlaylistDetails.coverUrl = ''">
                                    <div v-else
                                        class="flex h-full w-full items-center justify-center bg-gradient-to-br from-slate-100 to-slate-200 text-slate-300">
                                        <div class="i-fa6-solid-music text-2xl"></div>
                                    </div>
                                </div>
                                <div class="min-w-0 flex-1">
                                    <div class="flex items-center gap-2 mb-2">
                                        <span
                                            class="inline-flex items-center rounded-full px-2.5 py-1 text-[11px] font-bold text-white"
                                            :class="statusColor(currentPlaylistDetails.status)">
                                            {{ statusText(currentPlaylistDetails.status) }}
                                        </span>
                                        <span
                                            class="inline-flex items-center rounded-full border border-primary/10 bg-[rgba(51,91,241,0.05)] px-2.5 py-1 text-[11px] font-semibold text-primary">
                                            {{ visibilityText(currentPlaylistDetails) }}
                                        </span>
                                    </div>
                                    <h3 class="truncate text-base font-bold text-gray-800">{{
                                        currentPlaylistDetails.title }}</h3>
                                    <p class="mt-1.5 text-sm leading-6 text-gray-500 line-clamp-2">
                                        {{ currentPlaylistDetails.description || '该歌单未填写简介，请重点检查标题和曲目内容是否存在违规信息。' }}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div class="flex-1 overflow-y-auto px-5 py-4 space-y-4 min-h-0">
                            <section class="rounded-2xl border border-gray-100 bg-white p-4">
                                <h4 class="text-xs font-bold uppercase tracking-widest text-gray-400">审核摘要</h4>
                                <div class="mt-3 grid grid-cols-2 gap-3 text-sm">
                                    <div class="rounded-xl bg-gray-50 px-3 py-2">
                                        <div class="text-[11px] text-gray-400">创建者</div>
                                        <div class="mt-1 font-semibold text-gray-800">{{
                                            currentPlaylistDetails.user?.username || 'Unknown' }}</div>
                                    </div>
                                    <div class="rounded-xl bg-gray-50 px-3 py-2">
                                        <div class="text-[11px] text-gray-400">歌曲总数</div>
                                        <div class="mt-1 font-semibold text-gray-800">{{ currentPlaylistSongs.length }}
                                        </div>
                                    </div>
                                    <div class="rounded-xl bg-gray-50 px-3 py-2">
                                        <div class="text-[11px] text-gray-400">播放量</div>
                                        <div class="mt-1 font-semibold text-gray-800">{{
                                            formatNumber(currentPlaylistDetails.playCount || 0) }}</div>
                                    </div>
                                    <div class="rounded-xl bg-gray-50 px-3 py-2">
                                        <div class="text-[11px] text-gray-400">创建时间</div>
                                        <div class="mt-1 font-semibold text-gray-800">{{
                                            formatOptionalDate(currentPlaylistDetails.createdAt) }}</div>
                                    </div>
                                </div>
                            </section>

                            <section class="rounded-2xl border border-gray-100 bg-white p-4">
                                <div class="flex items-center justify-between">
                                    <h4 class="text-xs font-bold uppercase tracking-widest text-gray-400">风险提示</h4>
                                    <span class="text-[11px] text-gray-400">仅做人工审核辅助</span>
                                </div>
                                <div class="mt-3 space-y-2">
                                    <div class="rounded-xl border border-gray-100 px-3 py-2 text-sm text-gray-600">
                                        标题长度：{{ currentPlaylistDetails.title?.length || 0 }} 字
                                    </div>
                                    <div class="rounded-xl border border-gray-100 px-3 py-2 text-sm text-gray-600">
                                        简介状态：{{ currentPlaylistDetails.description ? '已填写简介，需检查文本是否合规' :
                                        '未填写简介，重点检查标题与歌曲内容' }}
                                    </div>
                                    <div class="rounded-xl border border-gray-100 px-3 py-2 text-sm text-gray-600">
                                        歌曲内容：当前收录 {{ currentPlaylistSongs.length }} 首歌，建议抽查前几首的标题与艺人信息
                                    </div>
                                </div>
                            </section>

                            <section class="rounded-2xl border border-gray-100 bg-white p-4">
                                <div class="flex items-center justify-between">
                                    <h4 class="text-xs font-bold uppercase tracking-widest text-gray-400">歌单歌曲</h4>
                                    <span class="text-[11px] text-gray-400">用于判断歌单内容是否异常</span>
                                </div>
                                <div v-if="currentPlaylistSongs.length === 0"
                                    class="mt-4 rounded-xl bg-gray-50 px-3 py-6 text-center text-sm text-gray-400">
                                    暂无歌曲
                                </div>
                                <div v-else class="mt-3 space-y-2">
                                    <div v-for="(song, index) in currentPlaylistSongs.slice(0, 6)" :key="song.id"
                                        class="flex items-center justify-between rounded-xl border border-gray-100 px-3 py-2">
                                        <div class="min-w-0">
                                            <div class="truncate text-sm font-medium text-gray-800">{{ index + 1 }}. {{
                                                song.title }}</div>
                                            <div class="truncate text-[11px] text-gray-400">{{ song.artist }}</div>
                                        </div>
                                        <div class="shrink-0 text-[11px] font-mono text-gray-400">{{
                                            formatDuration(song.duration) }}</div>
                                    </div>
                                </div>
                            </section>

                            <section class="rounded-2xl border border-gray-100 bg-white p-4">
                                <h4 class="text-xs font-bold uppercase tracking-widest text-gray-400">审核处理</h4>
                                <textarea v-model="rejectReason" rows="4"
                                    class="mt-3 w-full resize-y rounded-xl border border-gray-200 bg-gray-50 px-3 py-3 text-sm text-gray-700 outline-none transition focus:border-primary focus:bg-white focus:ring-2 focus:ring-primary/10"
                                    placeholder="如需驳回，请填写具体原因，例如：标题低俗、简介违规、歌单内容不符合平台规范。"></textarea>
                                <div v-if="currentPlaylistDetails.rejectReason"
                                    class="mt-3 rounded-xl bg-red-50 px-3 py-2 text-xs text-red-500">
                                    历史驳回原因：{{ currentPlaylistDetails.rejectReason }}
                                </div>
                                <div class="mt-4 grid grid-cols-3 gap-2">
                                    <button @click="handleApprove(currentPlaylistDetails)"
                                        class="inline-flex items-center justify-center gap-1.5 rounded-xl bg-primary px-3 py-2 text-sm font-medium text-white transition hover:opacity-90">
                                        <div class="i-fa6-solid-check text-xs"></div>
                                        通过
                                    </button>
                                    <button @click="confirmReject"
                                        class="inline-flex items-center justify-center gap-1.5 rounded-xl bg-gray-900 px-3 py-2 text-sm font-medium text-white transition hover:bg-black">
                                        <div class="i-fa6-solid-ban text-xs"></div>
                                        驳回
                                    </button>
                                    <button @click="handleDelete(currentPlaylistDetails)"
                                        class="inline-flex items-center justify-center gap-1.5 rounded-xl border border-red-200 bg-red-50 px-3 py-2 text-sm font-medium text-red-600 transition hover:bg-red-100">
                                        <div class="i-fa6-solid-trash text-xs"></div>
                                        删除
                                    </button>
                                </div>
                            </section>
                        </div>
                    </div>

                    <div v-else class="flex flex-1 items-center justify-center px-8 text-center">
                        <div>
                            <div
                                class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-[rgba(51,91,241,0.05)] text-primary">
                                <div class="i-fa6-solid-shield-halved text-2xl"></div>
                            </div>
                            <div class="text-sm font-semibold text-gray-700">选择一条歌单开始审核</div>
                            <p class="mt-2 text-xs leading-5 text-gray-400">右侧会展示标题、简介、歌曲明细和审核处理区，适合连续审核不同可见性的歌单。</p>
                        </div>
                    </div>

                    <Transition name="fade">
                        <div
                            v-if="detailLoading"
                            class="absolute inset-0 z-10 flex items-center justify-center bg-white/70 backdrop-blur-[1px] pointer-events-none"
                        >
                            <div class="flex items-center gap-3 rounded-full border border-gray-100 bg-white px-4 py-2 text-sm text-gray-500 shadow-sm">
                                <div class="i-fa6-solid-spinner animate-spin text-primary"></div>
                                正在更新详情
                            </div>
                        </div>
                    </Transition>
                </aside>
            </div>
        </div>

        <ConfirmModal
            v-model="showConfirmModal"
            :type="confirmType"
            :title="confirmTitle"
            :content="confirmContent"
            :loading="confirmLoading"
            :confirm-text="confirmText"
            @confirm="handleConfirmAction"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { resolveUrl, formatDate, formatDuration, showToast } from '@/utils/common'
import httpClient from '@/utils/request'
import Pagination from '@/components/ui/Pagination.vue'
import ConfirmModal from '@/components/ui/ConfirmModal.vue'

interface User {
    username: string
    avatarUrl: string
}

interface Playlist {
    id: number
    title: string
    description: string
    coverUrl: string
    isPublic?: boolean
    is_public?: boolean
    playCount: number
    songCount: number
    status: number
    rejectReason: string
    createdAt?: string
    user?: User
}

interface Song {
    id: number
    title: string
    artist: string
    duration: number
}

const playlists = ref<Playlist[]>([])
const loading = ref(false)
const detailLoading = ref(false)
const currentFilter = ref(0)
const pageIndex = ref(1)
const pageSize = ref(6)
const total = ref(0)
const searchQuery = ref('')
const selectedPlaylistId = ref<number | null>(null)
const currentPlaylistDetails = ref<Playlist | null>(null)
const currentPlaylistSongs = ref<Song[]>([])
const rejectReason = ref('')
const showConfirmModal = ref(false)
const confirmLoading = ref(false)
const confirmType = ref<'info' | 'danger'>('info')
const confirmTitle = ref('')
const confirmContent = ref('')
const confirmText = ref('确定')
const confirmAction = ref<(() => Promise<void>) | null>(null)

const filterTabs = [
    { label: '待审核', value: 0, activeClass: 'bg-primary text-white shadow-sm' },
    { label: '全部歌单', value: -1, activeClass: 'bg-gray-900 text-white shadow-sm' },
    { label: '已通过', value: 1, activeClass: 'bg-emerald-600 text-white shadow-sm' },
    { label: '已驳回', value: 2, activeClass: 'bg-gray-700 text-white shadow-sm' },
]


const fetchPlaylists = async (options?: { syncSelection?: boolean }) => {
    loading.value = true
    try {
        const res = await httpClient.get('/admin/playlists', {
            params: {
                status: currentFilter.value,
                pageIndex: pageIndex.value,
                pageSize: pageSize.value,
                keyword: searchQuery.value || undefined
            },
        })

        playlists.value = res.data.playlists || []
        total.value = res.data.total || 0

        const selectedStillExists = playlists.value.some(item => item.id === selectedPlaylistId.value)
        if (options?.syncSelection) {
            const nextSelected = playlists.value.find(item => item.id === selectedPlaylistId.value) || playlists.value[0]
            if (nextSelected) {
                await openDetailsPanel(nextSelected)
            } else {
                selectedPlaylistId.value = null
                currentPlaylistDetails.value = null
                currentPlaylistSongs.value = []
                rejectReason.value = ''
            }
        } else if (!selectedStillExists) {
            selectedPlaylistId.value = null
            currentPlaylistDetails.value = null
            currentPlaylistSongs.value = []
            rejectReason.value = ''
        }
    } catch (error) {
        console.error('Failed to fetch playlists:', error)
        showToast('获取歌单列表失败', 'error')
    } finally {
        loading.value = false
    }
}

const openDetailsPanel = async (playlist: Playlist) => {
    selectedPlaylistId.value = playlist.id
    detailLoading.value = true
    try {
        const res = await httpClient.get(`/admin/playlists/${playlist.id}`)
        currentPlaylistDetails.value = res.data.playlist
        currentPlaylistSongs.value = res.data.songs || []
        rejectReason.value = currentPlaylistDetails.value?.rejectReason || ''
    } catch (error) {
        console.error('Failed to fetch playlist details', error)
        showToast('获取详情失败', 'error')
    } finally {
        detailLoading.value = false
    }
}

const changePage = (page: number) => {
    pageIndex.value = page
    fetchPlaylists()
}

const handleSearch = () => {
    pageIndex.value = 1
    fetchPlaylists({ syncSelection: true })
}

watch(currentFilter, () => {
    pageIndex.value = 1
    fetchPlaylists({ syncSelection: true })
})

const handleApprove = (playlist: Playlist) => {
    confirmTitle.value = '确认通过'
    confirmContent.value = `确定通过歌单 "${playlist.title}" 吗？`
    confirmType.value = 'info'
    confirmText.value = '确认通过'
    confirmAction.value = async () => {
        await httpClient.put(`/admin/playlists/${playlist.id}/status`, {
            status: 1,
            rejectReason: '',
        })
        showToast('审核通过成功')
        rejectReason.value = ''
        await fetchPlaylists()
        showConfirmModal.value = false
    }
    showConfirmModal.value = true
}

const confirmReject = async () => {
    if (!currentPlaylistDetails.value) return
    if (!rejectReason.value.trim()) {
        showToast('请先填写驳回原因', 'error')
        return
    }

    try {
        await httpClient.put(`/admin/playlists/${currentPlaylistDetails.value.id}/status`, {
            status: 2,
            rejectReason: rejectReason.value.trim(),
        })
        showToast('已驳回该歌单')
        await fetchPlaylists()
    } catch (error) {
        showToast('操作失败', 'error')
    }
}

const handleDelete = (playlist: Playlist) => {
    confirmTitle.value = '确认删除'
    confirmContent.value = `确定永久删除歌单 "${playlist.title}" 吗？此操作无法撤销。`
    confirmType.value = 'danger'
    confirmText.value = '确认删除'
    confirmAction.value = async () => {
        await httpClient.delete(`/admin/playlists/${playlist.id}`)
        showToast('删除成功')
        await fetchPlaylists()
        showConfirmModal.value = false
    }
    showConfirmModal.value = true
}

const handleConfirmAction = async () => {
    if (!confirmAction.value) return
    confirmLoading.value = true
    try {
        await confirmAction.value()
    } catch (error) {
        showToast(confirmType.value === 'danger' ? '删除失败' : '操作失败', 'error')
    } finally {
        confirmLoading.value = false
    }
}

const statusColor = (status: number) => {
    switch (status) {
        case 0:
            return 'bg-primary'
        case 1:
            return 'bg-emerald-600'
        case 2:
            return 'bg-gray-700'
        default:
            return 'bg-gray-500'
    }
}

const statusText = (status: number) => {
    switch (status) {
        case 0:
            return '待审核'
        case 1:
            return '已通过'
        case 2:
            return '已驳回'
        default:
            return '未知'
    }
}

const visibilityText = (playlist: Playlist | null) => {
    if (!playlist) return '歌单'
    const isPublic = playlist.isPublic ?? playlist.is_public
    return isPublic ? '公开歌单' : '私密歌单'
}

const formatOptionalDate = (value?: string) => {
    return value ? formatDate(value) : '-'
}

const formatNumber = (num: number) => {
    if (num >= 10000) {
        return (num / 10000).toFixed(1) + 'w'
    }
    return String(num)
}

onMounted(() => {
    fetchPlaylists({ syncSelection: true })
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
