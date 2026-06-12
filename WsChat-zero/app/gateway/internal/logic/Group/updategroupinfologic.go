package Group

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupInfoLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewUpdateGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupInfoLogic {
	return &UpdateGroupInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateGroupInfoLogic) UpdateGroupInfo(req *types.UpdateGroupInfoReq) (resp *types.CommonResp, err error) {
	r, e := l.svcCtx.FriendClient.UpdateGroupInfo(l.ctx, &friendpb.UpdateGroupInfoRequest{GroupId: req.GroupId, Name: req.Name, Avatar: req.Avatar, Notice: req.Notice, AddMode: req.AddMode})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
