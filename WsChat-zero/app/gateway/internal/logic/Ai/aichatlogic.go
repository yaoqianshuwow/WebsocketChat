package Ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type AiChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAiChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatLogic {
	return &AiChatLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

type aiChatReq struct {
	Message   string `json:"message"`
	SessionID string `json:"sessionId"`
}

type aiChatResp struct {
	Reply string `json:"reply"`
}

func (l *AiChatLogic) AiChat(req *types.AiChatReq) (*types.AiChatResp, error) {
	userId := l.ctx.Value("userId").(int64)

	// 1. 保存用户消息到 message 表
	userMsgId := time.Now().UnixNano()
	_, err := l.svcCtx.MsgClient.SendMessage(l.ctx, &msgpb.SendMessageRequest{
		SenderId:   userId,
		ReceiverId: userId,
		ChatType:   3,
		MsgType:    1,
		Content:    req.Message,
		SessionId:  0,
		SenderName: "",
	})
	if err != nil {
		logx.Errorf("save user AI message failed: %v", err)
	}

	// 2. 调用 ppt-agent
	pptURL := l.svcCtx.Config.PptAgentURL
	if pptURL == "" {
		pptURL = "http://127.0.0.1:8123"
	}

	body, _ := json.Marshal(aiChatReq{
		Message:   req.Message,
		SessionID: fmt.Sprintf("user_%d", userId),
	})

	resp, err := http.Post(pptURL+"/api/chat/send", "application/json", bytes.NewReader(body))
	if err != nil {
		return &types.AiChatResp{Code: 1, Message: "AI服务连接失败: " + err.Error()}, nil
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var aiResp aiChatResp
	if err := json.Unmarshal(respBytes, &aiResp); err != nil {
		return &types.AiChatResp{Code: 1, Message: "AI响应解析失败"}, nil
	}

	// 3. 保存 AI 回复到 message 表
	_, _ = l.svcCtx.MsgClient.SendMessage(l.ctx, &msgpb.SendMessageRequest{
		SenderId:   -1,
		ReceiverId: userId,
		ChatType:   3,
		MsgType:    1,
		Content:    aiResp.Reply,
		SessionId:  0,
		SenderName: "PPT助手",
	})
	_ = userMsgId

	return &types.AiChatResp{Code: 0, Message: "ok", Reply: aiResp.Reply}, nil
}
