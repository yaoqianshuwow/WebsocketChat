// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package UserInfo

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/UserInfo"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DisableUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := UserInfo.NewDisableUsersLogic(r.Context(), svcCtx)
		resp, err := l.DisableUsers(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
