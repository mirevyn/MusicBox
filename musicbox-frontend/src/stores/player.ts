import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { toggleSongLike, getSongLikeStatus } from "@/api/songs";

// --- 类型定义 ---

export interface Song {
  songId: number | string;
  title: string;
  artist: string;
  album?: string;
  coverUrl: string;
  fileUrl: string;
  lyricUrl?: string;
  duration?: number;
  [key: string]: any;
}

// 定义简单的后端响应接口
interface LikeResponse {
  code?: number;
  msg?: string;
  data?: {
    isLiked: boolean;
  };
  [key: string]: any;
}

export const PLAY_MODE = {
  SEQUENCE: 0,
  LOOP: 1,
  RANDOM: 2,
};

export const usePlayerStore = defineStore("player", () => {
  // --- 响应式状态 ---
  const playList = ref<Song[]>([]);
  const currentSongIndex = ref(-1);
  const playing = ref(false);
  const mode = ref(PLAY_MODE.SEQUENCE);
  const volume = ref(0.8);
  const showPlayList = ref(false);
  const likedSongs = ref<Set<string>>(new Set());
  const playlistSource = ref("播放列表");

  // --- 计算属性 ---
  const currentSong = computed<Song | null>(() => {
    if (
      currentSongIndex.value < 0 ||
      currentSongIndex.value >= playList.value.length
    ) {
      return null;
    }
    return playList.value[currentSongIndex.value] || null;
  });

  const isLiked = computed(() => {
    if (!currentSong.value) return false;
    return likedSongs.value.has(String(currentSong.value.songId));
  });

  // --- 内部辅助函数 (Helper) ---

  // 发送全局 Toast 消息
  const notify = (message: string, type: "success" | "error" = "success") => {
    window.dispatchEvent(
      new CustomEvent("play-toast", { detail: { message, type } })
    );
  };

  // 检查是否登录
  const checkLogin = (): boolean => {
    const token = localStorage.getItem("token");
    if (!token) {
      notify("请先登录", "error");
      return false;
    }
    return true;
  };

  // --- 业务动作 ---

  function setPlayList(songs: Song[], source?: string) {
    playList.value = songs;
    if (source) {
      playlistSource.value = source;
    }
    if (songs.length > 0 && currentSongIndex.value === -1) {
      currentSongIndex.value = 0;
    }
    if (songs.length === 0) {
      currentSongIndex.value = -1;
      playing.value = false;
    }
  }

  function playByIndex(index: number) {
    if (index >= 0 && index < playList.value.length) {
      currentSongIndex.value = index;
      playing.value = true;
    }
  }

  function playSong(song: Song) {
    const index = playList.value.findIndex((s) => s.songId === song.songId);
    if (index !== -1) {
      playByIndex(index);
    } else {
      playList.value.unshift(song);
      playByIndex(0);
    }
  }

  function removeFromPlayList(songId: number | string) {
    const indexToRemove = playList.value.findIndex((s) => s.songId === songId);
    if (indexToRemove === -1) return;

    if (indexToRemove === currentSongIndex.value) {
      if (playList.value.length > 1) {
        playNext(false);
      } else {
        clearPlayList();
        return;
      }
    }

    if (indexToRemove < currentSongIndex.value) {
      currentSongIndex.value--;
    }

    playList.value.splice(indexToRemove, 1);
  }

  function reorderPlayList(fromIndex: number, toIndex: number) {
    const list = playList.value;
    if (
      fromIndex === toIndex ||
      fromIndex < 0 ||
      toIndex < 0 ||
      fromIndex >= list.length ||
      toIndex >= list.length
    ) {
      return;
    }

    const currentBeforeMove = currentSong.value;
    const currentSongId = currentBeforeMove?.songId;
    const [movedSong] = list.splice(fromIndex, 1);
    if (!movedSong) return;

    list.splice(toIndex, 0, movedSong);

    if (!currentBeforeMove) {
      if (currentSongIndex.value >= list.length) {
        currentSongIndex.value = list.length - 1;
      }
      return;
    }

    const nextCurrentIndexByReference = list.findIndex(
      (song) => song === currentBeforeMove
    );
    const nextCurrentIndex =
      nextCurrentIndexByReference !== -1
        ? nextCurrentIndexByReference
        : list.findIndex((song) => song.songId === currentSongId);

    if (nextCurrentIndex !== -1) {
      currentSongIndex.value = nextCurrentIndex;
    }
  }

  function clearPlayList() {
    playList.value = [];
    currentSongIndex.value = -1;
    playing.value = false;
  }

  function togglePlay() {
    if (!currentSong.value) return;
    playing.value = !playing.value;
  }

  function playNext(isEndOfSong = false) {
    if (playList.value.length === 0) return;
    if (playList.value.length === 1) {
      if (!playing.value) playing.value = true;
      return;
    }

    if (isEndOfSong && mode.value === PLAY_MODE.LOOP) {
      playing.value = true;
      return;
    }

    let newIndex = currentSongIndex.value;
    if (mode.value === PLAY_MODE.RANDOM) {
      do {
        newIndex = Math.floor(Math.random() * playList.value.length);
      } while (
        newIndex === currentSongIndex.value &&
        playList.value.length > 1
      );
    } else {
      newIndex = (newIndex + 1) % playList.value.length;
    }

    currentSongIndex.value = newIndex;
    playing.value = true;
  }

  function playPrev() {
    if (playList.value.length === 0) return;
    if (playList.value.length === 1) {
      playing.value = true;
      return;
    }

    let newIndex = currentSongIndex.value;
    if (mode.value === PLAY_MODE.RANDOM) {
      do {
        newIndex = Math.floor(Math.random() * playList.value.length);
      } while (
        newIndex === currentSongIndex.value &&
        playList.value.length > 1
      );
    } else {
      newIndex = (newIndex - 1 + playList.value.length) % playList.value.length;
    }

    currentSongIndex.value = newIndex;
    playing.value = true;
  }

  function togglePlayMode() {
    mode.value = (mode.value + 1) % 3;
  }

  function setVolume(vol: number) {
    volume.value = Math.max(0, Math.min(1, vol));
  }

  function togglePlayList() {
    showPlayList.value = !showPlayList.value;
  }

  function updateSongInfo(newSong: Song) {
    const index = playList.value.findIndex((s) => s.songId === newSong.songId);
    if (index !== -1) {
      const target = playList.value[index];
      if (!target) return;
      target.title = newSong.title;
      target.artist = newSong.artist;
      target.album = newSong.album;
      target.coverUrl = newSong.coverUrl;
      target.lyricUrl = newSong.lyricUrl;
      if (target.fileUrl !== newSong.fileUrl) {
        target.fileUrl = newSong.fileUrl;
      }
    }
  }

  function addToQueue(song: Song) {
    const exists = playList.value.some((s) => s.songId === song.songId);
    if (exists) {
      notify(`《${song.title}》已在播放列表中`);
      return;
    }
    playList.value.push(song);
    notify(`已添加到播放队列`);
  }

  async function toggleLike() {
    if (!currentSong.value) return;
    if (!checkLogin()) return;
    const id = String(currentSong.value.songId);
    const isCurrentlyLiked = likedSongs.value.has(id);
    const newSet = new Set(likedSongs.value);
    if (isCurrentlyLiked) {
      newSet.delete(id);
    } else {
      newSet.add(id);
    }
    likedSongs.value = newSet;
    try {
      await toggleSongLike(id);
    } catch (error) {
      console.error("点赞操作失败:", error);
      const revertSet = new Set(likedSongs.value);
      if (isCurrentlyLiked) {
        revertSet.add(id);
      } else {
        revertSet.delete(id);
      }
      likedSongs.value = revertSet;
      notify("操作失败，请检查网络", "error");
    }
  }

  async function syncCurrentSongLikeStatus() {
    if (!currentSong.value) return;
    const id = String(currentSong.value.songId);
    const token = localStorage.getItem("token");
    if (!token) {
      if (likedSongs.value.has(id)) {
        const newSet = new Set(likedSongs.value);
        newSet.delete(id);
        likedSongs.value = newSet;
      }
      return;
    }
    try {
      const res = (await getSongLikeStatus(id)) as unknown as LikeResponse;
      const isLikedServer = res?.data?.isLiked === true;
      const newSet = new Set(likedSongs.value);
      if (isLikedServer) {
        newSet.add(id);
      } else {
        newSet.delete(id);
      }
      likedSongs.value = newSet;
    } catch (error) {
      console.warn("获取歌曲点赞状态失败", error);
    }
  }

  function reset() {
    playList.value = [];
    currentSongIndex.value = -1;
    playing.value = false;
    mode.value = PLAY_MODE.SEQUENCE;
    showPlayList.value = false;
    likedSongs.value = new Set();
    playlistSource.value = "播放列表";
  }

  return {
    playList,
    currentSongIndex,
    playing,
    mode,
    volume,
    showPlayList,
    likedSongs,
    playlistSource,
    currentSong,
    isLiked,
    setPlayList,
    playByIndex,
    playSong,
    removeFromPlayList,
    reorderPlayList,
    clearPlayList,
    reset,
    togglePlay,
    playNext,
    playPrev,
    togglePlayMode,
    setVolume,
    togglePlayList,
    toggleLike,
    syncCurrentSongLikeStatus,
    updateSongInfo,
    addToQueue,
  };
});
