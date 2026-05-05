package v1

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"musicbox-backend/internal/config"
	"musicbox-backend/internal/global"
	"musicbox-backend/internal/service"
	"musicbox-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return isAllowedWebSocketOrigin(r)
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"musicbox.jwt"},
}

type NotificationController struct{}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

// @Summary WebSocket 通知
// @Description 管理员通过 WebSocket 接收系统实时通知（如新歌单审核通知）
// @Tags 系统模块
// @Param Sec-WebSocket-Protocol header string true "musicbox.jwt,<JWT 访问令牌>"
// @Success 101 {string} string "Switching Protocols"
// @Router /ws/admin/notifications [get]
func (ctrl *NotificationController) AdminNotifications(c *gin.Context) {
	token := websocketToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "缺少 token"})
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Token 无效或已过期"})
		return
	}

	role := claims.Role
	if global.DB != nil {
		user, err := service.GetActiveUserForClaims(claims)
		if err != nil {
			status := http.StatusUnauthorized
			message := "Token 无效或已过期"
			if errors.Is(err, service.ErrUserDisabled) {
				status = http.StatusForbidden
				message = "账号已被禁用"
			} else if !errors.Is(err, service.ErrNotFound) {
				status = http.StatusInternalServerError
				message = "认证状态校验失败"
			}
			c.JSON(status, gin.H{"code": status, "msg": message})
			return
		}
		role = user.Role
	}
	if strings.ToLower(role) != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "无管理员权限"})
		return
	}

	adminID := claims.UserID
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	hub := service.GetNotificationService()
	client := hub.RegisterAdmin(adminID, ws)
	hub.PushPendingPlaylistsToClient(client)
	defer func() {
		hub.UnregisterAdmin(adminID, client)
		_ = client.Close()
	}()

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func websocketToken(c *gin.Context) string {
	if authHeader := c.GetHeader("Authorization"); authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	for _, part := range strings.Split(c.GetHeader("Sec-WebSocket-Protocol"), ",") {
		part = strings.TrimSpace(part)
		if part == "" || strings.EqualFold(part, "musicbox.jwt") {
			continue
		}
		return part
	}

	return ""
}

func isAllowedWebSocketOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}

	if config.Conf != nil {
		for _, allowed := range config.Conf.CORS.AllowedOrigins {
			if strings.EqualFold(strings.TrimRight(allowed, "/"), strings.TrimRight(origin, "/")) {
				return true
			}
		}
	}

	originURL, err := url.Parse(origin)
	if err != nil || originURL.Host == "" {
		return false
	}
	return strings.EqualFold(originURL.Host, requestHost(r))
}

func requestHost(r *http.Request) string {
	if host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); host != "" {
		if comma := strings.Index(host, ","); comma >= 0 {
			host = host[:comma]
		}
		return strings.TrimSpace(host)
	}
	return r.Host
}
