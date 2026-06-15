package Message

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRecentMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRecentMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecentMessagesLogic {
	return &GetRecentMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRecentMessagesLogic) GetRecentMessages(req *types.GetRecentMessagesReq) (resp *types.MessageListResp, err error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	r, e := l.svcCtx.MsgClient.GetRecentMessages(l.ctx, &msgpb.GetRecentMessagesRequest{
		SessionId: req.SessionId,
		Limit:     limit,
	})
	if e != nil || r.Code != 0 {
		return &types.MessageListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}

	data := make([]types.MessageVo, 0)
	for _, m := range r.Data {
		data = append(data, types.MessageVo{
			MsgId:      m.Id,
			SenderId:   m.SenderId,
			ReceiverId: m.ReceiverId,
			MsgType:    m.MsgType,
			Content:    m.Content,
			FileUrl:    m.FileUrl,
			FileName:   m.FileName,
			FileSize:   m.FileSize,
			CreatedAt:  m.CreatedAt,
			SendName:   m.SenderName,
			SendAvatar: m.SenderAvatar,
		})
	}
	return &types.MessageListResp{Code: 0, Message: "ok", Data: data, Total: r.Total}, nil
}
