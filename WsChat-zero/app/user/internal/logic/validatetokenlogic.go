package logic

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateTokenLogic) ValidateToken(in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(in.Token, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(l.svcCtx.Config.JwtAuth.AccessSecret), nil
		})

	if err != nil || !token.Valid {
		return &pb.ValidateTokenResponse{Code: 1, Message: "token 无效"}, nil
	}

	uid, ok := claims["user_id"].(float64)
	if !ok {
		return &pb.ValidateTokenResponse{Code: 1, Message: "token 格式错误"}, nil
	}
	userId := int64(uid)

	// 校验 Redis 中的 token
	val, _ := l.svcCtx.Redis.Get(l.ctx, fmt.Sprintf("token:%d", userId)).Result()
	if val != in.Token {
		return &pb.ValidateTokenResponse{Code: 1, Message: "token 已失效"}, nil
	}

	return &pb.ValidateTokenResponse{
		Code:    0,
		Message: "ok",
		UserId:  userId,
	}, nil
}
