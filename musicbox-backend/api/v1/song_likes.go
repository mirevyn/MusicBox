package v1

import (
	"errors"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SongLikesController 处理歌曲点赞相关的请求
type SongLikesController struct {
	Service *service.SongLikesService
}

// NewSongLikesController 创建并返回一个新的 SongLikesController 实例
func NewSongLikesController() *SongLikesController {
	return &SongLikesController{
		Service: &service.SongLikesService{},
	}
}

// @Summary 切换点赞状态
// @Description 对指定的歌曲进行点赞或取消点赞
// @Tags 收藏与点赞
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param songId path uint true "歌曲 ID"
// @Success 200 {object} utils.Response "操作成功"
// @Failure 404 {object} utils.Response "歌曲不存在"
// @Router /song-likes/{songId} [post]
func (slc *SongLikesController) ToggleSongLike(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return
	}
	userID := val.(uint)

	songIDStr := c.Param("songId")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		utils.Result(c, 400, nil, "无效的歌曲ID")
		return
	}

	isLiked, err := slc.Service.ToggleSongLike(userID, uint(songID))
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			utils.Result(c, 404, nil, "歌曲不存在")
		} else if errors.Is(err, service.ErrUserNotFound) {
			utils.Result(c, 404, nil, "用户不存在")
		} else {
			utils.Result(c, 500, nil, "处理点赞失败")
		}
		return
	}

	msg := "取消点赞成功"
	if isLiked {
		msg = "点赞成功"
	}

	utils.Result(c, 200, gin.H{"isLiked": isLiked}, msg)
}

// @Summary 查询点赞状态
// @Description 检查当前登录用户是否已点赞指定歌曲
// @Tags 收藏与点赞
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param songId path uint true "歌曲 ID"
// @Success 200 {object} utils.Response "查询成功"
// @Router /song-likes/{songId} [get]
func (slc *SongLikesController) GetLikeStatus(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return
	}
	userID := val.(uint)

	songIDStr := c.Param("songId")
	songID, err := strconv.ParseUint(songIDStr, 10, 32)
	if err != nil {
		utils.Result(c, 400, nil, "无效的歌曲ID")
		return
	}

	isLiked, err := slc.Service.GetLikeStatus(userID, uint(songID))
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			utils.Result(c, 404, nil, "歌曲不存在")
		} else {
			utils.Result(c, 500, nil, "查询失败: "+err.Error())
		}
		return
	}

	utils.Result(c, 200, gin.H{"isLiked": isLiked}, "查询成功")
}

// @Summary 获取我的收藏
// @Description 获取当前用户已点赞的所有歌曲列表
// @Tags 收藏与点赞
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]model.Song} "查询成功"
// @Router /song-likes [get]
func (slc *SongLikesController) GetMyLikedSongs(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return
	}
	userID := val.(uint)

	songs, err := slc.Service.GetLikedSongsByUser(userID)
	if err != nil {
		utils.Result(c, 500, nil, "查询收藏列表失败")
		return
	}

	utils.Result(c, 200, songs, "查询成功")
}
