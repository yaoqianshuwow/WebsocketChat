package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const serverURL = "http://120.77.251.18"
const serverWS = "ws://120.77.251.18"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test/ws_test_server.go <single|group>")
		os.Exit(1)
	}
	mode := os.Args[1]

	adminToken := login("admin", "123456")
	u1Token := login("u1", "111111")

	if mode == "single" {
		testSingleChat(adminToken, u1Token)
	} else {
		testGroupChat(adminToken, u1Token)
	}
}

func testSingleChat(adminToken, u1Token string) {
	// Ensure session
	openSession(adminToken, 2, 1, "u1")

	adminConn := connectWS(adminToken)
	defer adminConn.Close()
	u1Conn := connectWS(u1Token)
	defer u1Conn.Close()

	time.Sleep(1 * time.Second)

	tag := fmt.Sprintf("single-%d", time.Now().UnixNano())
	msg1 := "admin->u1 " + tag

	// Send from admin to u1 (single chat, session 2, chatType 1)
	sendWS(adminConn, 2, 2, 1, msg1)
	ev1 := readWS(u1Conn, 5)
	if ev1 != nil && ev1.Content == msg1 {
		fmt.Println("✅ u1 received single message:", ev1.Content)
	} else if ev1 != nil {
		fmt.Println("❌ u1 got wrong message:", ev1.Content)
	} else {
		fmt.Println("❌ u1 did not receive message (timeout)")
	}

	time.Sleep(1 * time.Second)
	// Verify in history
	list := getMessageList(adminToken, 2)
	found := false
	for _, m := range list {
		if m.Content == msg1 {
			found = true
			break
		}
	}
	if found {
		fmt.Println("✅ Message found in history")
	} else {
		fmt.Println("❌ Message NOT in history")
	}
}

func testGroupChat(adminToken, u1Token string) {
	// Get group ID
	groups := loadMyGroup(adminToken)
	var groupID int64
	for _, g := range groups {
		if g.Name == "test-group" {
			groupID = g.GroupID
			break
		}
	}
	if groupID == 0 {
		// Create
		groupID = createGroup(adminToken, "test-group")
		joinGroup(u1Token, groupID)
	}
	fmt.Printf("Using group ID: %d\n", groupID)

	adminConn := connectWS(adminToken)
	defer adminConn.Close()
	u1Conn := connectWS(u1Token)
	defer u1Conn.Close()

	time.Sleep(1 * time.Second)

	tag := fmt.Sprintf("group-%d", time.Now().UnixNano())
	msg1 := "admin group " + tag

	// Send from admin to group (chatType=2)
	sendWS(adminConn, groupID, 0, 2, msg1)
	ev1 := readWS(u1Conn, 5)
	if ev1 != nil && ev1.ChatType == 2 && ev1.Content == msg1 {
		fmt.Println("✅ u1 received group message:", ev1.Content)
	} else if ev1 != nil {
		fmt.Printf("❌ u1 got wrong - type=%d content=%s\n", ev1.ChatType, ev1.Content)
	} else {
		fmt.Println("❌ u1 did not receive group message (timeout)")
	}

	time.Sleep(1 * time.Second)
	// Verify in group history
	list := getGroupMessageList(adminToken, groupID)
	found := false
	for _, m := range list {
		if m.Content == msg1 {
			found = true
			break
		}
	}
	if found {
		fmt.Println("✅ Group message found in history")
	} else {
		fmt.Println("❌ Group message NOT in history")
	}
}

func login(user, pass string) string {
	var resp struct {
		Code    int32  `json:"code"`
		Token   string `json:"token"`
		Message string `json:"message"`
	}
	postJSON("/api/v1/login", map[string]string{"username": user, "password": pass}, "", &resp)
	if resp.Code != 0 {
		fmt.Printf("Login failed for %s: %s\n", user, resp.Message)
		os.Exit(1)
	}
	fmt.Printf("Login ok: %s\n", user)
	return resp.Token
}

func openSession(token string, peerID int64, sessionType int32, name string) {
	var resp struct {
		Code int32 `json:"code"`
	}
	postJSON("/api/v1/session/openSession", map[string]any{
		"peerId": peerID, "sessionType": sessionType, "sessionName": name,
	}, token, &resp)
}

func createGroup(token, name string) int64 {
	var resp struct {
		Code    int32 `json:"code"`
		GroupID int64 `json:"group_id"`
	}
	postJSON("/api/v1/group/createGroup", map[string]any{"groupName": name}, token, &resp)
	return resp.GroupID
}

func joinGroup(token string, groupID int64) {
	var resp struct {
		Code int32 `json:"code"`
	}
	postJSON("/api/v1/group/joinGroup", map[string]any{"groupId": groupID}, token, &resp)
}

type GroupInfo struct {
	GroupID int64  `json:"group_id"`
	Name    string `json:"name"`
}

func loadMyGroup(token string) []GroupInfo {
	var resp struct {
		Code int32      `json:"code"`
		Data []GroupInfo `json:"data"`
	}
	postJSON("/api/v1/group/loadMyGroup", map[string]any{}, token, &resp)
	return resp.Data
}

type MessageItem struct {
	Content string `json:"content"`
}

func getMessageList(token string, sessionID int64) []struct {
	Content string `json:"content"`
} {
	var resp struct {
		Code int32 `json:"code"`
		Data []struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	postJSON("/api/v1/message/getMessageList", map[string]any{
		"sessionId": sessionID, "page": 1, "size": 20,
	}, token, &resp)
	return resp.Data
}

func getGroupMessageList(token string, groupID int64) []struct {
	Content string `json:"content"`
} {
	var resp struct {
		Code int32 `json:"code"`
		Data []struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	postJSON("/api/v1/message/getGroupMessageList", map[string]any{
		"groupId": groupID, "page": 1, "size": 20,
	}, token, &resp)
	return resp.Data
}

func connectWS(token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(serverWS+"/wss?token="+token, nil)
	if err != nil {
		fmt.Printf("WS connect failed: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func sendWS(conn *websocket.Conn, receiverID, sessionID int64, chatType int32, content string) {
	err := conn.WriteJSON(map[string]any{
		"type": "text",
		"data": map[string]any{
			"content":     content,
			"receiver_id": receiverID,
			"chat_type":   chatType,
			"msg_type":    1,
			"session_id":  sessionID,
		},
	})
	if err != nil {
		fmt.Printf("WS send failed: %v\n", err)
	}
	fmt.Printf("Sent: %s (receiver=%d, type=%d, session=%d)\n", content, receiverID, chatType, sessionID)
}

type WSMessage struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
	ChatType int32  `json:"chatType"`
}

func readWS(conn *websocket.Conn, timeoutSec int) *WSMessage {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutSec) * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var event struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if json.Unmarshal(raw, &event) != nil || event.Type == "heartbeat" {
			continue
		}
		if event.Type != "message:new" {
			continue
		}
		var msg WSMessage
		if json.Unmarshal(event.Data, &msg) != nil {
			return nil
		}
		return &msg
	}
}

func postJSON(path string, payload any, token string, out any) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, serverURL+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Request failed %s: %v\n", path, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Printf("HTTP %d for %s: %s\n", resp.StatusCode, path, string(raw))
		os.Exit(1)
	}
	json.Unmarshal(raw, out)
}
