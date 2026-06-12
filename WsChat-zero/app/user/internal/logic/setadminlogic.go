package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAdminLogic {
	return &SetAdminLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetAdminLogic) SetAdmin(in *pb.SetAdminRequest) (*pb.CommonResponse, error) {
	if err := l.svcCtx.DB.Model(&model.UserInfo{}).Where("id = ?", in.UserId).Update("role", in.Role).Error; err != nil {
		logx.Errorf("set admin error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "操作失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "设置成功"}, nil
}
