package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoListLogic {
	return &GetGroupInfoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupInfoListLogic) GetGroupInfoList(in *pb.GetContactListRequest) (*pb.GroupListResponse, error) {
	// 通过用户ID查找用户加入的所有群组
	var memberGroups []model.GroupMember
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&memberGroups).Error; err != nil {
		return &pb.GroupListResponse{Code: 1, Message: "查询失败"}, nil
	}

	if len(memberGroups) == 0 {
		return &pb.GroupListResponse{Code: 0, Message: "ok", Data: []*pb.GroupInfo{}}, nil
	}

	var groupIds []int64
	for _, mg := range memberGroups {
		groupIds = append(groupIds, mg.GroupId)
	}

	var groups []model.GroupInfo
	l.svcCtx.DB.Where("id IN ?", groupIds).Find(&groups)

	var data []*pb.GroupInfo
	for i := range groups {
		data = append(data, groupModelToProto(&groups[i]))
	}
	return &pb.GroupListResponse{Code: 0, Message: "ok", Data: data}, nil
}
