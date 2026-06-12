package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbleUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbleUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbleUsersLogic {
	return &AbleUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AbleUsersLogic) AbleUsers(in *pb.AbleUsersRequest) (*pb.CommonResponse, error) {
	if err := l.svcCtx.DB.Model(&model.UserInfo{}).Where("id IN ?", in.UserIds).Update("status", 0).Error; err != nil {
		logx.Errorf("able users error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "操作失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "启用成功"}, nil
}
