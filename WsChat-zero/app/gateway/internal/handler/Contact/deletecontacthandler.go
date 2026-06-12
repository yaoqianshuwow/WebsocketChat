// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Contact

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Contact"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteContactHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ContactInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := Contact.NewDeleteContactLogic(r.Context(), svcCtx)
		resp, err := l.DeleteContact(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
