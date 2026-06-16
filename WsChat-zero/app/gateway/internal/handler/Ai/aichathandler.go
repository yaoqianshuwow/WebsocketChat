package Ai

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Ai"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AiChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AiChatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.OkJson(w, &types.AiChatResp{Code: 1, Message: "参数错误"})
			return
		}
		l := Ai.NewAiChatLogic(r.Context(), svcCtx)
		resp, err := l.AiChat(&req)
		if err != nil {
			httpx.OkJson(w, &types.AiChatResp{Code: 1, Message: err.Error()})
			return
		}
		httpx.OkJson(w, resp)
	}
}
