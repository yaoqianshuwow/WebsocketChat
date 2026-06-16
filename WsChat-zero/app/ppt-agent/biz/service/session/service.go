package session

import "context"

// ISessionService handles persisted chat sessions/messages.
type ISessionService interface {
	AppendMessage(ctx context.Context, sessionID, role, content string) error
	GetHistory(ctx context.Context, sessionID string) ([]Message, SessionInfo, error)
	ListSessions(ctx context.Context) ([]SessionInfo, error)
}

// Message stored chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Time    string `json:"time"`
}

// SessionInfo stored session summary.
type SessionInfo struct {
	SessionID    string `json:"sessionId"`
	Title        string `json:"title"`
	MessageCount int    `json:"messageCount"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
