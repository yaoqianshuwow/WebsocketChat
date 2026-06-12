package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserInfoLogic) GetUserInfo(req *types.GetUserInfoReq) (resp *types.UserInfoResp, err error) {
	uid := req.UserId
	if uid == 0 {
		uid = l.ctx.Value("userId").(int64)
	}
	r, e := l.svcCtx.UserClient.GetUserInfo(l.ctx, &userpb.GetUserInfoRequest{UserId: uid})
	if e != nil || r.Code != 0 {
		return &types.UserInfoResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}
	return &types.UserInfoResp{
		Code: 0, Message: "ok",
		UserId: r.Data.Id, Username: r.Data.Username, Nickname: r.Data.Nickname,
		Avatar: r.Data.Avatar, Sex: r.Data.Sex, Age: r.Data.Age, Bio: r.Data.Bio,
		Phone: r.Data.Phone, Status: r.Data.Status, Role: r.Data.Role,
	}, nil
}
