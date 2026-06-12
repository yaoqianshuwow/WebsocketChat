package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyContactLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewApplyContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyContactLogic {
	return &ApplyContactLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *ApplyContactLogic) ApplyContact(req *types.ApplyContactReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.ApplyContact(l.ctx, &friendpb.ApplyContactRequest{FromId: uid, ToId: req.ToId, Remark: req.Remark})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
