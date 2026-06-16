package node

import (
	"context"
	"encoding/json"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var contentImageModel *llm.ImageModelWrapper

func InitContentImageCollectorNode(model *llm.ImageModelWrapper) {
	contentImageModel = model
}

func NewContentImageCollectorNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 内容配图收集")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil || wc.SlidesJSON == "" {
			return map[string]any{}, nil
		}

		var slides []map[string]any
		if err := json.Unmarshal([]byte(wc.SlidesJSON), &slides); err != nil {
			var wrapper struct {
				Slides []map[string]any `json:"slides"`
			}
			if e2 := json.Unmarshal([]byte(wc.SlidesJSON), &wrapper); e2 != nil || len(wrapper.Slides) == 0 {
				return map[string]any{}, nil
			}
			slides = wrapper.Slides
		}

		var images []state.ImageInfo
		for _, s := range slides {
			needImg, _ := s["needImage"].(bool)
			desc, _ := s["imageDescription"].(string)
			if !needImg || desc == "" {
				continue
			}
			title, _ := s["title"].(string)

			// 先尝试 agnes-image 生成，失败则搜索兜底
			img, err := GenerateImage(ctx, desc, title, "content")
			if err != nil {
				logger.Warnf("配图生成失败，走搜索兜底 [%s]: %v", title, err)
				// 搜索兜底：用描述作为关键词搜索图片
				fallbackImg, fallbackErr := searchFallbackImage(ctx, desc, title)
				if fallbackErr != nil {
					logger.Errorf("搜索兜底也失败 [%s]: %v", title, fallbackErr)
					images = append(images, state.ImageInfo{
						PageTitle: title, Description: desc,
					})
					continue
				}
				images = append(images, *fallbackImg)
				logger.Infof("搜索兜底完成 [%s]: %s", title, fallbackImg.LocalPath)
				continue
			}
			images = append(images, *img)
			logger.Infof("内容配图完成 [%s]: %s", title, img.LocalPath)
		}

		wc.ContentImage = images
		state.NotifyStepCompleted(wc, "内容配图")
		return map[string]any{}, nil
	})
}

func ContentImageCollectorStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}
