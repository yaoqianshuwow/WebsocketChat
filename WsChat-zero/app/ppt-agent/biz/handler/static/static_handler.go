package static

import (
	"context"
	"path/filepath"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"ppt-agent/pkg/myfile"
)

// ServeFrontend 提供首页（页面1：对话输入页）
func ServeFrontend(ctx context.Context, c *app.RequestContext) {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		c.String(consts.StatusInternalServerError, "project root not found")
		return
	}
	c.File(filepath.Join(root, "frontend", "index.html"))
}

// ServeWorkspace 提供工作台（页面2：对话+PPT预览）
func ServeWorkspace(ctx context.Context, c *app.RequestContext) {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		c.String(consts.StatusInternalServerError, "project root not found")
		return
	}
	c.File(filepath.Join(root, "frontend", "workspace.html"))
}

// ServeFootballNews 提供足球新闻页面
func ServeFootballNews(ctx context.Context, c *app.RequestContext) {
	root, err := myfile.GetProjectRoot()
	if err != nil {
		c.String(consts.StatusInternalServerError, "project root not found")
		return
	}
	c.File(filepath.Join(root, "frontend", "football-news.html"))
}
