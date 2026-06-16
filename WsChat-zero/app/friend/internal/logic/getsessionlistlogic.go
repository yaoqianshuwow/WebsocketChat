package logic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/internal/svc"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"
)

type GetSessionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionListLogic {
	return &GetSessionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSessionListLogic) GetSessionList(in *pb.GetSessionListRequest) (*pb.SessionListResponse, error) {
	var sessions []model.Session
	query := l.svcCtx.DB.Where("user_id = ?", in.UserId)
	if in.SessionType <= 0 || in.SessionType == 1 {
		query = query.Or("session_type = ? AND peer_id = ?", 1, in.UserId)
	} else if in.SessionType > 1 {
		query = query.Where("session_type = ?", in.SessionType)
	}
	result := query.Order("updated_at DESC").Find(&sessions)
	if result.Error != nil {
		return &pb.SessionListResponse{Code: 1, Message: "查询失败"}, nil
	}

	byKey := make(map[string]model.Session, len(sessions))
	for i := range sessions {
		key := sessionDedupKey(in.UserId, &sessions[i])
		existing, ok := byKey[key]
		if !ok || shouldPreferSession(&sessions[i], &existing) {
			byKey[key] = sessions[i]
		}
	}

	// 补充会话头像
	avatarCache := make(map[string]string)
	data := make([]*pb.Session, 0, len(byKey))
	for _, session := range byKey {
		s := session
		if s.Avatar == "" {
			cacheKey := fmt.Sprintf("%d_%d", s.SessionType, s.PeerId)
			if cached, ok := avatarCache[cacheKey]; ok {
				s.Avatar = cached
			} else {
				var avatar string
				if s.SessionType == 1 {
					l.svcCtx.DB.Raw("SELECT avatar FROM user_info WHERE id = ?", s.PeerId).Scan(&avatar)
				} else if s.SessionType == 2 {
					l.svcCtx.DB.Raw("SELECT avatar FROM group_info WHERE id = ?", s.PeerId).Scan(&avatar)
				}
				if avatar != "" {
					avatarCache[cacheKey] = avatar
					s.Avatar = avatar
				}
			}
		}
		data = append(data, sessionModelToProto(&s))
	}
	return &pb.SessionListResponse{Code: 0, Message: "ok", Data: data}, nil
}

func sessionDedupKey(currentUserID int64, s *model.Session) string {
	if s.SessionType == 1 {
		a := s.UserId
		b := s.PeerId
		if currentUserID > 0 && currentUserID == s.PeerId {
			a = s.PeerId
			b = s.UserId
		}
		if a > b {
			a, b = b, a
		}
		return fmt.Sprintf("%d:%d:%d", s.SessionType, a, b)
	}
	return fmt.Sprintf("%d:%d", s.SessionType, s.PeerId)
}

func shouldPreferSession(candidate, existing *model.Session) bool {
	if candidate.SessionType != 1 || existing.SessionType != 1 {
		return candidate.UpdatedAt.After(existing.UpdatedAt)
	}

	candidateCanonical := candidate.UserId <= candidate.PeerId
	existingCanonical := existing.UserId <= existing.PeerId
	switch {
	case candidateCanonical && !existingCanonical:
		return true
	case !candidateCanonical && existingCanonical:
		return false
	default:
		return candidate.UpdatedAt.After(existing.UpdatedAt)
	}
}
