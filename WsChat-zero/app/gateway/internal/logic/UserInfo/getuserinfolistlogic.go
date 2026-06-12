package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoListLogic {
	return &GetUserInfoListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserInfoListLogic) GetUserInfoList(req *types.GetUserInfoListReq) (resp *types.UserInfoListResp, err error) {
	r, e := l.svcCtx.UserClient.GetUserInfoList(l.ctx, &userpb.GetUserInfoListRequest{UserIds: req.UserIds})
	if e != nil || r.Code != 0 {
		return &types.UserInfoListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}
	data := make([]types.UserInfoResp, 0)
	for _, u := range r.Data {
		data = append(data, types.UserInfoResp{
			UserId: u.Id, Username: u.Username, Nickname: u.Nickname,
			Avatar: u.Avatar, Sex: u.Sex, Age: u.Age, Bio: u.Bio,
			Phone: u.Phone, Status: u.Status, Role: u.Role,
		})
	}
	return &types.UserInfoListResp{Code: 0, Message: "ok", Data: data}, nil
}
