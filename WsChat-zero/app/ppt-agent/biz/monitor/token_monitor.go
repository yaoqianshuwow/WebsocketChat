package monitor

import (
	"fmt"
	"sync"
	"time"
)

// TokenRecord Token使用记录
type TokenRecord struct {
	TaskID    string    `json:"taskId"`
	Topic     string    `json:"topic"`
	PromptTokens   int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens     int `json:"totalTokens"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  string   `json:"duration"`
}

// TokenMonitor Token使用监控
type TokenMonitor struct {
	mu      sync.Mutex
	records map[string]*TokenRecord
}

func NewTokenMonitor() *TokenMonitor {
	return &TokenMonitor{
		records: make(map[string]*TokenRecord),
	}
}

// StartTask 开始记录任务
func (m *TokenMonitor) StartTask(taskID, topic string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records[taskID] = &TokenRecord{
		TaskID:    taskID,
		Topic:     topic,
		StartTime: time.Now(),
	}
}

// RecordTokens 记录Token使用
func (m *TokenMonitor) RecordTokens(taskID string, promptTokens, completionTokens int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if record, ok := m.records[taskID]; ok {
		record.PromptTokens += promptTokens
		record.CompletionTokens += completionTokens
		record.TotalTokens += promptTokens + completionTokens
	}
}

// EndTask 结束任务记录
func (m *TokenMonitor) EndTask(taskID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if record, ok := m.records[taskID]; ok {
		record.EndTime = time.Now()
		record.Duration = time.Since(record.StartTime).Round(time.Second).String()
	}
}

// GetSummary 获取汇总信息
func (m *TokenMonitor) GetSummary() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalTokens := 0
	taskCount := 0
	for _, r := range m.records {
		totalTokens += r.TotalTokens
		taskCount++
	}
	return fmt.Sprintf("总任务数: %d, 总Token消耗: %d", taskCount, totalTokens)
}

// GetTaskRecord 获取单个任务记录
func (m *TokenMonitor) GetTaskRecord(taskID string) *TokenRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.records[taskID]
}

// GetAllRecords 获取所有记录
func (m *TokenMonitor) GetAllRecords() map[string]*TokenRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make(map[string]*TokenRecord)
	for k, v := range m.records {
		result[k] = v
	}
	return result
}
