// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/config"
	Ai "github.com/your-org/ws-chat-zero/app/gateway/internal/handler/Ai"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/handler"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// AI chat (auth required)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{ctx.Auth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/api/v1/ai/chat",
					Handler: Ai.AiChatHandler(ctx),
				},
			}...,
		),
	)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
