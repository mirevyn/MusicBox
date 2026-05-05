package main

// @title MusicBox API
// @version 1.0
// @description MusicBox 音乐盒后端 API 文档

// @BasePath /api
// @Schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description JWT Token，格式：Bearer <token>

import (
	"context"
	"log"
	"musicbox-backend/internal/config"
	"musicbox-backend/internal/global"
	"musicbox-backend/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	if config.Conf.Server.GinMode != "" {
		gin.SetMode(config.Conf.Server.GinMode)
	}

	// 初始化数据库
	global.InitDB()

	// 初始化路由
	r := router.InitRouter()

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:              config.Conf.Server.HTTPAddr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 启动服务 (非阻塞)
	go func() {
		log.Printf("服务启动在 %s (GIN_MODE=%s)", config.Conf.Server.HTTPAddr, gin.Mode())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅关闭：监听中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	// 5秒超时关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭异常: %v", err)
	}
	log.Println("服务已安全关闭")
}
