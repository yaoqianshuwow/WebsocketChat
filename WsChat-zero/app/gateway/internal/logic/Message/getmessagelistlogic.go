package Message

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageListLogic {
	return &GetMessageListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMessageListLogic) GetMessageList(req *types.GetMessageListReq) (resp *types.MessageListResp, err error) {
	r, e := l.svcCtx.MsgClient.GetMessageList(l.ctx, &msgpb.GetMessageListRequest{SessionId: req.SessionId, Page: req.Page, Size: req.Size, BeforeId: req.BeforeId})
	if e != nil || r.Code != 0 { return &types.MessageListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	currentUserId := l.ctx.Value("userId").(int64)
	data := make([]types.MessageVo, 0)
	for _, m := range r.Data {
		vo := types.MessageVo{
			MsgId: m.Id, SenderId: m.SenderId, ReceiverId: m.ReceiverId,
			MsgType: m.MsgType, Content: m.Content, FileUrl: m.FileUrl,
			FileName: m.FileName, FileSize: m.FileSize, CreatedAt: m.CreatedAt,
			SendName: m.SenderName, SendAvatar: m.SenderAvatar,
		}
		if m.SenderId == currentUserId {
			vo.Mine = true
		}
		data = append(data, vo)
	}
	return &types.MessageListResp{Code: 0, Message: "ok", Data: data, Total: r.Total}, nil
}
