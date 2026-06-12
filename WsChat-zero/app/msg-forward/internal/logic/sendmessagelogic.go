package logic

import (
	"context"
	"encoding/json"
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

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	payload, _ := json.Marshal(map[string]interface{}{
		"msg_id":      time.Now().UnixNano(),
		"sender_id":   in.SenderId,
		"receiver_id": in.ReceiverId,
		"chat_type":   in.ChatType,
		"msg_type":    in.MsgType,
		"content":     in.Content,
		"file_url":    in.FileUrl,
		"file_size":   in.FileSize,
		"file_name":   in.FileName,
		"session_id":  in.SessionId,
		"created_at":  time.Now().Unix(),
	})

	msg := kafka.Message{Key: []byte(in.Content), Value: payload}
	if err := l.svcCtx.KafkaWriter.WriteMessages(l.ctx, msg); err != nil {
		logx.Errorf("kafka write error: %v", err)
		return &pb.SendMessageResponse{Code: -1, Message: "消息发送失败"}, nil
	}

	logx.Infof("message forwarded to kafka: sender=%d, type=%d", in.SenderId, in.MsgType)
	return &pb.SendMessageResponse{Code: 0, Message: "ok", MsgId: time.Now().UnixNano()}, nil
}
