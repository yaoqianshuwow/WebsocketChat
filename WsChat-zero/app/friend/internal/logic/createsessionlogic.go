package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
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
	session := model.Session{
		UserId:      in.UserId,
		PeerId:      in.PeerId,
		SessionType: in.SessionType,
		SessionName: in.SessionName,
	}
	if err := l.svcCtx.DB.Create(&session).Error; err != nil {
		return &pb.SessionResponse{Code: 1, Message: "创建会话失败:" + err.Error()}, nil
	}
	return &pb.SessionResponse{Code: 0, Message: "ok", Data: sessionModelToProto(&session)}, nil
}
