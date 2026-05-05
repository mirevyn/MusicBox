package v1

import (
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
)

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

// @Summary 获取系统设置
// @Description 管理员获取系统全局配置（如是否开放注册）
// @Tags 系统模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "获取成功"
// @Router /admin/settings [get]
func (sc *SystemController) GetSettings(c *gin.Context) {
	allowRegister, _ := service.GetSystemService().GetSetting("allow_register")
	utils.Result(c, 200, gin.H{
		"allowRegister": allowRegister == "true",
	}, "获取成功")
}

// UpdateSettingsRequest 定义更新系统设置的请求体
type UpdateSettingsRequest struct {
	AllowRegister bool `json:"allowRegister"`
}

// @Summary 更新系统设置
// @Description 管理员修改系统全局配置
// @Tags 系统模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body UpdateSettingsRequest true "设置信息"
// @Success 200 {object} utils.Response "保存成功"
// @Router /admin/settings [put]
func (sc *SystemController) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误")
		return
	}

	val := "false"
	if req.AllowRegister {
		val = "true"
	}

	err := service.GetSystemService().UpdateSetting("allow_register", val)
	if err != nil {
		utils.Result(c, 500, nil, "保存失败")
		return
	}

	utils.Result(c, 200, nil, "保存成功")
}

// @Summary 获取公开配置
// @Description 获取访客可见的系统公开配置
// @Tags 系统模块
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "获取成功"
// @Router /config [get]
func (sc *SystemController) GetPublicConfig(c *gin.Context) {
	utils.Result(c, 200, gin.H{
		"allowRegister": service.GetSystemService().IsRegisterAllowed(),
	}, "获取成功")
}
