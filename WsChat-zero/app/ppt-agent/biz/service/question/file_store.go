package question

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

type FileQuestionService struct {
	mu   sync.Mutex
	path string
}

type questionFile struct {
	Items []Question `json:"items"`
}

func NewFileQuestionService() IQuestionService {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		root = "."
	}
	return &FileQuestionService{
		path: filepath.Join(root, "storage", "assistant", "questions.json"),
	}
}

func (s *FileQuestionService) Save(ctx context.Context, sessionID, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		sessionID = "workspace"
	}
	item := Question{
		SessionID: sessionID,
		Content:   strings.TrimSpace(content),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	data.Items = append(data.Items, item)
	return s.save(data)
}

func (s *FileQuestionService) List(ctx context.Context, sessionID string) ([]Question, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return nil, err
	}

	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		sessionID = "workspace"
	}

	items := make([]Question, 0)
	for _, item := range data.Items {
		if item.SessionID == sessionID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt > items[j].CreatedAt
	})
	return items, nil
}

func (s *FileQuestionService) load() (*questionFile, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return &questionFile{}, nil
		}
		return nil, err
	}
	var data questionFile
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *FileQuestionService) save(data *questionFile) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal question store: %w", err)
	}
	if err := os.WriteFile(tmp, raw, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
