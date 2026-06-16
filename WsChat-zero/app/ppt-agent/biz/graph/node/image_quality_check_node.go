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

var imageCheckModel *llm.ChatModelWrapper

func InitImageQualityCheckNode(model *llm.ChatModelWrapper) {
	imageCheckModel = model
}

func NewImageQualityCheckNode() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		logger.Info("执行节点: 图片质检(agnes过滤)")

		gs := state.GenGraphState(ctx)
		wc := state.GetContext(gs)
		if wc == nil {
			return map[string]any{}, nil
		}

		// 汇总所有图片信息
		type imgInfo struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			Path    string `json:"path"`
			Desc    string `json:"desc"`
		}
		var allImages []imgInfo

		for _, img := range wc.ContentImage {
			allImages = append(allImages, imgInfo{"配图", img.PageTitle, img.LocalPath, img.Description})
		}
		for _, img := range wc.Illustrations {
			allImages = append(allImages, imgInfo{"插画", img.PageTitle, img.LocalPath, img.Description})
		}
		for _, img := range wc.Diagrams {
			allImages = append(allImages, imgInfo{"图表", img.PageTitle, img.LocalPath, img.Description})
		}
		for _, img := range wc.Logos {
			allImages = append(allImages, imgInfo{"Logo", img.PageTitle, img.LocalPath, img.Description})
		}

		totalBefore := len(allImages)
		if totalBefore == 0 {
			logger.Info("无图片需要质检")
			return map[string]any{}, nil
		}

		// 构建图片摘要给 agnes 评估
		var imgSummaries []string
		for i, img := range allImages {
			status := "有文件"
			if img.Path == "" {
				status = "无文件"
			}
			imgSummaries = append(imgSummaries,
				fmt.Sprintf("[%d] type=%s title=%s desc=%s status=%s",
					i+1, img.Type, img.Title, truncateStr(img.Desc, 50), status))
		}

		systemPrompt := `You are an image quality reviewer.
Review each image entry. Remove entries that:
- Have no local file (status=无文件)
- Have vague or empty descriptions
- Are clearly irrelevant to their page title

Return JSON: {"keep_indices":[int], "removed":[{"index":int,"reason":str}]}
Indices are 1-based. Keep only high-quality, relevant images.`

		userMsg := fmt.Sprintf("Review these PPT images:\n%s\n\nDecide which to keep.",
			strings.Join(imgSummaries, "\n"))

		messages := []*schema.Message{
			{Role: schema.System, Content: systemPrompt},
			{Role: schema.User, Content: userMsg},
		}

		resp, err := imageCheckModel.Generate(ctx, messages)
		if err != nil {
			logger.Errorf("图片质检失败，保留全部: %v", err)
			return map[string]any{}, nil
		}

		var result struct {
			KeepIndices []int `json:"keep_indices"`
			Removed     []struct {
				Index  int    `json:"index"`
				Reason string `json:"reason"`
			} `json:"removed"`
		}
		if err := json.Unmarshal([]byte(resp.Content), &result); err != nil {
			logger.Warnf("图片质检结果解析失败，保留全部: %v", err)
			return map[string]any{}, nil
		}

		// 过滤
		keepSet := make(map[int]bool)
		for _, idx := range result.KeepIndices {
			if idx >= 1 && idx <= len(allImages) {
				keepSet[idx-1] = true
			}
		}

		filterByIndex := func(imgs []state.ImageInfo, startIdx int) ([]state.ImageInfo, int) {
			var kept []state.ImageInfo
			for i, img := range imgs {
				globalIdx := startIdx + i
				if keepSet[globalIdx] {
					kept = append(kept, img)
				}
			}
			return kept, startIdx + len(imgs)
		}

		idx := 0
		wc.ContentImage, idx = filterByIndex(wc.ContentImage, idx)
		wc.Illustrations, idx = filterByIndex(wc.Illustrations, idx)
		wc.Diagrams, idx = filterByIndex(wc.Diagrams, idx)
		wc.Logos, _ = filterByIndex(wc.Logos, idx)

		totalAfter := len(wc.ContentImage) + len(wc.Illustrations) +
			len(wc.Diagrams) + len(wc.Logos)
		removed := totalBefore - totalAfter

		logger.Infof("图片质检完成: %d→%d 张 (移除 %d)", totalBefore, totalAfter, removed)
		for _, r := range result.Removed {
			logger.Debugf("  移除 [%d]: %s", r.Index, r.Reason)
		}

		state.NotifyStepCompleted(wc, "图片质检")
		return map[string]any{}, nil
	})
}

func ImageQualityCheckStatePostHandler(ctx context.Context, output map[string]any, gs *state.GraphState) (map[string]any, error) {
	return output, nil
}

func truncateStr(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n]) + "..."
	}
	return s
}
