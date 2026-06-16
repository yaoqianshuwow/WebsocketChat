package Session

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetSessionListLogic) GetSessionList() (resp *types.SessionListResp, err error) {
	uid := l.ctx.Value("userId").(int64)
	r, e := l.svcCtx.FriendClient.GetSessionList(l.ctx, &friendpb.GetSessionListRequest{UserId: uid})
	if e != nil || r.Code != 0 {
		if r == nil {
			return &types.SessionListResp{Code: 1, Message: "获取会话失败"}, nil
		}
		return &types.SessionListResp{Code: r.GetCode(), Message: r.GetMessage()}, nil
	}

	// 收集需要补充头像的 peer IDs
	singlePeerIDs := make([]int64, 0, len(r.Data))
	groupPeerIDs := make([]int64, 0, len(r.Data))
	for _, s := range r.Data {
		if s.SessionType == 1 {
			singlePeerIDs = append(singlePeerIDs, sessionPeerID(uid, s.UserId, s.PeerId))
		} else if s.SessionType == 2 {
			groupPeerIDs = append(groupPeerIDs, s.PeerId)
		}
	}

	// 批量获取用户信息（单聊头像）
	userMap := make(map[int64]*types.UserInfoResp)
	if len(singlePeerIDs) > 0 {
		ur, ue := l.svcCtx.UserClient.GetUserInfoList(l.ctx, &userpb.GetUserInfoListRequest{UserIds: singlePeerIDs})
		if ue == nil && ur != nil && ur.Code == 0 {
			for _, u := range ur.Data {
				userMap[u.Id] = &types.UserInfoResp{
					UserId: u.Id, Username: u.Username, Nickname: u.Nickname,
					Avatar: u.Avatar, Sex: u.Sex, Age: u.Age, Bio: u.Bio,
					Phone: u.Phone, Status: u.Status, Role: u.Role,
				}
			}
		}
	}

	// 批量获取群组信息（群聊头像）
	groupMap := make(map[int64]*types.GroupInfoResp)
	if len(groupPeerIDs) > 0 {
		for _, gid := range groupPeerIDs {
			gr, ge := l.svcCtx.FriendClient.GetGroupInfo(l.ctx, &friendpb.GetGroupInfoRequest{GroupId: gid})
			if ge == nil && gr != nil && gr.Code == 0 && gr.Data != nil {
				groupMap[gr.Data.Id] = &types.GroupInfoResp{
					GroupId: gr.Data.Id, Name: gr.Data.Name, Avatar: gr.Data.Avatar,
				}
			}
		}
	}

	data := make([]types.SessionVo, 0, len(r.Data))
	for _, s := range r.Data {
		peerID := s.PeerId
		if s.SessionType == 1 {
			peerID = sessionPeerID(uid, s.UserId, s.PeerId)
		}
		vo := types.SessionVo{
			SessionId:      s.Id,
			PeerId:         peerID,
			SessionType:    s.SessionType,
			SessionName:    s.SessionName,
			LastMsgContent: s.LastMsgContent,
			LastMsgTime:    s.LastMsgTime,
			UnreadCount:    s.UnreadCount,
		}
		if s.SessionType == 1 {
			if u := userMap[peerID]; u != nil {
				if u.Nickname != "" {
					vo.SessionName = u.Nickname
				} else if u.Username != "" {
					vo.SessionName = u.Username
				}
				vo.Avatar = u.Avatar
			}
		} else if s.SessionType == 2 {
			if g := groupMap[s.PeerId]; g != nil {
				if vo.SessionName == "" {
					vo.SessionName = g.Name
				}
				vo.Avatar = g.Avatar
			}
		}
		if vo.SessionName == "" {
			vo.SessionName = "会话 " + formatInt64(s.PeerId)
		}
		data = append(data, vo)
	}
	return &types.SessionListResp{Code: 0, Message: "ok", Data: data}, nil
}

func formatInt64(v int64) string {
	if v == 0 {
		return "0"
	}
	sign := ""
	if v < 0 {
		sign = "-"
		v = -v
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		buf = append(buf, byte('0'+v%10))
		v /= 10
	}
	if sign != "" {
		buf = append(buf, '-')
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

func sessionPeerID(currentUserID, userID, peerID int64) int64 {
	if currentUserID == userID {
		return peerID
	}
	return userID
}
