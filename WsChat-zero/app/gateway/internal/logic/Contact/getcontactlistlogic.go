package Contact

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetContactListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetContactListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContactListLogic {
	return &GetContactListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetContactListLogic) GetContactList() (resp *types.ContactListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetContactList(l.ctx, &friendpb.GetContactListRequest{UserId: uid})
	if e != nil || r.Code != 0 {
		if r == nil {
			return &types.ContactListResp{Code: 1, Message: "获取联系人失败", Data: make([]types.ContactVo, 0)}, nil
		}
		return &types.ContactListResp{Code: r.GetCode(), Message: r.GetMessage(), Data: make([]types.ContactVo, 0)}, nil
	}

	userIDs := make([]int64, 0, len(r.Data))
	for _, c := range r.Data {
		userIDs = append(userIDs, c.ContactId)
	}

	userMap := make(map[int64]*types.UserInfoResp)
	if len(userIDs) > 0 {
		ur, ue := l.svcCtx.UserClient.GetUserInfoList(l.ctx, &userpb.GetUserInfoListRequest{UserIds: userIDs})
		if ue == nil && ur != nil && ur.Code == 0 {
			for _, u := range ur.Data {
				userMap[u.Id] = &types.UserInfoResp{
					UserId:   u.Id,
					Username: u.Username,
					Nickname: u.Nickname,
					Avatar:   u.Avatar,
					Sex:      u.Sex,
					Age:      u.Age,
					Bio:      u.Bio,
					Phone:    u.Phone,
					Status:   u.Status,
					Role:     u.Role,
				}
			}
		}
	}

	data := make([]types.ContactVo, 0, len(r.Data))
	for _, c := range r.Data {
		vo := types.ContactVo{
			ContactId:   c.ContactId,
			ContactType: c.ContactType,
			Status:      c.Status,
		}
		if u := userMap[c.ContactId]; u != nil {
			vo.Nickname = u.Nickname
			if vo.Nickname == "" {
				vo.Nickname = u.Username
			}
			vo.Avatar = u.Avatar
		}
		data = append(data, vo)
	}
	return &types.ContactListResp{Code: 0, Message: "ok", Data: data}, nil
}
