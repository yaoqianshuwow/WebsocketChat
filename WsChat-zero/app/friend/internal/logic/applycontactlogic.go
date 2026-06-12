package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyContactLogic {
	return &ApplyContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyContactLogic) ApplyContact(in *pb.ApplyContactRequest) (*pb.CommonResponse, error) {
	// 检查是否已经是好友
	var existing model.Contact
	if l.svcCtx.DB.Where("user_id = ? AND contact_id = ?", in.FromId, in.ToId).First(&existing).Error == nil {
		return &pb.CommonResponse{Code: 1, Message: "已经是好友"}, nil
	}

	// 检查是否已存在待处理的申请
	var existingApply model.ContactApply
	if l.svcCtx.DB.Where("from_id = ? AND to_id = ? AND status = 0", in.FromId, in.ToId).First(&existingApply).Error == nil {
		return &pb.CommonResponse{Code: 1, Message: "已发送过申请，请等待处理"}, nil
	}

	apply := model.ContactApply{
		FromId:    in.FromId,
		ToId:      in.ToId,
		ApplyType: 1,
		Remark:    in.Remark,
		Status:    0,
	}
	if err := l.svcCtx.DB.Create(&apply).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "申请失败:" + err.Error()}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "申请已发送"}, nil
}
