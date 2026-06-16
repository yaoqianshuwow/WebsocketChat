package node

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var postCheckModel *llm.ChatModelWrapper

func InitPostQualityCheckNode(model *llm.ChatModelWrapper) {
	postCheckModel = model
}

func NewPostQualityCheckNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 成品终检")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil {
			return map[string]any{}, nil
		}

		summary := fmt.Sprintf("Topic: %s\nStyle: %s\nFile: %s\nImages: content=%d illust=%d diagram=%d logo=%d",
			wc.Topic, wc.Style, wc.PptFilePath,
			len(wc.ContentImage), len(wc.Illustrations), len(wc.Diagrams), len(wc.Logos))

		systemPrompt := `You are a PPT final quality inspector. Review the generation summary.
Return JSON: {"passed":bool, "summary":string, "suggestions":[string]}`

		userMsg := fmt.Sprintf("Review this PPT generation result:\n%s", summary)

		messages := []*schema.Message{
			{Role: schema.System, Content: systemPrompt},
			{Role: schema.User, Content: userMsg},
		}

		resp, err := postCheckModel.Generate(ctx, messages)
		if err != nil {
			logger.Errorf("终检失败: %v", err)
			return map[string]any{}, nil
		}

		logger.Infof("终检完成: %.100s", resp.Content)
		return map[string]any{}, nil
	})
}

func PostQualityCheckStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}
