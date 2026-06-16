package question

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"ppt-agent/biz/model/api"
	"ppt-agent/biz/model/vo"
	questionSvc "ppt-agent/biz/service/question"
	pkg "ppt-agent/pkg/errors"
)

type QuestionHandler struct {
	service questionSvc.IQuestionService
}

func NewQuestionHandler(service questionSvc.IQuestionService) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (h *QuestionHandler) Save(ctx context.Context, c *app.RequestContext) {
	var req vo.QuestionSaveRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("请求参数格式错误")))
		return
	}
	if req.Content == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("内容不能为空")))
		return
	}
	if err := h.service.Save(ctx, req.SessionID, req.Content); err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}
	c.JSON(consts.StatusOK, api.NewSuccessResponse(map[string]any{"saved": true}))
}

func (h *QuestionHandler) List(ctx context.Context, c *app.RequestContext) {
	sessionID := string(c.QueryArgs().Peek("sessionId"))
	result, err := h.service.List(ctx, sessionID)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}
	c.JSON(consts.StatusOK, api.NewSuccessResponse(map[string]any{
		"questions": result,
	}))
}
