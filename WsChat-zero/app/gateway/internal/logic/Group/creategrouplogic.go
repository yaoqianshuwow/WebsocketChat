package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateGroupLogic) CreateGroup(req *types.CreateGroupReq) (resp *types.GroupInfoResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.CreateGroup(l.ctx, &friendpb.CreateGroupRequest{GroupName: req.GroupName, OwnerId: uid, MemberIds: req.MemberIds})
	if e != nil || r.Code != 0 { return &types.GroupInfoResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	g := r.Data
	return &types.GroupInfoResp{Code: 0, Message: "ok", GroupId: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId, MemberCount: g.MemberCount, AddMode: g.AddMode, Status: g.Status, Notice: g.Notice}, nil
}
