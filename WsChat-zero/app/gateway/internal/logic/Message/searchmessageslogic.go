package Message

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchMessagesLogic {
	return &SearchMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchMessagesLogic) SearchMessages(req *types.SearchMessagesReq) (resp *types.MessageListResp, err error) {
	r, e := l.svcCtx.MsgClient.SearchMessages(l.ctx, &msgpb.SearchMessagesRequest{
		Keyword:   req.Keyword,
		SenderId:  req.SenderId,
		SessionId: req.SessionId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Page:      req.Page,
		Size:      req.Size,
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
