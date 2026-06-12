package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUsersLogic {
	return &DeleteUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUsersLogic) DeleteUsers(in *pb.DeleteUsersRequest) (*pb.CommonResponse, error) {
	if err := l.svcCtx.DB.Where("id IN ?", in.UserIds).Delete(&model.UserInfo{}).Error; err != nil {
		logx.Errorf("delete users error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "删除失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "删除成功"}, nil
}
