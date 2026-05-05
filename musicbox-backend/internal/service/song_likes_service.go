package service

import (
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
)

// SongLikesService 歌曲点赞服务
type SongLikesService struct{}

// ToggleSongLike 切换用户对歌曲的点赞状态
func (s *SongLikesService) ToggleSongLike(userID uint, songID uint) (bool, error) {
	// 检查歌曲是否存在
	var songExists int64
	if err := global.DB.Model(&model.Song{}).Where("id = ?", songID).Count(&songExists).Error; err != nil {
		return false, err
	}
	if songExists == 0 {
		return false, ErrSongNotFound
	}

	// 检查用户是否存在
	var userExists int64
	if err := global.DB.Model(&model.User{}).Where("id = ?", userID).Count(&userExists).Error; err != nil {
		return false, err
	}
	if userExists == 0 {
		return false, ErrUserNotFound
	}

	like := model.SongLike{UserID: userID, SongID: songID}
	var count int64

	// 检查点赞记录是否存在
	if err := global.DB.Model(&model.SongLike{}).Where("user_id = ? AND song_id = ?", userID, songID).Count(&count).Error; err != nil {
		return false, err
	}

	// 执行切换逻辑
	if count == 0 {
		// 记录不存在，执行点赞
		if err := global.DB.Create(&like).Error; err != nil {
			return false, err
		}
		return true, nil // 点赞成功
	} else {
		// 记录已存在，执行取消点赞
		if err := global.DB.Where("user_id = ? AND song_id = ?", userID, songID).Delete(&like).Error; err != nil {
			return false, err
		}
		return false, nil // 取消点赞成功
	}
}

// GetLikeStatus 查询用户是否已点赞某首歌曲。
func (s *SongLikesService) GetLikeStatus(userID uint, songID uint) (bool, error) {
	// 检查歌曲是否存在
	var songExists int64
	if err := global.DB.Model(&model.Song{}).Where("id = ?", songID).Count(&songExists).Error; err != nil {
		return false, err
	}
	if songExists == 0 {
		return false, ErrSongNotFound
	}

	// 查询点赞记录
	var count int64
	if err := global.DB.Model(&model.SongLike{}).Where("user_id = ? AND song_id = ?", userID, songID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetLikedSongsByUser 获取用户点赞的所有歌曲列表
func (s *SongLikesService) GetLikedSongsByUser(userID uint) ([]model.Song, error) {
	var likedSongs []model.Song

	// 使用 Join 查询
	err := global.DB.
		Joins("JOIN song_likes ON songs.id = song_likes.song_id").
		Where("song_likes.user_id = ?", userID).
		Order("song_likes.liked_at DESC").
		Find(&likedSongs).Error

	if err != nil {
		return nil, err
	}

	return likedSongs, nil
}