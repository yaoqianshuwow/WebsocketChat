package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/schema"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/model/vo"
	sessionSvc "ppt-agent/biz/service/session"
	chatSvc "ppt-agent/biz/service/chat"
)

// ChatService 对话服务实现
type ChatService struct {
	agentModel   *llm.ChatModelWrapper
	sessionStore sessionSvc.ISessionService
}

func NewChatService(agentModel *llm.ChatModelWrapper, sessionStore sessionSvc.ISessionService) chatSvc.IChatService {
	return &ChatService{agentModel: agentModel, sessionStore: sessionStore}
}

func (s *ChatService) SendMessage(ctx context.Context, req *vo.ChatRequest) (*vo.ChatResponse, error) {
	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = "workspace"
	}
	if strings.TrimSpace(req.Message) == "" {
		return &vo.ChatResponse{Reply: "请输入消息"}, nil
	}

	systemPrompt := `你是一个 AI PPT 助手，具有完整的对话记忆能力。
重要：你可以访问本次对话的完整历史记录，包括用户之前的所有问题和你的所有回复。
请基于对话历史理解上下文，继续对话或回答用户的问题。

规则：
1. 如果用户想制作/生成/创建 PPT，提取 topic（PPT主题）和 style（风格，默认"简约"），设置 tool="generate_ppt"
2. 如果用户只是普通对话或追问，基于历史记录继续回答，只返回 reply
3. 当用户要求总结、回顾或询问之前的对话时，请基于对话历史如实回答，不要说无法访问历史

返回格式（严格 JSON，不要 markdown）：
- 需要生成PPT时 {"reply":"...","tool":"generate_ppt","toolArg":{"topic":"...","style":"..."}}
- 普通对话时 {"reply":"..."}`

	// 加载对话历史实现记忆功能
	messages := []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
	}
	if s.sessionStore != nil {
		history, _, _ := s.sessionStore.GetHistory(ctx, sessionID)
		for _, h := range history {
			role := schema.User
			if h.Role == "assistant" {
				role = schema.Assistant
			}
			messages = append(messages, &schema.Message{Role: role, Content: h.Content})
		}
	}
	messages = append(messages, &schema.Message{Role: schema.User, Content: req.Message})

	resp, err := s.agentModel.Generate(ctx, messages)
	if err != nil {
		logger.Errorf("chat llm call failed: %v", err)
		reply := fmt.Sprintf("抱歉，我暂时无法回复：%v", err)
		_ = s.persistConversation(ctx, sessionID, req.Message, reply)
		return &vo.ChatResponse{Reply: reply}, nil
	}

	parsed, parseErr := parseResponse(resp.Content)
	if parseErr == nil {
		_ = s.persistConversation(ctx, sessionID, req.Message, parsed.Reply)
	}
	return parsed, parseErr
}

// parseResponse 解析 LLM 返回的 JSON
func parseResponse(body string) (*vo.ChatResponse, error) {
	body = strings.TrimSpace(body)
	body = strings.TrimPrefix(body, "```json")
	body = strings.TrimPrefix(body, "```")
	body = strings.TrimSuffix(body, "```")
	body = strings.TrimSpace(body)

	var resp vo.ChatResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return &vo.ChatResponse{Reply: body}, nil
	}

	if resp.Reply == "" {
		resp.Reply = body
	}
	return &resp, nil
}

func (s *ChatService) persistConversation(ctx context.Context, sessionID, userMsg, assistantMsg string) error {
	if s.sessionStore == nil {
		return nil
	}
	if err := s.sessionStore.AppendMessage(ctx, sessionID, "user", userMsg); err != nil {
		logger.Errorf("append user message failed: %v", err)
	}
	if err := s.sessionStore.AppendMessage(ctx, sessionID, "assistant", assistantMsg); err != nil {
		logger.Errorf("append assistant message failed: %v", err)
	}
	return nil
}
