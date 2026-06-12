package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GroupIdReq) (resp *types.GroupMemberListResp, err error) {
	r, e := l.svcCtx.FriendClient.GetGroupMemberList(l.ctx, &friendpb.GetGroupMemberListRequest{GroupId: req.GroupId})
	if e != nil || r.Code != 0 { return &types.GroupMemberListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	var members []types.MemberVo
	for _, m := range r.Data { members = append(members, types.MemberVo{UserId: m.UserId, Nickname: m.Nickname, Role: m.Role, Avatar: ""}) }
	return &types.GroupMemberListResp{Code: 0, Message: "ok", MemberList: members}, nil
}
