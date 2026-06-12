package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisableUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableUsersLogic {
	return &DisableUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DisableUsersLogic) DisableUsers(in *pb.DisableUsersRequest) (*pb.CommonResponse, error) {
	if err := l.svcCtx.DB.Model(&model.UserInfo{}).Where("id IN ?", in.UserIds).Update("status", 1).Error; err != nil {
		logx.Errorf("disable users error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "操作失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "禁用成功"}, nil
}
