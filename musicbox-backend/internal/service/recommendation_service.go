package service

import (
	"sort"
	"strings"
	"time"

	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"

	"gorm.io/gorm"
)

type RecommendationService struct{}

// RecordPlayHistory 将一次成功播放落库，供推荐服务做偏好计算
func (s *RecommendationService) RecordPlayHistory(userID, songID uint, duration int, source string) error {
	var count int64
	if err := global.DB.Model(&model.Song{}).Where("id = ?", songID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return ErrSongNotFound
	}

	history := model.PlayHistory{
		UserID:   userID,
		SongID:   songID,
		Duration: duration,
		Source:   strings.TrimSpace(source),
	}

	if err := global.DB.Create(&history).Error; err != nil {
		return err
	}
	go GetNotificationService().PushDashboardRefresh()
	return nil
}

// GetDailyRecommendations 先用播放历史建模，再叠加点赞偏好，最后用新歌补齐结果集
func (s *RecommendationService) GetDailyRecommendations(userID uint, limit int) ([]model.Song, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	historySongs, err := s.getRecentHistorySongs(userID, 80)
	if err != nil {
		return nil, err
	}

	likedSongs, err := (&SongLikesService{}).GetLikedSongsByUser(userID)
	if err != nil {
		return nil, err
	}

	artistWeights := map[string]float64{}
	albumWeights := map[string]float64{}
	excludedIDs := map[uint]struct{}{}

	for i, song := range historySongs {
		// 越新的播放记录权重越高，避免旧习惯长期主导结果
		recencyWeight := float64(maxInt(1, 8-(i/10)))
		addPreferenceWeight(artistWeights, song.Artist, recencyWeight*2.6)
		addPreferenceWeight(albumWeights, song.Album, recencyWeight*1.2)
		if i < 24 {
			excludedIDs[song.ID] = struct{}{}
		}
	}

	for _, song := range likedSongs {
		// 点赞是强偏好信号，权重高于普通播放记录
		addPreferenceWeight(artistWeights, song.Artist, 6)
		addPreferenceWeight(albumWeights, song.Album, 3)
		excludedIDs[song.ID] = struct{}{}
	}

	candidates, err := s.getCandidateSongs(limit*6, excludedIDs, artistWeights, albumWeights)
	if err != nil {
		return nil, err
	}

	scored := scoreSongs(candidates, artistWeights, albumWeights)
	recommendations := pickTopSongs(scored, limit)

	if len(recommendations) < limit {
		// 偏好样本不足时，用未听过的新歌补齐，避免页面空洞
		fallbackSongs, err := s.getFallbackSongs(limit*2, excludedIDs, recommendations)
		if err == nil {
			recommendations = append(recommendations, fallbackSongs...)
		}
		
		// 如果还是不够 limit，说明库里绝大部分歌都被 excludedIDs 过滤掉了
		// 退化为无视过滤，强行拿最近更新的歌补足
		if len(recommendations) < limit {
			var raw []model.Song
			global.DB.Model(&model.Song{}).Order("upload_at DESC").Limit(limit * 2).Find(&raw)
			
			existsMap := make(map[uint]struct{}, len(recommendations))
			for _, r := range recommendations {
				existsMap[r.ID] = struct{}{}
			}
			
			for _, r := range raw {
				if _, ok := existsMap[r.ID]; !ok {
					recommendations = append(recommendations, r)
					existsMap[r.ID] = struct{}{}
					if len(recommendations) >= limit {
						break
					}
				}
			}
		}

		if len(recommendations) > limit {
			recommendations = recommendations[:limit]
		}
	}

	return recommendations, nil
}

func (s *RecommendationService) getRecentHistorySongs(userID uint, limit int) ([]model.Song, error) {
	var songs []model.Song
	err := global.DB.
		Table("songs").
		Joins("JOIN play_histories ON play_histories.song_id = songs.id").
		Where("play_histories.user_id = ?", userID).
		Order("play_histories.played_at DESC").
		Limit(limit).
		Find(&songs).Error
	return songs, err
}

func (s *RecommendationService) getCandidateSongs(
	limit int,
	excludedIDs map[uint]struct{},
	artistWeights map[string]float64,
	albumWeights map[string]float64,
) ([]model.Song, error) {
	query := global.DB.Model(&model.Song{})

	excludeList := mapKeysUint(excludedIDs)
	if len(excludeList) > 0 {
		query = query.Where("id NOT IN ?", excludeList)
	}

	topArtists := topPreferenceKeys(artistWeights, 5)
	topAlbums := topPreferenceKeys(albumWeights, 5)

	if len(topArtists) > 0 || len(topAlbums) > 0 {
		// genre 已移除，这里只保留歌手和专辑两个稳定维度做候选召回
		query = query.Scopes(buildPreferenceQuery(topArtists, topAlbums))
	}

	var songs []model.Song
	err := query.Order("upload_at DESC").Limit(limit).Find(&songs).Error
	if err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *RecommendationService) getFallbackSongs(limit int, excludedIDs map[uint]struct{}, selected []model.Song) ([]model.Song, error) {
	for _, song := range selected {
		excludedIDs[song.ID] = struct{}{}
	}

	query := global.DB.Model(&model.Song{})
	excludeList := mapKeysUint(excludedIDs)
	if len(excludeList) > 0 {
		query = query.Where("id NOT IN ?", excludeList)
	}

	var songs []model.Song
	err := query.Order("upload_at DESC").Limit(limit).Find(&songs).Error
	return songs, err
}

func addPreferenceWeight(target map[string]float64, key string, weight float64) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}
	target[key] += weight
}

type weightedKey struct {
	Key    string
	Weight float64
}

type scoredSong struct {
	Song  model.Song
	Score float64
}

func topPreferenceKeys(weights map[string]float64, limit int) []string {
	if len(weights) == 0 || limit <= 0 {
		return nil
	}

	items := make([]weightedKey, 0, len(weights))
	for key, weight := range weights {
		items = append(items, weightedKey{Key: key, Weight: weight})
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Weight == items[j].Weight {
			return items[i].Key < items[j].Key
		}
		return items[i].Weight > items[j].Weight
	})

	if len(items) > limit {
		items = items[:limit]
	}

	keys := make([]string, 0, len(items))
	for _, item := range items {
		keys = append(keys, item.Key)
	}
	return keys
}

func buildPreferenceQuery(artists, albums []string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db
		hasCondition := false

		if len(artists) > 0 {
			if hasCondition {
				query = query.Or("artist IN ?", artists)
			} else {
				query = query.Where("artist IN ?", artists)
				hasCondition = true
			}
		}
		if len(albums) > 0 {
			if hasCondition {
				query = query.Or("album IN ?", albums)
			} else {
				query = query.Where("album IN ?", albums)
			}
		}

		return query
	}
}

func scoreSongs(candidates []model.Song, artistWeights, albumWeights map[string]float64) []scoredSong {
	now := time.Now()
	scored := make([]scoredSong, 0, len(candidates))

	for _, song := range candidates {
		// 歌手权重最高，其次是专辑；新歌会拿到一个轻微的时效加分
		score := artistWeights[strings.TrimSpace(song.Artist)]*3.4 +
			albumWeights[strings.TrimSpace(song.Album)]*1.2

		if !song.UploadAt.IsZero() {
			days := now.Sub(song.UploadAt).Hours() / 24
			if days < 30 {
				score += (30 - days) / 10
			}
		}

		scored = append(scored, scoredSong{Song: song, Score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			if scored[i].Song.UploadAt.Equal(scored[j].Song.UploadAt) {
				return scored[i].Song.ID > scored[j].Song.ID
			}
			return scored[i].Song.UploadAt.After(scored[j].Song.UploadAt)
		}
		return scored[i].Score > scored[j].Score
	})

	return scored
}

func pickTopSongs(scored []scoredSong, limit int) []model.Song {
	if len(scored) == 0 || limit <= 0 {
		return nil
	}

	if len(scored) > limit {
		scored = scored[:limit]
	}

	result := make([]model.Song, 0, len(scored))
	for _, item := range scored {
		result = append(result, item.Song)
	}
	return result
}

func mapKeysUint(source map[uint]struct{}) []uint {
	if len(source) == 0 {
		return nil
	}
	keys := make([]uint, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	return keys
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
