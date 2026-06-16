package graph

import (
	"context"
	"fmt"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/eino/compose"

	"ppt-agent/biz/graph/node"
	"ppt-agent/biz/graph/state"
)

// CreatePptWorkflow 创建完整 DAG 工作流
//
//                    ┌→ pre_quality_check ─┐
//                    │  (deepseek:过滤文案)  │
//                    └────────┬────────────┘
//                             │
//              ┌──────────────┼──────────────┐
//              ▼              ▼              ▼
//   content_image   illustration   diagram    logo
//   (agnes-image    (agnes-image   (agnes-image (agnes-image
//    +搜索兜底)      插画+描述)     图表)         Logo)
//              │              │              │
//              └──────────────┼──────────────┘
//                             ▼
//                    image_quality_check
//                   (deepseek:过滤不合格图片)
//                             │
//                             ▼
//                       ppt_merge
//                     (python-pptx)
//                             │
//                             ▼
//                     post_quality_check
//                    (deepseek:质检成品)
//                             │
//                             ▼
//                            END
//
func CreatePptWorkflow(ctx context.Context) (compose.Runnable[map[string]any, map[string]any], error) {
	g := compose.NewGraph[map[string]any, map[string]any](
		compose.WithGenLocalState(state.GenGraphState),
	)

	// 0. PPT计划 (agnes-2.0: 生成20页文案)
	g.AddLambdaNode("ppt_plan", node.NewPptPlanNode(),
		compose.WithNodeName("PPT计划"),
		compose.WithStatePostHandler(node.PptPlanStatePostHandler))

	// 1. 预检 (deepseek: 检查文案质量)
	g.AddLambdaNode("pre_quality_check", node.NewPreQualityCheckNode(),
		compose.WithNodeName("文案预检"),
		compose.WithStatePostHandler(node.PreQualityCheckStatePostHandler))

	// 2. 并行收集 (全部用 agnes-image)
	g.AddLambdaNode("content_image_collector", node.NewContentImageCollectorNode(),
		compose.WithNodeName("内容配图收集"),
		compose.WithStatePostHandler(node.ContentImageCollectorStatePostHandler))

	g.AddLambdaNode("illustration_collector", node.NewIllustrationCollectorNode(),
		compose.WithNodeName("插画收集"),
		compose.WithStatePostHandler(node.IllustrationCollectorStatePostHandler))

	g.AddLambdaNode("diagram_collector", node.NewDiagramCollectorNode(),
		compose.WithNodeName("图表生成"),
		compose.WithStatePostHandler(node.DiagramCollectorStatePostHandler))

	g.AddLambdaNode("logo_collector", node.NewLogoCollectorNode(),
		compose.WithNodeName("Logo生成"),
		compose.WithStatePostHandler(node.LogoCollectorStatePostHandler))

	// 3. 图片质检 (过滤不合格图片)
	g.AddLambdaNode("image_quality_check", node.NewImageQualityCheckNode(),
		compose.WithNodeName("图片质检"),
		compose.WithStatePostHandler(node.ImageQualityCheckStatePostHandler))

	// 4. 合并 (python-pptx)
	g.AddLambdaNode("ppt_merge", node.NewPptMergeNode(),
		compose.WithNodeName("PPT合并"),
		compose.WithStatePostHandler(node.PptMergeStatePostHandler))

	// 4. 终检 (deepseek: 检查成品)
	g.AddLambdaNode("post_quality_check", node.NewPostQualityCheckNode(),
		compose.WithNodeName("成品终检"),
		compose.WithStatePostHandler(node.PostQualityCheckStatePostHandler))

	// DAG 连线
	g.AddEdge(compose.START, "ppt_plan")
	g.AddEdge("ppt_plan", "pre_quality_check")
	g.AddEdge("pre_quality_check", "content_image_collector")
	g.AddEdge("pre_quality_check", "illustration_collector")
	g.AddEdge("pre_quality_check", "diagram_collector")
	g.AddEdge("pre_quality_check", "logo_collector")
	g.AddEdge("content_image_collector", "image_quality_check")
	g.AddEdge("illustration_collector", "image_quality_check")
	g.AddEdge("diagram_collector", "image_quality_check")
	g.AddEdge("logo_collector", "image_quality_check")
	g.AddEdge("image_quality_check", "ppt_merge")
	g.AddEdge("ppt_merge", "post_quality_check")
	g.AddEdge("post_quality_check", compose.END)

	runnable, err := g.Compile(ctx, compose.WithGraphName("PPT生成工作流"))
	if err != nil {
		return nil, fmt.Errorf("编译工作流失败: %w", err)
	}
	return runnable, nil
}

func ExecutePptWorkflow(ctx context.Context, topic, style string) (*state.WorkFlowContext, error) {
	runnable, err := CreatePptWorkflow(ctx)
	if err != nil {
		return nil, err
	}

	wc := &state.WorkFlowContext{Topic: topic, Style: style}
	logger.Infof("开始执行PPT工作流: topic=%s style=%s", topic, style)

	ctx = state.WithWorkflowContext(ctx, wc)
	_, err = runnable.Invoke(ctx, map[string]any{})
	if err != nil {
		return wc, fmt.Errorf("工作流失败: %w", err)
	}

	logger.Infof("工作流完成: filePath=%s", wc.PptFilePath)
	return wc, nil
}
