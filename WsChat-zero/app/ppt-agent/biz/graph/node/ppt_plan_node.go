package node

import (
	"context"
	"fmt"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"
	openaiModel "github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var planAgentModel *llm.ChatModelWrapper

func InitPptPlanNode(model *llm.ChatModelWrapper) {
	planAgentModel = model
}

func NewPptPlanNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: PPT计划")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil {
			return nil, fmt.Errorf("工作流上下文为空")
		}

		systemPrompt := strings.Join([]string{
			"You are a professional PPT content planner. Generate rich, detailed content.",
			"Return ONLY valid JSON. Do NOT include markdown.",
			"Output: {\"slides\":[...]}",
			"Each slide: {\"title\":\"\",\"content\":\"\",\"needImage\":bool,\"imageDescription\":\"\",\"chart\":\"\"}",
			"content must be 100-300 Chinese characters of detailed bullet points and explanations.",
			"needImage=true means this slide needs a photograph/image for visual impact.",
			"imageDescription must be a specific image search query in English (e.g. 'giant panda bamboo forest eating').",
			"chart: \"bar\"|\"pie\"|\"line\"|\"\" (empty if no chart).",
			"IMPORTANT: Generate exactly 20 slides with comprehensive content.",
			"Slides structure: 1=cover, 2=agenda/outline, 3-18=detailed content sections,",
			"19=summary/key-takeaways, 20=thank-you/Q&A.",
			"Every 3rd slide should have needImage=true for visual breaks.",
			"Each content slide must have substantial text with 5-8 bullet points or paragraphs.",
		}, " ")

		userMsg := fmt.Sprintf("Topic: %s\nStyle: %s\nGenerate 20 slides with rich Chinese content. Cover the topic comprehensively with data, examples, and actionable insights.", wc.Topic, wc.Style)

		messages := []*schema.Message{
			{Role: schema.System, Content: systemPrompt},
			{Role: schema.User, Content: userMsg},
		}

		resp, err := planAgentModel.Generate(ctx, messages,
			openaiModel.WithExtraFields(map[string]any{
				"response_format": map[string]any{"type": "json_object"},
			}))
		if err != nil {
			logger.Errorf("PPT计划失败: %v", err)
			return map[string]any{"nodeName": "ppt_plan", "error": err.Error()}, nil
		}

		logger.Infof("PPT计划完成, slides长度=%d", len(resp.Content))

		return map[string]any{
			"nodeName":   "ppt_plan",
			"slidesJSON": resp.Content,
		}, nil
	})
}

func PptPlanStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	wc := state.GetContext(gs)
	if wc != nil {
		if sj, ok := output["slidesJSON"].(string); ok {
			wc.SlidesJSON = sj
		}
		state.NotifyStepCompleted(wc, "PPT计划")
	}
	return output, nil
}
