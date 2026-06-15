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
	groupBaseURL = "http://localhost:8888"
	groupWSURL   = "ws://localhost:8888"
)

type groupLoginResp struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token"`
	UserID  int64  `json:"user_id"`
}

type groupListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		GroupID int64  `json:"group_id"`
		Name    string `json:"name"`
	} `json:"data"`
}

type groupMessageListResp struct {
	Code int32 `json:"code"`
	Data []struct {
		Content  string `json:"content"`
		SenderID int64  `json:"senderId"`
	} `json:"data"`
}

type groupWSEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type groupWSMessage struct {
	SenderID   int64  `json:"senderId"`
	ReceiverID int64  `json:"receiverId"`
	ChatType   int32  `json:"chatType"`
	Content    string `json:"content"`
}

func main() {
	u1 := groupLogin("u1", "111111")
	u2 := groupLogin("u2", "111111")

	groupID := groupEnsureTestGroup(u1.Token, u2.Token)

	u1Conn := groupConnectWS(u1.Token)
	defer u1Conn.Close()
	u2Conn := groupConnectWS(u2.Token)
	defer u2Conn.Close()

	before1 := len(groupGetMessageList(u1.Token, groupID))
	before2 := len(groupGetMessageList(u2.Token, groupID))

	tag := fmt.Sprintf("group-%d", time.Now().UnixNano())
	msg1 := "u1 group " + tag
	msg2 := "u2 group " + tag

	groupSend(u1Conn, groupID, msg1)
	ev1 := groupReadMessage(u2Conn)
	groupSend(u2Conn, groupID, msg2)
	ev2 := groupReadMessage(u1Conn)

	after1 := groupWaitMessages(u1.Token, groupID, []string{msg1, msg2})
	after2 := groupWaitMessages(u2.Token, groupID, []string{msg1, msg2})

	mustGroup(len(after1) >= before1+2, "u1 群消息数量没有增长到预期值")
	mustGroup(len(after2) >= before2+2, "u2 群消息数量没有增长到预期值")
	mustGroup(groupContains(after1, msg1) && groupContains(after1, msg2), "u1 群消息列表缺少测试消息")
	mustGroup(groupContains(after2, msg1) && groupContains(after2, msg2), "u2 群消息列表缺少测试消息")
	mustGroup(ev1 != nil && ev1.ChatType == 2 && ev1.Content == msg1, "u2 没有实时收到 u1 群消息")
	mustGroup(ev2 != nil && ev2.ChatType == 2 && ev2.Content == msg2, "u1 没有实时收到 u2 群消息")

	fmt.Printf("群聊校验通过，groupId=%d\n", groupID)
}

func groupEnsureTestGroup(u1Token, u2Token string) int64 {
	groupID := groupFindByName(u1Token, "u1-u2-test-group")
	if groupID == 0 {
		var createResp map[string]any
		groupPost("/api/v1/group/createGroup", map[string]any{"groupName": "u1-u2-test-group"}, u1Token, &createResp)
		groupID = int64(createResp["group_id"].(float64))
	}

	groupPost("/api/v1/group/joinGroup", map[string]any{"groupId": groupID}, u2Token, &map[string]any{})
	return groupID
}

func groupFindByName(token, keyword string) int64 {
	var resp groupListResp
	groupPost("/api/v1/group/searchGroupList", map[string]any{"keyword": keyword}, token, &resp)
	for _, item := range resp.Data {
		if item.Name == keyword {
			return item.GroupID
		}
	}
	return 0
}

func groupLogin(user, pass string) *groupLoginResp {
	var resp groupLoginResp
	groupPost("/api/v1/login", map[string]any{"username": user, "password": pass}, "", &resp)
	mustGroup(resp.Code == 0, "登录失败: "+resp.Message)
	return &resp
}

func groupGetMessageList(token string, groupID int64) []struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	var resp groupMessageListResp
	groupPost("/api/v1/message/getGroupMessageList", map[string]any{
		"groupId": groupID,
		"page":    1,
		"size":    20,
	}, token, &resp)
	mustGroup(resp.Code == 0, "拉取群聊消息失败")
	return resp.Data
}

func groupConnectWS(token string) *websocket.Conn {
	conn, _, err := websocket.DefaultDialer.Dial(groupWSURL+"/wss?token="+token, nil)
	if err != nil {
		panic(err)
	}
	return conn
}

func groupSend(conn *websocket.Conn, groupID int64, content string) {
	err := conn.WriteJSON(map[string]any{
		"type": "text",
		"data": map[string]any{
			"content":     content,
			"receiver_id": groupID,
			"chat_type":   2,
			"msg_type":    1,
			"session_id":  0,
		},
	})
	if err != nil {
		panic(err)
	}
}

func groupReadMessage(conn *websocket.Conn) *groupWSMessage {
	_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return nil
		}
		var event groupWSEvent
		if json.Unmarshal(raw, &event) != nil || event.Type == "heartbeat" {
			continue
		}
		if event.Type != "message:new" {
			continue
		}
		var msg groupWSMessage
		if json.Unmarshal(event.Data, &msg) != nil {
			return nil
		}
		return &msg
	}
}

func groupPost(path string, payload any, token string, out any) {
	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, groupBaseURL+path, bytes.NewReader(body))
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

func groupContains(list []struct {
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

func groupWaitMessages(token string, groupID int64, contents []string) []struct {
	Content  string `json:"content"`
	SenderID int64  `json:"senderId"`
} {
	for i := 0; i < 10; i++ {
		list := groupGetMessageList(token, groupID)
		ok := true
		for _, content := range contents {
			if !groupContains(list, content) {
				ok = false
				break
			}
		}
		if ok {
			return list
		}
		time.Sleep(800 * time.Millisecond)
	}
	return groupGetMessageList(token, groupID)
}

func mustGroup(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}
