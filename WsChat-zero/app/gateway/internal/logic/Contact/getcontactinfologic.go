package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactInfoLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetContactInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactInfoLogic {
	return &GetContactInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetContactInfoLogic) GetContactInfo(req *types.ContactInfoReq) (resp *types.ContactListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetContactInfo(l.ctx, &friendpb.GetContactInfoRequest{UserId: uid, ContactId: req.ContactId})
	if e != nil || r.Code != 0 { return &types.ContactListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	c := r.Data
	return &types.ContactListResp{Code: 0, Message: "ok", Data: []types.ContactVo{{ContactId: c.ContactId, ContactType: c.ContactType, Status: c.Status}}}, nil
}
