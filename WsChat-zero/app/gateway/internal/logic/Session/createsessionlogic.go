package Session

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSessionLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateSessionLogic) CreateSession(req *types.CreateSessionReq) (resp *types.SessionResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.CreateSession(l.ctx, &friendpb.CreateSessionRequest{
		UserId: uid, PeerId: req.PeerId, SessionType: req.SessionType, SessionName: req.SessionName,
	})
	if e != nil || r.Code != 0 { return &types.SessionResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	return &types.SessionResp{Code: 0, Message: "ok", SessionId: r.Data.Id, PeerId: r.Data.PeerId, SessionType: r.Data.SessionType, SessionName: r.Data.SessionName}, nil
}
