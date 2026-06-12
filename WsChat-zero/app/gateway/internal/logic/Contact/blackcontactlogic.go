package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type BlackContactLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewBlackContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlackContactLogic {
	return &BlackContactLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *BlackContactLogic) BlackContact(req *types.ContactInfoReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.BlackContact(l.ctx, &friendpb.BlackContactRequest{UserId: uid, ContactId: req.ContactId})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
