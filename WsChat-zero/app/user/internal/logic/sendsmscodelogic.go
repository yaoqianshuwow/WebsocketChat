package logic

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsCodeLogic {
	return &SendSmsCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSmsCodeLogic) SendSmsCode(in *pb.SendSmsCodeRequest) (*pb.CommonResponse, error) {
	// 模拟生成 6 位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	key := fmt.Sprintf("sms_code:%s", in.Phone)

	// 存入 Redis，5 分钟过期
	if err := l.svcCtx.Redis.Set(l.ctx, key, code, 5*time.Minute).Err(); err != nil {
		logx.Errorf("send sms code error: %v", err)
		return &pb.CommonResponse{Code: -1, Message: "发送失败"}, nil
	}

	logx.Infof("SMS code sent to %s: %s", in.Phone, code)
	return &pb.CommonResponse{Code: 0, Message: "验证码已发送"}, nil
}
