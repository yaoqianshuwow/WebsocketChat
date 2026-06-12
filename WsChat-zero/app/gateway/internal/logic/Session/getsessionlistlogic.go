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

	peerIDs := make([]int64, 0, len(r.Data))
	for _, s := range r.Data {
		if s.SessionType == 1 {
			peerIDs = append(peerIDs, s.PeerId)
		}
	}

	userMap := make(map[int64]*types.UserInfoResp)
	if len(peerIDs) > 0 {
		ur, ue := l.svcCtx.UserClient.GetUserInfoList(l.ctx, &userpb.GetUserInfoListRequest{UserIds: peerIDs})
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

	data := make([]types.SessionVo, 0, len(r.Data))
	for _, s := range r.Data {
		vo := types.SessionVo{
			SessionId:      s.Id,
			PeerId:         s.PeerId,
			SessionType:    s.SessionType,
			SessionName:    s.SessionName,
			LastMsgContent: s.LastMsgContent,
			LastMsgTime:    s.LastMsgTime,
			UnreadCount:    s.UnreadCount,
		}
		if s.SessionType == 1 {
			if u := userMap[s.PeerId]; u != nil {
				if u.Nickname != "" {
					vo.SessionName = u.Nickname
				} else if u.Username != "" {
					vo.SessionName = u.Username
				}
				vo.Avatar = u.Avatar
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
