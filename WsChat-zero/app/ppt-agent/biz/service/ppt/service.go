package ppt

import (
	"context"
	"ppt-agent/biz/model/vo"
)

// IPptService PPT 生成服务接口
type IPptService interface {
	// GeneratePpt 根据主题生成PPT
	GeneratePpt(ctx context.Context, req *vo.PptGenerateRequest) (*vo.PptGenerateResponse, error)
	// GetPptStatus 查询PPT生成状态
	GetPptStatus(ctx context.Context, taskID string) (*vo.PptGenerateResponse, error)
	// DownloadPpt 下载PPT文件
	DownloadPpt(ctx context.Context, taskID string) (*vo.PptDownloadResponse, error)
	// GetPptPreview 预览PPT内容（JSON格式）
	GetPptPreview(ctx context.Context, taskID string) (*vo.PptPreviewResponse, error)
	// ListTasks 列出所有生成任务
	ListTasks(ctx context.Context) ([]*vo.PptTaskItem, error)
}
