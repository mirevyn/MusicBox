package service

import (
	"errors"
	"fmt"
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/model"
	"musicbox-backend/utils"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// GetSongs 根据过滤条件检索歌曲分页列表。
func GetSongs(keyword string, pageIndex, pageSize int, sortBy string, order string) ([]model.Song, int64, error) {
	var songs []model.Song

	// 创建查询构建器
	query := global.DB.Model(&model.Song{})

	// 应用过滤器
	if keyword != "" {
		// 在歌曲名或艺术家名中搜索
		query = query.Where("title LIKE ? OR artist LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 排序字段白名单校验
	allowedFields := map[string]bool{
		"upload_at": true,
		"id":        true,
		"title":     true,
		"artist":    true,
		"duration":  true,
	}

	// 默认排序
	finalSort := "upload_at"
	if allowedFields[sortBy] {
		finalSort = sortBy
	}

	// 排序方向校验
	finalOrder := "desc"
	if order == "asc" {
		finalOrder = "asc"
	}

	query = query.Order(songOrderClause(finalSort, finalOrder))

	// 使用通用的分页函数
	total, err := utils.Paginate(query, pageIndex, pageSize, &songs)
	if err != nil {
		return nil, 0, err
	}

	return songs, total, nil
}

func songOrderClause(sortBy, order string) string {
	orderClause := sortBy + " " + order
	if sortBy != "id" {
		orderClause += ", id " + order
	}
	return orderClause
}

// CreateSong 创建一个新的歌曲记录到数据库。
func CreateSong(song *model.Song) (*model.Song, error) {
	result := global.DB.Create(song)
	if result.Error != nil {
		return nil, result.Error
	}
	go GetNotificationService().PushDashboardRefresh()
	return song, nil
}

// UpdateSong 更新数据库中的歌曲记录。
func UpdateSong(id uint, updates map[string]interface{}) (*model.Song, error) {
	var song model.Song
	// 查找现有歌曲
	if err := global.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// 执行更新
	if err := global.DB.Model(&song).Updates(updates).Error; err != nil {
		return nil, err
	}

	go GetNotificationService().PushDashboardRefresh()

	// 返回更新后的歌曲
	return &song, nil
}

// GetSongByID 根据ID获取歌曲。
func GetSongByID(id uint) (*model.Song, error) {
	var song model.Song
	if err := global.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &song, nil
}

// DeleteSong 根据ID删除歌曲及其关联文件。
func DeleteSong(id uint) error {
	var song model.Song
	// 查找歌曲
	if err := global.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// 从数据库删除歌曲记录
	if err := global.DB.Delete(&song).Error; err != nil {
		return err
	}

	// 删除关联文件（忽略错误，因为文件可能不存在，且数据库记录已经删除）
	if song.FileUrl != "" {
		_ = utils.DeleteFile(song.FileUrl)
	}
	if song.CoverUrl != "" {
		_ = utils.DeleteFile(song.CoverUrl)
	}
	if song.LyricUrl != "" {
		_ = utils.DeleteFile(song.LyricUrl)
	}

	go GetNotificationService().PushDashboardRefresh()

	return nil
}

// ExportAllSongs 导出数据库中所有歌曲到指定目录。
// 返回 (成功数量, 总数量, 压缩包路径, 错误)
func ExportAllSongs(targetDir string) (int, int, string, error) {
	var songs []model.Song
	if err := global.DB.Find(&songs).Error; err != nil {
		return 0, 0, "", err
	}

	total := len(songs)
	if total == 0 {
		return 0, 0, "", nil
	}

	// 准备目录结构
	exportBase := filepath.Join(targetDir, "MusicBox_Export_"+time.Now().Format("20060102150405"))
	subDirs := []string{"music", "lyrics", "covers"}
	for _, d := range subDirs {
		if err := utils.EnsureDir(filepath.Join(exportBase, d)); err != nil {
			return 0, total, "", err
		}
	}

	// 导出文件
	successCount := 0
	for _, song := range songs {
		if err := exportSongRecord(&song, exportBase); err == nil {
			successCount++
		}
	}

	// 压缩目录
	zipPath := exportBase + ".zip"
	if err := utils.ZipDir(exportBase, zipPath); err != nil {
		return successCount, total, "", fmt.Errorf("压缩失败: %w", err)
	}

	return successCount, total, zipPath, nil
}

// exportSongRecord 内部使用的单条歌曲导出逻辑
func exportSongRecord(song *model.Song, exportBase string) error {
	// 定义导出函数：复制文件并保持格式 {标题}-{艺术家}.{后缀}
	exportFile := func(originalPath, subDir, suffix string) error {
		if originalPath == "" {
			return nil
		}

		// 检查路径是否存在
		if _, err := os.Stat(originalPath); os.IsNotExist(err) {
			return nil // 文件不存在则跳过
		}

		newFileName := fmt.Sprintf("%s-%s.%s", song.Title, song.Artist, suffix)
		destPath := filepath.Join(exportBase, subDir, newFileName)

		return utils.CopyFile(originalPath, destPath)
	}

	// 导出音频
	ext := filepath.Ext(song.FileUrl)
	if ext == "" {
		ext = "mp3"
	} else {
		ext = ext[1:]
	}
	_ = exportFile(song.FileUrl, "music", ext)

	// 导出封面
	if song.CoverUrl != "" {
		coverExt := filepath.Ext(song.CoverUrl)
		if coverExt != "" {
			coverExt = coverExt[1:]
		} else {
			coverExt = "jpg"
		}
		_ = exportFile(song.CoverUrl, "covers", coverExt)
	}

	// 导出歌词
	if song.LyricUrl != "" {
		lyricExt := filepath.Ext(song.LyricUrl)
		if lyricExt != "" {
			lyricExt = lyricExt[1:]
		} else {
			lyricExt = "lrc"
		}
		_ = exportFile(song.LyricUrl, "lyrics", lyricExt)
	}

	return nil
}
