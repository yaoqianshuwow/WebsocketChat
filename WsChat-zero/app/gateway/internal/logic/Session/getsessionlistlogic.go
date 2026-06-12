package Session

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetSessionListLogic) GetSessionList(req *types.CommonResp) (resp *types.SessionListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetSessionList(l.ctx, &friendpb.GetSessionListRequest{UserId: uid})
	if e != nil || r.Code != 0 { return &types.SessionListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	var data []types.SessionVo
	for _, s := range r.Data { data = append(data, types.SessionVo{SessionId: s.Id, PeerId: s.PeerId, SessionType: s.SessionType, SessionName: s.SessionName, LastMsgContent: s.LastMsgContent, LastMsgTime: s.LastMsgTime, UnreadCount: s.UnreadCount}) }
	return &types.SessionListResp{Code: 0, Message: "ok", Data: data}, nil
}
