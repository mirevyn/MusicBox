package v1

import (
	"errors"
	"fmt"
	"musicbox-backend/internal/config"
	"musicbox-backend/internal/model"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"os"
)

// @Summary 获取歌曲列表
// @Description 支持按关键字搜索和分页排序
// @Tags 歌曲管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键字"
// @Param pageIndex query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param sortBy query string false "排序字段" default(upload_at)
// @Param order query string false "排序方式" Enums(asc, desc) default(desc)
// @Success 200 {object} utils.Response "成功"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	keyword := c.Query("keyword")
	pageIndex, err := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	if err != nil || pageIndex < 1 {
		pageIndex = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}

	sortBy := c.DefaultQuery("sortBy", "upload_at")
	order := c.DefaultQuery("order", "desc")

	songs, total, err := service.GetSongs(keyword, pageIndex, pageSize, sortBy, order)
	if err != nil {
		utils.Result(c, 500, nil, "获取歌曲失败")
		return
	}

	utils.Result(c, 200, gin.H{
		"list":      songs,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	}, "成功")
}

// @Summary 获取歌曲详情
// @Description 根据歌曲 ID 获取歌曲完整信息
// @Tags 歌曲管理
// @Accept json
// @Produce json
// @Param id path uint true "歌曲 ID"
// @Success 200 {object} utils.Response{data=model.Song} "成功"
// @Failure 400 {object} utils.Response "无效的 ID"
// @Failure 404 {object} utils.Response "歌曲未找到"
// @Router /songs/{id} [get]
func GetSongByID(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的歌曲ID")
		return
	}

	song, err := service.GetSongByID(uint(songID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(c, 404, nil, "歌曲未找到")
		} else {
			utils.Result(c, 500, nil, "获取失败")
		}
		return
	}

	utils.Result(c, 200, song, "成功")
}

// @Summary 上传歌曲
// @Description 上传音频文件及其相关信息（需管理员权限）
// @Tags 歌曲管理
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "歌曲标题"
// @Param artist formData string true "艺术家"
// @Param album formData string false "专辑"
// @Param file formData file true "歌曲音频文件"
// @Param cover formData file false "封面图片"
// @Param lyric formData file false "歌词文件"
// @Success 201 {object} utils.Response{data=model.Song} "上传成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 500 {object} utils.Response "服务器内部错误"
// @Router /songs/upload [post]
func UploadSong(c *gin.Context) {
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	album := c.PostForm("album")

	if title == "" || artist == "" {
		utils.Result(c, 400, nil, "标题和艺术家是必填项")
		return
	}

	songFile, err := c.FormFile("file")
	if err != nil {
		utils.Result(c, 400, nil, "歌曲文件是必传的")
		return
	}

	songPath, err := service.SaveUploadedFile(c, songFile, config.Conf.Upload.SongPath)
	if err != nil {
		resultUploadError(c, err, "无法保存歌曲文件")
		return
	}
	uploadedPaths := []string{songPath}

	var coverPath, lyricPath string
	if coverFile, err := c.FormFile("cover"); err == nil {
		coverPath, err = service.SaveUploadedFile(c, coverFile, config.Conf.Upload.CoverPath)
		if err != nil {
			cleanupUploadedFiles(uploadedPaths...)
			resultUploadError(c, err, "无法保存封面文件")
			return
		}
		uploadedPaths = append(uploadedPaths, coverPath)
	}
	if lyricFile, err := c.FormFile("lyric"); err == nil {
		lyricPath, err = service.SaveUploadedFile(c, lyricFile, config.Conf.Upload.LyricPath)
		if err != nil {
			cleanupUploadedFiles(uploadedPaths...)
			resultUploadError(c, err, "无法保存歌词文件")
			return
		}
		uploadedPaths = append(uploadedPaths, lyricPath)
	}

	duration, err := utils.GetAudioDuration(songPath)
	if err != nil {
		fmt.Printf("获取音频时长错误: %v\n", err)
	}

	newSong := &model.Song{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
		FileUrl:  songPath,
		CoverUrl: coverPath,
		LyricUrl: lyricPath,
	}

	createdSong, err := service.CreateSong(newSong)
	if err != nil {
		cleanupUploadedFiles(uploadedPaths...)
		utils.Result(c, 500, nil, "无法创建歌曲记录: "+err.Error())
		return
	}

	utils.Result(c, 201, createdSong, "歌曲上传成功")
}

// @Summary 更新歌曲
// @Description 修改现有歌曲的信息或替换相关文件（需管理员权限）
// @Tags 歌曲管理
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "歌曲 ID"
// @Param title formData string false "歌曲标题"
// @Param artist formData string false "艺术家"
// @Param album formData string false "专辑"
// @Param file formData file false "新歌曲音频文件"
// @Param cover formData file false "新封面图片"
// @Param lyric formData file false "新歌词文件"
// @Success 200 {object} utils.Response{data=model.Song} "更新成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 404 {object} utils.Response "歌曲未找到"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的歌曲ID")
		return
	}

	existingSong, err := service.GetSongByID(uint(songID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(c, 404, nil, "歌曲未找到")
			return
		}
		utils.Result(c, 500, nil, "获取失败")
		return
	}

	updates := make(map[string]interface{})
	newUploadedPaths := make([]string, 0, 3)
	oldResourcePaths := make([]string, 0, 3)
	if title := c.PostForm("title"); title != "" {
		updates["title"] = title
	}
	if artist := c.PostForm("artist"); artist != "" {
		updates["artist"] = artist
	}
	if album := c.PostForm("album"); album != "" {
		updates["album"] = album
	}

	if songFile, err := c.FormFile("file"); err == nil {
		newSongPath, err := service.SaveUploadedFile(c, songFile, config.Conf.Upload.SongPath)
		if err != nil {
			resultUploadError(c, err, "无法保存新歌曲文件")
			return
		}
		newUploadedPaths = append(newUploadedPaths, newSongPath)
		if existingSong.FileUrl != "" {
			oldResourcePaths = append(oldResourcePaths, existingSong.FileUrl)
		}
		updates["file_url"] = newSongPath
		if newDuration, err := utils.GetAudioDuration(newSongPath); err == nil {
			updates["duration"] = newDuration
		}
	}

	if coverFile, err := c.FormFile("cover"); err == nil {
		newCoverPath, err := service.SaveUploadedFile(c, coverFile, config.Conf.Upload.CoverPath)
		if err != nil {
			cleanupUploadedFiles(newUploadedPaths...)
			resultUploadError(c, err, "无法保存新封面文件")
			return
		}
		newUploadedPaths = append(newUploadedPaths, newCoverPath)
		if existingSong.CoverUrl != "" {
			oldResourcePaths = append(oldResourcePaths, existingSong.CoverUrl)
		}
		updates["cover_url"] = newCoverPath
	}

	if lyricFile, err := c.FormFile("lyric"); err == nil {
		newLyricPath, err := service.SaveUploadedFile(c, lyricFile, config.Conf.Upload.LyricPath)
		if err != nil {
			cleanupUploadedFiles(newUploadedPaths...)
			resultUploadError(c, err, "无法保存新歌词文件")
			return
		}
		newUploadedPaths = append(newUploadedPaths, newLyricPath)
		if existingSong.LyricUrl != "" {
			oldResourcePaths = append(oldResourcePaths, existingSong.LyricUrl)
		}
		updates["lyric_url"] = newLyricPath
	}

	updatedSong, err := service.UpdateSong(uint(songID), updates)
	if err != nil {
		cleanupUploadedFiles(newUploadedPaths...)
		utils.Result(c, 500, nil, "无法更新歌曲记录: "+err.Error())
		return
	}
	cleanupUploadedFiles(oldResourcePaths...)

	utils.Result(c, 200, updatedSong, "歌曲更新成功")
}

// @Summary 删除歌曲
// @Description 根据 ID 删除歌曲记录及相关文件（需管理员权限）
// @Tags 歌曲管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌曲 ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 404 {object} utils.Response "歌曲未找到"
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	idParam := c.Param("id")
	songID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的歌曲ID")
		return
	}

	err = service.DeleteSong(uint(songID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(c, 404, nil, "歌曲未找到")
			return
		}
		utils.Result(c, 500, nil, "无法删除歌曲: "+err.Error())
		return
	}

	utils.Result(c, 200, nil, "歌曲删除成功")
}

// @Summary 一键下载所有歌曲包
// @Description 整理数据库中所有歌曲（含歌词、封面）并打包为 Zip 直接下载
// @Tags 歌曲管理
// @Produce application/octet-stream
// @Success 200 {file} file "歌曲压缩包"
// @Failure 404 {object} utils.Response "没有可导出的歌曲"
// @Failure 500 {object} utils.Response "导出失败"
// @Router /songs/export [post]
func ExportSongs(c *gin.Context) {
	// 1. 创建临时根目录用于存放导出文件
	tempRoot, err := os.MkdirTemp("", "musicbox_export_*")
	if err != nil {
		utils.Result(c, 500, nil, "无法创建临时目录: "+err.Error())
		return
	}
	// 确保函数结束（包括下载完成）后清理所有临时文件
	defer os.RemoveAll(tempRoot)

	// 调用 Service 执行导出和压缩
	successCount, total, zipPath, err := service.ExportAllSongs(tempRoot)
	if err != nil {
		utils.Result(c, 500, nil, "导出过程出错: "+err.Error())
		return
	}

	if total == 0 {
		utils.Result(c, 404, nil, "数据库中没有歌曲可供导出")
		return
	}

	// 检查压缩包是否生成
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		utils.Result(c, 500, nil, "生成的压缩包不存在")
		return
	}

	// 发送文件
	fmt.Printf("开始下载: %s (共 %d/%d 首)\n", zipPath, successCount, total)
	c.FileAttachment(zipPath, "MusicBox_Export.zip")
}

func cleanupUploadedFiles(paths ...string) {
	for _, filePath := range paths {
		if filePath == "" {
			continue
		}
		_ = utils.DeleteFile(filePath)
	}
}

func resultUploadError(c *gin.Context, err error, message string) {
	status := 500
	if errors.Is(err, service.ErrInvalidUpload) {
		status = 400
	}
	utils.Result(c, status, nil, message+": "+err.Error())
}
