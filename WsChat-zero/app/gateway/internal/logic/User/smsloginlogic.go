package User

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type SmsLoginLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewSmsLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsLoginLogic {
	return &SmsLoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SmsLoginLogic) SmsLogin(req *types.SmsLoginReq) (resp *types.LoginResp, err error) {
	r, e := l.svcCtx.UserClient.SmsLogin(l.ctx, &userpb.SmsLoginRequest{Phone: req.Phone, Code: req.Code})
	if e != nil || r.Code != 0 { return &types.LoginResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	return &types.LoginResp{Code: 0, Message: "ok", Token: r.Token, UserId: r.UserInfo.GetId(), Nickname: r.UserInfo.GetNickname(), Avatar: r.UserInfo.GetAvatar()}, nil
}
