package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PassContactApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPassContactApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassContactApplyLogic {
	return &PassContactApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PassContactApplyLogic) PassContactApply(in *pb.PassContactApplyRequest) (*pb.CommonResponse, error) {
	var apply model.ContactApply
	if err := l.svcCtx.DB.Where("id = ?", in.ApplyId).First(&apply).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "申请记录不存在"}, nil
	}

	if apply.Status != 0 {
		return &pb.CommonResponse{Code: 1, Message: "申请已处理"}, nil
	}

	apply.Status = in.Status

	tx := l.svcCtx.DB.Begin()

	// 更新申请状态
	if err := tx.Save(&apply).Error; err != nil {
		tx.Rollback()
		return &pb.CommonResponse{Code: 1, Message: "操作失败"}, nil
	}

	if in.Status == 1 { // 同意
		// 双方添加为好友
		contact1 := model.Contact{
			UserId:      apply.FromId,
			ContactId:   apply.ToId,
			ContactType: 1,
			Status:      0,
		}
		contact2 := model.Contact{
			UserId:      apply.ToId,
			ContactId:   apply.FromId,
			ContactType: 1,
			Status:      0,
		}
		if err := tx.Create(&contact1).Error; err != nil {
			tx.Rollback()
			return &pb.CommonResponse{Code: 1, Message: "添加好友失败"}, nil
		}
		if err := tx.Create(&contact2).Error; err != nil {
			tx.Rollback()
			return &pb.CommonResponse{Code: 1, Message: "添加好友失败"}, nil
		}
	}

	tx.Commit()
	return &pb.CommonResponse{Code: 0, Message: "操作成功"}, nil
}
