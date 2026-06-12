package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type PassContactApplyLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewPassContactApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PassContactApplyLogic {
	return &PassContactApplyLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PassContactApplyLogic) PassContactApply(req *types.PassContactApplyReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.PassContactApply(l.ctx, &friendpb.PassContactApplyRequest{ApplyId: req.ApplyId, UserId: uid, Status: req.Status})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
