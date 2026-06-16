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

const (
	singleBaseURL = "http://localhost:8888"
	singleWSURL   = "ws://localhost:8888"
)

type singleLoginResp struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

type singleRegisterResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	UserID  int64  `json:"user_id"`
}

type singleSessionsResp struct {
	Code int32 `json:"code"`
	Data []struct {
		SessionID   int64 `json:"sessionId"`
		PeerID      int64 `json:"peerId"`
		SessionType int32 `json:"sessionType"`
	} `json:"data"`
}

type singleMessageListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		MsgID    int64  `json:"msgId"`
		Content  string `json:"content"`
		SenderID int64  `json:"senderId"`
	} `json:"data"`
}

type singleWSEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type singleWSMessage struct {
	SessionID int64  `json:"sessionId"`
	SenderID  int64  `json:"senderId"`
	Content   string `json:"content"`
}

func main() {
	tag := fmt.Sprintf("%d-%d", time.Now().UnixNano(), os.Getpid())
	u1 := singleEnsureAccount("u1-"+tag, "111111")
	u2 := singleEnsureAccount("u2-"+tag, "111111")

	u1Session := singleEnsureSession(u1.Token, u2.UserID, u2.Nickname)
	u2Session := singleEnsureSession(u2.Token, u1.UserID, u1.Nickname)

	u1Conn := singleConnectWS(u1.Token)
	defer u1Conn.Close()
	u2Conn := singleConnectWS(u2.Token)
	defer u2Conn.Close()

	before1 := len(singleGetMessageList(u1.Token, u1Session))
	before2 := len(singleGetMessageList(u2.Token, u2Session))

	msg1 := "u1->u2 single-" + tag
	msg2 := "u2->u1 single-" + tag

	singleSend(u1Conn, u2.UserID, u1Session, msg1)
	ev1 := singleReadMessage(u2Conn)
	singleSend(u2Conn, u1.UserID, u2Session, msg2)
	ev2 := singleReadMessage(u1Conn)

	afterMsgs1 := singleWaitMessages(u1.Token, u1Session, []string{msg1, msg2})
	afterMsgs2 := singleWaitMessages(u2.Token, u2Session, []string{msg1, msg2})

	fmt.Printf("u1 session=%d before=%d after=%d\n", u1Session, before1, len(afterMsgs1))
	fmt.Printf("u2 session=%d before=%d after=%d\n", u2Session, before2, len(afterMsgs2))

	must(singleContains(afterMsgs1, msg1) && singleContains(afterMsgs1, msg2), "u1 消息列表未同时包含两条测试消息")
	must(singleContains(afterMsgs2, msg1) && singleContains(afterMsgs2, msg2), "u2 消息列表未同时包含两条测试消息")
	must(ev1 != nil && ev1.Content == msg1, "u2 没有实时收到 u1 消息")
	must(ev2 != nil && ev2.Content == msg2, "u1 没有实时收到 u2 消息")

	fmt.Println("单聊校验通过")
}

func singleEnsureAccount(user, pass string) *singleLoginResp {
	var regResp singleRegisterResp
	singlePost("/api/v1/register", map[string]any{
		"username": user,
		"password": pass,
		"nickname": user,
	}, "", &regResp)
	if regResp.Code != 0 && regResp.Message != "用户名已存在" {
		panic("注册失败: " + regResp.Message)
	}
	return singleLogin(user, pass)
}

func singleLogin(user, pass string) *singleLoginResp {
	var resp singleLoginResp
	singlePost("/api/v1/login", map[string]any{"username": user, "password": pass}, "", &resp)
	must(resp.Code == 0, "登录失败: "+resp.Message)
	return &resp
}

func singleEnsureSession(token string, peerID int64, peerName string) int64 {
	var openResp map[string]any
	singlePost("/api/v1/session/openSession", map[string]any{
		"peerId":      peerID,
		"sessionType": 1,
		"sessionName": peerName,
	}, token, &openResp)

	var sessions singleSessionsResp
	singlePost("/api/v1/session/getUserSessionList", map[string]any{"sessionType": 0}, token, &sessions)
	for _, item := range sessions.Data {
		if item.SessionType == 1 && item.PeerID == peerID {
			return item.SessionID
		}
	}
	panic("未找到单聊 session")
}

func singleGetMessageList(token string, sessionID int64) []struct {
	MsgID    int64  `json:"msgId"`
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	var resp singleMessageListResp
	singlePost("/api/v1/message/getMessageList", map[string]any{
		"sessionId": sessionID,
		"page":      1,
		"size":      20,
	}, token, &resp)
	must(resp.Code == 0, "拉取单聊消息失败")
	return resp.Data
}

func singleConnectWS(token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(singleWSURL+"/wss?token="+token, nil)
	if err != nil {
		panic(err)
	}
	return conn
}

func singleSend(conn *websocket.Conn, receiverID, sessionID int64, content string) {
	err := conn.WriteJSON(map[string]any{
		"type": "text",
		"data": map[string]any{
			"content":     content,
			"receiver_id": receiverID,
			"chat_type":   1,
			"msg_type":    1,
			"session_id":  sessionID,
		},
	})
	if err != nil {
		panic(err)
	}
}

func singleReadMessage(conn *websocket.Conn) *singleWSMessage {
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var event singleWSEvent
		if json.Unmarshal(raw, &event) != nil || event.Type == "heartbeat" {
			continue
		}
		if event.Type != "message:new" {
			continue
		}
		var msg singleWSMessage
		if json.Unmarshal(event.Data, &msg) != nil {
			return nil
		}
		return &msg
	}
}

func singlePost(path string, payload any, token string, out any) {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, singleBaseURL+path, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		panic(string(raw))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		panic(err)
	}
}

func singleContains(list []struct {
	MsgID    int64  `json:"msgId"`
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
}, content string) bool {
	for _, item := range list {
		if item.Content == content {
			return true
		}
	}
	return false
}

func singleWaitMessages(token string, sessionID int64, contents []string) []struct {
	MsgID    int64  `json:"msgId"`
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	for i := 0; i < 10; i++ {
		list := singleGetMessageList(token, sessionID)
		ok := true
		for _, content := range contents {
			if !singleContains(list, content) {
				ok = false
				break
			}
		}
		if ok {
			return list
		}
		time.Sleep(800 * time.Millisecond)
	}
	return singleGetMessageList(token, sessionID)
}

func must(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
