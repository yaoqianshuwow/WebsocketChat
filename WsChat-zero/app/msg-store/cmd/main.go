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
	logx.MustSetup(c.Log)

	ctx := svc.NewServiceContext(c)

	// 启动 Kafka 消费者（消息存储）
	go consumer.StartMessageConsumer(ctx)

	logx.Info("msg-store service started")
	select {}
}
