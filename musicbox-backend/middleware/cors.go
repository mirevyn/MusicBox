package middleware

import (
	"musicbox-backend/internal/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 处理跨域请求的中间件
func CORS() gin.HandlerFunc {
	allowedOrigins := config.Conf.CORS.AllowedOrigins
	if len(allowedOrigins) == 0 {
		// 同源部署下无需显式挂载 CORS 中间件。
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,

		// 允许的请求方法
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},

		// 允许的请求头
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Range", //允许 Range 头，支持音频拖动
		},

		// 暴露给前端的响应头
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Content-Range", //暴露内容范围
			"Accept-Ranges", //暴露支持范围请求
		},

		// 允许携带凭证
		AllowCredentials: true,

		// 预检请求缓存时间
		MaxAge: 12 * time.Hour,
	})
}
