package middleware

import (
	"errors"
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth 是一个中间件，用于验证请求头中的 JWT Token。
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Header 中的 Authorization
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			utils.Result(c, 401, nil, "请求未携带Token，无权访问")
			c.Abort()
			return
		}

		// 按惯例 Token 格式为 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Result(c, 401, nil, "Token格式错误")
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Result(c, 401, nil, "Token无效或已过期")
			c.Abort()
			return
		}

		userID := claims.UserID
		username := claims.Username
		role := claims.Role
		if global.DB != nil {
			user, err := service.GetActiveUserForClaims(claims)
			if err != nil {
				status := 401
				message := "Token无效或已过期"
				if errors.Is(err, service.ErrUserDisabled) {
					status = 403
					message = "账号已被禁用"
				} else if !errors.Is(err, service.ErrNotFound) {
					status = 500
					message = "认证状态校验失败"
				}
				utils.Result(c, status, nil, message)
				c.Abort()
				return
			}
			userID = user.ID
			username = user.Username
			role = user.Role
		}

		// 将当前用户信息存在上下文，供后续 Controller 使用。
		c.Set("userID", userID)
		c.Set("username", username)
		c.Set("role", role)

		c.Next()
	}
}
