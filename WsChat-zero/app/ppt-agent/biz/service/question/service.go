package question

import "context"

// IQuestionService handles persisted prompt/question records.
type IQuestionService interface {
	Save(ctx context.Context, sessionID, content string) error
	List(ctx context.Context, sessionID string) ([]Question, error)
}

// Question is a persisted prompt.
type Question struct {
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}
