import { type Song } from "@/stores/player";
import type { SongDTO } from "@/types/api";

const BASE_URL = import.meta.env.VITE_BACKEND_URL || "";

/**
 * 统一处理 URL 拼接
 * 解决后端返回相对路径、多余斜杠等问题
 */
export const resolveUrl = (url?: string): string => {
  if (!url) return "";
  if (url.startsWith("http") || url.startsWith("blob")) return url;
  // 去除开头的斜杠，防止双斜杠
  const cleanUrl = url.startsWith("/") ? url.substring(1) : url;
  // 拼接并正则清理多余的 // (保留协议后的 //)
  return `${BASE_URL}/${cleanUrl}`.replace(/([^:])\/\/+/g, "$1/");
};

export const resolveCover = (url?: string): string => {
  return resolveUrl(url);
};

/**
 * 格式化时长 (秒 -> mm:ss)
 * @param seconds 秒数
 */
export const formatDuration = (seconds: number): string => {
  if (!seconds || isNaN(seconds)) return "00:00";
  const min = Math.floor(seconds / 60);
  const sec = Math.floor(seconds % 60);
  return `${min.toString().padStart(2, "0")}:${sec
    .toString()
    .padStart(2, "0")}`;
};

/**
 * 格式化日期 (YYYY年MM月DD日)
 */
export const formatDate = (dateStr: string): string => {
  if (!dateStr) return "";
  return new Date(dateStr).toLocaleDateString("zh-CN", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
};



/**
 * 发送全局 Toast 消息
 */
export const showToast = (
  message: string,
  type: "success" | "error" | "info" | "warning" = "success"
) => {
  window.dispatchEvent(
    new CustomEvent("play-toast", { detail: { message, type } })
  );
};

/**
 * 统一转换后端 SongDTO -> 前端播放器 Song 结构。
 */
export const formatSong = (song: SongDTO): Song | null => {
  if (!song) return null;
  try {
    if (!song.id || !song.title || !song.artist || !song.fileUrl) {
      return null;
    }
    return {
      songId: song.id,
      title: song.title,
      artist: song.artist,
      album: song.album,
      duration: song.duration,
      coverUrl: resolveUrl(song.coverUrl),
      lyricUrl: resolveUrl(song.lyricUrl),
      fileUrl: resolveUrl(song.fileUrl),
    } as Song;
  } catch (e) {
    console.warn("Song format error:", e);
    return null;
  }
};

export const formatSongs = (songs: SongDTO[] = []): Song[] => {
  return songs.map(formatSong).filter((song): song is Song => song !== null);
};
