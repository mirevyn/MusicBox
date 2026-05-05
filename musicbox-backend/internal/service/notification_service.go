package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"musicbox-backend/internal/global"

	"github.com/gorilla/websocket"
)

type AdminNotification struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	TargetID  uint   `json:"targetId,omitempty"`
	Route     string `json:"route,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type NotificationClient struct {
	adminID uint
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func NewNotificationClient(adminID uint, conn *websocket.Conn) *NotificationClient {
	return &NotificationClient{
		adminID: adminID,
		conn:    conn,
	}
}

func (c *NotificationClient) WriteMessage(messageType int, payload []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteMessage(messageType, payload)
}

func (c *NotificationClient) Close() error {
	return c.conn.Close()
}

type NotificationService struct {
	mu      sync.RWMutex
	clients map[uint]map[*NotificationClient]struct{}
}

var notificationHub = &NotificationService{
	clients: make(map[uint]map[*NotificationClient]struct{}),
}

func GetNotificationService() *NotificationService {
	return notificationHub
}

func (s *NotificationService) RegisterAdmin(adminID uint, conn *websocket.Conn) *NotificationClient {
	client := NewNotificationClient(adminID, conn)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.clients[adminID] == nil {
		s.clients[adminID] = make(map[*NotificationClient]struct{})
	}
	s.clients[adminID][client] = struct{}{}

	return client
}

func (s *NotificationService) UnregisterAdmin(adminID uint, client *NotificationClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	connections, exists := s.clients[adminID]
	if !exists {
		return
	}

	delete(connections, client)
	if len(connections) == 0 {
		delete(s.clients, adminID)
	}
}

func (s *NotificationService) BroadcastToAdmins(notification AdminNotification) {
	payload, err := json.Marshal(notification)
	if err != nil {
		return
	}

	for _, client := range s.adminClients() {
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			s.closeClient(client)
		}
	}
}

func (s *NotificationService) adminClients() []*NotificationClient {
	s.mu.RLock()
	defer s.mu.RUnlock()

	clients := make([]*NotificationClient, 0)
	for _, connections := range s.clients {
		for client := range connections {
			clients = append(clients, client)
		}
	}
	return clients
}

func (s *NotificationService) closeClient(client *NotificationClient) {
	s.UnregisterAdmin(client.adminID, client)
	_ = client.Close()
}

func (s *NotificationService) PushPlaylistPendingNotification(title string, playlistID uint) {
	s.BroadcastToAdmins(AdminNotification{
		ID:        time.Now().Format("20060102150405.000000000"),
		Type:      "playlist_pending_review",
		Title:     "有新的歌单待审核",
		Content:   "歌单《" + title + "》已进入待审核队列",
		TargetID:  playlistID,
		Route:     "/admin/playlists",
		CreatedAt: time.Now().Format(time.RFC3339),
	})
}

// PushDashboardRefresh 通知所有在线管理员刷新仪表盘数据。
// 该消息为轻量信号，前端收到后自行重新拉取 REST API。
func (s *NotificationService) PushDashboardRefresh() {
	s.BroadcastToAdmins(AdminNotification{
		ID:        "dashboard_" + time.Now().Format("20060102150405.000000000"),
		Type:      "dashboard_refresh",
		Title:     "数据已更新",
		Content:   "仪表盘数据已发生变化",
		CreatedAt: time.Now().Format(time.RFC3339),
	})
}

// PushPendingPlaylistsToClient 在管理员 WS 连接建立时，将数据库中所有待审核歌单
// 逐条推送给该连接，确保管理员不会遗漏离线期间产生的待审核通知。
func (s *NotificationService) PushPendingPlaylistsToClient(client *NotificationClient) {
	var pending []struct {
		ID    uint   `gorm:"column:id"`
		Title string `gorm:"column:title"`
	}

	if err := global.DB.Table("playlists").
		Select("id, title").
		Where("is_public = ? AND status = ?", true, 0).
		Order("created_at DESC").
		Limit(20).
		Scan(&pending).Error; err != nil {
		return
	}

	for _, p := range pending {
		notification := AdminNotification{
			ID:        fmt.Sprintf("pending_%d", p.ID),
			Type:      "playlist_pending_review",
			Title:     "待审核歌单",
			Content:   "歌单《" + p.Title + "》需要审核",
			TargetID:  p.ID,
			Route:     "/admin/playlists",
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		payload, err := json.Marshal(notification)
		if err != nil {
			continue
		}
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			s.closeClient(client)
			return
		}
	}
}
