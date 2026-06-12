package User

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsCodeLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewSendSmsCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsCodeLogic {
	return &SendSmsCodeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SendSmsCodeLogic) SendSmsCode(req *types.SmsCodeReq) (resp *types.CommonResp, err error) {
	r, e := l.svcCtx.UserClient.SendSmsCode(l.ctx, &userpb.SendSmsCodeRequest{Phone: req.Phone})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
