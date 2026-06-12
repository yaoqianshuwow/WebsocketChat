package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *pb.UpdateUserInfoRequest) (*pb.CommonResponse, error) {
	updates := map[string]interface{}{}
	if in.Nickname != "" {
		updates["nickname"] = in.Nickname
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}
	if in.Sex != "" {
		updates["sex"] = in.Sex
	}
	if in.Age > 0 {
		updates["age"] = in.Age
	}
	if in.Bio != "" {
		updates["bio"] = in.Bio
	}

	if len(updates) == 0 {
		return &pb.CommonResponse{Code: 1, Message: "无更新内容"}, nil
	}

	if err := l.svcCtx.DB.Model(&model.UserInfo{}).Where("id = ?", in.UserId).Updates(updates).Error; err != nil {
		logx.Errorf("update user error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "更新失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "更新成功"}, nil
}
