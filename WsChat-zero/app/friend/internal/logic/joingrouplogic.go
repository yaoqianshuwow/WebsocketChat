package logic

import (
	"context"
	"time"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type JoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *JoinGroupLogic) JoinGroup(in *pb.DismissGroupRequest) (*pb.CommonResponse, error) {
	var group model.GroupInfo
	if err := l.svcCtx.DB.Where("id = ? AND status = 0", in.GroupId).First(&group).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "群组不存在或已解散"}, nil
	}

	var member model.GroupMember
	if err := l.svcCtx.DB.Where("group_id = ? AND user_id = ?", in.GroupId, in.UserId).First(&member).Error; err == nil {
		return &pb.CommonResponse{Code: 1, Message: "已在群组中"}, nil
	}

	now := time.Now()
	newMember := model.GroupMember{
		GroupId:   in.GroupId,
		UserId:    in.UserId,
		Role:      0,
		JoinedAt:  now,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := l.svcCtx.DB.Create(&newMember).Error; err != nil {
		return &pb.CommonResponse{Code: -1, Message: "加入失败"}, nil
	}

	if err := l.svcCtx.DB.Model(&model.GroupInfo{}).
		Where("id = ?", in.GroupId).
		UpdateColumn("member_count", gorm.Expr("member_count + 1")).Error; err != nil {
		logx.Errorf("update group member count failed: groupId=%d err=%v", in.GroupId, err)
	}

	return &pb.CommonResponse{Code: 0, Message: "ok"}, nil
}
