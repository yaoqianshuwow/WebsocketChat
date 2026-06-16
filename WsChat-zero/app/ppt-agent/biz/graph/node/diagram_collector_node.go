package node

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var diagramModel *llm.ImageModelWrapper

func InitDiagramCollectorNode(model *llm.ImageModelWrapper) {
	diagramModel = model
}

func NewDiagramCollectorNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 图表生成")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil || wc.SlidesJSON == "" {
			return map[string]any{}, nil
		}

		var slides []map[string]any
		if err := json.Unmarshal([]byte(wc.SlidesJSON), &slides); err != nil {
			var wrapper struct{ Slides []map[string]any `json:"slides"` }
			if e2 := json.Unmarshal([]byte(wc.SlidesJSON), &wrapper); e2 != nil || len(wrapper.Slides) == 0 {
				return map[string]any{}, nil
			}
			slides = wrapper.Slides
		}

		var images []state.ImageInfo
		for i, s := range slides {
			chart, _ := s["chart"].(string)
			title, _ := s["title"].(string)
			content, _ := s["content"].(string)
			if chart == "" {
				continue
			}

			chartDesc := fmt.Sprintf("%s chart about: %s. Data context: %s", chart, title, truncate(content, 100))
			img, err := GenerateImage(ctx, chartDesc, fmt.Sprintf("chart_%d_%s", i, title), "diagram")
			if err != nil {
				logger.Errorf("图表生成失败 [%s]: %v", title, err)
				images = append(images, state.ImageInfo{
					PageTitle: title, Description: chartDesc,
				})
				continue
			}
			images = append(images, *img)
			logger.Infof("图表完成 [%s] type=%s: %s", title, chart, img.LocalPath)
		}

		wc.Diagrams = images
		state.NotifyStepCompleted(wc, "图表")
		return map[string]any{}, nil
	})
}

func DiagramCollectorStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}
