package node

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var logoModel *llm.ImageModelWrapper

func InitLogoCollectorNode(model *llm.ImageModelWrapper) {
	logoModel = model
}

func NewLogoCollectorNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: Logo生成")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)

		var images []state.ImageInfo
		if wc != nil && wc.Topic != "" {
			logoDesc := fmt.Sprintf("Professional logo for presentation about: %s. Minimalist, clean design, suitable for slide cover", wc.Topic)
			img, err := GenerateImage(ctx, logoDesc, "cover_logo", "logo")
			if err != nil {
				logger.Errorf("Logo生成失败: %v", err)
				images = append(images, state.ImageInfo{
					PageTitle: "封面Logo", Description: logoDesc,
				})
			} else {
				images = append(images, *img)
				logger.Infof("Logo完成: %s", img.LocalPath)
			}
			wc.Logos = images
			state.NotifyStepCompleted(wc, "Logo")
		}
		return map[string]any{}, nil
	})
}

func LogoCollectorStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}
