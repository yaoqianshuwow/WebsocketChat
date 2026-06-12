package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetContactInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactInfoLogic {
	return &GetContactInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetContactInfoLogic) GetContactInfo(in *pb.GetContactInfoRequest) (*pb.ContactInfoResponse, error) {
	var contact model.Contact
	result := l.svcCtx.DB.Where("user_id = ? AND contact_id = ?", in.UserId, in.ContactId).First(&contact)
	if result.Error != nil {
		return &pb.ContactInfoResponse{Code: 1, Message: "联系人不存在"}, nil
	}
	return &pb.ContactInfoResponse{Code: 0, Message: "ok", Data: contactModelToProto(&contact)}, nil
}
