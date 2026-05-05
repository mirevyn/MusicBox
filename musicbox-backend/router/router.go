package router

import (
	v1 "musicbox-backend/api/v1"
	"musicbox-backend/docs"
	"musicbox-backend/middleware"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var swaggerInfoMu sync.Mutex

func InitRouter() *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.MaxMultipartMemory = 16 << 20

	r.Use(middleware.CORS())
	r.Static("/uploads", "./uploads")
	r.GET("/healthz", v1.Healthz)

	// Swagger 路由
	r.GET("/swagger/*any", swaggerHandler())

	api := r.Group("/api")
	{

		// 认证模块
		auth := api.Group("/auth")
		{
			auth.POST("/login", v1.Login)
			auth.POST("/register", v1.Register)
		}

		// 系统配置 (公开)
		systemCtrl := v1.NewSystemController()
		api.GET("/config", systemCtrl.GetPublicConfig)

		// 歌曲模块
		songs := api.Group("/songs")
		{
			songs.GET("", v1.GetSongs)
			songs.GET("/:id", v1.GetSongByID)
			protected := songs.Group("")
			protected.Use(middleware.JWTAuth(), middleware.AdminAuth())
			{
				protected.POST("/upload", v1.UploadSong)
				protected.POST("/export", v1.ExportSongs)
				protected.PUT("/:id", v1.UpdateSong)
				protected.DELETE("/:id", v1.DeleteSong)
			}
		}

		// 用户个人中心
		user := api.Group("/user")
		user.Use(middleware.JWTAuth())
		{
			user.GET("/profile", v1.GetMyProfile)
			user.PUT("/password", v1.ChangePassword)
			user.POST("/avatar", v1.UploadAvatar)
		}

		// 收藏与点赞
		songLikesCtrl := v1.NewSongLikesController()
		songLikes := api.Group("/song-likes")
		songLikes.Use(middleware.JWTAuth())
		{
			songLikes.GET("", songLikesCtrl.GetMyLikedSongs)
			songLikes.POST("/:songId", songLikesCtrl.ToggleSongLike)
			songLikes.GET("/:songId", songLikesCtrl.GetLikeStatus)
		}

		// 私享推荐模块
		recommendationCtrl := v1.NewRecommendationController()
		recommendations := api.Group("/recommendations")
		recommendations.Use(middleware.JWTAuth())
		{
			recommendations.GET("/daily", recommendationCtrl.GetDailyRecommendations)
		}

		// 播放历史模块
		playHistories := api.Group("/play-histories")
		playHistories.Use(middleware.JWTAuth())
		{
			playHistories.POST("", recommendationCtrl.RecordPlayHistory)
		}

		// AI 助手模块
		aiCtrl := v1.NewAIController()
		ai := api.Group("/ai")
		ai.Use(middleware.JWTAuth())
		{
			ai.GET("/status", aiCtrl.GetStatus)
			ai.POST("/status", aiCtrl.GetStatus)
			ai.POST("/chat/stream", aiCtrl.ChatStream)
		}

		// 管理员模块
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
		{
			admin.GET("/dashboard/stats", v1.GetDashboardStats)
			admin.GET("/dashboard/analytics", v1.GetDashboardAnalytics)
			admin.GET("/dashboard/export", v1.ExportDashboardReport)
			admin.GET("/users", v1.GetUsers)
			admin.GET("/users/:id", v1.GetUserByID)
			admin.POST("/users", v1.CreateUser)
			admin.PUT("/users/:id", v1.UpdateUser)
			admin.DELETE("/users/:id", v1.DeleteUser)

			// 歌单管理
			admin.GET("/playlists", v1.GetAdminPlaylists)
			admin.GET("/playlists/:id", v1.GetAdminPlaylistDetails)
			admin.PUT("/playlists/:id/status", v1.UpdatePlaylistStatus)
			admin.DELETE("/playlists/:id", v1.DeleteAdminPlaylist)

			// 管理系统设置
			admin.GET("/settings", systemCtrl.GetSettings)
			admin.PUT("/settings", systemCtrl.UpdateSettings)
		}

		// 歌单模块
		playlistCtrl := v1.NewPlaylistController()
		api.GET("/playlists/recommended", playlistCtrl.GetRecommendedPlaylists)
		api.GET("/playlists/search", playlistCtrl.SearchPlaylists)
		api.GET("/playlists/:id", playlistCtrl.GetPlaylistDetails)

		playlists := api.Group("/playlists")
		playlists.Use(middleware.JWTAuth())
		{
			playlists.POST("", playlistCtrl.CreatePlaylist)
			playlists.GET("/my", playlistCtrl.GetMyPlaylists)
			playlists.PUT("/:id", playlistCtrl.UpdatePlaylist)
			playlists.DELETE("/:id", playlistCtrl.DeletePlaylist)
			playlists.POST("/:id/songs", playlistCtrl.AddSongToPlaylist)
			playlists.DELETE("/:id/songs/:songId", playlistCtrl.RemoveSongFromPlaylist)
		}
	}

	// WebSocket 通知模块
	notificationCtrl := v1.NewNotificationController()
	ws := r.Group("/ws")
	{
		ws.GET("/admin/notifications", notificationCtrl.AdminNotifications)
	}

	return r
}

func swaggerHandler() gin.HandlerFunc {
	handler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	return func(c *gin.Context) {
		swaggerInfoMu.Lock()
		defer swaggerInfoMu.Unlock()

		docs.SwaggerInfo.Host = c.Request.Host
		docs.SwaggerInfo.Schemes = []string{requestScheme(c.Request)}
		handler(c)
	}
}

func requestScheme(r *http.Request) string {
	if proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); proto != "" {
		if comma := strings.Index(proto, ","); comma >= 0 {
			proto = proto[:comma]
		}
		proto = strings.ToLower(strings.TrimSpace(proto))
		if proto == "http" || proto == "https" {
			return proto
		}
	}

	if r.TLS != nil {
		return "https"
	}
	return "http"
}
