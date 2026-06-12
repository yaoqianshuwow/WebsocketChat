package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUsersLogic {
	return &SearchUsersLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchUsersLogic) SearchUsers(req *types.SearchUsersReq) (resp *types.SearchUsersResp, err error) {
	r, e := l.svcCtx.UserClient.SearchUsers(l.ctx, &userpb.SearchUsersRequest{
		Keyword: req.Keyword, Page: req.Page, Size: req.Size,
	})
	if e != nil || r.Code != 0 {
		return &types.SearchUsersResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}
	var data []types.UserInfoResp
	for _, u := range r.Data {
		data = append(data, types.UserInfoResp{
			UserId: u.Id, Username: u.Username, Nickname: u.Nickname,
			Avatar: u.Avatar, Sex: u.Sex, Age: u.Age, Bio: u.Bio,
			Phone: u.Phone, Status: u.Status, Role: u.Role,
		})
	}
	return &types.SearchUsersResp{Code: 0, Message: "ok", Data: data, Total: r.Total}, nil
}
