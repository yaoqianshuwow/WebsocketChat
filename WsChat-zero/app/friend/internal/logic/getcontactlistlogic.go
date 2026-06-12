package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetContactListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactListLogic {
	return &GetContactListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetContactListLogic) GetContactList(in *pb.GetContactListRequest) (*pb.ContactListResponse, error) {
	var contacts []model.Contact
	result := l.svcCtx.DB.Where("user_id = ? AND status = 0", in.UserId).Find(&contacts)
	if result.Error != nil {
		return &pb.ContactListResponse{Code: 1, Message: "查询失败:" + result.Error.Error()}, nil
	}

	var data []*pb.ContactInfo
	for i := range contacts {
		data = append(data, contactModelToProto(&contacts[i]))
	}
	return &pb.ContactListResponse{Code: 0, Message: "ok", Data: data}, nil
}
