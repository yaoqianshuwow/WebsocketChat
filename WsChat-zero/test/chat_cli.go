package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultBaseURL = "http://localhost:8888"
	defaultWSURL   = "ws://localhost:8888"
)

type loginResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
	UserID  int64  `json:"user_id"`
}

type commonResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

type searchUsersResp struct {
	Code int32 `json:"code"`
	Data []struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	} `json:"data"`
}

type sessionsResp struct {
	Code int32 `json:"code"`
	Data []struct {
		SessionID   int64  `json:"sessionId"`
		PeerID      int64  `json:"peerId"`
		SessionType int32  `json:"sessionType"`
		SessionName string `json:"sessionName"`
	} `json:"data"`
}

type wsEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type wsIncoming struct {
	SessionID  int64  `json:"sessionId"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	ChatType   int32  `json:"chatType"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
	SendName   string `json:"sendName"`
}

func main() {
	var (
		baseURL   = flag.String("base", defaultBaseURL, "HTTP base URL")
		wsURL     = flag.String("ws", defaultWSURL, "WebSocket base URL")
		user      = flag.String("user", "", "login username")
		pass      = flag.String("pass", "111111", "login password")
		mode      = flag.String("mode", "single", "single or group")
		peer      = flag.String("peer", "", "peer username for single chat")
		groupID   = flag.Int64("group-id", 0, "group ID for group chat")
		sessionID = flag.Int64("session-id", 0, "override session ID for single chat")
	)
	flag.Parse()

	if *user == "" {
		exitf("必须传入 -user")
	}
	if *mode == "single" && *peer == "" {
		exitf("单聊模式必须传入 -peer")
	}
	if *mode == "group" && *groupID <= 0 {
		exitf("群聊模式必须传入 -group-id")
	}

	me := login(*baseURL, *user, *pass)
	fmt.Printf("登录成功: user=%s userId=%d mode=%s\n", *user, me.UserID, *mode)

	targetID := *groupID
	sendSessionID := *sessionID
	if *mode == "single" {
		peerUserID := lookupUserID(*baseURL, me.Token, *peer)
		if peerUserID == 0 {
			exitf("未找到用户 %s", *peer)
		}
		targetID = peerUserID
		if sendSessionID == 0 {
			sendSessionID = ensureSingleSession(*baseURL, me.Token, peerUserID, *peer)
		}
		fmt.Printf("单聊对象: %s (%d), sessionId=%d\n", *peer, peerUserID, sendSessionID)
	} else {
		fmt.Printf("群聊对象: groupId=%d\n", targetID)
	}

	conn := connectWS(*wsURL, me.Token)
	defer conn.Close()
	fmt.Println("WebSocket 已连接，输入内容后回车即可发送，输入 /quit 退出。")

	done := make(chan struct{})
	go readLoop(conn, done)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("[%s] > ", *user)
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "/quit" {
			break
		}

		msg := map[string]any{
			"type": "text",
			"data": map[string]any{
				"content":     line,
				"receiver_id": targetID,
				"chat_type":   chatType(*mode),
				"msg_type":    1,
				"session_id":  sendSessionID,
			},
		}
		if err := conn.WriteJSON(msg); err != nil {
			exitf("发送失败: %v", err)
		}
	}

	_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"), time.Now().Add(time.Second))
	<-done
}

func chatType(mode string) int32 {
	if mode == "group" {
		return 2
	}
	return 1
}

func readLoop(conn *websocket.Conn, done chan<- struct{}) {
	defer close(done)
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("\n连接已关闭: %v\n", err)
			return
		}

		var event wsEnvelope
		if err := json.Unmarshal(raw, &event); err != nil {
			fmt.Printf("\n收到非标准消息: %s\n", string(raw))
			continue
		}
		if event.Type == "heartbeat" {
			continue
		}
		if event.Type != "message:new" {
			fmt.Printf("\n收到事件 %s: %s\n", event.Type, string(event.Data))
			continue
		}

		var msg wsIncoming
		if err := json.Unmarshal(event.Data, &msg); err != nil {
			fmt.Printf("\n解析消息失败: %v\n", err)
			continue
		}
		fmt.Printf("\n[收到] sender=%d name=%s chatType=%d session=%d content=%s\n", msg.SenderID, msg.SendName, msg.ChatType, msg.SessionID, msg.Content)
	}
}

func login(baseURL, user, pass string) *loginResp {
	body := map[string]string{"username": user, "password": pass}
	var resp loginResp
	postJSON(baseURL, "/api/v1/login", body, "", &resp)
	if resp.Code != 0 {
		exitf("登录失败: %s", resp.Message)
	}
	return &resp
}

func lookupUserID(baseURL, token, keyword string) int64 {
	var resp searchUsersResp
	postJSON(baseURL, "/api/v1/user/searchUsers", map[string]any{
		"keyword": keyword,
		"page":    1,
		"size":    20,
	}, token, &resp)
	if resp.Code != 0 {
		exitf("查询用户失败")
	}
	for _, user := range resp.Data {
		if user.Username == keyword {
			return user.UserID
		}
	}
	return 0
}

func ensureSingleSession(baseURL, token string, peerID int64, peerName string) int64 {
	var openResp commonResp
	postJSON(baseURL, "/api/v1/session/openSession", map[string]any{
		"peerId":      peerID,
		"sessionType": 1,
		"sessionName": peerName,
	}, token, &openResp)
	if openResp.Code != 0 {
		exitf("打开会话失败: %s", openResp.Message)
	}

	var list sessionsResp
	postJSON(baseURL, "/api/v1/session/getUserSessionList", map[string]any{
		"sessionType": 0,
	}, token, &list)
	for _, session := range list.Data {
		if session.SessionType == 1 && session.PeerID == peerID {
			return session.SessionID
		}
	}
	exitf("未找到单聊 session")
	return 0
}

func connectWS(wsURL, token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL+"/wss?token="+token, nil)
	if err != nil {
		exitf("连接 WebSocket 失败: %v", err)
	}
	return conn
}

func postJSON(baseURL, path string, payload any, token string, out any) {
	body, err := json.Marshal(payload)
	if err != nil {
		exitf("序列化请求失败: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		exitf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		exitf("请求失败 %s: %v", path, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		exitf("读取响应失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		exitf("HTTP %d %s: %s", resp.StatusCode, path, string(raw))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		exitf("解析响应失败 %s: %v\n原始响应: %s", path, err, string(raw))
	}
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
