package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmsLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSmsLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsLoginLogic {
	return &SmsLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SmsLoginLogic) SmsLogin(in *pb.SmsLoginRequest) (*pb.LoginResponse, error) {
	// 验证短信验证码
	key := fmt.Sprintf("sms_code:%s", in.Phone)
	storedCode, err := l.svcCtx.Redis.Get(l.ctx, key).Result()
	if err != nil || storedCode != in.Code {
		return &pb.LoginResponse{Code: 1, Message: "验证码错误或已过期"}, nil
	}
	l.svcCtx.Redis.Del(l.ctx, key)

	// 查找用户
	var user model.UserInfo
	if l.svcCtx.DB.Where("phone = ?", in.Phone).First(&user).Error != nil {
		return &pb.LoginResponse{Code: 1, Message: "手机号未注册"}, nil
	}

	// 生成 token
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"exp":     time.Now().Add(time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))

	l.svcCtx.Redis.Set(l.ctx, fmt.Sprintf("token:%d", user.Id), tokenStr,
		time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)*time.Second)

	return &pb.LoginResponse{
		Code:    0,
		Message: "登录成功",
		Token:   tokenStr,
		UserInfo: &pb.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			Nickname:  user.Nickname,
			Status:    int32(user.Status),
			Role:      int32(user.Role),
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}
