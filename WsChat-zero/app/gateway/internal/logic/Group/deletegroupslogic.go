package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteGroupsLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewDeleteGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteGroupsLogic {
	return &DeleteGroupsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteGroupsLogic) DeleteGroups(req *types.IdsReq) (resp *types.CommonResp, err error) {
	r, e := l.svcCtx.FriendClient.DeleteGroups(l.ctx, &friendpb.DeleteGroupsRequest{GroupIds: req.Ids})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
