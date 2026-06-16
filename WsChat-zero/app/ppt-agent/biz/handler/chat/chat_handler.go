package chat

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"ppt-agent/biz/model/api"
	"ppt-agent/biz/model/vo"
	chatSvc "ppt-agent/biz/service/chat"
	sessionSvc "ppt-agent/biz/service/session"
	pkg "ppt-agent/pkg/errors"
)

type ChatHandler struct {
	chatService    chatSvc.IChatService
	sessionService sessionSvc.ISessionService
}

func NewChatHandler(chatService chatSvc.IChatService, sessionService sessionSvc.ISessionService) *ChatHandler {
	return &ChatHandler{chatService: chatService, sessionService: sessionService}
}

func (h *ChatHandler) SendMessage(ctx context.Context, c *app.RequestContext) {
	var req vo.ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("请求参数格式错误")))
		return
	}

	if strings.TrimSpace(req.Message) == "" {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.ParamsError.WithMessage("消息不能为空")))
		return
	}

	result, err := h.chatService.SendMessage(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(result))
}

func (h *ChatHandler) History(ctx context.Context, c *app.RequestContext) {
	sessionID := string(c.QueryArgs().Peek("sessionId"))
	messages, info, err := h.sessionService.GetHistory(ctx, sessionID)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}

	c.JSON(consts.StatusOK, api.NewSuccessResponse(map[string]any{
		"session":  info,
		"messages": messages,
	}))
}

func (h *ChatHandler) Sessions(ctx context.Context, c *app.RequestContext) {
	sessions, err := h.sessionService.ListSessions(ctx)
	if err != nil {
		c.JSON(consts.StatusOK, api.NewErrorResponse[any](pkg.SystemError.WithMessage(err.Error())))
		return
	}
	c.JSON(consts.StatusOK, api.NewSuccessResponse(map[string]any{
		"sessions": sessions,
	}))
}
