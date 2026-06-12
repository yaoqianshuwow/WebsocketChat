package UserInfo

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.UserClient.UpdateUserInfo(l.ctx, &userpb.UpdateUserInfoRequest{
		UserId: uid, Nickname: req.Nickname, Avatar: req.Avatar,
		Sex: req.Sex, Age: req.Age, Bio: req.Bio,
	})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: "更新失败"}, nil
	}
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
