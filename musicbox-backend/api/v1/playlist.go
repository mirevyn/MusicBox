package v1

import (
	"errors"
	"musicbox-backend/internal/config"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PlaylistController 处理歌单相关的HTTP请求
type PlaylistController struct {
	playlistService *service.PlaylistService
}

// NewPlaylistController 创建一个新的PlaylistController实例
func NewPlaylistController() *PlaylistController {
	return &PlaylistController{
		playlistService: &service.PlaylistService{},
	}
}

// 从 Context 中安全获取 UserID
func getAuthUserID(ctx *gin.Context) (uint, bool) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.Result(ctx, 401, nil, "未经授权: 无法获取用户信息")
		return 0, false
	}
	uid, ok := userID.(uint)
	if !ok {
		utils.Result(ctx, 401, nil, "未经授权: 用户ID格式错误")
		return 0, false
	}
	return uid, true
}

// 可选解析当前登录用户
func getOptionalAuthUserID(ctx *gin.Context) uint {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return 0
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return 0
	}

	claims, err := utils.ParseToken(parts[1])
	if err != nil {
		return 0
	}

	return claims.UserID
}

// 解析 URL 中的 ID 参数
func parseParamID(ctx *gin.Context, key string) (uint, bool) {
	idStr := ctx.Param(key)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Result(ctx, 400, nil, "无效的ID参数")
		return 0, false
	}
	return uint(id), true
}

// @Summary 创建歌单
// @Description 创建一个新的个人歌单
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "歌单标题"
// @Param description formData string false "歌单描述"
// @Param isPublic formData bool false "是否公开" default(true)
// @Param coverFile formData file false "封面文件"
// @Success 200 {object} utils.Response{data=model.Playlist} "创建成功"
// @Failure 401 {object} utils.Response "未授权"
// @Router /playlists [post]
func (c *PlaylistController) CreatePlaylist(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}

	// 1. 获取文本参数
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")

	isPublic := true
	if isPublicStr, exists := ctx.GetPostForm("isPublic"); exists {
		parsed, err := strconv.ParseBool(isPublicStr)
		if err != nil {
			utils.Result(ctx, 400, nil, "isPublic 参数错误")
			return
		}
		isPublic = parsed
	}

	// 简单校验
	if title == "" {
		utils.Result(ctx, 400, nil, "标题不能为空")
		return
	}
	if len(title) > 100 {
		utils.Result(ctx, 400, nil, "标题长度不能超过100")
		return
	}

	// 处理封面图片上传 (可选)
	var coverUrl string
	file, err := ctx.FormFile("coverFile")
	if err == nil {
		coverUrl, err = service.SaveUploadedFile(ctx, file, config.Conf.Upload.CoverPath)
		if err != nil {
			resultUploadError(ctx, err, "封面保存失败")
			return
		}
	}

	// 调用 Service
	playlist, err := c.playlistService.CreatePlaylist(uid, title, description, coverUrl, isPublic)
	if err != nil {
		cleanupUploadedFiles(coverUrl)
		utils.Result(ctx, 500, nil, "创建歌单失败")
		return
	}

	utils.Result(ctx, 200, playlist, "歌单创建成功")
}

// @Summary 获取我的歌单
// @Description 获取当前登录用户创建的所有歌单
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=[]model.Playlist} "获取成功"
// @Router /playlists/my [get]
func (c *PlaylistController) GetMyPlaylists(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}

	playlists, err := c.playlistService.GetMyPlaylists(uid)
	if err != nil {
		utils.Result(ctx, 500, nil, "获取歌单列表失败")
		return
	}

	utils.Result(ctx, 200, playlists, "获取成功")
}

// @Summary 获取推荐歌单
// @Description 首页展示的热门公开歌单
// @Tags 歌单管理
// @Accept json
// @Produce json
// @Param pageIndex query int false "页码" default(1)
// @Param limit query int false "每页数量" default(4)
// @Success 200 {object} utils.Response "成功"
// @Router /playlists/recommended [get]
func (c *PlaylistController) GetRecommendedPlaylists(ctx *gin.Context) {
	pageIndex, err := strconv.Atoi(ctx.DefaultQuery("pageIndex", "1"))
	if err != nil || pageIndex < 1 {
		pageIndex = 1
	}

	pageSizeRaw := ctx.Query("pageSize")
	if pageSizeRaw == "" {
		pageSizeRaw = ctx.DefaultQuery("limit", "4")
	}
	pageSize, err := strconv.Atoi(pageSizeRaw)
	if err != nil || pageSize < 1 {
		pageSize = 4
	}

	playlists, total, err := c.playlistService.GetRecommendedPlaylists(pageIndex, pageSize)
	if err != nil {
		utils.Result(ctx, 500, nil, "获取推荐歌单失败")
		return
	}

	utils.Result(ctx, 200, gin.H{
		"playlists": playlists,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
		"hasMore":   int64(pageIndex*pageSize) < total,
	}, "获取成功")
}

// @Summary 搜索歌单
// @Description 根据关键字搜索所有公开歌单
// @Tags 歌单管理
// @Accept json
// @Produce json
// @Param keyword query string false "关键字"
// @Param pageIndex query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response "成功"
// @Router /playlists/search [get]
func (c *PlaylistController) SearchPlaylists(ctx *gin.Context) {
	keyword := ctx.Query("keyword")

	pageIndex, err := strconv.Atoi(ctx.DefaultQuery("pageIndex", "1"))
	if err != nil || pageIndex < 1 {
		pageIndex = 1
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	playlists, total, err := c.playlistService.SearchPlaylists(keyword, pageIndex, pageSize)
	if err != nil {
		utils.Result(ctx, 500, nil, "搜索歌单失败")
		return
	}

	utils.Result(ctx, 200, gin.H{
		"playlists": playlists,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
		"hasMore":   int64(pageIndex*pageSize) < total,
	}, "获取成功")
}

// @Summary 更新歌单
// @Description 修改歌单的基本信息或封面
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path uint true "歌单 ID"
// @Param title formData string false "歌单标题"
// @Param description formData string false "歌单描述"
// @Param isPublic formData bool false "是否公开"
// @Param coverFile formData file false "新封面文件"
// @Success 200 {object} utils.Response "更新成功"
// @Failure 403 {object} utils.Response "无权修改"
// @Failure 404 {object} utils.Response "歌单不存在"
// @Router /playlists/{id} [put]
func (c *PlaylistController) UpdatePlaylist(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}

	playlistID, ok := parseParamID(ctx, "id")
	if !ok {
		return
	}

	updateInput := service.PlaylistUpdateInput{ID: playlistID}

	if title, exists := ctx.GetPostForm("title"); exists {
		title = strings.TrimSpace(title)
		if title == "" {
			utils.Result(ctx, 400, nil, "标题不能为空")
			return
		}
		if len([]rune(title)) > 100 {
			utils.Result(ctx, 400, nil, "标题长度不能超过100")
			return
		}
		updateInput.Title = &title
	}

	if description, exists := ctx.GetPostForm("description"); exists {
		if len([]rune(description)) > 500 {
			utils.Result(ctx, 400, nil, "描述长度不能超过500")
			return
		}
		updateInput.Description = &description
	}

	if isPublicStr, exists := ctx.GetPostForm("isPublic"); exists {
		isPublic, err := strconv.ParseBool(isPublicStr)
		if err != nil {
			utils.Result(ctx, 400, nil, "isPublic 参数错误")
			return
		}
		updateInput.IsPublic = &isPublic
	}

	// 处理封面
	file, err := ctx.FormFile("coverFile")
	if err == nil {
		newCoverUrl, err := service.SaveUploadedFile(ctx, file, config.Conf.Upload.CoverPath)
		if err != nil {
			resultUploadError(ctx, err, "封面保存失败")
			return
		}
		updateInput.CoverURL = &newCoverUrl
	}

	if updateInput.Title == nil && updateInput.Description == nil && updateInput.IsPublic == nil && updateInput.CoverURL == nil {
		utils.Result(ctx, 400, nil, "未提供任何要更新的信息")
		return
	}

	err = c.playlistService.UpdatePlaylist(uid, updateInput)
	if err != nil {
		if updateInput.CoverURL != nil {
			cleanupUploadedFiles(*updateInput.CoverURL)
		}
		if errors.Is(err, service.ErrForbiddenAccess) {
			utils.Result(ctx, 403, nil, "您无权修改此歌单")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(ctx, 404, nil, "歌单不存在")
			return
		}
		utils.Result(ctx, 500, nil, "更新歌单失败")
		return
	}

	utils.Result(ctx, 200, nil, "歌单更新成功")
}

// @Summary 删除歌单
// @Description 删除指定的个人歌单
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Success 200 {object} utils.Response "删除成功"
// @Failure 403 {object} utils.Response "无权删除"
// @Router /playlists/{id} [delete]
func (c *PlaylistController) DeletePlaylist(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}

	playlistID, ok := parseParamID(ctx, "id")
	if !ok {
		return
	}

	err := c.playlistService.DeletePlaylist(playlistID, uid)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenAccess) {
			utils.Result(ctx, 403, nil, "您无权删除此歌单")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(ctx, 404, nil, "歌单不存在")
			return
		}
		utils.Result(ctx, 500, nil, "删除歌单失败")
		return
	}

	utils.Result(ctx, 200, nil, "歌单删除成功")
}

type AddSongReq struct {
	SongID uint `json:"songId" binding:"required"`
}

// @Summary 向歌单添加歌曲
// @Description 将指定歌曲加入到你的歌单中
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Param request body AddSongReq true "歌曲信息"
// @Success 200 {object} utils.Response "添加成功"
// @Failure 400 {object} utils.Response "参数错误或已存在"
// @Router /playlists/{id}/songs [post]
func (c *PlaylistController) AddSongToPlaylist(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}
	playlistID, ok := parseParamID(ctx, "id")
	if !ok {
		return
	}

	var req AddSongReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Result(ctx, 400, nil, "参数错误")
		return
	}

	err := c.playlistService.AddSongToPlaylist(uid, playlistID, req.SongID)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenAccess) {
			utils.Result(ctx, 403, nil, "无权操作")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(ctx, 404, nil, "歌单或歌曲不存在")
			return
		}
		if errors.Is(err, service.ErrDuplicateEntry) {
			utils.Result(ctx, 409, nil, "歌曲已存在于歌单中")
			return
		}
		utils.Result(ctx, 500, nil, "添加失败: "+err.Error())
		return
	}

	utils.Result(ctx, 200, nil, "添加成功")
}

// @Summary 从歌单移除歌曲
// @Description 将歌曲从指定歌单中移除
// @Tags 歌单管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Param songId path uint true "歌曲 ID"
// @Success 200 {object} utils.Response "移除成功"
// @Failure 403 {object} utils.Response "无权操作"
// @Router /playlists/{id}/songs/{songId} [delete]
func (c *PlaylistController) RemoveSongFromPlaylist(ctx *gin.Context) {
	uid, ok := getAuthUserID(ctx)
	if !ok {
		return
	}
	playlistID, ok := parseParamID(ctx, "id")
	if !ok {
		return
	}
	songID, ok := parseParamID(ctx, "songId")
	if !ok {
		return
	}

	err := c.playlistService.RemoveSongFromPlaylist(uid, playlistID, songID)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenAccess) {
			utils.Result(ctx, 403, nil, "无权操作")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(ctx, 404, nil, "资源不存在")
			return
		}
		utils.Result(ctx, 500, nil, "移除失败")
		return
	}
	utils.Result(ctx, 200, nil, "移除成功")
}

// @Summary 获取歌单详情
// @Description 获取歌单信息及其包含的所有歌曲列表
// @Tags 歌单管理
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Success 200 {object} utils.Response "成功"
// @Failure 403 {object} utils.Response "私密歌单无权访问"
// @Router /playlists/{id} [get]
func (c *PlaylistController) GetPlaylistDetails(ctx *gin.Context) {
	uid := getOptionalAuthUserID(ctx)
	playlistID, ok := parseParamID(ctx, "id")
	if !ok {
		return
	}

	playlist, songs, err := c.playlistService.GetPlaylistDetails(playlistID, uid)
	if err != nil {
		if errors.Is(err, service.ErrForbiddenAccess) {
			utils.Result(ctx, 403, nil, "私密歌单，无权访问")
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(ctx, 404, nil, "歌单不存在")
			return
		}
		utils.Result(ctx, 500, nil, "获取失败")
		return
	}

	utils.Result(ctx, 200, gin.H{
		"playlist": playlist,
		"songs":    songs,
	}, "获取成功")
}
