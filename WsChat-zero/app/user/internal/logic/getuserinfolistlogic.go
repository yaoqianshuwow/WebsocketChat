package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoListLogic {
	return &GetUserInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoListLogic) GetUserInfoList(in *pb.GetUserInfoListRequest) (*pb.UserInfoListResponse, error) {
	var users []model.UserInfo
	l.svcCtx.DB.Where("id IN ?", in.UserIds).Find(&users)

	var data []*pb.UserInfo
	for i := range users {
		data = append(data, modelToProto(&users[i]))
	}
	return &pb.UserInfoListResponse{Code: 0, Message: "ok", Data: data}, nil
}
