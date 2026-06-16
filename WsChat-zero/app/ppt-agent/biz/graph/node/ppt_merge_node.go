package node

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"

	"ppt-agent/biz/core"
	"ppt-agent/biz/graph/state"
)

var mergePptGen *core.PptGenerator

func InitPptMergeNode(generator *core.PptGenerator) {
	mergePptGen = generator
}

func NewPptMergeNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: PPT合并")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil {
			return nil, fmt.Errorf("工作流上下文为空")
		}

		totalImages := len(wc.ContentImage) + len(wc.Illustrations) +
			len(wc.Diagrams) + len(wc.Logos)
		logger.Infof("共收集 %d 张图片: content=%d, illust=%d, diagram=%d, logo=%d",
			totalImages, len(wc.ContentImage), len(wc.Illustrations),
			len(wc.Diagrams), len(wc.Logos))

		if wc.SlidesJSON == "" {
			return map[string]any{"nodeName": "ppt_merge", "error": "幻灯片内容为空"}, nil
		}

		filePath, err := mergePptGen.CreatePptFile(ctx, wc.SlidesJSON, wc.Style, wc.Topic)
		if err != nil {
			logger.Errorf("PPT生成失败: %v", err)
			return map[string]any{"nodeName": "ppt_merge", "error": err.Error()}, nil
		}

		logger.Infof("PPT合并完成: %s", filePath)
		return map[string]any{
			"nodeName": "ppt_merge",
			"filePath": filePath,
		}, nil
	})
}

func PptMergeStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	wc := state.GetContext(gs)
	if wc != nil {
		if fp, ok := output["filePath"].(string); ok {
			wc.PptFilePath = fp
		}
		if err, ok := output["error"].(string); ok {
			wc.ErrorMessage = err
		}
		state.NotifyStepCompleted(wc, "PPT合并")
	}
	return output, nil
}
