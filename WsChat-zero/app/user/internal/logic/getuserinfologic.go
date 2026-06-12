package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoRequest) (*pb.UserInfoResponse, error) {
	var user model.UserInfo
	if l.svcCtx.DB.First(&user, in.UserId).Error != nil {
		return &pb.UserInfoResponse{Code: 1, Message: "用户不存在"}, nil
	}
	return &pb.UserInfoResponse{
		Code:    0,
		Message: "ok",
		Data:    modelToProto(&user),
	}, nil
}
