import request from "../utils/request";
import type { ApiResponse, Playlist, SongDTO } from "@/types/api";

// 歌单详情响应类型
export interface PlaylistDetailResponse extends Playlist {
  songs: SongDTO[];
}

export interface PlaylistPageResponse {
  playlists: Playlist[];
  total: number;
  pageIndex: number;
  pageSize: number;
  hasMore?: boolean;
}

/**
 * 获取我的歌单列表
 */
export function fetchMyPlaylists(): Promise<ApiResponse<Playlist[]>> {
  return request.get("/playlists/my");
}

/**
 * 获取首页热门推荐歌单
 * 仅返回公开且审核通过的歌单
 */
export function fetchRecommendedPlaylists(params?: {
  limit?: number;
  pageIndex?: number;
  pageSize?: number;
}): Promise<ApiResponse<PlaylistPageResponse>> {
  return request.get("/playlists/recommended", {
    params,
  });
}

/**
 * 搜索公开且审核通过的歌单
 */
export function searchPlaylists(params: {
  keyword: string;
  pageIndex?: number;
  pageSize?: number;
}): Promise<ApiResponse<PlaylistPageResponse>> {
  return request.get("/playlists/search", { params });
}

/**
 * 获取歌单详情（包含歌曲列表）
 * @param id 歌单ID
 */
export function getPlaylistDetails(id: number | string) {
  return request.get(`/playlists/${id}`);
}

/**
 * 创建歌单
 * @param data 包含标题、描述、公开状态、封面文件的对象
 */
export function createPlaylist(data: {
  title: string;
  description?: string;
  isPublic?: boolean; // 默认为 true (后端逻辑)
  coverFile?: File;
}) {
  const formData = new FormData();
  formData.append("title", data.title);

  if (data.description) {
    formData.append("description", data.description);
  }

  // 处理布尔值：FormData 只能传字符串
  // 后端默认是 "true"，只有传 "false" 才会设为私有
  if (data.isPublic !== undefined) {
    formData.append("isPublic", String(data.isPublic));
  }

  if (data.coverFile) {
    formData.append("coverFile", data.coverFile);
  }

  return request.post("/playlists", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

/**
 * 更新歌单
 * @param id 歌单ID
 * @param data 更新的数据 (可选)
 */
export function updatePlaylist(
  id: number | string,
  data: {
    title?: string;
    description?: string;
    isPublic?: boolean;
    coverFile?: File;
  }
) {
  const formData = new FormData();

  if (data.title) {
    formData.append("title", data.title);
  }

  if (data.description !== undefined) {
    formData.append("description", data.description);
  }

  // 注意：后端 Update 逻辑中，如果不传 isPublic 可能会默认为 false
  // 建议在调用更新时，最好明确传入当前的 isPublic 状态
  if (data.isPublic !== undefined) {
    formData.append("isPublic", String(data.isPublic));
  }

  if (data.coverFile) {
    formData.append("coverFile", data.coverFile);
  }

  return request.put(`/playlists/${id}`, formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
}

/**
 * 删除歌单
 * @param id 歌单ID
 */
export function deletePlaylist(id: number | string) {
  return request.delete(`/playlists/${id}`);
}

/**
 * 添加歌曲到歌单
 * @param playlistId 歌单ID
 * @param songId 歌曲ID
 */
export function addSongToPlaylist(
  playlistId: number | string,
  songId: number | string
) {
  return request.post(`/playlists/${playlistId}/songs`, { songId });
}

/**
 * 从歌单移除歌曲
 * @param playlistId 歌单ID
 * @param songId 歌曲ID
 */
export function removeSongFromPlaylist(
  playlistId: number | string,
  songId: number | string
) {
  return request.delete(`/playlists/${playlistId}/songs/${songId}`);
}
