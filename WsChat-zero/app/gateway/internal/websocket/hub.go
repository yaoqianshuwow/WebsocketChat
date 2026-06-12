package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

// Client 代表一个 WebSocket 连接
type Client struct {
	UserId   int64
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	mu       sync.RWMutex
	clients  map[int64]*Client  // userId -> Client
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
			h.mu.Unlock()
			logx.Infof("ws client registered: userId=%d, total=%d", client.UserId, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if current, ok := h.clients[client.UserId]; ok && current == client {
				delete(h.clients, client.UserId)
				close(client.Send)
			}
			h.mu.Unlock()
			logx.Infof("ws client unregistered: userId=%d, total=%d", client.UserId, len(h.clients))
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser 向指定用户发送消息
func (h *Hub) SendToUser(userId int64, data interface{}) error {
	h.mu.RLock()
	client, ok := h.clients[userId]
	h.mu.RUnlock()
	if !ok {
		return nil // 用户不在线
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	select {
	case client.Send <- msg:
	default:
		logx.Errorf("ws send buffer full: userId=%d", userId)
	}
	return nil
}

// BroadcastToUsers 向多个用户广播消息
func (h *Hub) BroadcastToUsers(userIds []int64, data interface{}) {
	for _, uid := range userIds {
		_ = h.SendToUser(uid, data)
	}
}

// IsOnline 检查用户是否在线
func (h *Hub) IsOnline(userId int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userId]
	return ok
}

// OnlineCount 获取在线人数
func (h *Hub) OnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
