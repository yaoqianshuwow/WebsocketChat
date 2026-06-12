package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetContactListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactListLogic {
	return &GetContactListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetContactListLogic) GetContactList(req *types.CommonResp) (resp *types.ContactListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetContactList(l.ctx, &friendpb.GetContactListRequest{UserId: uid})
	if e != nil || r.Code != 0 {
		return &types.ContactListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}
	var data []types.ContactVo
	for _, c := range r.Data {
		data = append(data, types.ContactVo{ContactId: c.ContactId, ContactType: c.ContactType, Status: c.Status})
	}
	return &types.ContactListResp{Code: 0, Message: "ok", Data: data}, nil
}
