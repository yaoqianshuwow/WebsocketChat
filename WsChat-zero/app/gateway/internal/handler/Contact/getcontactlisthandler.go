// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Contact

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Contact"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetContactListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Contact.NewGetContactListLogic(r.Context(), svcCtx)
		resp, err := l.GetContactList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
