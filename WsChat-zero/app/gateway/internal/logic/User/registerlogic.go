package User

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	rpcResp, err := l.svcCtx.UserClient.Register(l.ctx, &userpb.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Phone:    req.Phone,
		Nickname: req.Nickname,
	})
	if err != nil || rpcResp.Code != 0 {
		return &types.RegisterResp{
			Code:    rpcResp.GetCode(),
			Message: rpcResp.GetMessage(),
		}, nil
	}

	return &types.RegisterResp{
		Code:    0,
		Message: "注册成功",
		UserId:  rpcResp.UserId,
	}, nil
}
