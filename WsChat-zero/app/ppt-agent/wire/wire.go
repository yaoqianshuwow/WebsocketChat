//go:build wireinject

package wire

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	hertzConfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/wire"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/core"
	"ppt-agent/biz/graph/node"
	chatHandler "ppt-agent/biz/handler/chat"
	pptHandler "ppt-agent/biz/handler/ppt"
	questionHandler "ppt-agent/biz/handler/question"
	chatLogic "ppt-agent/biz/logic/chat"
	pptLogic "ppt-agent/biz/logic/ppt"
	"ppt-agent/biz/monitor"
	"ppt-agent/biz/router"
	questionLogic "ppt-agent/biz/service/question"
	sessionLogic "ppt-agent/biz/service/session"
	"ppt-agent/config"
)

var configSet = wire.NewSet(config.InitConfig)
var monitorSet = wire.NewSet(monitor.NewTokenMonitor)
var coreSet = wire.NewSet(core.NewPptGenerator)

type LLMWrappers struct {
	TextModel  *llm.ChatModelWrapper  // deepseek-v4-flash — 预检+终检
	AgentModel *llm.ChatModelWrapper  // agnes-2.0-flash — 文案生成+对话
	ImageModel *llm.ImageModelWrapper // agnes-image-2.0-flash — 配图/插画/图表/Logo
}

func InitLLMs(cfg *config.Config) *LLMWrappers {
	return &LLMWrappers{
		TextModel:  llm.NewTextModel(cfg),
		AgentModel: llm.NewAgentModel(cfg),
		ImageModel: llm.NewImageModel(cfg),
	}
}

var llmSet = wire.NewSet(InitLLMs)

// ProvideAgentModel 暴露 AgentModel 给 ChatService
func ProvideAgentModel(wrappers *LLMWrappers) *llm.ChatModelWrapper {
	return wrappers.AgentModel
}

var agentModelSet = wire.NewSet(ProvideAgentModel)

// ProvideSessionService 创建 MySQL 会话存储（对话存 message 表）
func ProvideSessionService(cfg *config.Config) (sessionLogic.ISessionService, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName,
	)
	// 改为 use wechat_database for message table
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/wechat_database?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
	)
	return sessionLogic.NewMySQLSessionService(dsn)
}

var sessionStoreSet = wire.NewSet(ProvideSessionService)

type NodeInitializer struct{}

func InitAllNodes(wrappers *LLMWrappers, pptGen *core.PptGenerator) *NodeInitializer {
	// agnes-2.0: 文案
	node.InitPptPlanNode(wrappers.AgentModel)
	// deepseek: 预检+终检
	node.InitPreQualityCheckNode(wrappers.TextModel)
	node.InitPostQualityCheckNode(wrappers.TextModel)
	// agnes-2.0: 图片质检
	node.InitImageQualityCheckNode(wrappers.AgentModel)
	// agnes-image: 全部配图/插画/图表/Logo
	node.InitContentImageCollectorNode(wrappers.ImageModel)
	node.InitIllustrationCollectorNode(wrappers.ImageModel)
	node.InitDiagramCollectorNode(wrappers.ImageModel)
	node.InitLogoCollectorNode(wrappers.ImageModel)
	// python-pptx: 生成文件
	node.InitPptMergeNode(pptGen)
	logger.Info("所有工作流节点已初始化")
	return &NodeInitializer{}
}

var nodeInitSet = wire.NewSet(InitAllNodes)
var serviceSet = wire.NewSet(pptLogic.NewPptService, chatLogic.NewChatService, questionLogic.NewFileQuestionService)
var handlerSet = wire.NewSet(pptHandler.NewPptHandler, chatHandler.NewChatHandler, questionHandler.NewQuestionHandler)

func CustomRecoveryHandler(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
	logger.Errorf("panic recovered: %v\n%s", err, stack)
	c.JSON(consts.StatusOK, map[string]interface{}{"code": 500, "message": fmt.Sprintf("%v", err)})
	c.Abort()
}

func InitServer(serverConfig *config.Config, pptHandler *pptHandler.PptHandler, chatHandler *chatHandler.ChatHandler, questionHandler *questionHandler.QuestionHandler, _ *NodeInitializer) *server.Hertz {
	serverOpts := []hertzConfig.Option{
		server.WithHostPorts(":" + strconv.Itoa(serverConfig.Server.Port)),
		server.WithBasePath(serverConfig.Server.ContextPath),
	}
	h := server.New(serverOpts...)
	h.Use(recovery.Recovery(recovery.WithRecoveryHandler(CustomRecoveryHandler)))
	router.RegisterRoutes(h, pptHandler, chatHandler, questionHandler)
	return h
}

func InitializeApp() (*server.Hertz, error) {
	panic(wire.Build(
		configSet, monitorSet, coreSet, llmSet, agentModelSet,
		nodeInitSet, serviceSet, handlerSet, InitServer,
	))
}
