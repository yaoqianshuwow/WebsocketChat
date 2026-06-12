package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteContactLogic {
	return &DeleteContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteContactLogic) DeleteContact(in *pb.DeleteContactRequest) (*pb.CommonResponse, error) {
	// 双向删除好友关系
	tx := l.svcCtx.DB.Begin()
	if err := tx.Where("user_id = ? AND contact_id = ?", in.UserId, in.ContactId).Delete(&model.Contact{}).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}
	if err := tx.Where("user_id = ? AND contact_id = ?", in.ContactId, in.UserId).Delete(&model.Contact{}).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}
	tx.Commit()
	return &pb.CommonResponse{Code: 0, Message: "删除成功"}, nil
}
