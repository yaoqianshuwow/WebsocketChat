package Contact

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContactInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactInfoLogic {
	return &GetContactInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetContactInfoLogic) GetContactInfo(req *types.ContactInfoReq) (resp *types.ContactListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetContactInfo(l.ctx, &friendpb.GetContactInfoRequest{UserId: uid, ContactId: req.ContactId})
	if e != nil || r.Code != 0 {
		if r == nil {
			return &types.ContactListResp{Code: 1, Message: "获取联系人失败"}, nil
		}
		return &types.ContactListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}

	c := r.Data
	vo := types.ContactVo{
		ContactId:   c.ContactId,
		ContactType: c.ContactType,
		Status:      c.Status,
	}
	ur, ue := l.svcCtx.UserClient.GetUserInfoList(l.ctx, &userpb.GetUserInfoListRequest{UserIds: []int64{c.ContactId}})
	if ue == nil && ur != nil && ur.Code == 0 && len(ur.Data) > 0 {
		u := ur.Data[0]
		vo.Nickname = u.Nickname
		if vo.Nickname == "" {
			vo.Nickname = u.Username
		}
		vo.Avatar = u.Avatar
	}
	return &types.ContactListResp{Code: 0, Message: "ok", Data: []types.ContactVo{vo}}, nil
}
