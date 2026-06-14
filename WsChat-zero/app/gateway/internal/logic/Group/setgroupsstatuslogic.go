package Group

import (
	"context"

	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupsStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetGroupsStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupsStatusLogic {
	return &SetGroupsStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *SetGroupsStatusLogic) SetGroupsStatus(req *types.GroupIdReq) (resp *types.CommonResp, err error) {
	groupResp, err := l.svcCtx.FriendClient.GetGroupInfo(l.ctx, &friendpb.GetGroupInfoRequest{GroupId: req.GroupId})
	if err != nil {
		return &types.CommonResp{Code: -1, Message: err.Error()}, nil
	}
	if groupResp == nil || groupResp.Code != 0 || groupResp.Data == nil {
		if groupResp == nil {
			return &types.CommonResp{Code: -1, Message: "获取群组状态失败"}, nil
		}
		return &types.CommonResp{Code: groupResp.Code, Message: groupResp.Message}, nil
	}

	nextStatus := int32(1)
	if groupResp.Data.Status == 1 {
		nextStatus = 0
	}

	r, e := l.svcCtx.FriendClient.SetGroupsStatus(l.ctx, &friendpb.SetGroupsStatusRequest{
		GroupIds: []int64{req.GroupId},
		Status:   nextStatus,
	})
	if e != nil {
		return &types.CommonResp{Code: -1, Message: e.Error()}, nil
	}
	return &types.CommonResp{Code: r.Code, Message: r.Message}, nil
}
