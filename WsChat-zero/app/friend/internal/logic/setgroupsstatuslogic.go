package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupsStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupsStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupsStatusLogic {
	return &SetGroupsStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetGroupsStatusLogic) SetGroupsStatus(in *pb.SetGroupsStatusRequest) (*pb.CommonResponse, error) {
	if err := l.svcCtx.DB.Model(&model.GroupInfo{}).Where("id IN ?", in.GroupIds).Update("status", in.Status).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "操作失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "操作成功"}, nil
}
