package Message

import (
	"context"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMessageListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewGetGroupMessageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMessageListLogic {
	return &GetGroupMessageListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupMessageListLogic) GetGroupMessageList(req *types.GetGroupMessageReq) (resp *types.MessageListResp, err error) {
	r, e := l.svcCtx.MsgClient.GetGroupMessageList(l.ctx, &msgpb.GetGroupMessageListRequest{GroupId: req.GroupId, Page: req.Page, Size: req.Size, BeforeId: req.BeforeId})
	if e != nil || r.Code != 0 { return &types.MessageListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil }
	data := make([]types.MessageVo, 0)
	for _, m := range r.Data { data = append(data, types.MessageVo{MsgId: m.Id, SenderId: m.SenderId, ReceiverId: m.ReceiverId, MsgType: m.MsgType, Content: m.Content, FileUrl: m.FileUrl, FileName: m.FileName, FileSize: m.FileSize, CreatedAt: m.CreatedAt, SendName: m.SenderName, SendAvatar: m.SenderAvatar}) }
	return &types.MessageListResp{Code: 0, Message: "ok", Data: data, Total: r.Total}, nil
}
