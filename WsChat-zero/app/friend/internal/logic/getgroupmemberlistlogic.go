package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(in *pb.GetGroupMemberListRequest) (*pb.GroupMemberListResponse, error) {
	var members []model.GroupMember
	result := l.svcCtx.DB.Where("group_id = ?", in.GroupId).Order("role DESC, joined_at ASC").Find(&members)
	if result.Error != nil {
		return &pb.GroupMemberListResponse{Code: 1, Message: "查询失败"}, nil
	}

	var data []*pb.GroupMember
	for i := range members {
		data = append(data, memberModelToProto(&members[i]))
	}
	return &pb.GroupMemberListResponse{Code: 0, Message: "ok", Data: data}, nil
}
