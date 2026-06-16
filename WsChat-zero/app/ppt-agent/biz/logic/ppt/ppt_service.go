package ppt

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/gopkg/util/logger"

	"ppt-agent/biz/graph"
	"ppt-agent/biz/model/vo"
	"ppt-agent/biz/monitor"
	"ppt-agent/biz/service/ppt"
	"ppt-agent/pkg/snowflake"
)

type PptService struct {
	monitor *monitor.TokenMonitor
	tasks   sync.Map
}

type taskStatus struct {
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Topic     string    `json:"topic"`
	Style     string    `json:"style"`
	FilePath  string    `json:"filePath,omitempty"`
	Slides    string    `json:"slides,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type slidePayload struct {
	Title            string `json:"title"`
	Content          string `json:"content"`
	NeedImage        bool   `json:"needImage"`
	ImageDescription string `json:"imageDescription,omitempty"`
	Chart            string `json:"chart,omitempty"`
}

func NewPptService(
	monitor *monitor.TokenMonitor,
) ppt.IPptService {
	return &PptService{
		monitor: monitor,
	}
}

func (s *PptService) GeneratePpt(ctx context.Context, req *vo.PptGenerateRequest) (*vo.PptGenerateResponse, error) {
	if strings.TrimSpace(req.Topic) == "" {
		return nil, fmt.Errorf("PPT topic cannot be empty")
	}
	if strings.TrimSpace(req.Style) == "" {
		req.Style = "简约"
	}

	taskID := snowflake.NextStrID()
	s.tasks.Store(taskID, &taskStatus{
		Status:    "processing",
		Message:   "generating slides",
		Topic:     req.Topic,
		Style:     req.Style,
		CreatedAt: time.Now(),
	})

	go s.runGeneration(taskID, *req)

	return &vo.PptGenerateResponse{
		TaskID:  taskID,
		Status:  "processing",
		Message: "generation started",
	}, nil
}

func (s *PptService) runGeneration(taskID string, req vo.PptGenerateRequest) {
	ctx := context.Background()
	logger.Infof("start generating ppt: topic=%s, style=%s, taskID=%s", req.Topic, req.Style, taskID)

	s.monitor.StartTask(taskID, req.Topic)
	defer s.monitor.EndTask(taskID)

	// 使用 DAG 工作流生成 PPT
	wc, err := graph.ExecutePptWorkflow(ctx, req.Topic, req.Style)
	if err != nil {
		s.tasks.Store(taskID, &taskStatus{
			Status:    "failed",
			Message:   fmt.Sprintf("workflow failed: %v", err),
			Topic:     req.Topic,
			Style:     req.Style,
			CreatedAt: time.Now(),
		})
		logger.Errorf("ppt workflow failed: taskID=%s err=%v", taskID, err)
		return
	}

	if wc.ErrorMessage != "" {
		s.tasks.Store(taskID, &taskStatus{
			Status:    "failed",
			Message:   wc.ErrorMessage,
			Topic:     req.Topic,
			Style:     req.Style,
			Slides:    wc.SlidesJSON,
			CreatedAt: time.Now(),
		})
		logger.Errorf("ppt workflow error: taskID=%s err=%s", taskID, wc.ErrorMessage)
		return
	}

	s.tasks.Store(taskID, &taskStatus{
		Status:    "completed",
		Message:   "ppt generated",
		Topic:     req.Topic,
		Style:     req.Style,
		FilePath:  wc.PptFilePath,
		Slides:    wc.SlidesJSON,
		CreatedAt: time.Now(),
	})
	logger.Infof("ppt generated successfully: taskID=%s filePath=%s", taskID, wc.PptFilePath)
}

func (s *PptService) GetPptStatus(ctx context.Context, taskID string) (*vo.PptGenerateResponse, error) {
	val, ok := s.tasks.Load(taskID)
	if !ok {
		return &vo.PptGenerateResponse{
			TaskID:  taskID,
			Status:  "not_found",
			Message: "task not found",
		}, nil
	}

	ts := val.(*taskStatus)
	return &vo.PptGenerateResponse{
		TaskID:  taskID,
		Status:  ts.Status,
		Message: ts.Message,
		Topic:   ts.Topic,
		Style:   ts.Style,
	}, nil
}

func (s *PptService) DownloadPpt(ctx context.Context, taskID string) (*vo.PptDownloadResponse, error) {
	val, ok := s.tasks.Load(taskID)
	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	ts := val.(*taskStatus)
	if ts.Status != "completed" {
		return nil, fmt.Errorf("task is not completed: %s", ts.Status)
	}

	return &vo.PptDownloadResponse{
		FileURL:  ts.FilePath,
		FileName: fmt.Sprintf("ppt_%s.pptx", taskID),
	}, nil
}

func (s *PptService) GetPptPreview(ctx context.Context, taskID string) (*vo.PptPreviewResponse, error) {
	val, ok := s.tasks.Load(taskID)
	if !ok {
		return nil, fmt.Errorf("task not found")
	}

	ts := val.(*taskStatus)
	if strings.TrimSpace(ts.Slides) == "" {
		return &vo.PptPreviewResponse{
			TaskID:    taskID,
			Status:    ts.Status,
			TotalPage: 0,
			Slides:    []*vo.SlideInfo{},
		}, nil
	}

	// 尝试解析顶层 slides 数组或嵌套 {"slides": [...]}
	var slides []slidePayload
	if err := json.Unmarshal([]byte(ts.Slides), &slides); err != nil {
		var wrapper struct {
			Slides []slidePayload `json:"slides"`
		}
		if e2 := json.Unmarshal([]byte(ts.Slides), &wrapper); e2 != nil {
			return nil, fmt.Errorf("parse preview slides failed: %w", err)
		}
		slides = wrapper.Slides
	}

	result := make([]*vo.SlideInfo, 0, len(slides))
	for _, slide := range slides {
		result = append(result, &vo.SlideInfo{
			Title:    slide.Title,
			Content:  slide.Content,
			ImageURL: slide.ImageDescription,
			Chart:    slide.Chart,
		})
	}

	return &vo.PptPreviewResponse{
		TaskID:    taskID,
		Status:    ts.Status,
		TotalPage: len(result),
		Slides:    result,
	}, nil
}

func (s *PptService) ListTasks(ctx context.Context) ([]*vo.PptTaskItem, error) {
	var items []*vo.PptTaskItem
	s.tasks.Range(func(key, value any) bool {
		taskID := key.(string)
		ts := value.(*taskStatus)
		items = append(items, &vo.PptTaskItem{
			TaskID:    taskID,
			Topic:     ts.Topic,
			Style:     ts.Style,
			Status:    ts.Status,
			CreatedAt: ts.CreatedAt.Format("2006-01-02 15:04:05"),
		})
		return true
	})
	// 按创建时间降序排列
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].CreatedAt < items[j].CreatedAt {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	return items, nil
}
