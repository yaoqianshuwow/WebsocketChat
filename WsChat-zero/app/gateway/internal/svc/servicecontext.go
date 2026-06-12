package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/config"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/middleware"
	filepb "github.com/your-org/ws-chat-zero/app/file/fileservice"
	friendpb "github.com/your-org/ws-chat-zero/app/friend/friendservice"
	msgpb "github.com/your-org/ws-chat-zero/app/msg-forward/messageservice"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
)

type ServiceContext struct {
	Config       config.Config
	Auth         rest.Middleware
	UserClient   userpb.UserService
	MsgClient    msgpb.MessageService
	FriendClient friendpb.FriendService
	FileClient   filepb.FileService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userCli := userpb.NewUserService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.UserRpc.Target,
	}))
	msgCli := msgpb.NewMessageService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.MsgRpc.Target,
	}))
	friendCli := friendpb.NewFriendService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.FriendRpc.Target,
	}))
	fileCli := filepb.NewFileService(zrpc.MustNewClient(zrpc.RpcClientConf{
		Target: c.FileRpc.Target,
	}))

	ctx := &ServiceContext{
		Config:       c,
		UserClient:   userCli,
		MsgClient:    msgCli,
		FriendClient: friendCli,
		FileClient:   fileCli,
	}
	ctx.Auth = middleware.NewAuthMiddleware(userCli).Handle
	return ctx
}
