package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz 返回轻量健康探针结果，供容器编排和反向代理检查使用。
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
