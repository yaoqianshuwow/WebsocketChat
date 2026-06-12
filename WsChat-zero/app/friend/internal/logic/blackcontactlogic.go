package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BlackContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlackContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlackContactLogic {
	return &BlackContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BlackContactLogic) BlackContact(in *pb.BlackContactRequest) (*pb.CommonResponse, error) {
	var contact model.Contact
	result := l.svcCtx.DB.Where("user_id = ? AND contact_id = ?", in.UserId, in.ContactId).First(&contact)
	if result.Error != nil {
		return &pb.CommonResponse{Code: 1, Message: "联系人不存在"}, nil
	}
	contact.Status = 1 // 黑名单
	if err := l.svcCtx.DB.Save(&contact).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "拉黑失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "已拉黑"}, nil
}
