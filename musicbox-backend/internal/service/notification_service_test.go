package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestBroadcastToAdminsSerializesConcurrentWrites(t *testing.T) {
	serverConnCh := make(chan *websocket.Conn, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("websocket upgrade failed: %v", err)
			return
		}
		serverConnCh <- conn

		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	readerConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("websocket dial failed: %v", err)
	}
	defer readerConn.Close()

	var serverConn *websocket.Conn
	select {
	case serverConn = <-serverConnCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for server websocket connection")
	}

	notifications := &NotificationService{
		clients: make(map[uint]map[*NotificationClient]struct{}),
	}
	client := notifications.RegisterAdmin(1, serverConn)
	defer func() {
		notifications.UnregisterAdmin(1, client)
		_ = client.Close()
	}()

	const broadcastCount = 64
	readErrCh := make(chan error, 1)
	go func() {
		for i := 0; i < broadcastCount; i++ {
			if err := readerConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
				readErrCh <- err
				return
			}

			_, payload, err := readerConn.ReadMessage()
			if err != nil {
				readErrCh <- fmt.Errorf("read broadcast %d: %w", i, err)
				return
			}

			var notification AdminNotification
			if err := json.Unmarshal(payload, &notification); err != nil {
				readErrCh <- fmt.Errorf("decode broadcast %d: %w", i, err)
				return
			}
		}
		readErrCh <- nil
	}()

	var wg sync.WaitGroup
	panicCh := make(chan any, broadcastCount)
	for i := 0; i < broadcastCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() {
				if recovered := recover(); recovered != nil {
					panicCh <- recovered
				}
			}()

			notifications.BroadcastToAdmins(AdminNotification{
				ID:        fmt.Sprintf("broadcast_%d", i),
				Type:      "dashboard_refresh",
				Title:     "数据已更新",
				Content:   "仪表盘数据已发生变化",
				CreatedAt: time.Now().Format(time.RFC3339),
			})
		}(i)
	}

	wg.Wait()
	close(panicCh)

	for recovered := range panicCh {
		t.Fatalf("BroadcastToAdmins panicked during concurrent writes: %v", recovered)
	}

	if err := <-readErrCh; err != nil {
		t.Fatalf("reading broadcasts failed: %v", err)
	}
}
