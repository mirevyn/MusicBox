package v1

import (
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 管理员获取歌单列表
// @Description 管理员分页查询系统中的所有歌单
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param keyword query string false "关键字"
// @Param status query int false "状态 (0-待审核, 1-已通过, 2-已驳回)" default(-1)
// @Param pageIndex query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response "成功"
// @Router /admin/playlists [get]
func GetAdminPlaylists(c *gin.Context) {
	keyword := c.Query("keyword")
	statusStr := c.DefaultQuery("status", "-1")

	pageIndex, _ := strconv.Atoi(c.DefaultQuery("pageIndex", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status, _ := strconv.Atoi(statusStr)

	playlists, total, err := service.GetAdminPlaylists(keyword, status, pageIndex, pageSize)
	if err != nil {
		utils.Result(c, 500, nil, "获取歌单失败: "+err.Error())
		return
	}

	utils.Result(c, 200, gin.H{
		"playlists": playlists,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	}, "成功")
}

// UpdatePlaylistStatusReq 更新状态请求体
type UpdatePlaylistStatusReq struct {
	Status       int    `json:"status" binding:"oneof=1 2"` // 1-通过, 2-驳回
	RejectReason string `json:"rejectReason"`
}

// @Summary 审核歌单
// @Description 管理员对歌单进行审核通过或驳回操作
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Param request body UpdatePlaylistStatusReq true "审核信息"
// @Success 200 {object} utils.Response "操作成功"
// @Router /admin/playlists/{id}/status [put]
func UpdatePlaylistStatus(c *gin.Context) {
	idParam := c.Param("id")
	playlistID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的ID")
		return
	}

	var req UpdatePlaylistStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误: "+err.Error())
		return
	}

	if req.Status == 2 && req.RejectReason == "" {
		utils.Result(c, 400, nil, "驳回时必须填写理由")
		return
	}

	err = service.UpdatePlaylistStatus(uint(playlistID), req.Status, req.RejectReason)
	if err != nil {
		utils.Result(c, 500, nil, "操作失败: "+err.Error())
		return
	}

	utils.Result(c, 200, nil, "操作成功")
}

// @Summary 管理员删除歌单
// @Description 管理员强制删除指定的歌单
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Success 200 {object} utils.Response "删除成功"
// @Router /admin/playlists/{id} [delete]
func DeleteAdminPlaylist(c *gin.Context) {
	idParam := c.Param("id")
	playlistID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的ID")
		return
	}

	err = service.DeletePlaylistByAdmin(uint(playlistID))
	if err != nil {
		utils.Result(c, 500, nil, "删除失败: "+err.Error())
		return
	}

	utils.Result(c, 200, nil, "删除成功")
}

// @Summary 管理员获取歌单详情
// @Description 管理员查看歌单的详细信息及包含的歌曲
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "歌单 ID"
// @Success 200 {object} utils.Response "获取成功"
// @Router /admin/playlists/{id} [get]
func GetAdminPlaylistDetails(c *gin.Context) {
	idParam := c.Param("id")
	playlistID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的ID")
		return
	}

	playlist, songs, err := service.GetAdminPlaylistDetails(uint(playlistID))
	if err != nil {
		utils.Result(c, 500, nil, "获取详情失败: "+err.Error())
		return
	}

	utils.Result(c, 200, gin.H{
		"playlist": playlist,
		"songs":    songs,
	}, "获取成功")
}
