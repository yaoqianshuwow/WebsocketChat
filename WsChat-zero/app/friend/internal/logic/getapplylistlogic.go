package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyListLogic {
	return &GetApplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApplyListLogic) GetApplyList(in *pb.GetContactListRequest) (*pb.ApplyListResponse, error) {
	var applies []model.ContactApply
	result := l.svcCtx.DB.Where("to_id = ?", in.UserId).Order("created_at DESC").Find(&applies)
	if result.Error != nil {
		return &pb.ApplyListResponse{Code: 1, Message: "查询失败"}, nil
	}

	var data []*pb.ContactApply
	for i := range applies {
		data = append(data, applyModelToProto(&applies[i]))
	}
	return &pb.ApplyListResponse{Code: 0, Message: "ok", Data: data}, nil
}
