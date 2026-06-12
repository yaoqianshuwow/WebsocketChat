package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupInfoLogic) GetGroupInfo(req *types.GroupIdReq) (resp *types.GroupInfoResp, err error) {
	r, e := l.svcCtx.FriendClient.GetGroupInfo(l.ctx, &friendpb.GetGroupInfoRequest{GroupId: req.GroupId})
	if e != nil || r.Code != 0 { return &types.GroupInfoResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	g := r.Data
	return &types.GroupInfoResp{Code: 0, Message: "ok", GroupId: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId, MemberCount: g.MemberCount, AddMode: g.AddMode, Status: g.Status, Notice: g.Notice}, nil
}
