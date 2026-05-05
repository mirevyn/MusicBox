<template>
    <div class="flex flex-col">
        <!-- 骨架屏 Loading -->
        <div v-if="loading" class="space-y-6 animate-pulse">
            <div class="h-64 bg-gray-200/50 rounded-[2rem]"></div>
            <div class="space-y-3">
                <div v-for="n in 6" :key="n" class="h-16 bg-gray-200/50 rounded-xl"></div>
            </div>
        </div>

        <div v-else class="space-y-6 animate-fade-in">
            <!-- 审核驳回 Banner -->
            <div v-if="isOwner && playlist.isPublic && playlist.status === 2"
                class="bg-red-50 border border-red-200 rounded-[2rem] p-4 md:p-6 flex flex-col md:flex-row gap-4 items-start md:items-center justify-between shadow-sm">
                <div class="flex gap-4">
                    <div class="w-10 h-10 rounded-full bg-red-100 flex items-center justify-center shrink-0 text-red-500">
                        <div class="i-fa6-solid-triangle-exclamation text-lg"></div>
                    </div>
                    <div>
                        <h3 class="font-bold text-red-800">此歌单未通过审核</h3>
                        <p class="text-sm text-red-600 mt-1 whitespace-pre-wrap">
                            原因：{{ playlist.rejectReason || '内容不符合社区规范' }}
                        </p>
                    </div>
                </div>
                <button @click="handleEdit"
                    class="shrink-0 px-5 py-2.5 bg-red-600 text-white font-bold rounded-full shadow-lg shadow-red-500/30 hover:bg-red-700 hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-2">
                    <div class="i-fa6-solid-pen"></div> 去修改
                </button>
            </div>

            <!-- 1. 头部信息卡片 -->
            <div
                class="relative flex flex-col md:flex-row gap-6 items-center md:items-end bg-white/60 backdrop-blur-md p-6 md:p-8 rounded-[2rem] border border-white shadow-sm overflow-hidden">

                <!-- 背景装饰 -->
                <div
                    class="absolute top-0 right-0 w-64 h-64 bg-blue-200/20 rounded-full blur-3xl -translate-y-1/2 translate-x-1/3 pointer-events-none">
                </div>

                <!-- 封面区 -->
                <div class="relative group shrink-0 z-10">
                    <div
                        class="w-40 h-40 md:w-48 md:h-48 rounded-3xl flex items-center justify-center text-white shadow-xl shadow-blue-500/10 transform group-hover:scale-105 transition-transform duration-500 overflow-hidden bg-gray-100">
                        <img v-if="playlist.coverUrl && !coverImgError" :src="resolveUrl(playlist.coverUrl)"
                            class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                            alt="Cover" @error="coverImgError = true" />
                        <div v-else
                            class="w-full h-full bg-gradient-to-br from-blue-400 to-cyan-500 flex items-center justify-center">
                            <div class="i-fa6-solid-music text-5xl md:text-6xl drop-shadow-md opacity-80"></div>
                        </div>
                    </div>
                </div>

                <!-- 信息区 -->
                <div class="flex-1 text-center md:text-left space-y-4 w-full z-10">
                    <div>
                        <div class="flex items-center justify-center md:justify-start gap-2 mb-2">
                            <span class="px-2 py-0.5 rounded-md text-[10px] font-bold uppercase tracking-wider border"
                                :class="playlist.isPublic ? 'border-blue-200 text-blue-500 bg-blue-50' : 'border-gray-200 text-gray-500 bg-gray-50'">
                                {{ playlist.isPublic ? 'Public' : 'Private' }}
                            </span>
                            <span class="text-xs font-bold text-gray-400 uppercase tracking-widest">Playlist</span>
                        </div>
                        <h1 class="text-3xl md:text-4xl font-extrabold text-gray-800 tracking-tight leading-tight">
                            {{ playlist.title }}
                        </h1>
                        <p v-if="playlist.description" class="text-sm text-gray-500 mt-2 line-clamp-2 max-w-2xl">
                            {{ playlist.description }}
                        </p>
                    </div>

                    <!-- 用户信息与元数据 -->
                    <div class="flex items-center justify-center md:justify-start gap-3 text-sm text-gray-500">
                        <div
                            class="flex items-center gap-1.5 p-1 pr-3 bg-white/50 rounded-full border border-white/50 shadow-sm">
                            <div v-if="playlist.user?.avatarUrl && !avatarImgError" class="w-6 h-6 rounded-full overflow-hidden">
                                <img :src="resolveUrl(playlist.user.avatarUrl)" class="w-full h-6 object-cover"
                                    @error="avatarImgError = true" />
                            </div>
                            <div v-else
                                class="w-6 h-6 rounded-full bg-gradient-to-br from-blue-400 to-indigo-500 text-white flex items-center justify-center text-xs font-bold shadow-sm">
                                {{ (playlist.user?.username || 'U').charAt(0).toUpperCase() }}
                            </div>
                            <span class="font-bold text-gray-700 text-xs">{{ playlist.user?.nickname ||
                                playlist.user?.username }}</span>
                        </div>
                        <span class="text-gray-300">|</span>
                        <span>{{ formatDate(playlist.createdAt ?? '') }}</span>
                        <span class="text-gray-300">|</span>
                        <span>{{ songs.length }} 首歌曲</span>
                    </div>

                    <!-- 操作按钮组 -->
                    <div class="flex gap-3 justify-center md:justify-start pt-2">
                        <button @click="playAll" :disabled="songs.length === 0"
                            class="px-8 py-2.5 bg-blue-500 text-white rounded-full font-bold shadow-lg shadow-blue-500/30 hover:bg-blue-600 hover:shadow-blue-500/40 hover:-translate-y-0.5 active:scale-95 transition-all flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed">
                            <div class="i-fa6-solid-play"></div> 播放全部
                        </button>

                        <!-- 所有者操作 -->
                        <template v-if="isOwner">
                            <button @click="handleEdit"
                                class="px-5 py-2.5 bg-white text-gray-700 border border-gray-200 rounded-full font-medium hover:bg-gray-50 hover:border-gray-300 transition-colors flex items-center gap-2 shadow-sm">
                                <div class="i-fa6-solid-pen text-sm"></div> 编辑
                            </button>
                            <button @click="handleDelete"
                                class="w-10 h-10 flex items-center justify-center bg-white text-gray-500 border border-gray-200 rounded-full hover:bg-red-50 hover:text-red-500 hover:border-red-200 transition-colors shadow-sm"
                                title="删除歌单">
                                <div class="i-fa6-solid-trash text-sm"></div>
                            </button>
                        </template>
                    </div>
                </div>
            </div>

            <!-- 2. 歌曲列表区域 -->
            <SongList :is-empty="formattedSongs.length === 0" :count="formattedSongs.length" theme="blue" action-label="操作">
                <SongRow v-for="(song, index) in formattedSongs" :key="song.songId" :song="song" :index="index"
                    :is-active="isCurrentSong(song)" :is-playing="!!player.playing" theme="blue" :show-remove="isOwner"
                    @play="playSong" @add-to-queue="player.addToQueue" @remove="handleRemoveSong" @add-to-playlist="onAddToPlaylist" />

                <template #empty>
                    <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                        <div class="i-fa6-solid-record-vinyl text-3xl opacity-30"></div>
                    </div>
                    <p class="text-sm">这个歌单还没有添加歌曲</p>
                </template>
            </SongList>
        </div>
        
        <!-- 编辑弹窗 -->
        <PlaylistModal
            v-model="showEditModal"
            mode="edit"
            :initial-data="playlist"
            @success="handleEditSuccess"
        />

        <!-- 添加到歌单弹窗 -->
        <AddToPlaylistModal
            v-model="showAddModal"
            :song-id="currentSongId"
            @create-new="openCreateModal"
        />

        <!-- 新建歌单弹窗 (复用 PlaylistModal) -->
        <PlaylistModal
            v-model="showCreateModal"
            mode="create"
            @success="handleCreateSuccess"
        />
        
        <!-- 确认弹窗 -->
        <ConfirmModal
            v-model="confirmState.show"
            :title="confirmState.title"
            :content="confirmState.content"
            :type="confirmState.type"
            :confirm-text="confirmState.confirmText"
            :loading="confirmState.loading"
            @confirm="confirmState.onConfirm"
        />
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { getPlaylistDetails, deletePlaylist, removeSongFromPlaylist } from '@/api/playlists';
import { usePlayerStore, type Song } from '@/stores/player';
import { useUserStore } from '@/stores/user';
import { resolveUrl, formatDate, showToast, formatSongs } from '@/utils/common';
import type { SongDTO } from '@/types/api';
import SongList from '@/components/song/SongList.vue';
import SongRow from '@/components/song/SongRow.vue';
import PlaylistModal from '@/components/playlist/PlaylistModal.vue';
import AddToPlaylistModal from '@/components/playlist/AddToPlaylistModal.vue';
import ConfirmModal from '@/components/ui/ConfirmModal.vue';


// 组件内类型定义
interface PlaylistUser {
    id: number;
    username: string;
    nickname?: string;
    avatarUrl?: string;
}

interface PlaylistDetail {
    id: number;
    title: string;
    description?: string;
    coverUrl?: string;
    isPublic: boolean;
    userId: number;
    songCount: number;
    createdAt?: string;
    status?: number;
    rejectReason?: string;
    user?: PlaylistUser;
}

const route = useRoute();
const router = useRouter();
const player = usePlayerStore();
const userStore = useUserStore();

// 状态
const loading = ref(true);
const playlist = ref<PlaylistDetail>({
    id: 0,
    title: '',
    isPublic: true,
    userId: 0,
    songCount: 0
});
const songs = ref<SongDTO[]>([]);
const showEditModal = ref(false);

const coverImgError = ref(false);
const avatarImgError = ref(false);

watch(() => playlist.value.coverUrl, () => {
    coverImgError.value = false;
});
watch(() => playlist.value.user?.avatarUrl, () => {
    avatarImgError.value = false;
});

// 确认弹窗状态
const confirmState = ref({
    show: false,
    title: '',
    content: '',
    type: 'info' as 'info' | 'danger',
    confirmText: '确定',
    loading: false,
    onConfirm: () => {}
});

// 计算属性
const formattedSongs = computed(() => {
    return formatSongs(songs.value);
});

const isOwner = computed(() => {
    return !!(userStore.user?.id && playlist.value.userId === userStore.user.id);
});

const isCurrentSong = (song: Song): boolean => {
    return player.currentSong?.songId === song.songId;
};

// 生命周期
onMounted(async () => {
    const id = route.params.id as string;
    if (!id) {
        router.push('/library');
        return;
    }

    // 确保用户状态已加载
    if (!userStore.user) {
        await userStore.fetchProfile();
    }
    
    await loadPlaylistData(id);
});

const loadPlaylistData = async (id: string) => {
    loading.value = true;
    try {
        const res = await getPlaylistDetails(id) as unknown as { code: number; data: { playlist: PlaylistDetail; songs: SongDTO[] } };
        if (res.code === 200 && res.data) {
            playlist.value = res.data.playlist;
            songs.value = res.data.songs || [];
        }
    } catch (error) {
        console.error("加载歌单失败", error);
        showToast("加载歌单失败", "error");
    } finally {
        loading.value = false;
    }
};

// 方法
const playAll = () => {
    if (formattedSongs.value.length === 0) return;
    
    player.setPlayList([...formattedSongs.value], playlist.value.title);
    player.playByIndex(0);
    showToast(`开始播放: ${playlist.value.title}`);
};

const playSong = (index: number) => {
    // 检查是否是同一个播放列表
    const isSameList = player.playList.length === formattedSongs.value.length && 
        player.playList[0]?.songId === formattedSongs.value[0]?.songId;
    
    if (!isSameList) {
        player.setPlayList([...formattedSongs.value], playlist.value.title);
    }
    player.playByIndex(index);
};



const handleDelete = () => {
    confirmState.value = {
        show: true,
        title: '删除歌单',
        content: '确定要删除这个歌单吗？此操作不可恢复，请谨慎操作。',
        type: 'danger',
        confirmText: '确认删除',
        loading: false,
        onConfirm: async () => {
            confirmState.value.loading = true;
            try {
                await deletePlaylist(playlist.value.id);
                showToast("歌单删除成功");
                router.push('/library');
            } catch (error) {
                console.error(error);
                showToast("删除失败", "error");
            } finally {
                confirmState.value.loading = false;
                confirmState.value.show = false;
            }
        }
    };
};

const handleEdit = () => {
    showEditModal.value = true;
};

const handleEditSuccess = () => {
    const id = route.params.id as string;
    if (id) loadPlaylistData(id);
};

const handleRemoveSong = (songId: string | number) => {
    confirmState.value = {
        show: true,
        title: '移除歌曲',
        content: '确定将这首歌从歌单移除吗？',
        type: 'danger',
        confirmText: '确认移除',
        loading: false,
        onConfirm: async () => {
            confirmState.value.loading = true;
            try {
                await removeSongFromPlaylist(playlist.value.id, songId);
                songs.value = songs.value.filter(s => s.id !== songId);
                playlist.value.songCount = Math.max(0, playlist.value.songCount - 1);
                showToast("移除成功");
            } catch (error) {
                console.error(error);
                showToast("移除失败", "error");
            } finally {
                confirmState.value.loading = false;
                confirmState.value.show = false;
            }
        }
    };
};

// 添加到歌单的处理函数
const showAddModal = ref(false);
const showCreateModal = ref(false);
const currentSongId = ref<number | string | null>(null);

const onAddToPlaylist = (song: Song) => {
    currentSongId.value = song.songId || song.id;
    showAddModal.value = true;
};

const openCreateModal = () => {
    showAddModal.value = false;
    showCreateModal.value = true;
};

const handleCreateSuccess = () => {
    showToast('歌单创建成功，请重新选择添加到歌单');
    showAddModal.value = true;
};
</script>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(15px);
    }

    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
