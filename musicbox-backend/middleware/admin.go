package middleware

import (
	"musicbox-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuth 是一个中间件，用于检查用户是否具有 'Admin' 角色
// 在 JWTAuth 中间件之后使用
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取角色 (由 JWTAuth 中间件设置)
		roleVal, exists := c.Get("role")
		if !exists {
			utils.Result(c, 403, nil, "无法获取用户角色信息")
			c.Abort() // 必须终止请求
			return
		}

		// 检查角色是否为 'Admin'
		role, ok := roleVal.(string)
		if !ok || !strings.EqualFold(role, "Admin") {
			utils.Result(c, 403, nil, "权限不足，需要管理员身份")
			c.Abort() // 必须终止请求
			return
		}

		// 用户是管理员，继续执行下一个处理程序
		c.Next()
	}
}
