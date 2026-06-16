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

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.Size
	if pageSize <= 0 {
		pageSize = 20
	}

	var messages []model.Message
	var total int64

	if in.StartTime > 0 || in.EndTime > 0 {
		if err := query.Order("id DESC").Find(&messages).Error; err != nil {
			return &pb.MessageListResponse{Code: 1, Message: "搜索失败"}, nil
		}

		filtered := make([]model.Message, 0, len(messages))
		for _, msg := range messages {
			ts := msg.CreatedAt.Unix()
			if in.StartTime > 0 && ts < in.StartTime {
				continue
			}
			if in.EndTime > 0 && ts > in.EndTime {
				continue
			}
			filtered = append(filtered, msg)
		}
		messages = filtered
		total = int64(len(messages))
	} else {
		query.Model(&model.Message{}).Count(&total)

		offset := int((page - 1) * pageSize)
		result := query.Order("id DESC").Offset(offset).Limit(int(pageSize)).Find(&messages)
		if result.Error != nil {
			return &pb.MessageListResponse{Code: 1, Message: "搜索失败"}, nil
		}
	}

	start := int((page - 1) * pageSize)
	if start > len(messages) {
		start = len(messages)
	}
	end := start + int(pageSize)
	if end > len(messages) {
		end = len(messages)
	}
	pageMessages := messages[start:end]

	var data []*pb.Message
	for i := len(pageMessages) - 1; i >= 0; i-- {
		data = append(data, messageModelToProto(&pageMessages[i]))
	}

	return &pb.MessageListResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
		Total:   total,
	}, nil
}
