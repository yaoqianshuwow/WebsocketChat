package state

import (
	"context"
	"sync"
)

type contextKey string

const workflowContextCtxKey contextKey = "workflow_context_ctx"

type StepCallback func(stepNumber int, currentStep string)

// WorkFlowContext PPT 生成工作流上下文
type WorkFlowContext struct {
	mu           sync.Mutex
	CurrentStep  string
	Topic        string
	Style        string
	Outline      string
	SlidesJSON   string     // LLM 生成的幻灯片 JSON
	ContentImage []ImageInfo // 内容配图
	Illustrations []ImageInfo // 插画
	Diagrams     []ImageInfo // 图表
	Logos        []ImageInfo // Logo
	PptFilePath  string
	ErrorMessage string
	StepCallback StepCallback
	StepCounter  int
}

// ImageInfo 图片资源信息
type ImageInfo struct {
	PageTitle   string `json:"pageTitle"`
	Description string `json:"description"`
	LocalPath   string `json:"localPath,omitempty"`
	URL         string `json:"url,omitempty"`
}

// GraphState 图状态
type GraphState struct {
	WorkFlowContext *WorkFlowContext
}

func GetContext(graphState *GraphState) *WorkFlowContext {
	if graphState == nil {
		return nil
	}
	return graphState.WorkFlowContext
}

func GenGraphState(ctx context.Context) *GraphState {
	workflowCtx, ok := ctx.Value(workflowContextCtxKey).(*WorkFlowContext)
	if !ok || workflowCtx == nil {
		workflowCtx = &WorkFlowContext{}
	}
	return &GraphState{
		WorkFlowContext: workflowCtx,
	}
}

func WithWorkflowContext(ctx context.Context, workflowCtx *WorkFlowContext) context.Context {
	return context.WithValue(ctx, workflowContextCtxKey, workflowCtx)
}

func NotifyStepCompleted(workflowCtx *WorkFlowContext, currentStep string) {
	if workflowCtx == nil {
		return
	}
	workflowCtx.mu.Lock()
	workflowCtx.StepCounter++
	workflowCtx.CurrentStep = currentStep
	callback := workflowCtx.StepCallback
	workflowCtx.mu.Unlock()

	if callback != nil {
		callback(workflowCtx.StepCounter, currentStep)
	}
}
