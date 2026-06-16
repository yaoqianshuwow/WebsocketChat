package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"

	"ppt-agent/biz/handler"
	chatHandler "ppt-agent/biz/handler/chat"
	pptHandler "ppt-agent/biz/handler/ppt"
	questionHandler "ppt-agent/biz/handler/question"
	"ppt-agent/biz/handler/static"
)

// RegisterRoutes registers all HTTP routes.
func RegisterRoutes(
	r *server.Hertz,
	pptH *pptHandler.PptHandler,
	chatH *chatHandler.ChatHandler,
	questionH *questionHandler.QuestionHandler,
) {
	r.GET("/", static.ServeFrontend)
	r.GET("/workspace", static.ServeWorkspace)
	r.GET("/football-news.html", static.ServeFootballNews)
	r.GET("/api/ping", handler.Ping)

	pptRoute := r.Group("/api/ppt")
	{
		pptRoute.POST("/generate", pptH.GeneratePpt)
		pptRoute.GET("/status", pptH.GetPptStatus)
		pptRoute.GET("/download/:taskId", pptH.DownloadPpt)
		pptRoute.GET("/preview/:taskId", pptH.PreviewPpt)
		pptRoute.GET("/list", pptH.ListTasks)
	}

	chatRoute := r.Group("/api/chat")
	{
		chatRoute.POST("/send", chatH.SendMessage)
		chatRoute.GET("/history", chatH.History)
		chatRoute.GET("/sessions", chatH.Sessions)
	}

	questionRoute := r.Group("/api/question")
	{
		questionRoute.POST("/save", questionH.Save)
		questionRoute.GET("/list", questionH.List)
	}

	monitorRoute := r.Group("/api/monitor")
	{
		monitorRoute.GET("/tokens", pptH.GetTokenUsage)
	}
}
