package Session

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSessionLogic {
	return &DeleteSessionLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteSessionLogic) DeleteSession(req *types.DeleteSessionReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)

	// 先查会话类型
	sessionResp, se := l.svcCtx.FriendClient.GetSessionList(l.ctx, &friendpb.GetSessionListRequest{UserId: uid})
	if se != nil || sessionResp == nil || sessionResp.Code != 0 {
		return &types.CommonResp{Code: 1, Message: "查询会话失败"}, nil
	}
	isGroup := false
	for _, s := range sessionResp.Data {
		if s.Id == req.SessionId {
			isGroup = s.SessionType == 2
			break
		}
	}

	// 删除会话
	r, e := l.svcCtx.FriendClient.DeleteSession(l.ctx, &friendpb.DeleteSessionRequest{SessionId: req.SessionId, UserId: uid})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: e.Error()}, nil
	}
	if r.Code != 0 {
		return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
	}

	// 单聊：软删除自己的消息；群聊：不删消息
	if !isGroup {
		_, me := l.svcCtx.MsgClient.DeleteSessionMessages(l.ctx, &msgpb.GetRecentMessagesRequest{
			SessionId: req.SessionId,
			Limit:     0,
		})
		if me != nil {
			logx.Errorf("delete session messages failed: sessionId=%d err=%v", req.SessionId, me)
		}
	}

	return &types.CommonResp{Code: 0, Message: "删除成功"}, nil
}
