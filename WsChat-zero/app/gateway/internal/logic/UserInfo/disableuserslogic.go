package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisableUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDisableUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisableUsersLogic {
	return &DisableUsersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DisableUsersLogic) DisableUsers(req *types.IdsReq) (resp *types.CommonResp, err error) {
	r, e := l.svcCtx.UserClient.DisableUsers(l.ctx, &userpb.DisableUsersRequest{UserIds: req.Ids})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: e.Error()}, nil
	}
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
