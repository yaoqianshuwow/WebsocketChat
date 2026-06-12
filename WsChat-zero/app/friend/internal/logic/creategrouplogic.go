package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateGroupLogic) CreateGroup(in *pb.CreateGroupRequest) (*pb.GroupInfoResponse, error) {
	now := time.Now()
	group := model.GroupInfo{
		Name:      in.GroupName,
		OwnerId:   in.OwnerId,
		AddMode:   2, // 默认直接加入
		Status:    0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx := l.svcCtx.DB.Begin()

	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		return &pb.GroupInfoResponse{Code: 1, Message: "创建群组失败"}, nil
	}

	// 添加群主为管理员
	ownerMember := model.GroupMember{
		GroupId:  group.Id,
		UserId:   in.OwnerId,
		Role:     2, // 群主
		Nickname: "",
		JoinedAt: now,
	}
	if err := tx.Create(&ownerMember).Error; err != nil {
		tx.Rollback()
		return &pb.GroupInfoResponse{Code: 1, Message: "添加群主失败"}, nil
	}

	// 添加其他成员
	for _, memberId := range in.MemberIds {
		member := model.GroupMember{
			GroupId:  group.Id,
			UserId:   memberId,
			Role:     0, // 普通成员
			Nickname: "",
			JoinedAt: now,
		}
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return &pb.GroupInfoResponse{Code: 1, Message: "添加成员失败"}, nil
		}
	}

	// 更新成员数
	group.MemberCount = int32(1 + len(in.MemberIds))
	tx.Model(&group).Update("member_count", group.MemberCount)

	tx.Commit()
	return &pb.GroupInfoResponse{Code: 0, Message: "ok", Data: groupModelToProto(&group)}, nil
}
