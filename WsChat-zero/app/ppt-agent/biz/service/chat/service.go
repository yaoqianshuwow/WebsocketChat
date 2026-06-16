package chat

import (
	"context"
	"ppt-agent/biz/model/vo"
)

// IChatService 对话服务接口
type IChatService interface {
	// SendMessage 发送对话消息，返回 AI 回复
	SendMessage(ctx context.Context, req *vo.ChatRequest) (*vo.ChatResponse, error)
}
