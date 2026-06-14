package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type TextMessageData struct {
	Content    string `json:"content"`
	ReceiverId int64  `json:"receiver_id"`
	ChatType   int32  `json:"chat_type"`
	MsgType    int32  `json:"msg_type"`
	SessionId  int64  `json:"session_id"`
	FileUrl    string `json:"file_url"`
	FileSize   int64  `json:"file_size"`
	FileName   string `json:"file_name"`
}

type MessageEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type MessageEventData struct {
	SessionId   int64  `json:"sessionId"`
	SenderId    int64  `json:"senderId"`
	ReceiverId  int64  `json:"receiverId"`
	ChatType    int32  `json:"chatType"`
	MsgType     int32  `json:"msgType"`
	Content     string `json:"content"`
	FileUrl     string `json:"fileUrl,omitempty"`
	FileName    string `json:"fileName,omitempty"`
	FileSize    int64  `json:"fileSize,omitempty"`
	SendName    string `json:"sendName,omitempty"`
	SendAvatar  string `json:"sendAvatar,omitempty"`
}

func WSHandler(w http.ResponseWriter, r *http.Request, serverCtx *svc.ServiceContext) {
	token := strings.TrimPrefix(r.URL.Query().Get("token"), "Bearer ")
	if token == "" {
		httpx.OkJson(w, map[string]interface{}{
			"code":    401,
			"message": "token 不能为空",
		})
		return
	}

	userId, err := parseUserIDFromToken(token, serverCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		logx.Errorf("ws auth failed: token=%s err=%v", token[:min(len(token), 10)], err)
		httpx.OkJson(w, map[string]interface{}{
			"code":    401,
			"message": "token 无效",
		})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Errorf("ws upgrade error: %v", err)
		return
	}

	hub := GetHub()
	client := &Client{
		UserId: userId,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.Register(client)
	go client.WritePump()
	go client.ReadPump(func(uid int64, raw []byte) {
		handleIncomingMessage(context.Background(), serverCtx, hub, uid, raw)
	})
}

func fetchSenderInfo(ctx context.Context, serverCtx *svc.ServiceContext, uid int64) (string, string) {
	resp, err := serverCtx.UserClient.GetUserInfo(ctx, &userpb.GetUserInfoRequest{UserId: uid})
	if err != nil || resp.Code != 0 || resp.Data == nil {
		return "", ""
	}
	return resp.Data.Nickname, resp.Data.Avatar
}

func handleIncomingMessage(ctx context.Context, serverCtx *svc.ServiceContext, hub *Hub, uid int64, raw []byte) {
	var wsMsg WsMessage
	if err := json.Unmarshal(raw, &wsMsg); err != nil {
		logx.Errorf("ws unmarshal error: %v", err)
		return
	}

	switch wsMsg.Type {
	case "heartbeat":
		return
	case "text":
		var payload TextMessageData
		if err := json.Unmarshal(wsMsg.Data, &payload); err != nil {
			logx.Errorf("ws text payload error: %v", err)
			return
		}
		if payload.MsgType == 0 {
			payload.MsgType = 1
		}
		if payload.ChatType == 0 {
			payload.ChatType = 1
		}
		if payload.ReceiverId <= 0 {
			logx.Errorf("ws text payload invalid: sender=%d receiver=%d session=%d", uid, payload.ReceiverId, payload.SessionId)
			return
		}

		// 获取发送者信息
		senderName, senderAvatar := fetchSenderInfo(ctx, serverCtx, uid)

		if payload.ChatType == 1 {
			receiverSessionId := payload.SessionId
			peerSessionId, err := ensurePeerSession(ctx, serverCtx.FriendClient, uid, payload.ReceiverId)
			if err != nil {
				logx.Errorf("ensure peer session failed: sender=%d receiver=%d err=%v", uid, payload.ReceiverId, err)
				return
			}
			receiverSessionId = peerSessionId

			if _, err := serverCtx.MsgClient.SendMessage(ctx, &msgpb.SendMessageRequest{
				SenderId:     uid,
				ReceiverId:   payload.ReceiverId,
				ChatType:     payload.ChatType,
				MsgType:      payload.MsgType,
				Content:      payload.Content,
				FileUrl:      payload.FileUrl,
				FileSize:     payload.FileSize,
				FileName:     payload.FileName,
				SessionId:    payload.SessionId,
				SenderName:   senderName,
				SenderAvatar: senderAvatar,
			}); err != nil {
				logx.Errorf("send message failed for sender session: %v", err)
				return
			}

			_ = hub.SendToUser(uid, MessageEvent{
				Type: "message:new",
				Data: MessageEventData{
					SessionId:  payload.SessionId,
					SenderId:   uid,
					ReceiverId: payload.ReceiverId,
					ChatType:   payload.ChatType,
					MsgType:    payload.MsgType,
					Content:    payload.Content,
					FileUrl:    payload.FileUrl,
					FileName:   payload.FileName,
					FileSize:   payload.FileSize,
					SendName:   senderName,
					SendAvatar: senderAvatar,
				},
			})
			_ = hub.SendToUser(payload.ReceiverId, MessageEvent{
				Type: "message:new",
				Data: MessageEventData{
					SessionId:  receiverSessionId,
					SenderId:   uid,
					ReceiverId: payload.ReceiverId,
					ChatType:   payload.ChatType,
					MsgType:    payload.MsgType,
					Content:    payload.Content,
					FileUrl:    payload.FileUrl,
					FileName:   payload.FileName,
					FileSize:   payload.FileSize,
					SendName:   senderName,
					SendAvatar: senderAvatar,
				},
			})
			return
		}

		groupName, sessionByUser, err := ensureGroupSessions(ctx, serverCtx.FriendClient, uid, payload.ReceiverId)
		if err != nil {
			logx.Errorf("ensure group sessions failed: sender=%d group=%d err=%v", uid, payload.ReceiverId, err)
			return
		}

		senderSessionId := sessionByUser[uid]
		if senderSessionId <= 0 {
			logx.Errorf("group sender session missing: sender=%d group=%d", uid, payload.ReceiverId)
			return
		}

		if _, err := serverCtx.MsgClient.SendMessage(ctx, &msgpb.SendMessageRequest{
			SenderId:     uid,
			ReceiverId:   payload.ReceiverId,
			ChatType:     payload.ChatType,
			MsgType:      payload.MsgType,
			Content:      payload.Content,
			FileUrl:      payload.FileUrl,
			FileSize:     payload.FileSize,
			FileName:     payload.FileName,
			SessionId:    senderSessionId,
			SenderName:   senderName,
			SenderAvatar: senderAvatar,
		}); err != nil {
			logx.Errorf("send group message failed: sender=%d group=%d err=%v", uid, payload.ReceiverId, err)
			return
		}

		logx.Infof("group message enqueued: sender=%d group=%d session=%d name=%s", uid, payload.ReceiverId, senderSessionId, groupName)
		for userId, sessionId := range sessionByUser {
			_ = hub.SendToUser(userId, MessageEvent{
				Type: "message:new",
				Data: MessageEventData{
					SessionId:  sessionId,
					SenderId:   uid,
					ReceiverId: payload.ReceiverId,
					ChatType:   payload.ChatType,
					MsgType:    payload.MsgType,
					Content:    payload.Content,
					FileUrl:    payload.FileUrl,
					FileName:   payload.FileName,
					FileSize:   payload.FileSize,
					SendName:   senderName,
					SendAvatar: senderAvatar,
				},
			})
		}
	default:
		logx.Infof("ws message ignored: user=%d type=%s", uid, wsMsg.Type)
	}
}

func ensurePeerSession(ctx context.Context, friendClient friendpb.FriendService, senderId, receiverId int64) (int64, error) {
	resp, err := friendClient.CreateSession(ctx, &friendpb.CreateSessionRequest{
		UserId:      receiverId,
		PeerId:      senderId,
		SessionType: 1,
	})
	if err != nil {
		return 0, err
	}
	if resp == nil || resp.Code != 0 || resp.Data == nil {
		if resp == nil {
			return 0, fmt.Errorf("peer session response is nil")
		}
		return 0, fmt.Errorf("peer session create failed: code=%d message=%s", resp.Code, resp.Message)
	}
	return resp.Data.Id, nil
}

func ensureGroupSessions(ctx context.Context, friendClient friendpb.FriendService, senderId, groupId int64) (string, map[int64]int64, error) {
	groupResp, err := friendClient.GetGroupInfo(ctx, &friendpb.GetGroupInfoRequest{GroupId: groupId})
	if err != nil {
		return "", nil, err
	}
	if groupResp == nil || groupResp.Code != 0 || groupResp.Data == nil {
		if groupResp == nil {
			return "", nil, fmt.Errorf("group info response is nil")
		}
		return "", nil, fmt.Errorf("group info query failed: code=%d message=%s", groupResp.Code, groupResp.Message)
	}

	memberResp, err := friendClient.GetGroupMemberList(ctx, &friendpb.GetGroupMemberListRequest{GroupId: groupId})
	if err != nil {
		return "", nil, err
	}
	if memberResp == nil || memberResp.Code != 0 {
		if memberResp == nil {
			return "", nil, fmt.Errorf("group member response is nil")
		}
		return "", nil, fmt.Errorf("group member query failed: code=%d message=%s", memberResp.Code, memberResp.Message)
	}

	sessionByUser := make(map[int64]int64, len(memberResp.Data)+1)
	foundSender := false
	for _, member := range memberResp.Data {
		if member.UserId == senderId {
			foundSender = true
		}
		sessionResp, err := friendClient.CreateSession(ctx, &friendpb.CreateSessionRequest{
			UserId:      member.UserId,
			PeerId:      groupId,
			SessionType: 2,
			SessionName: groupResp.Data.Name,
		})
		if err != nil {
			return "", nil, err
		}
		if sessionResp == nil || sessionResp.Code != 0 || sessionResp.Data == nil {
			if sessionResp == nil {
				return "", nil, fmt.Errorf("group session response is nil")
			}
			return "", nil, fmt.Errorf("group session create failed: code=%d message=%s", sessionResp.Code, sessionResp.Message)
		}
		sessionByUser[member.UserId] = sessionResp.Data.Id
	}

	if !foundSender {
		return "", nil, fmt.Errorf("sender not in group")
	}

	return groupResp.Data.Name, sessionByUser, nil
}

func parseUserIDFromToken(tokenString, secret string) (int64, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || token == nil || !token.Valid {
		return 0, fmt.Errorf("invalid jwt")
	}

	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("missing user_id")
	}
	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return 0, fmt.Errorf("token expired")
	}
	return int64(uid), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
