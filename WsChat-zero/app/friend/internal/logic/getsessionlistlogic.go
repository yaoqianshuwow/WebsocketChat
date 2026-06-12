package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionListLogic) GetSessionList(in *pb.GetSessionListRequest) (*pb.SessionListResponse, error) {
	var sessions []model.Session
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)
	if in.SessionType > 0 {
		query = query.Where("session_type = ?", in.SessionType)
	}
	result := query.Order("updated_at DESC").Find(&sessions)
	if result.Error != nil {
		return &pb.SessionListResponse{Code: 1, Message: "查询失败"}, nil
	}

	var data []*pb.Session
	for i := range sessions {
		data = append(data, sessionModelToProto(&sessions[i]))
	}
	return &pb.SessionListResponse{Code: 0, Message: "ok", Data: data}, nil
}
