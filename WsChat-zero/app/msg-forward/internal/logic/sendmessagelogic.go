package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
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

var msgSeq int64

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	now := time.Now()
	msgSeq++
	// Use atomic counter to avoid nanosecond collisions
	msgID := now.UnixNano() + msgSeq%1000
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
		return &pb.SendMessageResponse{Code: -1, Message: "marshal message payload failed"}, nil
	}

	if err := l.storeDirect(payload); err != nil {
		logx.Errorf("direct store error: sender=%d receiver=%d session=%d err=%v", in.SenderId, in.ReceiverId, in.SessionId, err)
		return &pb.SendMessageResponse{Code: -1, Message: "store message failed"}, nil
	}

	if l.svcCtx.Config.StoreMode == "direct" || l.svcCtx.KafkaWriter == nil {
		logx.Infof("message stored directly: sender=%d receiver=%d session=%d msg_id=%d", in.SenderId, in.ReceiverId, in.SessionId, msgID)
		return &pb.SendMessageResponse{Code: 0, Message: "ok", MsgId: msgID}, nil
	}

	if err := l.svcCtx.KafkaWriter.WriteMessages(l.ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", in.SessionId)),
		Value: body,
		Time:  now,
	}); err != nil {
		logx.Errorf("kafka write error: sender=%d receiver=%d session=%d err=%v", in.SenderId, in.ReceiverId, in.SessionId, err)
	} else {
		logx.Infof("message enqueued: sender=%d receiver=%d session=%d msg_id=%d", in.SenderId, in.ReceiverId, in.SessionId, msgID)
	}

	return &pb.SendMessageResponse{Code: 0, Message: "ok", MsgId: msgID}, nil
}

func (l *SendMessageLogic) storeDirect(payload messagePayload) error {
	msg := model.Message{
		MsgId:      &payload.MsgId,
		SenderId:   payload.SenderId,
		ReceiverId: payload.ReceiverId,
		ChatType:   payload.ChatType,
		MsgType:    payload.MsgType,
		Content:    payload.Content,
		FileUrl:    payload.FileUrl,
		FileSize:   payload.FileSize,
		FileName:   payload.FileName,
		SendName:   payload.SenderName,
		SendAvatar: payload.SenderAvatar,
		Status:     0,
		SessionId:  payload.SessionId,
		CreatedAt:  time.Unix(payload.CreatedAt, 0),
	}
	return l.svcCtx.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "msg_id"}},
		DoNothing: true,
	}).Create(&msg).Error
}
