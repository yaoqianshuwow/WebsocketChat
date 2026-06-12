package logic

import (
	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
	"github.com/your-org/ws-chat-zero/app/friend/pb/pb"
)

func contactModelToProto(c *model.Contact) *pb.ContactInfo {
	return &pb.ContactInfo{
		Id:          c.Id,
		UserId:      c.UserId,
		ContactId:   c.ContactId,
		ContactType: c.ContactType,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt.Unix(),
	}
}

func applyModelToProto(a *model.ContactApply) *pb.ContactApply {
	return &pb.ContactApply{
		Id:        a.Id,
		FromId:    a.FromId,
		ToId:      a.ToId,
		ApplyType: a.ApplyType,
		Remark:    a.Remark,
		Status:    a.Status,
		CreatedAt: a.CreatedAt.Unix(),
	}
}

func sessionModelToProto(s *model.Session) *pb.Session {
	return &pb.Session{
		Id:             s.Id,
		UserId:         s.UserId,
		PeerId:         s.PeerId,
		SessionType:    s.SessionType,
		SessionName:    s.SessionName,
		LastMsgId:      s.LastMsgId,
		LastMsgContent: s.LastMsgContent,
		LastMsgTime:    s.LastMsgTime,
		UnreadCount:    s.UnreadCount,
		CreatedAt:      s.CreatedAt.Unix(),
	}
}

func groupModelToProto(g *model.GroupInfo) *pb.GroupInfo {
	return &pb.GroupInfo{
		Id:          g.Id,
		Name:        g.Name,
		Avatar:      g.Avatar,
		OwnerId:     g.OwnerId,
		MemberCount: g.MemberCount,
		AddMode:     g.AddMode,
		Status:      g.Status,
		Notice:      g.Notice,
		CreatedAt:   g.CreatedAt.Unix(),
	}
}

func memberModelToProto(m *model.GroupMember) *pb.GroupMember {
	return &pb.GroupMember{
		Id:        m.Id,
		GroupId:   m.GroupId,
		UserId:    m.UserId,
		Role:      m.Role,
		Nickname:  m.Nickname,
		JoinedAt:  m.JoinedAt.Unix(),
	}
}
