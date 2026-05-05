package v1

import (
	"errors"
	"musicbox-backend/internal/config"
	"musicbox-backend/internal/model"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginRequest 定义登录请求参数结构体
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 定义注册请求参数结构体
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

func formatRegisterBindError(err error) string {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "RegisterRequest.Username") && strings.Contains(msg, "required"):
		return "用户名不能为空"
	case strings.Contains(msg, "RegisterRequest.Username") && strings.Contains(msg, "min"):
		return "用户名至少需要 3 个字符"
	case strings.Contains(msg, "RegisterRequest.Username") && strings.Contains(msg, "max"):
		return "用户名长度不能超过 30 个字符"
	case strings.Contains(msg, "RegisterRequest.Password") && strings.Contains(msg, "required"):
		return "密码不能为空"
	case strings.Contains(msg, "RegisterRequest.Password") && strings.Contains(msg, "min"):
		return "密码长度至少需要 6 位"
	default:
		return "注册参数错误"
	}
}

// ChangePasswordRequest 定义修改密码请求参数
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// @Summary 修改密码
// @Description 允许登录用户修改其登录密码
// @Tags 用户接口
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "修改密码请求参数"
// @Success 200 {object} utils.Response "密码修改成功"
// @Failure 401 {object} utils.Response "未登录或旧密码错误"
// @Router /user/password [put]
func ChangePassword(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return
	}
	userID := val.(uint)

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误: "+err.Error())
		return
	}

	err := service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrPasswordIncorrect) {
			utils.Result(c, 401, nil, "旧密码错误")
		} else {
			utils.Result(c, 500, nil, "修改失败: "+err.Error())
		}
		return
	}

	utils.Result(c, 200, nil, "密码修改成功")
}

// @Summary 上传头像
// @Description 用户上传自定义头像文件
// @Tags 用户接口
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "头像图片文件"
// @Success 200 {object} utils.Response "头像更新成功"
// @Failure 400 {object} utils.Response "未提供文件"
// @Router /user/avatar [post]
func UploadAvatar(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return
	}
	userID := val.(uint)

	file, err := c.FormFile("file")
	if err != nil {
		utils.Result(c, 400, nil, "未提供文件")
		return
	}

	// 限制文件大小
	if file.Size > 2*1024*1024 {
		utils.Result(c, 400, nil, "图片大小不能超过2MB")
		return
	}

	newAvatarURL, err := service.SaveUploadedFile(c, file, config.Conf.Upload.AvatarPath)
	if err != nil {
		resultUploadError(c, err, "保存头像文件失败")
		return
	}

	oldAvatarURL, err := service.UpdateAvatar(userID, newAvatarURL)
	if err != nil {
		cleanupUploadedFiles(newAvatarURL)
		utils.Result(c, 500, nil, "更新头像信息失败: "+err.Error())
		return
	}

	if oldAvatarURL != "" && oldAvatarURL != "avatar/default.jpg" {
		_ = utils.DeleteFile(oldAvatarURL)
	}

	utils.Result(c, 200, gin.H{"avatarUrl": newAvatarURL}, "头像更新成功")
}

// @Summary 用户注册
// @Description 允许新用户创建账号（如果系统配置允许）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} utils.Response "注册成功"
// @Failure 400 {object} utils.Response "参数错误"
// @Failure 403 {object} utils.Response "注册已关闭"
// @Failure 409 {object} utils.Response "用户名已存在"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	// 拦截逻辑：检查是否开启注册
	if !service.GetSystemService().IsRegisterAllowed() {
		utils.Result(c, 403, nil, "当前系统已关闭新用户注册，请联系管理员")
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, formatRegisterBindError(err))
		return
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Role:         "User",
		Status:       1,
	}

	if err := service.Register(&user); err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			utils.Result(c, 409, nil, "用户名已存在")
		} else {
			utils.Result(c, 500, nil, "注册失败: "+err.Error())
		}
		return
	}

	utils.Result(c, 201, nil, "注册成功")
}

// @Summary 用户登录
// @Description 验证用户名和密码并返回 JWT Token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} utils.Response "登录成功"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误: "+err.Error())
		return
	}

	user, token, expiration, err := service.Login(req.Username, req.Password)
	if err != nil {
		// 如果是账号禁用错误，明确返回
		if err.Error() == "账号已被禁用" {
			utils.Result(c, 403, nil, "账号已被禁用")
			return
		}
		// 统一登录失败的错误提示
		utils.Result(c, 401, nil, "用户名或密码错误")
		return
	}

	utils.Result(c, 200, gin.H{
		"token":      token,
		"expiration": expiration,
		"user":       user,
	}, "登录成功")
}

// @Summary 获取个人信息
// @Description 获取当前登录用户的详细档案
// @Tags 用户接口
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response{data=model.User} "获取成功"
// @Router /user/profile [get]
func GetMyProfile(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "未登录")
		return
	}

	userID := val.(uint)
	user, err := service.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(c, 404, nil, "用户不存在")
		} else {
			utils.Result(c, 500, nil, "获取用户信息失败")
		}
		return
	}

	utils.Result(c, 200, user, "success")
}
