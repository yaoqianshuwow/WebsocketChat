package session

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLSessionService stores conversations in the shared wechat_database.message table
type MySQLSessionService struct {
	db *sql.DB
}

func NewMySQLSessionService(dsn string) (*MySQLSessionService, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}
	return &MySQLSessionService{db: db}, nil
}

func (s *MySQLSessionService) AppendMessage(ctx context.Context, sessionID, role, content string) error {
	if role == "assistant" || role == "system" {
		return nil // only store user messages + AI replies; skip system
	}
	now := time.Now().UnixNano()
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO message (msg_id, sender_id, receiver_id, chat_type, msg_type, content, status, session_id, send_name, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		now, -1, 0, 3, 1, content, 0, 0, "用户",
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (s *MySQLSessionService) GetHistory(ctx context.Context, sessionID string) ([]Message, SessionInfo, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT content, send_name, created_at FROM message
		 WHERE chat_type = 3 AND status = 0
		 ORDER BY id ASC
		 LIMIT 50`,
	)
	if err != nil {
		return nil, SessionInfo{}, fmt.Errorf("query history: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var content, sendName string
		var createdAt time.Time
		if err := rows.Scan(&content, &sendName, &createdAt); err != nil {
			continue
		}
		role := "assistant"
		if sendName == "用户" || sendName == "" {
			role = "user"
		}
		messages = append(messages, Message{
			Role:    role,
			Content: content,
			Time:    createdAt.Format(time.RFC3339),
		})
	}
	return messages, SessionInfo{SessionID: sessionID, Title: "AI对话"}, nil
}

func (s *MySQLSessionService) ListSessions(ctx context.Context) ([]SessionInfo, error) {
	return []SessionInfo{{SessionID: "workspace", Title: "AI 工作台"}}, nil
}

func (s *MySQLSessionService) Close() error {
	return s.db.Close()
}

// Ensure MySQLSessionService implements ISessionService
var _ ISessionService = (*MySQLSessionService)(nil)
