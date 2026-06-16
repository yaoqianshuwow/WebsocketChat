package Contact

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetApplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyListLogic {
	return &GetApplyListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetApplyListLogic) GetApplyList() (resp *types.ApplyListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetApplyList(l.ctx, &friendpb.GetContactListRequest{UserId: uid})
	if e != nil || r.Code != 0 {
		if r == nil {
			return &types.ApplyListResp{Code: 1, Message: "获取申请列表失败", Data: make([]types.ApplyVo, 0)}, nil
		}
		return &types.ApplyListResp{Code: r.GetCode(), Message: r.GetMessage(), Data: make([]types.ApplyVo, 0)}, nil
	}

	userIDs := collectApplyFromIDs(r.Data)
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

	data := make([]types.ApplyVo, 0)
	for _, a := range r.Data {
		vo := types.ApplyVo{ApplyId: a.Id, FromId: a.FromId, Remark: a.Remark, Status: a.Status}
		if u := userMap[a.FromId]; u != nil {
			vo.Nickname = u.Nickname
			if vo.Nickname == "" {
				vo.Nickname = u.Username
			}
			vo.Avatar = u.Avatar
		}
		data = append(data, vo)
	}
	return &types.ApplyListResp{Code: 0, Message: "ok", Data: data}, nil
}

func collectApplyFromIDs(items []*friendpb.ContactApply) []int64 {
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.FromId)
	}
	return ids
}
