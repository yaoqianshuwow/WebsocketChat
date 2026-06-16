package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
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

type sessionListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		SessionID   int64  `json:"sessionId"`
		PeerID      int64  `json:"peerId"`
		SessionType int32  `json:"sessionType"`
		SessionName string `json:"sessionName"`
		UnreadCount int32  `json:"unreadCount"`
	} `json:"data"`
}

type contactListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		ContactID   int64  `json:"contactId"`
		ContactType int32  `json:"contactType"`
		Nickname    string `json:"nickname"`
		Avatar      string `json:"avatar"`
		Status      int32  `json:"status"`
	} `json:"data"`
}

type applyListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		ApplyID int64  `json:"applyId"`
		FromID  int64  `json:"fromId"`
		Status  int32  `json:"status"`
	} `json:"data"`
}

type messageListResp struct {
	Code    int32 `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		MsgID     int64  `json:"msgId"`
		SenderID  int64  `json:"senderId"`
		ReceiverID int64 `json:"receiverId"`
		MsgType   int32  `json:"msgType"`
		Content   string `json:"content"`
		FileURL   string `json:"fileUrl"`
		FileName  string `json:"fileName"`
		FileSize  int64  `json:"fileSize"`
		CreatedAt int64  `json:"createdAt"`
		SendName  string `json:"sendName"`
		SendAvatar string `json:"sendAvatar"`
	} `json:"data"`
	Total int64 `json:"total"`
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
		baseURL = flag.String("base", "http://localhost:8888", "gateway http base url")
		wsURL   = flag.String("ws", "ws://localhost:8888", "gateway ws base url")
		userA   = flag.String("user-a", "u1", "first username")
		passA   = flag.String("pass-a", "111111", "first password")
		userB   = flag.String("user-b", "u2", "second username")
		passB   = flag.String("pass-b", "111111", "second password")
	)
	flag.Parse()

	fmt.Println("== API smoke test ==")
	a := login(*baseURL, *userA, *passA)
	b := login(*baseURL, *userB, *passB)
	must(a.UserID > 0 && b.UserID > 0, "login should return user ids")

	peerBID := lookupUserID(*baseURL, a.Token, *userB)
	must(peerBID == b.UserID, "lookupUserID should resolve peer id, got=%d want=%d", peerBID, b.UserID)

	sessionA := ensureSingleSession(*baseURL, a.Token, b.UserID, *userB)
	sessionB := ensureSingleSession(*baseURL, b.Token, a.UserID, *userA)
	must(sessionA > 0 && sessionB > 0, "single sessions should exist")

	beforeSessions := listSessions(*baseURL, b.Token)
	beforeUnread := unreadForPeer(beforeSessions, a.UserID)

	aConn := wsConnect(*wsURL, a.Token)
	defer aConn.Close()
	bConn := wsConnect(*wsURL, b.Token)
	defer bConn.Close()

	tag := fmt.Sprintf("api-%d", time.Now().UnixNano())
	singleText := "single message " + tag
	wsSend(aConn, b.UserID, sessionA, 1, singleText)
	ev := wsReadMessage(bConn, 5)
	must(ev != nil && ev.Content == singleText, "ws echo to peer should match sent single message")

	waitMessages(*baseURL, a.Token, sessionA, []string{singleText}, 10)
	waitMessages(*baseURL, b.Token, sessionB, []string{singleText}, 10)

	searchResp := searchMessages(*baseURL, a.Token, sessionA, a.UserID, singleText)
	must(searchResp.Code == 0, "searchMessages should succeed")
	must(searchResp.Total > 0 && containsMessage(searchResp.Data, singleText), "searchMessages should find the inserted message")

	recentResp := getRecentMessages(*baseURL, b.Token, sessionB, 5)
	must(recentResp.Code == 0, "getRecentMessages should succeed")
	must(len(recentResp.Data) > 0 && recentResp.Data[0].Content != "", "getRecentMessages should return data")
	must(containsMessage(recentResp.Data, singleText), "getRecentMessages should include the inserted message")

	afterSessions := listSessions(*baseURL, b.Token)
	afterUnread := unreadForPeer(afterSessions, a.UserID)
	must(afterUnread >= beforeUnread, "session unread count should not decrease after incoming message")

	verifyDeleteAndRestore(*baseURL, a.Token, b.Token, a.UserID, b.UserID)

	fmt.Println("== API smoke test passed ==")
}

func verifyDeleteAndRestore(baseURL, aToken, bToken string, aID, bID int64) {
	fmt.Println("-- deleteContact / restore contact --")

	// Ensure friend relation exists before delete.
	if !containsContact(listContacts(baseURL, aToken), bID) {
		ensureFriend(baseURL, aToken, bToken, aID, bID)
	}
	must(containsContact(listContacts(baseURL, aToken), bID), "precondition: contact should exist")

	delResp := commonPost(baseURL, "/contact/deleteContact", map[string]any{"contactId": bID}, aToken)
	must(delResp.Code == 0, "deleteContact should succeed: %s", delResp.Message)
	must(!containsContact(listContacts(baseURL, aToken), bID), "contact should be removed after delete")

	ensureFriend(baseURL, aToken, bToken, aID, bID)
	must(containsContact(listContacts(baseURL, aToken), bID), "contact should be restored after re-apply")
}

func ensureFriend(baseURL, aToken, bToken string, aID, bID int64) {
	appResp := commonPost(baseURL, "/contact/applyContact", map[string]any{"toId": bID, "remark": "api-check friend restore"}, aToken)
	must(appResp.Code == 0, "applyContact should succeed: %s", appResp.Message)

	applies := listApplies(baseURL, bToken)
	applyID := int64(0)
	for _, item := range applies {
		if item.FromID == aID && item.Status == 0 {
			applyID = item.ApplyID
			break
		}
	}
	must(applyID > 0, "should find apply record from %d", aID)

	passResp := commonPost(baseURL, "/contact/passContactApply", map[string]any{"applyId": applyID, "status": 1}, bToken)
	must(passResp.Code == 0, "passContactApply should succeed: %s", passResp.Message)
}

func searchMessages(baseURL, token string, sessionID, senderID int64, keyword string) *messageListResp {
	var resp messageListResp
	postJSON(baseURL, "/message/searchMessages", map[string]any{
		"sessionId": sessionID,
		"senderId":  senderID,
		"keyword":   keyword,
		"page":      1,
		"size":      20,
	}, token, &resp)
	return &resp
}

func getRecentMessages(baseURL, token string, sessionID int64, limit int) *messageListResp {
	var resp messageListResp
	postJSON(baseURL, "/message/getRecentMessages", map[string]any{
		"sessionId": sessionID,
		"limit":     limit,
	}, token, &resp)
	return &resp
}

func lookupUserID(baseURL, token, keyword string) int64 {
	var resp struct {
		Code int32 `json:"code"`
		Data []struct {
			UserID int64 `json:"user_id"`
		} `json:"data"`
	}
	postJSON(baseURL, "/user/searchUsers", map[string]any{"keyword": keyword, "page": 1, "size": 10}, token, &resp)
	must(resp.Code == 0 && len(resp.Data) > 0, "searchUsers should find %s", keyword)
	return resp.Data[0].UserID
}

func ensureSingleSession(baseURL, token string, peerID int64, peerName string) int64 {
	commonPost(baseURL, "/session/openSession", map[string]any{"peerId": peerID, "sessionType": 1, "sessionName": peerName}, token)
	sessions := listSessions(baseURL, token)
	for _, s := range sessions {
		if s.SessionType == 1 && s.PeerID == peerID {
			return s.SessionID
		}
	}
	fatal("single session not found for peer=%d", peerID)
	return 0
}

func listSessions(baseURL, token string) []struct {
	SessionID   int64  `json:"sessionId"`
	PeerID      int64  `json:"peerId"`
	SessionType int32  `json:"sessionType"`
	SessionName string `json:"sessionName"`
	UnreadCount int32  `json:"unreadCount"`
} {
	var resp sessionListResp
	postJSON(baseURL, "/session/getUserSessionList", map[string]any{"sessionType": 0}, token, &resp)
	must(resp.Code == 0, "getUserSessionList should succeed")
	return resp.Data
}

func unreadForPeer(list []struct {
	SessionID   int64  `json:"sessionId"`
	PeerID      int64  `json:"peerId"`
	SessionType int32  `json:"sessionType"`
	SessionName string `json:"sessionName"`
	UnreadCount int32  `json:"unreadCount"`
}, peerID int64) int32 {
	for _, item := range list {
		if item.SessionType == 1 && item.PeerID == peerID {
			return item.UnreadCount
		}
	}
	return 0
}

func listContacts(baseURL, token string) []struct {
	ContactID   int64  `json:"contactId"`
	ContactType int32  `json:"contactType"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Status      int32  `json:"status"`
} {
	var resp contactListResp
	postJSON(baseURL, "/contact/getUserList", map[string]any{}, token, &resp)
	must(resp.Code == 0, "getUserList should succeed")
	return resp.Data
}

func listApplies(baseURL, token string) []struct {
	ApplyID int64 `json:"applyId"`
	FromID  int64 `json:"fromId"`
	Status  int32 `json:"status"`
} {
	var resp applyListResp
	postJSON(baseURL, "/contact/getNewContactList", map[string]any{}, token, &resp)
	must(resp.Code == 0, "getNewContactList should succeed")
	return resp.Data
}

func containsContact(list []struct {
	ContactID   int64  `json:"contactId"`
	ContactType int32  `json:"contactType"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Status      int32  `json:"status"`
}, contactID int64) bool {
	for _, item := range list {
		if item.ContactID == contactID {
			return true
		}
	}
	return false
}

func containsMessage(list []struct {
	MsgID      int64  `json:"msgId"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
	FileURL    string `json:"fileUrl"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  int64  `json:"createdAt"`
	SendName   string `json:"sendName"`
	SendAvatar string `json:"sendAvatar"`
}, keyword string) bool {
	for _, item := range list {
		if item.Content == keyword {
			return true
		}
	}
	return false
}

func commonPost(baseURL, path string, body any, token string) commonResp {
	var resp commonResp
	postJSON(baseURL, path, body, token, &resp)
	return resp
}

func postJSON(baseURL, path string, body any, token string, out any) {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/v1"+path, bytes.NewReader(buf))
	if err != nil {
		fatal("build request failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fatal("request %s failed: %v", path, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fatal("request %s status=%d body=%s", path, resp.StatusCode, string(raw))
	}
	if err := json.Unmarshal(raw, out); err != nil {
		fatal("decode %s failed: %v body=%s", path, err, string(raw))
	}
}

func login(baseURL, user, pass string) *loginResp {
	var resp loginResp
	postJSON(baseURL, "/login", map[string]any{"username": user, "password": pass}, "", &resp)
	must(resp.Code == 0, "%s login failed: %s", user, resp.Message)
	fmt.Printf("login ok: %s id=%d\n", user, resp.UserID)
	return &resp
}

func wsConnect(wsURL, token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL+"/wss?token="+token, nil)
	if err != nil {
		fatal("websocket connect failed: %v", err)
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
		fatal("ws send failed: %v", err)
	}
}

func wsReadMessage(conn *websocket.Conn, timeoutSec int) *wsIncoming {
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeoutSec) * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var env wsEnvelope
		if json.Unmarshal(raw, &env) != nil || env.Type == "heartbeat" {
			continue
		}
		if env.Type != "message:new" {
			continue
		}
		var msg wsIncoming
		if json.Unmarshal(env.Data, &msg) != nil {
			return nil
		}
		return &msg
	}
}

func waitMessages(baseURL, token string, sessionID int64, contents []string, maxWait int) []struct {
	MsgID      int64  `json:"msgId"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
	FileURL    string `json:"fileUrl"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  int64  `json:"createdAt"`
	SendName   string `json:"sendName"`
	SendAvatar string `json:"sendAvatar"`
} {
	for i := 0; i < maxWait; i++ {
		list := getMessageList(baseURL, token, sessionID)
		ok := true
		for _, c := range contents {
			if !containsMessage(list, c) {
				ok = false
				break
			}
		}
		if ok {
			return list
		}
		time.Sleep(800 * time.Millisecond)
	}
	return getMessageList(baseURL, token, sessionID)
}

func getMessageList(baseURL, token string, sessionID int64) []struct {
	MsgID      int64  `json:"msgId"`
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	MsgType    int32  `json:"msgType"`
	Content    string `json:"content"`
	FileURL    string `json:"fileUrl"`
	FileName   string `json:"fileName"`
	FileSize   int64  `json:"fileSize"`
	CreatedAt  int64  `json:"createdAt"`
	SendName   string `json:"sendName"`
	SendAvatar string `json:"sendAvatar"`
} {
	var resp messageListResp
	postJSON(baseURL, "/message/getMessageList", map[string]any{"sessionId": sessionID, "page": 1, "size": 50}, token, &resp)
	must(resp.Code == 0, "getMessageList should succeed: %s", resp.Message)
	return resp.Data
}

func must(ok bool, format string, args ...any) {
	if !ok {
		fatal(format, args...)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "FAIL: "+format+"\n", args...)
	os.Exit(1)
}
