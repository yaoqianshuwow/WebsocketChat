// WS 单聊功能测试脚本 — 验证双方消息列表一致
// go run test/ws_chat_check.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	BaseURL = "http://localhost:8888"
	WSUrl   = "ws://localhost:8888"
)

type LoginResp struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
	Token    string `json:"token"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UserInfoResp struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type WsMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MsgEventData struct {
	SessionId  int64  `json:"sessionId"`
	SenderId   int64  `json:"senderId"`
	ReceiverId int64  `json:"receiverId"`
	ChatType   int32  `json:"chatType"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
	SendName   string `json:"sendName"`
	SendAvatar string `json:"sendAvatar"`
}

type MessageVo struct {
	MsgId      int64  `json:"msgId"`
	SenderId   int64  `json:"senderId"`
	Content    string `json:"content"`
	SendName   string `json:"sendName"`
	SendAvatar string `json:"sendAvatar"`
	MsgType    int32  `json:"msgType"`
}

type MessageListResp struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    []MessageVo `json:"data"`
	Total   int64       `json:"total"`
}

type SessionVo struct {
	SessionId      int64  `json:"sessionId"`
	PeerId         int64  `json:"peerId"`
	SessionType    int32  `json:"sessionType"`
	SessionName    string `json:"sessionName"`
	LastMsgContent string `json:"lastMsgContent"`
	LastMsgTime    int64  `json:"lastMsgTime"`
}

type SessionsResp struct {
	Code int32       `json:"code"`
	Data []SessionVo `json:"data"`
}

func main() {
	fmt.Println("══════════════════════════════════════════════════")
	fmt.Println("  WS 单聊消息一致性测试")
	fmt.Println("══════════════════════════════════════════════════")

	// 1. 登录
	u1 := login("u1", "111111")
	u2 := login("u2", "111111")
	fmt.Printf("u1: id=%d name=%s\n", u1.UserId, u1.Nickname)
	fmt.Printf("u2: id=%d name=%s\n", u2.UserId, u2.Nickname)

	// 2. 获取双方会话
	u1Sess := getSessions(u1.Token)
	u2Sess := getSessions(u2.Token)

	var u1SessionId, u2SessionId int64
	for _, s := range u1Sess {
		if s.SessionType == 1 && s.PeerId == u2.UserId {
			u1SessionId = s.SessionId
			break
		}
	}
	for _, s := range u2Sess {
		if s.SessionType == 1 && s.PeerId == u1.UserId {
			u2SessionId = s.SessionId
			break
		}
	}
	fmt.Printf("u1↔u2 会话: u1视角 sessionId=%d, u2视角 sessionId=%d\n", u1SessionId, u2SessionId)

	// 3. 连接 WS
	u1WS := connectWS(u1.Token)
	defer u1WS.Close()
	u2WS := connectWS(u2.Token)
	defer u2WS.Close()
	time.Sleep(500 * time.Millisecond)

	// 4. 记录双方消息数量
	before1 := len(getMessageList(u1.Token, u1SessionId))
	before2 := len(getMessageList(u2.Token, u2SessionId))
	fmt.Printf("测试前 u1消息数=%d, u2消息数=%d\n", before1, before2)

	// 5. u1 发消息给 u2
	fmt.Println("\n─ u1 → u2 发消息 ─")
	tag := fmt.Sprintf("test-%d", time.Now().Unix())
	msg1 := fmt.Sprintf("u1 say: %s", tag)
	sendText(u1WS, u2.UserId, u1SessionId, msg1)
	time.Sleep(2 * time.Second)

	// 6. 检查 u2 实时收到
	ev := readWS(u2WS)
	if ev != nil && ev.Type == "message:new" {
		var d MsgEventData
		json.Unmarshal(ev.Data, &d)
		fmt.Printf("u2 WS收到: [%d]%s → %s\n", d.SenderId, d.SendName, d.Content)
		if d.Content == msg1 {
			fmt.Println("  ✅ u2 实时收到 u1 消息")
		}
	}

	// 7. u2 回复
	fmt.Println("\n─ u2 → u1 回复 ─")
	msg2 := fmt.Sprintf("u2 reply: %s", tag)
	sendText(u2WS, u1.UserId, u2SessionId, msg2)
	time.Sleep(2 * time.Second)

	ev = readWS(u1WS)
	if ev != nil && ev.Type == "message:new" {
		var d MsgEventData
		json.Unmarshal(ev.Data, &d)
		fmt.Printf("u1 WS收到: [%d]%s → %s\n", d.SenderId, d.SendName, d.Content)
		if d.Content == msg2 {
			fmt.Println("  ✅ u1 实时收到 u2 回复")
		}
	}

	// 8. 验证 API 列表
	time.Sleep(1 * time.Second)
	fmt.Println("\n─ API 验证 ─")

	msgs1 := getMessageList(u1.Token, u1SessionId)
	msgs2 := getMessageList(u2.Token, u2SessionId)

	fmt.Printf("\nu1 消息列表 (%d 条, sessionId=%d):\n", len(msgs1), u1SessionId)
	for _, m := range msgs1 {
		icon := map[int32]string{1: "💬", 2: "📎"}[m.MsgType]
		fmt.Printf("  %s [%d] %s: %s\n", icon, m.SenderId, m.SendName, truncate(m.Content, 40))
	}

	fmt.Printf("\nu2 消息列表 (%d 条, sessionId=%d):\n", len(msgs2), u2SessionId)
	for _, m := range msgs2 {
		icon := map[int32]string{1: "💬", 2: "📎"}[m.MsgType]
		fmt.Printf("  %s [%d] %s: %s\n", icon, m.SenderId, m.SendName, truncate(m.Content, 40))
	}

	// 9. 核心验证：双方消息数增长一致
	fmt.Println("\n─ 校验 ─")
	after1 := len(msgs1)
	after2 := len(msgs2)
	fmt.Printf("u1: %d → %d (+%d)\n", before1, after1, after1-before1)
	fmt.Printf("u2: %d → %d (+%d)\n", before2, after2, after2-before2)

	// 新消息应在双方列表中都出现
	found1 := false
	for _, m := range msgs1 {
		if m.Content == msg1 || m.Content == msg2 {
			found1 = true
			break
		}
	}
	found2 := false
	for _, m := range msgs2 {
		if m.Content == msg1 || m.Content == msg2 {
			found2 = true
			break
		}
	}

	allPass := true
	if !found1 {
		fmt.Println("  ❌ u1 消息列表未包含新测试消息!")
		allPass = false
	}
	if !found2 {
		fmt.Println("  ❌ u2 消息列表未包含新测试消息!")
		allPass = false
	}
	if after1-before1 >= 2 && after2-before2 >= 2 {
		fmt.Println("  ✅ 双方均新增 2+ 条消息")
	} else {
		fmt.Printf("  ⚠️ 消息增量: u1=%d, u2=%d (预期各 >=2)\n", after1-before1, after2-before2)
		allPass = false
	}

	fmt.Println()
	if allPass {
		fmt.Println("══════════════════════════════════════════════════")
		fmt.Println("  🎉 全部通过! 双方消息列表一致")
		fmt.Println("══════════════════════════════════════════════════")
	} else {
		fmt.Println("══════════════════════════════════════════════════")
		fmt.Println("  ⚠️ 部分校验未通过，请检查")
		fmt.Println("══════════════════════════════════════════════════")
	}
}

// ── HTTP ──

func login(u, p string) *LoginResp {
	b := must(json.Marshal(map[string]string{"username": u, "password": p}))
	r := post("/api/v1/login", b, "")
	var lr LoginResp
	mustV(json.Unmarshal(r, &lr))
	if lr.Code != 0 {
		panic(lr.Message)
	}
	return &lr
}

func getSessions(tok string) []SessionVo {
	r := post("/api/v1/session/getUserSessionList", must(json.Marshal(map[string]int32{"sessionType": 0})), tok)
	var sr SessionsResp
	mustV(json.Unmarshal(r, &sr))
	return sr.Data
}

func getMessageList(tok string, sid int64) []MessageVo {
	b := must(json.Marshal(map[string]interface{}{"sessionId": sid, "page": 1, "size": 20}))
	r := post("/api/v1/message/getMessageList", b, tok)
	var mlr MessageListResp
	mustV(json.Unmarshal(r, &mlr))
	if mlr.Code != 0 {
		return nil
	}
	return mlr.Data
}

func post(path string, body []byte, token string) []byte {
	req, _ := http.NewRequest("POST", BaseURL+path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("HTTP %d %s: %s", resp.StatusCode, path, string(b)))
	}
	return b
}

// ── WS ──

func connectWS(token string) *websocket.Conn {
	c, _, err := websocket.DefaultDialer.Dial(WSUrl+"/wss?token="+token, nil)
	if err != nil {
		panic(err)
	}
	return c
}

func sendText(conn *websocket.Conn, toId, sessId int64, content string) {
	conn.WriteJSON(WsMsg{Type: "text", Data: must(json.Marshal(map[string]interface{}{
		"content":     content,
		"receiver_id": toId,
		"chat_type":   1,
		"msg_type":    1,
		"session_id":  sessId,
	}))})
}

func readWS(conn *websocket.Conn) *WsMsg {
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		return nil
	}
	var m WsMsg
	if json.Unmarshal(raw, &m) != nil {
		return nil
	}
	if m.Type == "heartbeat" {
		return readWS(conn)
	}
	return &m
}

// ── util ──

func must(b []byte, e error) []byte {
	if e != nil {
		panic(e)
	}
	return b
}

func mustV(e error) {
	if e != nil {
		panic(e)
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "..."
}
