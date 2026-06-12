package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupsLogic {
	return &DeleteGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteGroupsLogic) DeleteGroups(in *pb.DeleteGroupsRequest) (*pb.CommonResponse, error) {
	tx := l.svcCtx.DB.Begin()
	// 删除群组
	if err := tx.Where("id IN ?", in.GroupIds).Delete(&model.GroupInfo{}).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}
	// 删除群成员
	if err := tx.Where("group_id IN ?", in.GroupIds).Delete(&model.GroupMember{}).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}
	tx.Commit()
	return &pb.CommonResponse{Code: 0, Message: "删除成功"}, nil
}
