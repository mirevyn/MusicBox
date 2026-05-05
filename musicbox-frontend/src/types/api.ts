// API 响应通用类型
export interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data?: T;
}

// 用户相关类型
export interface User {
  id: number;
  username: string;
  avatarUrl?: string;
  avatar_url?: string;
  role?: string;
  createdAt?: string;
}

export interface LoginDTO {
  username: string;
  password: string;
}

export interface RegisterDTO {
  username: string;
  password: string;
}

export interface ChangePasswordDTO {
  oldPassword: string;
  newPassword: string;
}

export interface LoginResponse {
  token: string;
  user?: User;
}

// 后端歌曲 DTO
export interface SongDTO {
  id: number | string;
  title: string;
  artist: string;
  album?: string;
  coverUrl?: string;
  fileUrl: string;
  lyricUrl?: string;
  duration?: number;
  uploadAt?: string;
}

// 旧组件迁移期间保留别名
export type RawSong = SongDTO;

export interface FetchSongsParams {
  page?: number;
  pageIndex?: number;
  pageSize?: number;
  keyword?: string;
  sortBy?: string;
  order?: "asc" | "desc";
}

export interface PaginatedResponse<T> {
  list: T[];
  total: number;
  pageIndex?: number;
  pageSize: number;
}

// 歌单相关类型
export interface Playlist {
  id: number;
  title: string;
  description?: string;
  coverUrl?: string;
  cover_url?: string;
  isPublic: boolean;
  is_public?: boolean;
  status?: number;
  rejectReason?: string;
  reject_reason?: string;
  songCount: number;
  song_count?: number;
  userId: number;
  createdAt?: string;
  user?: {
    id: number;
    username: string;
    nickname?: string;
    avatarUrl?: string;
  };
}

export interface CreatePlaylistDTO {
  title: string;
  description?: string;
  isPublic?: boolean;
}

export interface UpdatePlaylistDTO {
  title?: string;
  description?: string;
  isPublic?: boolean;
}

// 点赞相关类型
export interface LikeStatusResponse {
  isLiked: boolean;
}
