package vo

// ChatRequest 瀵硅瘽璇锋眰
type ChatRequest struct {
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

// ChatResponse 瀵硅瘽鍝嶅簲
type ChatResponse struct {
	Reply   string       `json:"reply"`
	Tool    string       `json:"tool,omitempty"`
	ToolArg *ChatToolArg `json:"toolArg,omitempty"`
}

// ChatToolArg 宸ュ叿鍙傛暟
type ChatToolArg struct {
	Topic string `json:"topic"`
	Style string `json:"style"`
}

// ChatMessage 瀵硅瘽娑堟伅璁板綍
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// ChatSessionSummary 会话摘要
type ChatSessionSummary struct {
	SessionID    string `json:"sessionId"`
	Title        string `json:"title"`
	MessageCount int    `json:"messageCount"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// ChatHistoryResponse 会话历史
type ChatHistoryResponse struct {
	SessionID string        `json:"sessionId"`
	Title     string        `json:"title"`
	Messages  []*ChatMessage `json:"messages"`
}

// QuestionRecord 问卷/主题记录
type QuestionRecord struct {
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

// QuestionSaveRequest 保存问卷/主题
type QuestionSaveRequest struct {
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
}
