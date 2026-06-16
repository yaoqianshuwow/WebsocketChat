package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"ppt-agent/pkg/myfile"
)

type FileSessionService struct {
	mu   sync.Mutex
	path string
}

type storedSession struct {
	SessionID    string    `json:"sessionId"`
	Title        string    `json:"title"`
	MessageCount int       `json:"messageCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Messages     []Message `json:"messages"`
}

type sessionFile struct {
	Sessions map[string]*storedSession `json:"sessions"`
}

func NewFileSessionService() ISessionService {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		root = "."
	}
	return &FileSessionService{
		path: filepath.Join(root, "storage", "assistant", "chat_sessions.json"),
	}
}

func (s *FileSessionService) AppendMessage(ctx context.Context, sessionID, role, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	sess, ok := data.Sessions[sessionID]
	if !ok {
		sess = &storedSession{
			SessionID: sessionID,
			Title:     deriveTitle(content),
			CreatedAt: now,
		}
		data.Sessions[sessionID] = sess
	}
	if strings.TrimSpace(sess.Title) == "" {
		sess.Title = deriveTitle(content)
	}
	sess.Messages = append(sess.Messages, Message{Role: role, Content: content, Time: now.Format(time.RFC3339)})
	sess.MessageCount = len(sess.Messages)
	sess.UpdatedAt = now

	return s.save(data)
}

func (s *FileSessionService) GetHistory(ctx context.Context, sessionID string) ([]Message, SessionInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return nil, SessionInfo{}, err
	}

	sess, ok := data.Sessions[sessionID]
	if !ok {
		return []Message{}, SessionInfo{SessionID: sessionID}, nil
	}
	history := make([]Message, len(sess.Messages))
	copy(history, sess.Messages)
	info := SessionInfo{
		SessionID:    sess.SessionID,
		Title:        sess.Title,
		MessageCount: sess.MessageCount,
		CreatedAt:    sess.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    sess.UpdatedAt.Format(time.RFC3339),
	}
	return history, info, nil
}

func (s *FileSessionService) ListSessions(ctx context.Context) ([]SessionInfo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return nil, err
	}

	items := make([]SessionInfo, 0, len(data.Sessions))
	for _, sess := range data.Sessions {
		items = append(items, SessionInfo{
			SessionID:    sess.SessionID,
			Title:        sess.Title,
			MessageCount: sess.MessageCount,
			CreatedAt:    sess.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    sess.UpdatedAt.Format(time.RFC3339),
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt > items[j].UpdatedAt
	})
	return items, nil
}

func (s *FileSessionService) load() (*sessionFile, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &sessionFile{Sessions: map[string]*storedSession{}}, nil
		}
		return nil, err
	}

	var data sessionFile
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	if data.Sessions == nil {
		data.Sessions = map[string]*storedSession{}
	}
	return &data, nil
}

func (s *FileSessionService) save(data *sessionFile) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session store: %w", err)
	}
	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func deriveTitle(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return "PPT 工作台"
	}
	runes := []rune(content)
	if len(runes) > 24 {
		return string(runes[:24]) + "..."
	}
	return content
}
