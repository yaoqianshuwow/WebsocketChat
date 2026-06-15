package main

import (
	"flag"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/your-org/ws-chat-zero/app/msg-store/internal/config"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/consumer"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/svc"
)

var configFile = flag.String("f", "etc/msg-store.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.MustSetup(c.RestConf.Log)
	// 手动初始化 Telemetry（因为 msg-store 不是 rest/zrpc server，不会自动初始化）
	c.RestConf.MustSetUp()

	ctx := svc.NewServiceContext(c)

	// 启动 Kafka 消费者，负责消息落库和会话同步
	go consumer.StartMessageConsumer(ctx)

	logx.Info("msg-store service started")
	select {}
}
