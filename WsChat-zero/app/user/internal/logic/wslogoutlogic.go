package logic

import (
	"context"
	"fmt"

	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type WsLogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWsLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsLogoutLogic {
	return &WsLogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WsLogoutLogic) WsLogout(in *pb.WsLogoutRequest) (*pb.CommonResponse, error) {
	key := fmt.Sprintf("token:%d", in.UserId)
	if err := l.svcCtx.Redis.Del(l.ctx, key).Err(); err != nil {
		logx.Errorf("ws logout error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "退出失败"}, nil
	}
	return &pb.CommonResponse{Code: 0, Message: "退出成功"}, nil
}
