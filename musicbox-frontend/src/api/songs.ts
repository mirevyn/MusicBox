import request from "../utils/request";
import type {
  ApiResponse,
  FetchSongsParams,
  PaginatedResponse,
  SongDTO,
  LikeStatusResponse,
} from "@/types/api";
import type { AxiosProgressEvent } from "axios";

// 获取歌曲列表（带分页）
export function fetchSongs(
  params: FetchSongsParams
): Promise<ApiResponse<PaginatedResponse<SongDTO>>> {
  return request.get("/songs", { params });
}

// 获取单个歌曲详情
export function getSongById(
  id: number | string
): Promise<ApiResponse<SongDTO>> {
  return request.get(`/songs/${id}`);
}

// 切换歌曲点赞状态
export function toggleSongLike(
  songId: number | string
): Promise<ApiResponse<LikeStatusResponse>> {
  return request.post(`/song-likes/${songId}`);
}

// 获取歌曲点赞状态
export function getSongLikeStatus(
  songId: number | string
): Promise<ApiResponse<LikeStatusResponse>> {
  return request.get(`/song-likes/${songId}`);
}

export function getUserLikedSongs(): Promise<ApiResponse<SongDTO[]>> {
  return request.get("/song-likes");
}

// 记录当前用户的播放历史，用于私享推荐计算
export function recordPlayHistory(data: {
  songId: number | string;
  duration?: number;
  source?: string;
}): Promise<ApiResponse<null>> {
  return request.post("/play-histories", data);
}

// 获取基于播放历史和点赞偏好的私享推荐歌曲
export function getDailyRecommendations(
  limit = 20
): Promise<ApiResponse<{ date: string; songs: SongDTO[] }>> {
  return request.get("/recommendations/daily", {
    params: { limit },
  });
}

// 上传歌曲（管理员功能）
export function uploadSong(
  formData: FormData,
  onProgress?: (percent: number) => void
): Promise<ApiResponse<SongDTO>> {
  return request.post("/songs/upload", formData, {
    onUploadProgress: onProgress
      ? (progressEvent: AxiosProgressEvent) => {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / (progressEvent.total || 1)
        );
        onProgress(percentCompleted);
      }
      : undefined,
  });
}

// 更新歌曲（管理员功能）
export function updateSong(
  songId: number | string,
  formData: FormData,
  onProgress?: (percent: number) => void
): Promise<ApiResponse<SongDTO>> {
  return request.put(`/songs/${songId}`, formData, {
    onUploadProgress: onProgress
      ? (progressEvent: AxiosProgressEvent) => {
        const percentCompleted = Math.round(
          (progressEvent.loaded * 100) / (progressEvent.total || 1)
        );
        onProgress(percentCompleted);
      }
      : undefined,
  });
}

// 删除歌曲（管理员功能）
export function deleteSong(
  songId: number | string
): Promise<ApiResponse<null>> {
  return request.delete(`/songs/${songId}`);
}

// 导出所有歌曲（管理员功能）
export function exportSongs(): Promise<Blob> {
  return request.post("/songs/export", {}, {
    responseType: 'blob'
  });
}
