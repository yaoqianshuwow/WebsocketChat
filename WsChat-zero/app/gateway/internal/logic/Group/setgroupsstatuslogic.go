package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupsStatusLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewSetGroupsStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupsStatusLogic {
	return &SetGroupsStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetGroupsStatusLogic) SetGroupsStatus(req *types.GroupIdReq) (resp *types.CommonResp, err error) {
	r, e := l.svcCtx.FriendClient.SetGroupsStatus(l.ctx, &friendpb.SetGroupsStatusRequest{GroupIds: []int64{req.GroupId}, Status: 1})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
