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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	Auth         rest.Middleware
	UserClient   userpb.UserService
	MsgClient    msgpb.MessageService
	FriendClient friendpb.FriendService
	FileClient   filepb.FileService
	DB           *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	userCli := userpb.NewUserService(zrpc.MustNewClient(buildClientConf(c.UserRpc)))
	msgCli := msgpb.NewMessageService(zrpc.MustNewClient(buildClientConf(c.MsgRpc)))
	friendCli := friendpb.NewFriendService(zrpc.MustNewClient(buildClientConf(c.FriendRpc)))
	fileCli := filepb.NewFileService(zrpc.MustNewClient(buildClientConf(c.FileRpc)))

	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	ctx := &ServiceContext{
		Config:       c,
		UserClient:   userCli,
		MsgClient:    msgCli,
		FriendClient: friendCli,
		FileClient:   fileCli,
		DB:           db,
	}
	ctx.Auth = middleware.NewAuthMiddleware(userCli).Handle
	return ctx
}

func buildClientConf(conf config.RpcConf) zrpc.RpcClientConf {
	if len(conf.Etcd.Hosts) > 0 && len(conf.Etcd.Key) > 0 {
		return zrpc.RpcClientConf{Etcd: conf.Etcd}
	}
	return zrpc.RpcClientConf{Target: conf.Target}
}
