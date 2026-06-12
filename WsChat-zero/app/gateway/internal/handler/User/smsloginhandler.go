// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package User

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/User"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SmsLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SmsLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := User.NewSmsLoginLogic(r.Context(), svcCtx)
		resp, err := l.SmsLogin(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
