package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMessageListLogic) GetMessageList(in *pb.GetMessageListRequest) (*pb.MessageListResponse, error) {
	var messages []model.Message
	query := l.svcCtx.DB.Where("session_id = ? AND status = 0", in.SessionId)

	if in.BeforeId > 0 {
		query = query.Where("id < ?", in.BeforeId)
	}

	pageSize := in.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	query.Model(&model.Message{}).Count(&total)
	if in.BeforeId <= 0 {
		pageSize = total
	}

	result := query.Order("id DESC").Limit(int(pageSize)).Find(&messages)
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
		Total:   total,
	}, nil
}
