package node

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"ppt-agent/biz/ai/llm"
	"ppt-agent/biz/graph/state"
)

var preCheckModel *llm.ChatModelWrapper

func InitPreQualityCheckNode(model *llm.ChatModelWrapper) {
	preCheckModel = model
}

// slideCheckItem 单页检测结果
type slideCheckItem struct {
	Index     int    `json:"index"`
	Title     string `json:"title"`
	Qualified bool   `json:"qualified"`
	Reason    string `json:"reason,omitempty"`
}

func NewPreQualityCheckNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 文案预检(过滤)")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil || wc.SlidesJSON == "" {
			return map[string]any{}, nil
		}

		// 解析 slides
		var slides []map[string]any
		if err := json.Unmarshal([]byte(wc.SlidesJSON), &slides); err != nil {
			var wrapper struct{ Slides []map[string]any `json:"slides"` }
			if e2 := json.Unmarshal([]byte(wc.SlidesJSON), &wrapper); e2 == nil {
				slides = wrapper.Slides
			}
		}

		if len(slides) == 0 {
			return map[string]any{}, nil
		}

		// 构建每页摘要供模型检查
		var pageSummaries []string
		for i, s := range slides {
			title, _ := s["title"].(string)
			content, _ := s["content"].(string)
			c := strings.TrimSpace(content)
			if len(c) > 80 {
				c = string([]rune(c)[:80]) + "..."
			}
			pageSummaries = append(pageSummaries, fmt.Sprintf("[%d] title=%s | content=%s", i+1, title, c))
		}

		systemPrompt := `You are a PPT slide quality filter.
Review each slide individually. Remove slides that are:
- Too short or empty content (<20 chars)
- Vague or generic without substance
- Duplicate of another slide
- Not relevant to the topic

Return JSON: {"qualified":bool, "keep_indices":[int], "removed":[{"index":int,"reason":str}]}
Keep at least 10 slides. Keep cover(1) and thank-you(last).`

		userMsg := fmt.Sprintf("Topic: %s\nStyle: %s\nTotal slides: %d\n\nSlides:\n%s\n\nReview each slide and decide which to keep.",
			wc.Topic, wc.Style, len(slides), strings.Join(pageSummaries, "\n"))

		messages := []*schema.Message{
			{Role: schema.System, Content: systemPrompt},
			{Role: schema.User, Content: userMsg},
		}

		resp, err := preCheckModel.Generate(ctx, messages)
		if err != nil {
			logger.Errorf("预检失败，保留全部: %v", err)
			return map[string]any{}, nil
		}

		// 解析检测结果
		var result struct {
			Qualified   bool `json:"qualified"`
			KeepIndices []int `json:"keep_indices"`
			Removed     []struct {
				Index  int    `json:"index"`
				Reason string `json:"reason"`
			} `json:"removed"`
		}
		if err := json.Unmarshal([]byte(resp.Content), &result); err != nil {
			logger.Warnf("预检结果解析失败，保留全部: %v", err)
			return map[string]any{}, nil
		}

		// 过滤：只保留合格的页（1-based → 0-based）
		if len(result.KeepIndices) == 0 && len(result.Removed) == 0 {
			return map[string]any{}, nil
		}

		keepSet := make(map[int]bool)
		for _, idx := range result.KeepIndices {
			if idx >= 1 && idx <= len(slides) {
				keepSet[idx-1] = true
			}
		}

		var keptSlides []map[string]any
		for i, s := range slides {
			if keepSet[i] {
				keptSlides = append(keptSlides, s)
			}
		}

		// 最少保留 10 页
		if len(keptSlides) < 10 {
			logger.Warnf("过滤后只剩 %d 页，低于下限，保留全部", len(keptSlides))
			return map[string]any{}, nil
		}

		// 更新 SlidesJSON
		filteredJSON, _ := json.Marshal(keptSlides)
		wc.SlidesJSON = string(filteredJSON)

		removedCount := len(slides) - len(keptSlides)
		logger.Infof("文案预检完成: %d→%d 页 (移除 %d 页不合格)",
			len(slides), len(keptSlides), removedCount)

		for _, r := range result.Removed {
			logger.Debugf("  移除 [%d]: %s", r.Index, r.Reason)
		}

		return map[string]any{}, nil
	})
}

func PreQualityCheckStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}
