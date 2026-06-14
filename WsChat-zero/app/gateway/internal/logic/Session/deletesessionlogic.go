package Session

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	// 删除 session 记录
	r, e := l.svcCtx.FriendClient.DeleteSession(l.ctx, &friendpb.DeleteSessionRequest{SessionId: req.SessionId, UserId: uid})
	if e != nil { return &types.CommonResp{Code: -1, Message: e.Error()}, nil }
	if r.Code != 0 { return &types.CommonResp{Code: r.Code, Message: r.Message}, nil }

	// 删除该会话下的所有消息（软删除置 status=1）
	if err := l.svcCtx.DB.Model(&struct{}{}).Table("message").
		Where("session_id = ?", req.SessionId).
		Update("status", 1).Error; err != nil {
		logx.Errorf("delete messages for session %d failed: %v", req.SessionId, err)
		// 不阻断，session 已删成功
	}

	return &types.CommonResp{Code: 0, Message: "删除成功"}, nil
}
