package Group

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *JoinGroupLogic) JoinGroup(req *types.GroupIdReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.JoinGroup(l.ctx, &friendpb.DismissGroupRequest{GroupId: req.GroupId, UserId: uid})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: e.Error()}, nil
	}
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
