package v1

import (
	"errors"
	"musicbox-backend/internal/model"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest 定义管理员创建用户的请求体
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // 非必填, 默认为 "User"
}

// UpdateUserRequest 定义管理员更新用户的请求体
type UpdateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Status   *int   `json:"status"`
}

// @Summary 获取仪表盘统计
// @Description 获取系统总用户数、音频数、播放历史等统计数据
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "成功"
// @Router /admin/dashboard/stats [get]
func GetDashboardStats(c *gin.Context) {
	stats, err := service.GetDashboardStats()
	if err != nil {
		utils.Result(c, 500, nil, "获取仪表盘数据失败")
		return
	}

	utils.Result(c, 200, stats, "成功")
}

// @Summary 获取分析数据
// @Description 获取用于图表展示的分析趋势数据
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.Response "成功"
// @Router /admin/dashboard/analytics [get]
func GetDashboardAnalytics(c *gin.Context) {
	analytics, err := service.GetDashboardAnalytics()
	if err != nil {
		utils.Result(c, 500, nil, "获取分析数据失败")
		return
	}

	utils.Result(c, 200, analytics, "成功")
}

// @Summary 导出仪表盘报告
// @Description 下载包含系统统计信息的 CSV 报告
// @Tags 管理员模块
// @Security ApiKeyAuth
// @Produce text/csv
// @Success 200 {file} file "CSV 报告文件"
// @Router /admin/dashboard/export [get]
func ExportDashboardReport(c *gin.Context) {
	report, err := service.BuildDashboardReportCSV()
	if err != nil {
		utils.Result(c, 500, nil, "导出报告失败")
		return
	}

	filename := "dashboard-report-" + time.Now().Format("20060102-150405") + ".csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(200, "text/csv; charset=utf-8", report)
}

// @Summary 获取用户列表
// @Description 管理员分页查询系统用户
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param keyword query string false "用户名关键字"
// @Param role query string false "角色过滤"
// @Param pageIndex query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} utils.Response "成功"
// @Router /admin/users [get]
func GetUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	role := c.Query("role")

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

	users, total, err := service.GetUsers(keyword, role, pageIndex, pageSize)
	if err != nil {
		utils.Result(c, 500, nil, "获取用户失败")
		return
	}

	utils.Result(c, 200, gin.H{
		"users":     users,
		"total":     total,
		"pageIndex": pageIndex,
		"pageSize":  pageSize,
	}, "成功")
}

// @Summary 获取用户详情
// @Description 根据 ID 获取指定用户的详细信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "用户 ID"
// @Success 200 {object} utils.Response{data=model.User} "成功"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /admin/users/{id} [get]
func GetUserByID(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的用户ID")
		return
	}

	user, err := service.GetUserByID(uint(userID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			utils.Result(c, 404, nil, "用户未找到")
		} else {
			utils.Result(c, 500, nil, "获取失败")
		}
		return
	}

	utils.Result(c, 200, user, "成功")
}

// @Summary 创建用户
// @Description 管理员手动创建新用户账号
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "创建用户请求体"
// @Success 201 {object} utils.Response{data=model.User} "成功"
// @Failure 409 {object} utils.Response "用户名已存在"
// @Router /admin/users [post]
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误: "+err.Error())
		return
	}

	user := &model.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Role:         req.Role,
		Status:       1, // 默认正常状态
	}

	if err := service.Register(user); err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			utils.Result(c, 409, nil, "用户名已存在")
		} else {
			utils.Result(c, 500, nil, "创建用户失败: "+err.Error())
		}
		return
	}

	utils.Result(c, 201, user, "用户创建成功")
}

// @Summary 更新用户
// @Description 管理员修改用户的基本信息、角色或状态
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "用户 ID"
// @Param request body UpdateUserRequest true "更新用户请求体"
// @Success 200 {object} utils.Response{data=model.User} "成功"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的用户ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Result(c, 400, nil, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		utils.Result(c, 400, nil, "未提供任何要更新的信息")
		return
	}

	actorUserID, ok := getAdminActorID(c)
	if !ok {
		return
	}

	updatedUser, err := service.UpdateUserByAdmin(actorUserID, uint(userID), updates)
	if err != nil {
		if errors.Is(err, service.ErrUsernameExists) {
			utils.Result(c, 409, nil, "用户名已存在")
		} else if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrUserNotFound) {
			utils.Result(c, 404, nil, "用户不存在")
		} else if errors.Is(err, service.ErrProtectedAdmin) || errors.Is(err, service.ErrSelfOperation) {
			utils.Result(c, 403, nil, err.Error())
		} else {
			utils.Result(c, 500, nil, "更新用户信息失败: "+err.Error())
		}
		return
	}

	utils.Result(c, 200, updatedUser, "用户信息更新成功")
}

// @Summary 删除用户
// @Description 管理员彻底删除系统中的用户账号
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path uint true "用户 ID"
// @Success 200 {object} utils.Response "成功"
// @Failure 404 {object} utils.Response "用户不存在"
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		utils.Result(c, 400, nil, "无效的用户ID")
		return
	}

	actorUserID, ok := getAdminActorID(c)
	if !ok {
		return
	}

	err = service.DeleteUser(actorUserID, uint(userID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrUserNotFound) {
			utils.Result(c, 404, nil, "用户未找到")
			return
		}
		if errors.Is(err, service.ErrProtectedAdmin) || errors.Is(err, service.ErrSelfOperation) {
			utils.Result(c, 403, nil, err.Error())
			return
		}
		utils.Result(c, 500, nil, "删除用户失败: "+err.Error())
		return
	}

	utils.Result(c, 200, nil, "用户删除成功")
}

func getAdminActorID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("userID")
	if !exists {
		utils.Result(c, 401, nil, "用户未登录")
		return 0, false
	}
	userID, ok := val.(uint)
	if !ok || userID == 0 {
		utils.Result(c, 401, nil, "用户身份无效")
		return 0, false
	}
	return userID, true
}
