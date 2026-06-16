package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/logger"
	"ppt-agent/wire"
)

func main() {
	h, err := wire.InitializeApp()
	if err != nil {
		panic(fmt.Errorf("注入失败: %w", err))
	}

	ctx := context.Background()
	_ = ctx

	// 优雅关闭
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		logger.Infof("正在关闭服务...")
		h.Close()
	}()

	logger.Infof("PPT Agent 服务启动完成")
	h.Spin()
}
