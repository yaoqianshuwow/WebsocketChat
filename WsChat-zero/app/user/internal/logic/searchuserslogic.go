package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUsersLogic {
	return &SearchUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUsersLogic) SearchUsers(in *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}

	var users []model.UserInfo
	var total int64

	query := l.svcCtx.DB.Model(&model.UserInfo{}).
		Where("username LIKE ? OR nickname LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	query.Count(&total)
	query.Offset(int((in.Page - 1) * in.Size)).Limit(int(in.Size)).Find(&users)

	var data []*pb.UserInfo
	for i := range users {
		data = append(data, modelToProto(&users[i]))
	}
	return &pb.SearchUsersResponse{Code: 0, Message: "ok", Data: data, Total: total}, nil
}
