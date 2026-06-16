package llm

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"ppt-agent/config"
)

// ChatModelWrapper 聊天模型封装
type ChatModelWrapper struct {
	*openai.ChatModel
	ModelName string
}

// ImageModelWrapper 图片模型封装（兼容 OpenAI 格式的图片生成API）
type ImageModelWrapper struct {
	BaseURL   string
	APIKey    string
	ModelName string
}

// NewTextModel 创建文本模型（deepseek-v4-flash）
func NewTextModel(cfg *config.Config) *ChatModelWrapper {
	ctx := context.Background()
	modelName := cfg.AI.TextModel.ModelName

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: cfg.AI.TextModel.BaseURL,
		Model:   modelName,
		APIKey:  cfg.AI.TextModel.APIKey,
	})
	if err != nil {
		panic(fmt.Errorf("创建文本模型失败: %w", err))
	}

	return &ChatModelWrapper{
		ChatModel: chatModel,
		ModelName: modelName,
	}
}

// NewAgentModel 创建Agent模型（agnes-2.0-flash，用于工具调用）
func NewAgentModel(cfg *config.Config) *ChatModelWrapper {
	ctx := context.Background()
	modelName := cfg.AI.AgentModel.ModelName

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: cfg.AI.AgentModel.BaseURL,
		Model:   modelName,
		APIKey:  cfg.AI.AgentModel.APIKey,
	})
	if err != nil {
		panic(fmt.Errorf("创建Agent模型失败: %w", err))
	}

	return &ChatModelWrapper{
		ChatModel: chatModel,
		ModelName: modelName,
	}
}

// NewImageModel 创建图片生成模型配置
func NewImageModel(cfg *config.Config) *ImageModelWrapper {
	return &ImageModelWrapper{
		BaseURL:   cfg.AI.ImageModel.BaseURL,
		APIKey:    cfg.AI.ImageModel.APIKey,
		ModelName: cfg.AI.ImageModel.ModelName,
	}
}

func (w *ChatModelWrapper) GetChatModel() *openai.ChatModel {
	return w.ChatModel
}

func (w *ChatModelWrapper) GetModelName() string {
	return w.ModelName
}
