package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetAdminLogic {
	return &SetAdminLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetAdminLogic) SetAdmin(req *types.IdsReq) (resp *types.CommonResp, err error) {
	if len(req.Ids) == 0 {
		return &types.CommonResp{Code: 1, Message: "参数错误"}, nil
	}
	r, e := l.svcCtx.UserClient.SetAdmin(l.ctx, &userpb.SetAdminRequest{
		UserId: req.Ids[0], Role: 1,
	})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: e.Error()}, nil
	}
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
