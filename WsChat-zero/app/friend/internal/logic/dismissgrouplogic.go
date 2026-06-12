package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DismissGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDismissGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DismissGroupLogic {
	return &DismissGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DismissGroupLogic) DismissGroup(in *pb.DismissGroupRequest) (*pb.CommonResponse, error) {
	var group model.GroupInfo
	if err := l.svcCtx.DB.First(&group, in.GroupId).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "群组不存在"}, nil
	}

	if group.OwnerId != in.UserId {
		return &pb.CommonResponse{Code: 1, Message: "只有群主才能解散群组"}, nil
	}

	tx := l.svcCtx.DB.Begin()
	// 删除所有成员
	if err := tx.Where("group_id = ?", in.GroupId).Delete(&model.GroupMember{}).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "解散失败"}, nil
	}
	// 删除群组
	if err := tx.Delete(&group).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "解散失败"}, nil
	}
	tx.Commit()
	return &pb.CommonResponse{Code: 0, Message: "群组已解散"}, nil
}
