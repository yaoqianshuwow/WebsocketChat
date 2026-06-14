package Group

import (
	"context"
	"strings"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SearchGroupListLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewSearchGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchGroupListLogic {
	return &SearchGroupListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchGroupListLogic) SearchGroupList(req *types.SearchGroupListReq) (resp *types.GroupListResp, err error) {
	r, e := l.svcCtx.FriendClient.GetGroupInfoList(l.ctx, &friendpb.GetContactListRequest{UserId: 0})
	if e != nil || r == nil || r.Code != 0 {
		if r == nil {
			return &types.GroupListResp{Code: -1, Message: e.Error()}, nil
		}
		return &types.GroupListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}

	keyword := strings.TrimSpace(strings.ToLower(req.Keyword))
	data := make([]types.GroupInfoResp, 0, len(r.Data))
	for _, g := range r.Data {
		if keyword != "" {
			name := strings.ToLower(g.Name)
			notice := strings.ToLower(g.Notice)
			if !strings.Contains(name, keyword) && !strings.Contains(notice, keyword) {
				continue
			}
		}
		data = append(data, types.GroupInfoResp{
			GroupId: g.Id, Name: g.Name, Avatar: g.Avatar, OwnerId: g.OwnerId,
			MemberCount: g.MemberCount, AddMode: g.AddMode, Status: g.Status, Notice: g.Notice,
		})
	}
	return &types.GroupListResp{Code: 0, Message: "ok", Data: data}, nil
}
