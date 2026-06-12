package Contact

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetApplyListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyListLogic {
	return &GetApplyListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetApplyListLogic) GetApplyList() (resp *types.ApplyListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetApplyList(l.ctx, &friendpb.GetContactListRequest{UserId: uid})
	if e != nil || r.Code != 0 { return &types.ApplyListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	data := make([]types.ApplyVo, 0)
	for _, a := range r.Data { data = append(data, types.ApplyVo{ApplyId: a.Id, FromId: a.FromId, Remark: a.Remark, Status: a.Status}) }
	return &types.ApplyListResp{Code: 0, Message: "ok", Data: data}, nil
}
