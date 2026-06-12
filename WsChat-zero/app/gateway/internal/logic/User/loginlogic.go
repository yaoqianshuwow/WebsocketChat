package User

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	rpcResp, err := l.svcCtx.UserClient.Login(l.ctx, &userpb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil || rpcResp.Code != 0 {
		return &types.LoginResp{
			Code:    rpcResp.GetCode(),
			Message: rpcResp.GetMessage(),
		}, nil
	}

	return &types.LoginResp{
		Code:     0,
		Message:  "登录成功",
		Token:    rpcResp.Token,
		UserId:   rpcResp.UserInfo.GetId(),
		Nickname: rpcResp.UserInfo.GetNickname(),
		Avatar:   rpcResp.UserInfo.GetAvatar(),
	}, nil
}
