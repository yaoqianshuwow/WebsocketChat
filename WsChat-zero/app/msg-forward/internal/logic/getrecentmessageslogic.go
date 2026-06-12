package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRecentMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecentMessagesLogic {
	return &GetRecentMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRecentMessagesLogic) GetRecentMessages(in *pb.GetRecentMessagesRequest) (*pb.MessageListResponse, error) {
	limit := in.Limit
	if limit <= 0 {
		limit = 50
	}

	var messages []model.Message
	result := l.svcCtx.DB.Where("session_id = ? AND status = 0", in.SessionId).
		Order("id DESC").
		Limit(int(limit)).
		Find(&messages)
	if result.Error != nil {
		return &pb.MessageListResponse{Code: 1, Message: "查询失败"}, nil
	}

	var data []*pb.Message
	for i := len(messages) - 1; i >= 0; i-- {
		data = append(data, messageModelToProto(&messages[i]))
	}

	return &pb.MessageListResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
		Total:   0,
	}, nil
}
