package Group

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GroupIdReq) (resp *types.GroupMemberListResp, err error) {
	r, e := l.svcCtx.FriendClient.GetGroupMemberList(l.ctx, &friendpb.GetGroupMemberListRequest{GroupId: req.GroupId})
	if e != nil || r.Code != 0 {
		if r == nil {
			return &types.GroupMemberListResp{Code: 1, Message: "获取群成员失败", MemberList: make([]types.MemberVo, 0)}, nil
		}
		return &types.GroupMemberListResp{Code: r.GetCode(), Message: r.GetMessage(), MemberList: make([]types.MemberVo, 0)}, nil
	}

	userIDs := collectGroupUserIDs(r.Data)
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

	members := make([]types.MemberVo, 0)
	for _, m := range r.Data {
		vo := types.MemberVo{UserId: m.UserId, Nickname: m.Nickname, Role: m.Role}
		if u := userMap[m.UserId]; u != nil {
			if vo.Nickname == "" {
				vo.Nickname = u.Nickname
				if vo.Nickname == "" {
					vo.Nickname = u.Username
				}
			}
			vo.Avatar = u.Avatar
		}
		members = append(members, vo)
	}
	return &types.GroupMemberListResp{Code: 0, Message: "ok", MemberList: members}, nil
}

func collectGroupUserIDs(items []*friendpb.GroupMember) []int64 {
	ids := make([]int64, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.UserId)
	}
	return ids
}
