package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateGroupInfoLogic) UpdateGroupInfo(in *pb.UpdateGroupInfoRequest) (*pb.CommonResponse, error) {
	var group model.GroupInfo
	if err := l.svcCtx.DB.First(&group, in.GroupId).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "群组不存在"}, nil
	}

	updates := map[string]interface{}{}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Avatar != "" {
		updates["avatar"] = in.Avatar
	}
	if in.Notice != "" {
		updates["notice"] = in.Notice
	}
	if in.AddMode > 0 {
		updates["add_mode"] = in.AddMode
	}

	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(&group).Updates(updates).Error; err != nil {
			return &pb.CommonResponse{Code: 1, Message: "更新失败"}, nil
		}
	}
	return &pb.CommonResponse{Code: 0, Message: "更新成功"}, nil
}
