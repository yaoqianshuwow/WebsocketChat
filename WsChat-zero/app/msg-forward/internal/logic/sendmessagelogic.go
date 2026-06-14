package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type messagePayload struct {
	MsgId        int64  `json:"msg_id"`
	SenderId     int64  `json:"sender_id"`
	ReceiverId   int64  `json:"receiver_id"`
	ChatType     int32  `json:"chat_type"`
	MsgType      int32  `json:"msg_type"`
	Content      string `json:"content"`
	FileUrl      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	FileName     string `json:"file_name"`
	SessionId    int64  `json:"session_id"`
	CreatedAt    int64  `json:"created_at"`
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	now := time.Now()
	msgID := now.UnixNano()
	payload := messagePayload{
		MsgId:        msgID,
		SenderId:     in.SenderId,
		ReceiverId:   in.ReceiverId,
		ChatType:     in.ChatType,
		MsgType:      in.MsgType,
		Content:      in.Content,
		FileUrl:      in.FileUrl,
		FileSize:     in.FileSize,
		FileName:     in.FileName,
		SessionId:    in.SessionId,
		CreatedAt:    now.Unix(),
		SenderName:   in.SenderName,
		SenderAvatar: in.SenderAvatar,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logx.Errorf("marshal message payload error: %v", err)
		return &pb.SendMessageResponse{Code: -1, Message: "消息打包失败"}, nil
	}

	if err := l.svcCtx.KafkaWriter.WriteMessages(l.ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", in.SessionId)),
		Value: body,
		Time:  now,
	}); err != nil {
		logx.Errorf("kafka write error: sender=%d receiver=%d session=%d err=%v", in.SenderId, in.ReceiverId, in.SessionId, err)
		return &pb.SendMessageResponse{Code: -1, Message: "消息发送失败"}, nil
	}

	logx.Infof("message enqueued: sender=%d receiver=%d session=%d msg_id=%d", in.SenderId, in.ReceiverId, in.SessionId, msgID)
	return &pb.SendMessageResponse{Code: 0, Message: "ok", MsgId: msgID}, nil
}
