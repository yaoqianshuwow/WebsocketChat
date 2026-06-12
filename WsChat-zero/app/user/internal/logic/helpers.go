package logic

import (
	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"
)

func modelToProto(u *model.UserInfo) *pb.UserInfo {
	return &pb.UserInfo{
		Id:        u.Id,
		Username:  u.Username,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Nickname:  u.Nickname,
		Sex:       u.Sex,
		Age:       int32(u.Age),
		Bio:       u.Bio,
		Status:    int32(u.Status),
		Role:      int32(u.Role),
		CreatedAt: u.CreatedAt.Unix(),
		UpdatedAt: u.UpdatedAt.Unix(),
	}
}
