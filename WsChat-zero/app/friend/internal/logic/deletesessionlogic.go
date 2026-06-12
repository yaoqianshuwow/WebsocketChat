package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSessionLogic) DeleteSession(in *pb.DeleteSessionRequest) (*pb.CommonResponse, error) {
	result := l.svcCtx.DB.Where("id = ? AND user_id = ?", in.SessionId, in.UserId).Delete(&model.Session{})
	if result.Error != nil {
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}
	if result.RowsAffected == 0 {
		return &pb.CommonResponse{Code: 1, Message: "会话不存在"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "删除成功"}, nil
}
