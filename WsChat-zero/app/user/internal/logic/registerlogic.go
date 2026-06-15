package logic

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if in.Username == "" {
		return &pb.RegisterResponse{Code: 1, Message: "用户名不能为空"}, nil
	}
	if in.Nickname == "" {
		return &pb.RegisterResponse{Code: 1, Message: "昵称不能为空"}, nil
	}
	var count int64
	l.svcCtx.DB.Model(&model.UserInfo{}).Where("username = ?", in.Username).Count(&count)
	if count > 0 {
		return &pb.RegisterResponse{Code: 1, Message: "用户名已存在"}, nil
	}

	pwd := fmt.Sprintf("%x", md5.Sum([]byte(in.Password)))
	user := model.UserInfo{
		Username:  in.Username,
		Password:  pwd,
		Phone:     in.Phone,
		Nickname:  in.Nickname,
		Avatar:    "/static/avatars/default.png",
		Status:    0,
		Role:      0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := l.svcCtx.DB.Create(&user).Error; err != nil {
		logx.Errorf("register error: %v", err)
		return &pb.RegisterResponse{Code: -1, Message: "注册失败"}, nil
	}

	// 同步用户到 ES（异步非阻塞，失败不影响注册）
	go l.indexUserToES(&user)

	return &pb.RegisterResponse{
		Code:    0,
		Message: "注册成功",
		UserId:  user.Id,
	}, nil
}

func (l *RegisterLogic) indexUserToES(user *model.UserInfo) {
	if l.svcCtx.ES == nil {
		return
	}

	doc := map[string]any{
		"user_id":    user.Id,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"status":     user.Status,
		"created_at": user.CreatedAt.Format(time.RFC3339),
	}

	body, _ := json.Marshal(doc)
	index := l.svcCtx.Config.ES.Index
	if index == "" {
		index = "users"
	}

	resp, err := l.svcCtx.ES.Index(index, bytes.NewReader(body),
		l.svcCtx.ES.Index.WithDocumentID(fmt.Sprintf("%d", user.Id)))
	if err != nil {
		logx.Errorf("ES index user error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logx.Errorf("ES index user error response: %s", resp.String())
		return
	}
	logx.Infof("user indexed to ES: id=%d username=%s", user.Id, user.Username)
}
