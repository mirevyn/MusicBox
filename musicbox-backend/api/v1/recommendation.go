package v1

import (
	"errors"
	"strconv"
	"time"

	"musicbox-backend/internal/service"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	Service *service.RecommendationService
}

func NewRecommendationController() *RecommendationController {
	return &RecommendationController{
		Service: &service.RecommendationService{},
	}
}

type RecordPlayHistoryReq struct {
	SongID   uint   `json:"songId" binding:"required"`
	Duration int    `json:"duration"`
	Source   string `json:"source"`
}

// @Summary 记录播放历史
// @Description 记录用户播放歌曲的行为，用于个性化推荐
// @Tags 推荐系统
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body RecordPlayHistoryReq true "播放历史信息"
// @Success 200 {object} utils.Response "记录成功"
// @Router /play-histories [post]
func (rc *RecommendationController) RecordPlayHistory(c *gin.Context) {
	uid, ok := getAuthUserID(c)
	if !ok {
		return
	}

	var req RecordPlayHistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误")
		return
	}

	if req.Duration < 0 {
		req.Duration = 0
	}

	err := rc.Service.RecordPlayHistory(uid, req.SongID, req.Duration, req.Source)
	if err != nil {
		if errors.Is(err, service.ErrSongNotFound) {
			utils.Result(c, 404, nil, "歌曲不存在")
			return
		}
		utils.Result(c, 500, nil, "记录播放历史失败")
		return
	}

	utils.Result(c, 200, nil, "记录成功")
}

// @Summary 获取每日推荐
// @Description 根据用户的听歌偏好生成每日推荐歌曲列表
// @Tags 推荐系统
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param limit query int false "推荐数量" default(20)
// @Success 200 {object} utils.Response "获取成功"
// @Router /recommendations/daily [get]
func (rc *RecommendationController) GetDailyRecommendations(c *gin.Context) {
	uid, ok := getAuthUserID(c)
	if !ok {
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}

	songs, err := rc.Service.GetDailyRecommendations(uid, limit)
	if err != nil {
		utils.Result(c, 500, nil, "获取每日推荐失败")
		return
	}

	utils.Result(c, 200, gin.H{
		"date":  time.Now().Format("2006-01-02"),
		"songs": songs,
	}, "获取成功")
}
