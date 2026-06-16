package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionMessagesLogic {
	return &DeleteSessionMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionMessagesLogic) DeleteSessionMessages(in *pb.GetRecentMessagesRequest) (*pb.SendMessageResponse, error) {
	result := l.svcCtx.DB.Model(&model.Message{}).
		Where("session_id = ? AND status = 0", in.SessionId).
		Update("status", 1)
	if result.Error != nil {
		return &pb.SendMessageResponse{Code: 1, Message: "删除消息失败"}, nil
	}
	return &pb.SendMessageResponse{Code: 0, Message: "ok"}, nil
}
