package websocket

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/your-org/ws-chat-zero/app/user/userservice"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

// WsMessage WebSocket 消息格式
type WsMessage struct {
	Type string          `json:"type"` // text, file, system, heartbeat
	Data json.RawMessage `json:"data"`
}

// WSHandler 处理 WebSocket 升级请求（带鉴权）
func WSHandler(w http.ResponseWriter, r *http.Request, userClient userservice.UserService) {
	// 1. 从请求参数中获取 token
	token := r.URL.Query().Get("token")
	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		httpx.OkJson(w, map[string]interface{}{
			"code":    401,
			"message": "token 不能为空",
		})
		return
	}

	// 2. 调用 user RPC 验证 token
	resp, err := userClient.ValidateToken(r.Context(), &userservice.ValidateTokenRequest{
		Token: token,
	})
	if err != nil || resp == nil || resp.Code != 0 {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		logx.Errorf("ws auth failed: token=%s, err=%s", token[:min(len(token), 10)], errMsg)
		httpx.OkJson(w, map[string]interface{}{
			"code":    401,
			"message": "token 无效",
		})
		return
	}

	userId := resp.UserId
	logx.Infof("ws auth success: user_id=%d", userId)

	// 3. 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Errorf("ws upgrade error: %v", err)
		return
	}

	// 4. 注册客户端
	hub := GetHub()
	client := &Client{
		UserId: userId,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.Register(client)

	// 5. 启动读写协程
	go client.WritePump()
	go client.ReadPump(func(uid int64, msg []byte) {
		var wsMsg WsMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			logx.Errorf("ws unmarshal error: %v", err)
			return
		}
		logx.Infof("ws message from user=%d, type=%s", uid, wsMsg.Type)
		// 这里可以转发消息到 Kafka / gRPC
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
