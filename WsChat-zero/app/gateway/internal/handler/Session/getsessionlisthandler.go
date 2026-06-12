// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Session

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Session"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetSessionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Session.NewGetSessionListLogic(r.Context(), svcCtx)
		resp, err := l.GetSessionList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
