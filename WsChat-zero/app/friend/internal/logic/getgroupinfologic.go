package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupInfoLogic) GetGroupInfo(in *pb.GetGroupInfoRequest) (*pb.GroupInfoResponse, error) {
	var group model.GroupInfo
	result := l.svcCtx.DB.First(&group, in.GroupId)
	if result.Error != nil {
		return &pb.GroupInfoResponse{Code: 1, Message: "群组不存在"}, nil
	}
	return &pb.GroupInfoResponse{Code: 0, Message: "ok", Data: groupModelToProto(&group)}, nil
}
