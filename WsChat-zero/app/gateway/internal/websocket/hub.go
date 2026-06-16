package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

// Client represents a single websocket connection.
type Client struct {
	UserId int64
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
}

// Hub keeps the latest websocket connection for each user.
type Hub struct {
	mu       sync.RWMutex
	clients  map[int64]*Client
	register chan *Client
	unregister chan *Client
}

var (
	GlobalHub *Hub
	once      sync.Once
)

func GetHub() *Hub {
	once.Do(func() {
		GlobalHub = &Hub{
			clients:    make(map[int64]*Client),
			register:   make(chan *Client, 256),
			unregister: make(chan *Client, 256),
		}
		go GlobalHub.run()
	})
	return GlobalHub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserId] = client
			total := len(h.clients)
			h.mu.Unlock()
			logx.Infof("ws client registered: userId=%d, total=%d", client.UserId, total)
		case client := <-h.unregister:
			h.mu.Lock()
			if current, ok := h.clients[client.UserId]; ok && current == client {
				delete(h.clients, client.UserId)
				close(client.Send)
			}
			total := len(h.clients)
			h.mu.Unlock()
			logx.Infof("ws client unregistered: userId=%d, total=%d", client.UserId, total)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) SendToUser(userId int64, data interface{}) error {
	h.mu.RLock()
	client, ok := h.clients[userId]
	h.mu.RUnlock()
	if !ok {
		logx.Infof("ws send skipped: userId=%d offline", userId)
		return nil
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	select {
	case client.Send <- msg:
		logx.Infof("ws send queued: userId=%d bytes=%d", userId, len(msg))
	default:
		logx.Errorf("ws send buffer full: userId=%d", userId)
	}
	return nil
}

func (h *Hub) BroadcastToUsers(userIds []int64, data interface{}) {
	for _, uid := range userIds {
		_ = h.SendToUser(uid, data)
	}
}

func (h *Hub) IsOnline(userId int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userId]
	return ok
}

func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
