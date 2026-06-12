package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateSessionLogic) CreateSession(in *pb.CreateSessionRequest) (*pb.SessionResponse, error) {
	var session model.Session
	err := l.svcCtx.DB.Where("user_id = ? AND peer_id = ? AND session_type = ?", in.UserId, in.PeerId, in.SessionType).
		First(&session).Error
	if err == nil {
		if in.SessionName != "" && in.SessionName != session.SessionName {
			if err := l.svcCtx.DB.Model(&model.Session{}).
				Where("id = ?", session.Id).
				Update("session_name", in.SessionName).Error; err != nil {
				return &pb.SessionResponse{Code: 1, Message: "更新会话失败:" + err.Error()}, nil
			}
			session.SessionName = in.SessionName
		}
		return &pb.SessionResponse{Code: 0, Message: "ok", Data: sessionModelToProto(&session)}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return &pb.SessionResponse{Code: 1, Message: "查询会话失败:" + err.Error()}, nil
	}

	session = model.Session{
		UserId:      in.UserId,
		PeerId:      in.PeerId,
		SessionType: in.SessionType,
		SessionName: in.SessionName,
	}
	if err := l.svcCtx.DB.Create(&session).Error; err != nil {
		dup := model.Session{}
		if err := l.svcCtx.DB.Where("user_id = ? AND peer_id = ? AND session_type = ?", in.UserId, in.PeerId, in.SessionType).
			First(&dup).Error; err == nil {
			return &pb.SessionResponse{Code: 0, Message: "ok", Data: sessionModelToProto(&dup)}, nil
		}
		return &pb.SessionResponse{Code: 1, Message: "创建会话失败:" + err.Error()}, nil
	}
	return &pb.SessionResponse{Code: 0, Message: "ok", Data: sessionModelToProto(&session)}, nil
}
