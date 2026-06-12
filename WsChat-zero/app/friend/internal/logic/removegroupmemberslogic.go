package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveGroupMembersLogic {
	return &RemoveGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveGroupMembersLogic) RemoveGroupMembers(in *pb.RemoveGroupMembersRequest) (*pb.CommonResponse, error) {
	var group model.GroupInfo
	if err := l.svcCtx.DB.First(&group, in.GroupId).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "群组不存在"}, nil
	}

	// 验证操作者是群主或管理员
	var operator model.GroupMember
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.OperatorId).First(&operator).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "操作者不在群组中"}, nil
	}
	if operator.Role < 1 {
		return &pb.CommonResponse{Code: 1, Message: "没有权限移除成员"}, nil
	}

	result := l.svcCtx.DB.Where("group_id = ? AND user_id IN ?", in.GroupId, in.MemberIds).Delete(&model.GroupMember{})
	if result.Error != nil {
		return &pb.CommonResponse{Code: 1, Message: "移除失败"}, nil
	}

	// 更新成员数
	l.svcCtx.DB.Model(&group).UpdateColumn("member_count", group.MemberCount-int32(len(in.MemberIds)))
	return &pb.CommonResponse{Code: 0, Message: "移除成功"}, nil
}
