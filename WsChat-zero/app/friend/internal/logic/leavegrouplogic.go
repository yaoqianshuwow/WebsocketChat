package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LeaveGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLeaveGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LeaveGroupLogic {
	return &LeaveGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LeaveGroupLogic) LeaveGroup(in *pb.LeaveGroupRequest) (*pb.CommonResponse, error) {
	var group model.GroupInfo
	if err := l.svcCtx.DB.First(&group, in.GroupId).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "群组不存在"}, nil
	}

	if group.OwnerId == in.UserId {
		return &pb.CommonResponse{Code: 1, Message: "群主不能退出群组，请转让或解散"}, nil
	}

	result := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).Delete(&model.GroupMember{})
	if result.Error != nil {
		return &pb.CommonResponse{Code: 1, Message: "退出失败"}, nil
	}
	if result.RowsAffected == 0 {
		return &pb.CommonResponse{Code: 1, Message: "不在该群组中"}, nil
	}

	// 更新成员数
	l.svcCtx.DB.Model(&group).UpdateColumn("member_count", group.MemberCount-1)
	return &pb.CommonResponse{Code: 0, Message: "已退出群组"}, nil
}
