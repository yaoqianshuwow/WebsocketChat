package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

type sessionLite struct {
	Id          int64 `gorm:"column:id"`
	UserId      int64 `gorm:"column:user_id"`
	UnreadCount int32 `gorm:"column:unread_count"`
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *SendMessageLogic) SendMessage(in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	now := time.Now()
	msg := model.Message{
		Uuid:       fmt.Sprintf("M%019d", now.UnixNano()),
		Type:       in.MsgType,
		Content:    in.Content,
		Url:        in.FileUrl,
		SendID:     formatID(in.SenderId),
		SendName:   formatID(in.SenderId),
		SendAvatar: "",
		ReceiveID:  formatID(in.ReceiverId),
		FileType:   "",
		FileName:   in.FileName,
		FileSize:   in.FileSize,
		Status:     0,
		CreatedAt:  now,
		SenderId:   in.SenderId,
		ReceiverId: in.ReceiverId,
		ChatType:   in.ChatType,
		MsgType:    in.MsgType,
		FileUrl:    in.FileUrl,
		SessionId:  in.SessionId,
	}

	if err := l.svcCtx.DB.Create(&msg).Error; err != nil {
		logx.Errorf("message insert error: %v", err)
		return &pb.SendMessageResponse{Code: -1, Message: "消息发送失败"}, nil
	}

	if err := l.updateSessionPreview(in.SessionId, in.SenderId, now.Unix(), in.Content); err != nil {
		logx.Errorf("update session preview error: session=%d err=%v", in.SessionId, err)
	}

	logx.Infof("message stored directly: sender=%d receiver=%d session=%d id=%d", in.SenderId, in.ReceiverId, in.SessionId, msg.Id)
	return &pb.SendMessageResponse{Code: 0, Message: "ok", MsgId: msg.Id}, nil
}

func (l *SendMessageLogic) updateSessionPreview(sessionId, senderId, sentAt int64, content string) error {
	var session sessionLite
	if err := l.svcCtx.DB.Table("session").Where("id = ?", sessionId).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	update := map[string]interface{}{
		"last_msg_content": content,
		"last_msg_time":    sentAt,
	}
	if session.UserId != senderId {
		update["unread_count"] = session.UnreadCount + 1
	}

	return l.svcCtx.DB.Table("session").Where("id = ?", sessionId).Updates(update).Error
}

func formatID(v int64) string {
	return fmt.Sprintf("%d", v)
}
