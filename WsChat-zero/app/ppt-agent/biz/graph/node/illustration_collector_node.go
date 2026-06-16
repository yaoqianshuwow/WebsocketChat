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

var illustModel *llm.ImageModelWrapper

func InitIllustrationCollectorNode(model *llm.ImageModelWrapper) {
	illustModel = model
}

func NewIllustrationCollectorNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 插画收集")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil || wc.SlidesJSON == "" {
			return map[string]any{}, nil
		}

		// 从 slides 中提取需要插画的场景
		var slides []map[string]any
		if err := json.Unmarshal([]byte(wc.SlidesJSON), &slides); err != nil {
			var wrapper struct{ Slides []map[string]any `json:"slides"` }
			if e2 := json.Unmarshal([]byte(wc.SlidesJSON), &wrapper); e2 != nil || len(wrapper.Slides) == 0 {
				return map[string]any{}, nil
			}
			slides = wrapper.Slides
		}

		var images []state.ImageInfo

		// 为每个需要配图的页面生成插画风格配图
		for i, s := range slides {
			needImg, _ := s["needImage"].(bool)
			desc, _ := s["imageDescription"].(string)
			if !needImg || desc == "" {
				continue
			}
			title, _ := s["title"].(string)
			illustDesc := fmt.Sprintf("Flat illustration style: %s", desc)

			img, err := GenerateImage(ctx, illustDesc, fmt.Sprintf("illust_%d_%s", i, title), "illustration")
			if err != nil {
				logger.Errorf("插画生成失败 [%s]: %v", title, err)
				images = append(images, state.ImageInfo{
					PageTitle: title, Description: illustDesc,
				})
				continue
			}
			images = append(images, *img)
			logger.Infof("插画完成 [%s]: %s", title, img.LocalPath)
		}

		wc.Illustrations = images
		state.NotifyStepCompleted(wc, "插画")
		return map[string]any{}, nil
	})
}

func IllustrationCollectorStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}
