package service

import (
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
	"musicbox-backend/utils"

	"gorm.io/gorm"
)

// PlaylistService 处理歌单相关的业务逻辑
type PlaylistService struct{}

type PlaylistUpdateInput struct {
	ID          uint
	Title       *string
	Description *string
	CoverURL    *string
	IsPublic    *bool
}

// CreatePlaylist 为用户创建一个新歌单。
func (s *PlaylistService) CreatePlaylist(userID uint, title, description, coverUrl string, isPublic bool) (*model.Playlist, error) {
	playlist := model.Playlist{
		UserID:      userID,
		Title:       title,
		Description: description,
		CoverUrl:    coverUrl,
		IsPublic:    isPublic,
	}
	if err := global.DB.Create(&playlist).Error; err != nil {
		return nil, err
	}

	// 当模型字段带有 default 标签时，GORM 在插入数据时可能会跳过零值字段。
	// 因此需要显式保存审核状态，确保公开歌单始终进入“待审核”状态。
	if isPublic {
		if err := global.DB.Model(&playlist).Updates(map[string]interface{}{
			"status":        0,
			"reject_reason": "",
		}).Error; err != nil {
			return nil, err
		}
		playlist.Status = 0
		playlist.RejectReason = ""
		GetNotificationService().PushPlaylistPendingNotification(playlist.Title, playlist.ID)
	}

	return &playlist, nil
}

// UpdatePlaylist 更新歌单信息
func (s *PlaylistService) UpdatePlaylist(userID uint, input PlaylistUpdateInput) error {
	var existingPlaylist model.Playlist

	if err := global.DB.First(&existingPlaylist, input.ID).Error; err != nil {
		return err
	}

	if existingPlaylist.UserID != userID {
		return ErrForbiddenAccess
	}
	wasPublic := existingPlaylist.IsPublic
	wasPending := existingPlaylist.Status == 0
	nextIsPublic := existingPlaylist.IsPublic
	if input.IsPublic != nil {
		nextIsPublic = *input.IsPublic
	}

	// 准备更新的数据
	updates := make(map[string]interface{})

	if input.Description != nil {
		updates["description"] = *input.Description
	}
	if input.IsPublic != nil {
		updates["is_public"] = *input.IsPublic
	}
	if input.Title != nil {
		updates["title"] = *input.Title
	}

	oldCoverURL := ""
	if input.CoverURL != nil {
		if existingPlaylist.CoverUrl != "" && existingPlaylist.CoverUrl != *input.CoverURL {
			oldCoverURL = existingPlaylist.CoverUrl
		}
		updates["cover_url"] = *input.CoverURL
	}

	if nextIsPublic && len(updates) > 0 {
		// 如果播放列表是公开的并且正在被更新，则需要重新审核
		updates["status"] = 0
		updates["reject_reason"] = ""
	}
	if len(updates) == 0 {
		return nil
	}
	if err := global.DB.Model(&existingPlaylist).Updates(updates).Error; err != nil {
		return err
	}
	if oldCoverURL != "" {
		_ = utils.DeleteFile(oldCoverURL)
	}

	var latestPlaylist model.Playlist
	if err := global.DB.First(&latestPlaylist, existingPlaylist.ID).Error; err != nil {
		return err
	}

	if (!wasPublic || !wasPending) && latestPlaylist.IsPublic && latestPlaylist.Status == 0 {
		GetNotificationService().PushPlaylistPendingNotification(latestPlaylist.Title, latestPlaylist.ID)
	}

	return nil
}

// DeletePlaylist 删除歌单
func (s *PlaylistService) DeletePlaylist(playlistID uint, userID uint) error {
	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}

	if playlist.UserID != userID {
		return ErrForbiddenAccess
	}

	if err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("playlist_id = ?", playlistID).Delete(&model.PlaylistSong{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&playlist).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if playlist.CoverUrl != "" {
		_ = utils.DeleteFile(playlist.CoverUrl)
	}

	return nil
}

// GetMyPlaylists 获取用户的所有歌单 (包含歌曲数量统计)
func (s *PlaylistService) GetMyPlaylists(userID uint) ([]model.Playlist, error) {
	type PlaylistResult struct {
		model.Playlist
		SongCount int `gorm:"column:song_count"`
	}

	var results []PlaylistResult

	err := global.DB.Table("playlists").
		Select("playlists.*, count(playlist_songs.song_id) as song_count").
		Joins("LEFT JOIN playlist_songs ON playlist_songs.playlist_id = playlists.id").
		Where("playlists.user_id = ?", userID).
		Group("playlists.id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	playlists := make([]model.Playlist, len(results))
	for i, r := range results {
		playlists[i] = r.Playlist
		playlists[i].SongCount = r.SongCount
	}

	return playlists, nil
}

// GetRecommendedPlaylists 获取首页推荐歌单，只返回公开且审核通过的数据。
func (s *PlaylistService) GetRecommendedPlaylists(pageIndex, pageSize int) ([]model.Playlist, int64, error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 4
	}
	if pageSize > 24 {
		pageSize = 24
	}

	query := global.DB.Model(&model.Playlist{}).
		Where("is_public = ? AND status = ?", true, 1)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var playlists []model.Playlist
	if err := global.DB.
		Preload("User").
		Where("is_public = ? AND status = ?", true, 1).
		Order("play_count DESC, created_at DESC, id DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&playlists).Error; err != nil {
		return nil, 0, err
	}

	for i := range playlists {
		var count int64
		if err := global.DB.Model(&model.PlaylistSong{}).
			Where("playlist_id = ?", playlists[i].ID).
			Count(&count).Error; err != nil {
			return nil, 0, err
		}
		playlists[i].SongCount = int(count)
	}

	return playlists, total, nil
}

// SearchPlaylists 通过关键词搜索公开且审核通过的歌单
func (s *PlaylistService) SearchPlaylists(keyword string, pageIndex, pageSize int) ([]model.Playlist, int64, error) {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	query := global.DB.Model(&model.Playlist{}).
		Where("is_public = ? AND status = ?", true, 1)

	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var playlists []model.Playlist
	if err := query.
		Preload("User").
		Order("play_count DESC, created_at DESC, id DESC").
		Offset((pageIndex - 1) * pageSize).
		Limit(pageSize).
		Find(&playlists).Error; err != nil {
		return nil, 0, err
	}

	// 统计歌单下歌曲数量
	for i := range playlists {
		var count int64
		if err := global.DB.Model(&model.PlaylistSong{}).
			Where("playlist_id = ?", playlists[i].ID).
			Count(&count).Error; err != nil {
			return nil, 0, err
		}
		playlists[i].SongCount = int(count)
	}

	return playlists, total, nil
}

// AddSongToPlaylist 将歌曲添加到歌单
func (s *PlaylistService) AddSongToPlaylist(userID, playlistID, songID uint) error {
	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}
	if playlist.UserID != userID {
		return ErrForbiddenAccess
	}

	var count int64
	global.DB.Model(&model.Song{}).Where("id = ?", songID).Count(&count)
	if count == 0 {
		return ErrNotFound
	}

	playlistSong := model.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
	}

	var existCount int64
	global.DB.Model(&model.PlaylistSong{}).
		Where("playlist_id = ? AND song_id = ?", playlistID, songID).
		Count(&existCount)
	if existCount > 0 {
		return ErrDuplicateEntry
	}

	return global.DB.Create(&playlistSong).Error
}

// RemoveSongFromPlaylist 从歌单中移除歌曲
func (s *PlaylistService) RemoveSongFromPlaylist(userID, playlistID, songID uint) error {
	var playlist model.Playlist
	if err := global.DB.First(&playlist, playlistID).Error; err != nil {
		return err
	}
	if playlist.UserID != userID {
		return ErrForbiddenAccess
	}

	result := global.DB.Where("playlist_id = ? AND song_id = ?", playlistID, songID).Delete(&model.PlaylistSong{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// GetPlaylistDetails 获取歌单详情
func (s *PlaylistService) GetPlaylistDetails(playlistID uint, currentUserID uint) (*model.Playlist, []model.Song, error) {
	var playlist model.Playlist

	if err := global.DB.Preload("User").First(&playlist, playlistID).Error; err != nil {
		return nil, nil, err
	}

	if playlist.UserID != currentUserID && (!playlist.IsPublic || playlist.Status != 1) {
		return nil, nil, ErrForbiddenAccess
	}

	// 增加播放量/热度
	global.DB.Model(&playlist).UpdateColumn("play_count", gorm.Expr("play_count + ?", 1))

	var songs []model.Song
	err := global.DB.
		Table("songs").
		Joins("JOIN playlist_songs ps ON songs.id = ps.song_id").
		Where("ps.playlist_id = ?", playlistID).
		Order("ps.sort_order ASC, ps.added_at DESC").
		Find(&songs).Error

	if err != nil {
		return nil, nil, err
	}

	playlist.SongCount = len(songs)

	return &playlist, songs, nil
}
