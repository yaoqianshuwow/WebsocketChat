package Message

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Message"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRecentMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRecentMessagesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Message.NewGetRecentMessagesLogic(r.Context(), svcCtx)
		resp, err := l.GetRecentMessages(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
