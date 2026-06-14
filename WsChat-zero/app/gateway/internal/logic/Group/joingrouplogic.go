package Group

import (
	"context"
	"time"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type JoinGroupLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinGroupLogic {
	return &JoinGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *JoinGroupLogic) JoinGroup(req *types.GroupIdReq) (resp *types.CommonResp, err error) {
	uid := l.ctx.Value("userId").(int64)

	// Check group exists and not dismissed
	var count int64
	l.svcCtx.DB.Table("group_info").Where("id = ? AND status = 0", req.GroupId).Count(&count)
	if count == 0 {
		return &types.CommonResp{Code: 1, Message: "群组不存在或已解散"}, nil
	}

	// Check not already member
	l.svcCtx.DB.Table("group_member").Where("group_id = ? AND user_id = ?", req.GroupId, uid).Count(&count)
	if count > 0 {
		return &types.CommonResp{Code: 1, Message: "已在群组中"}, nil
	}

	// Insert member
	now := time.Now()
	err = l.svcCtx.DB.Exec(
		"INSERT INTO group_member (group_id, user_id, role, joined_at, created_at, updated_at) VALUES (?,?,?,?,?,?)",
		req.GroupId, uid, 0, now, now, now,
	).Error
	if err != nil {
		return &types.CommonResp{Code: -1, Message: "加入失败"}, nil
	}

	// Update member count
	l.svcCtx.DB.Exec("UPDATE group_info SET member_count = member_count + 1 WHERE id = ?", req.GroupId)

	return &types.CommonResp{Code: 0, Message: "已加入群组"}, nil
}
