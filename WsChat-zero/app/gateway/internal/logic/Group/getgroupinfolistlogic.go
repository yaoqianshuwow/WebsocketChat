package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetGroupInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoListLogic {
	return &GetGroupInfoListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupInfoListLogic) GetGroupInfoList() (resp *types.GroupListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetGroupInfoList(l.ctx, &friendpb.GetContactListRequest{UserId: uid})
	if e != nil || r.Code != 0 { return &types.GroupListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	data := make([]types.GroupInfoResp, 0)
	for _, g := range r.Data { data = append(data, types.GroupInfoResp{GroupId: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId, MemberCount: g.MemberCount, AddMode: g.AddMode, Status: g.Status, Notice: g.Notice}) }
	return &types.GroupListResp{Code: 0, Message: "ok", Data: data}, nil
}
