package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/svc"
	"github.com/your-org/ws-chat-zero/app/msg-forward/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMessagesLogic {
	return &SearchMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchMessagesLogic) SearchMessages(in *pb.SearchMessagesRequest) (*pb.MessageListResponse, error) {
	var messages []model.Message
	query := l.svcCtx.DB.Where("status = 0")

	if in.Keyword != "" {
		query = query.Where("content LIKE ?", "%"+in.Keyword+"%")
	}
	if in.SenderId > 0 {
		query = query.Where("sender_id = ?", in.SenderId)
	}
	if in.SessionId > 0 {
		query = query.Where("session_id = ?", in.SessionId)
	}
	if in.StartTime > 0 {
		query = query.Where("created_at >= ?", in.StartTime)
	}
	if in.EndTime > 0 {
		query = query.Where("created_at <= ?", in.EndTime)
	}

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	var total int64
	query.Model(&model.Message{}).Count(&total)

	offset := int((page - 1) * pageSize)
	result := query.Order("id DESC").Offset(offset).Limit(int(pageSize)).Find(&messages)
	if result.Error != nil {
		return &pb.MessageListResponse{Code: 1, Message: "搜索失败"}, nil
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
