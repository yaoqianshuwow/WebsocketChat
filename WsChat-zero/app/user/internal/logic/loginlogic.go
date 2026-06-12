package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. 查询用户
	var user model.UserInfo
	pwd := md5Hash(in.Password)
	result := l.svcCtx.DB.Where("username = ? AND password = ?", in.Username, pwd).First(&user)
	if result.Error != nil {
		return &pb.LoginResponse{Code: 1, Message: "用户名或密码错误"}, nil
	}
	if user.Status == 1 {
		return &pb.LoginResponse{Code: 1, Message: "用户已被禁用"}, nil
	}

	// 2. 生成 JWT token
	token, err := l.generateToken(user.Id)
	if err != nil {
		logx.Errorf("generate token error: %v", err)
		return &pb.LoginResponse{Code: -1, Message: "系统错误"}, nil
	}

	// 3. 存储 token 到 Redis
	_ = l.svcCtx.Redis.Set(l.ctx, fmt.Sprintf("token:%d", user.Id), token,
		time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)*time.Second)

	return &pb.LoginResponse{
		Code:    0,
		Message: "登录成功",
		Token:   token,
		UserInfo: &pb.UserInfo{
			Id:        user.Id,
			Username:  user.Username,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			Nickname:  user.Nickname,
			Sex:       user.Sex,
			Age:       int32(user.Age),
			Bio:       user.Bio,
			Status:    int32(user.Status),
			Role:      int32(user.Role),
			CreatedAt: user.CreatedAt.Unix(),
			UpdatedAt: user.UpdatedAt.Unix(),
		},
	}, nil
}

func (l *LoginLogic) generateToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.JwtAuth.AccessSecret))
}

func md5Hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
