// chat_pair_check.go — 单聊 + 群聊全流程自动化测试
// 测试 u1 ↔ u2 单聊，和 u1+u2 群聊的消息收发、WebSocket 实时推送、消息持久化
//
// 使用方式：go run test/chat_pair_check.go
//
// 前置条件：
// 1. 所有服务在 localhost 上运行（docker 依赖 + 本地业务服务）
// 2. u1、u2 账号已注册（密码 111111）
// 3. gateway 监听 localhost:8888

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
	baseURL = "http://localhost:8888"
	wsURL   = "ws://localhost:8888"
)

// ── 数据结构 ──

type LoginResp struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
	Token    string `json:"token"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

type SessionsResp struct {
	Code int32 `json:"code"`
	Data []struct {
		SessionID   int64  `json:"sessionId"`
		PeerID      int64  `json:"peerId"`
		SessionType int32  `json:"sessionType"`
		SessionName string `json:"sessionName"`
	} `json:"data"`
}

type GroupSearchResp struct {
	Code int32 `json:"code"`
	Data []struct {
		GroupID int64  `json:"group_id"`
		Name    string `json:"name"`
	} `json:"data"`
}

type GroupCreateResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	GroupID int64  `json:"group_id"`
}

type MessageListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		MsgID    int64  `json:"msgId"`
		Content  string `json:"content"`
		SenderID int64  `json:"senderId"`
	} `json:"data"`
}

type GroupMessageListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		Content  string `json:"content"`
		SenderID int64  `json:"senderId"`
	} `json:"data"`
}

type WSEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type WSMessage struct {
	SessionID  int64  `json:"sessionId"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	ChatType   int32  `json:"chatType"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
}

// ── 全局状态 ──

var (
	u1, u2        *LoginResp
	u1SessionID   int64
	u2SessionID   int64
	groupID       int64
	testResults   []string
)

func main() {
	fmt.Println("═══ WsChat-zero 全功能双向测试 ═══")
	startTime := time.Now()

	// 1. 登录
	step("1/8 用户登录")
	u1 = apiLogin("u1", "111111")
	u2 = apiLogin("u2", "111111")
	pass("u1 login ok, id=%d", u1.UserID)
	pass("u2 login ok, id=%d", u2.UserID)

	// 2. 搜索用户
	step("2/8 搜索用户")
	searchUser(u1.Token, "u2")
	pass("u1 搜索 u2 成功")

	// 3. 单聊会话
	step("3/8 获取/创建单聊会话")
	u1SessionID = findSession(u1.Token, u2.UserID, 1)
	u2SessionID = findSession(u2.Token, u1.UserID, 1)
	pass("单聊会话已就绪: u1 session=%d, u2 session=%d", u1SessionID, u2SessionID)

	// 4. 单聊消息收发
	step("4/8 单聊 WebSocket 双向消息")
	testSingleChat()

	// 5. 创建群组并加人
	step("5/8 创建群组并添加成员")
	createGroupAndAddMember()

	// 6. 群聊消息收发
	step("6/8 群聊 WebSocket 双向消息")
	testGroupChat()

	// 7. 消息历史记录验证
	step("7/8 消息历史记录验证")
	verifyMessageHistory()

	// 8. 联系人/群组接口验证
	step("8/8 联系人 & 群组接口验证")
	testContactAndGroupAPIs()

	// 汇总
	elapsed := time.Since(startTime)
	fmt.Printf("\n═══════════════════════════════════\n")
	fmt.Printf("测试完成，耗时 %v\n", elapsed)
	fmt.Printf("共 %d 项检查\n", countTests())
	for _, r := range testResults {
		fmt.Println(r)
	}
	writeResultFile(elapsed)
}

// ── 单聊测试 ──

func testSingleChat() {
	u1Conn := wsConnect(u1.Token)
	defer u1Conn.Close()
	u2Conn := wsConnect(u2.Token)
	defer u2Conn.Close()

	time.Sleep(500 * time.Millisecond)

	tag := fmt.Sprintf("single-%d", time.Now().UnixNano())
	msg1 := "u1->u2 " + tag
	msg2 := "u2->u1 " + tag

	before1 := len(getMessageList(u1.Token, u1SessionID))
	before2 := len(getMessageList(u2.Token, u2SessionID))

	wsSend(u1Conn, u2.UserID, u1SessionID, 1, msg1)
	// u1 发消息后自己的 WS 会收到 echo，u1 读消息时跳过自己的 echo
	ev1 := wsReadMessage(u2Conn, 5)
	must(ev1 != nil && ev1.Content == msg1, "u2 应实时收到 u1 单聊消息: %s", msg1)

	wsSend(u2Conn, u1.UserID, u2SessionID, 1, msg2)
	ev2 := wsSkipSelf(u1Conn, u1.UserID, 5)
	must(ev2 != nil && ev2.Content == msg2, "u1 应实时收到 u2 单聊消息: %s", msg2)

	after1 := waitMessages(u1.Token, u1SessionID, []string{msg1, msg2}, 10)
	after2 := waitMessages(u2.Token, u2SessionID, []string{msg1, msg2}, 10)

	must(len(after1) >= before1+2, "u1 消息列表应包含两条新消息 (before=%d after=%d)", before1, len(after1))
	must(len(after2) >= before2+2, "u2 消息列表应包含两条新消息 (before=%d after=%d)", before2, len(after2))
	must(contains(after1, msg1) && contains(after1, msg2), "u1 列表应包含 msg1 和 msg2")
	must(contains(after2, msg1) && contains(after2, msg2), "u2 列表应包含 msg1 和 msg2")
	pass("单聊双向消息收发验证通过")
}

// ── 群聊测试 ──

func createGroupAndAddMember() {
	var searchResp GroupSearchResp
	apiPost("/api/v1/group/searchGroupList", map[string]any{"keyword": "auto-test-group"}, u1.Token, &searchResp)
	for _, g := range searchResp.Data {
		if g.Name == "auto-test-group" {
			groupID = g.GroupID
			break
		}
	}
	if groupID == 0 {
		var createResp GroupCreateResp
		apiPost("/api/v1/group/createGroup", map[string]any{"groupName": "auto-test-group"}, u1.Token, &createResp)
		must(createResp.Code == 0 && createResp.GroupID > 0, "创建群组失败: %s", createResp.Message)
		groupID = createResp.GroupID
	}

	var memberListResp struct {
		Code int32 `json:"code"`
		Data []struct {
			UserID int64 `json:"user_id"`
		} `json:"data"`
	}
	apiPost("/api/v1/group/getGroupMemberList", map[string]any{"groupId": groupID}, u1.Token, &memberListResp)
	u2InGroup := false
	for _, m := range memberListResp.Data {
		if m.UserID == u2.UserID {
			u2InGroup = true
			break
		}
	}
	if !u2InGroup {
		var joinResp struct {
			Code    int32  `json:"code"`
			Message string `json:"message"`
		}
		apiPost("/api/v1/group/joinGroup", map[string]any{"groupId": groupID}, u2.Token, &joinResp)
		// 已在群组中也算成功
		must(joinResp.Code == 0 || joinResp.Code == 1, "u2 加群失败: code=%d msg=%s", joinResp.Code, joinResp.Message)
	}
	pass("群组已就绪: groupId=%d, u1=%d, u2=%d", groupID, u1.UserID, u2.UserID)
}

func testGroupChat() {
	u1Conn := wsConnect(u1.Token)
	defer u1Conn.Close()
	u2Conn := wsConnect(u2.Token)
	defer u2Conn.Close()

	time.Sleep(500 * time.Millisecond)

	tag := fmt.Sprintf("group-%d", time.Now().UnixNano())
	msg1 := "u1 group " + tag
	msg2 := "u2 group " + tag

	before1 := len(getGroupMessageList(u1.Token, groupID))
	before2 := len(getGroupMessageList(u2.Token, groupID))

	wsSend(u1Conn, groupID, 0, 2, msg1)
	ev1 := wsReadMessage(u2Conn, 5)
	must(ev1 != nil && ev1.ChatType == 2 && ev1.Content == msg1, "u2 应实时收到 u1 群消息: %s", msg1)

	wsSend(u2Conn, groupID, 0, 2, msg2)
	ev2 := wsSkipSelf(u1Conn, u1.UserID, 5)
	must(ev2 != nil && ev2.ChatType == 2 && ev2.Content == msg2, "u1 应实时收到 u2 群消息: %s", msg2)

	after1 := waitGroupMessages(u1.Token, groupID, []string{msg1, msg2}, 10)
	after2 := waitGroupMessages(u2.Token, groupID, []string{msg1, msg2}, 10)

	must(len(after1) >= before1+2, "u1 群消息列表应增长 >=2 (before=%d after=%d)", before1, len(after1))
	must(len(after2) >= before2+2, "u2 群消息列表应增长 >=2 (before=%d after=%d)", before2, len(after2))
	must(containsGroup(after1, msg1) && containsGroup(after1, msg2), "u1 群聊列表应包含 msg1 和 msg2")
	must(containsGroup(after2, msg1) && containsGroup(after2, msg2), "u2 群聊列表应包含 msg1 和 msg2")
	pass("群聊双向消息收发验证通过")
}

// ── 消息历史验证 ──

func verifyMessageHistory() {
	tag := fmt.Sprintf("history-%d", time.Now().UnixNano())
	msg := "history " + tag

	u1Conn := wsConnect(u1.Token)
	defer u1Conn.Close()
	wsSend(u1Conn, u2.UserID, u1SessionID, 1, msg)
	time.Sleep(1500 * time.Millisecond)

	list := waitMessages(u1.Token, u1SessionID, []string{msg}, 10)
	must(contains(list, msg), "单聊消息历史应包含最新消息: %s", msg)

	gMsg := "group-history " + tag
	gConn := wsConnect(u1.Token)
	defer gConn.Close()
	wsSend(gConn, groupID, 0, 2, gMsg)
	time.Sleep(1500 * time.Millisecond)

	gList := waitGroupMessages(u2.Token, groupID, []string{gMsg}, 10)
	must(containsGroup(gList, gMsg), "群聊消息历史应包含最新消息: %s", gMsg)
	pass("消息历史记录验证通过")
}

// ── 联系人&群组API验证 ──

func testContactAndGroupAPIs() {
	var contactListResp struct {
		Code int32 `json:"code"`
		Data []struct {
			UserID int64  `json:"user_id"`
			Name   string `json:"name"`
		} `json:"data"`
	}
	apiPost("/api/v1/contact/getUserList", map[string]any{}, u1.Token, &contactListResp)
	must(contactListResp.Code == 0, "获取联系人列表失败: %d", contactListResp.Code)
	pass("联系人列表接口正常: 共 %d 个联系人", len(contactListResp.Data))

	var myGroupsResp struct {
		Code int32 `json:"code"`
		Data []struct {
			GroupID int64  `json:"group_id"`
			Name    string `json:"name"`
		} `json:"data"`
	}
	apiPost("/api/v1/group/loadMyGroup", map[string]any{}, u1.Token, &myGroupsResp)
	must(myGroupsResp.Code == 0, "获取群组列表失败: %d", myGroupsResp.Code)
	pass("群组列表接口正常: 共 %d 个群组", len(myGroupsResp.Data))

	var memberResp struct {
		Code int32 `json:"code"`
	}
	apiPost("/api/v1/group/getGroupMemberList", map[string]any{"groupId": groupID}, u1.Token, &memberResp)
	must(memberResp.Code == 0, "获取群成员列表失败: %d", memberResp.Code)
	pass("群成员列表接口正常")

	var sessResp struct {
		Code int32 `json:"code"`
	}
	apiPost("/api/v1/session/getUserSessionList", map[string]any{"sessionType": 0}, u1.Token, &sessResp)
	must(sessResp.Code == 0, "获取会话列表失败: %d", sessResp.Code)
	pass("会话列表接口正常")

	var userInfoResp struct {
		Code     int32  `json:"code"`
		Nickname string `json:"nickname"`
	}
	apiPost("/api/v1/user/getUserInfo", map[string]any{"userId": u2.UserID}, u1.Token, &userInfoResp)
	must(userInfoResp.Code == 0, "获取用户信息失败: %d", userInfoResp.Code)
	pass("用户信息接口正常")
}

// ── HTTP 辅助 ──

func apiPost(path string, payload any, token string, out any) {
	body, err := json.Marshal(payload)
	if err != nil {
		fatal("序列化失败: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(body))
	if err != nil {
		fatal("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fatal("请求失败 %s: %v", path, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fatal("HTTP %d %s: %s", resp.StatusCode, path, string(raw))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		fatal("解析响应失败 %s: %v\n原始: %s", path, err, string(raw))
	}
}

func apiLogin(user, pass string) *LoginResp {
	var resp LoginResp
	apiPost("/api/v1/login", map[string]any{"username": user, "password": pass}, "", &resp)
	must(resp.Code == 0, "%s 登录失败: %s", user, resp.Message)
	return &resp
}

func searchUser(token, keyword string) {
	var resp struct {
		Code int32 `json:"code"`
		Data []struct {
			UserID   int64  `json:"user_id"`
			Username string `json:"username"`
		} `json:"data"`
	}
	apiPost("/api/v1/user/searchUsers", map[string]any{"keyword": keyword, "page": 1, "size": 20}, token, &resp)
	must(resp.Code == 0 && len(resp.Data) > 0, "搜索用户 %s 失败", keyword)
}

func findSession(token string, peerID int64, sessionType int32) int64 {
	var openResp struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
	}
	apiPost("/api/v1/session/openSession", map[string]any{
		"peerId":      peerID,
		"sessionType": sessionType,
		"sessionName": fmt.Sprintf("peer-%d", peerID),
	}, token, &openResp)

	var list SessionsResp
	apiPost("/api/v1/session/getUserSessionList", map[string]any{"sessionType": 0}, token, &list)
	for _, s := range list.Data {
		if s.SessionType == sessionType && s.PeerID == peerID {
			return s.SessionID
		}
	}
	fatal("未找到 session: peerID=%d type=%d", peerID, sessionType)
	return 0
}

func getMessageList(token string, sessionID int64) []struct {
	MsgID    int64  `json:"msgId"`
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	var resp MessageListResp
	apiPost("/api/v1/message/getMessageList", map[string]any{
		"sessionId": sessionID, "page": 1, "size": 50,
	}, token, &resp)
	must(resp.Code == 0, "拉取消息列表失败: %d", resp.Code)
	return resp.Data
}

func getGroupMessageList(token string, gid int64) []struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	var resp GroupMessageListResp
	apiPost("/api/v1/message/getGroupMessageList", map[string]any{
		"groupId": gid, "page": 1, "size": 50,
	}, token, &resp)
	must(resp.Code == 0, "拉取群消息列表失败: %d", resp.Code)
	return resp.Data
}

// wsSkipSelf 从 WebSocket 连接中读取一条 message:new 事件，跳过自己的 echo
// 返回匹配的消息，超时返回 nil
func wsSkipSelf(conn *websocket.Conn, selfID int64, timeoutSec int) *WSMessage {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutSec) * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var event WSEvent
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
		if msg.SenderID == selfID {
			continue // 跳过自己的 echo
		}
		return &msg
	}
}

// ── WebSocket 辅助 ──

func wsConnect(token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL+"/wss?token="+token, nil)
	if err != nil {
		fatal("WebSocket 连接失败: %v", err)
	}
	return conn
}

func wsSend(conn *websocket.Conn, receiverID, sessionID int64, chatType int32, content string) {
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
		fatal("WS 发送失败: %v", err)
	}
}

func wsReadMessage(conn *websocket.Conn, timeoutSec int) *WSMessage {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutSec) * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var event WSEvent
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

// ── 轮询辅助 ──

func waitMessages(token string, sessionID int64, contents []string, maxWait int) []struct {
	MsgID    int64  `json:"msgId"`
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	for i := 0; i < maxWait; i++ {
		list := getMessageList(token, sessionID)
		ok := true
		for _, c := range contents {
			if !contains(list, c) {
				ok = false
				break
			}
		}
		if ok {
			return list
		}
		time.Sleep(800 * time.Millisecond)
	}
	return getMessageList(token, sessionID)
}

func waitGroupMessages(token string, gid int64, contents []string, maxWait int) []struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	for i := 0; i < maxWait; i++ {
		list := getGroupMessageList(token, gid)
		ok := true
		for _, c := range contents {
			if !containsGroup(list, c) {
				ok = false
				break
			}
		}
		if ok {
			return list
		}
		time.Sleep(800 * time.Millisecond)
	}
	return getGroupMessageList(token, gid)
}

func contains(list []struct {
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

func containsGroup(list []struct {
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

// ── 报告辅助 ──

var testCount = 0
var failCount = 0

func must(ok bool, format string, args ...any) {
	testCount++
	if !ok {
		failCount++
		msg := fmt.Sprintf(format, args...)
		fmt.Printf("  ✗ FAIL: %s\n", msg)
		testResults = append(testResults, fmt.Sprintf("  ✗ FAIL: %s", msg))
	} else {
		testResults = append(testResults, fmt.Sprintf("  ✓ %s", fmt.Sprintf(format, args...)))
	}
}

func step(name string) {
	fmt.Printf("\n--- %s ---\n", name)
}

func pass(format string, args ...any) {
	fmt.Printf("  ✓ %s\n", fmt.Sprintf(format, args...))
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "致命错误: "+format+"\n", args...)
	os.Exit(1)
}

func countTests() int {
	return testCount
}

func writeResultFile(elapsed time.Duration) {
	filename := fmt.Sprintf("test/chat_pair_result_%s.txt", time.Now().Format("150405"))
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("无法写结果文件: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintf(f, "WsChat-zero 全功能双向测试结果\n")
	fmt.Fprintf(f, "时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(f, "耗时: %v\n", elapsed)
	fmt.Fprintf(f, "检查项: %d, 失败: %d\n\n", testCount, failCount)
	for _, r := range testResults {
		fmt.Fprintf(f, "%s\n", r)
	}
	fmt.Printf("\n结果已写入: %s\n", filename)
	if failCount > 0 {
		fmt.Printf("⚠  %d 项检查失败，请查看日志\n", failCount)
		os.Exit(1)
	} else {
		fmt.Printf("✅ 全部 %d 项检查通过！\n", testCount)
	}
}
