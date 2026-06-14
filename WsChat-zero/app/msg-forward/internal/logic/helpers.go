package logic

import (
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"
)

func messageModelToProto(m *model.Message) *pb.Message {
	return &pb.Message{
		Id:           m.Id,
		SenderId:     m.SenderId,
		ReceiverId:   m.ReceiverId,
		ChatType:     m.ChatType,
		MsgType:      m.MsgType,
		Content:      m.Content,
		FileUrl:      m.FileUrl,
		FileSize:     m.FileSize,
		FileName:     m.FileName,
		Status:       m.Status,
		SessionId:    m.SessionId,
		CreatedAt:    m.CreatedAt.Unix(),
		SenderName:   m.SendName,
		SenderAvatar: m.SendAvatar,
	}
}
