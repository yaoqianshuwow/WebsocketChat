package ppt

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"ppt-agent/biz/model/api"
	"ppt-agent/biz/model/vo"
	"ppt-agent/biz/monitor"
	"ppt-agent/biz/service/ppt"
	pkg "ppt-agent/pkg/errors"
	"ppt-agent/pkg/myfile"
)

type PptHandler struct {
	pptService ppt.IPptService
	monitor    *monitor.TokenMonitor
}

func NewPptHandler(pptService ppt.IPptService, monitor *monitor.TokenMonitor) *PptHandler {
	return &PptHandler{
		pptService: pptService,
		monitor:    monitor,
	}
}

func (h *PptHandler) GeneratePpt(ctx context.Context, c *app.RequestContext) {
	var req vo.PptGenerateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("请求参数格式错误")))
		return
	}

	if strings.TrimSpace(req.Topic) == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("PPT 主题不能为空")))
		return
	}

	logger.Infof("receive generate ppt request: topic=%s, style=%s", req.Topic, req.Style)

	result, err := h.pptService.GeneratePpt(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(result))
}

func (h *PptHandler) GetPptStatus(ctx context.Context, c *app.RequestContext) {
	taskID := c.Query("taskId")
	if taskID == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("任务 ID 不能为空")))
		return
	}

	result, err := h.pptService.GetPptStatus(ctx, taskID)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(result))
}

func (h *PptHandler) DownloadPpt(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("任务 ID 不能为空")))
		return
	}

	downloadInfo, err := h.pptService.DownloadPpt(ctx, taskID)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	filePath := downloadInfo.FileURL
	if !filepath.IsAbs(filePath) {
		root, _ := myfile.GetProjectRoot()
		filePath = filepath.Join(root, "output", filepath.Base(filePath))
	}

	c.File(filePath)
}

func (h *PptHandler) GetTokenUsage(ctx context.Context, c *app.RequestContext) {
	summary := h.monitor.GetSummary()
	c.JSON(consts.StatusOK, api.NewSuccessResponse(summary))
}

func (h *PptHandler) PreviewPpt(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("taskId")
	if taskID == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("任务 ID 不能为空")))
		return
	}

	result, err := h.pptService.GetPptPreview(ctx, taskID)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(result))
}

func (h *PptHandler) ListTasks(ctx context.Context, c *app.RequestContext) {
	result, err := h.pptService.ListTasks(ctx)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(result))
}
