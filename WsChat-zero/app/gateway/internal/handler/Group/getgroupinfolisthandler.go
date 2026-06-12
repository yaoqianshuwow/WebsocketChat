// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package Group

import (
	"net/http"

	"github.com/your-org/ws-chat-zero/app/gateway/internal/logic/Group"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGroupInfoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := Group.NewGetGroupInfoListLogic(r.Context(), svcCtx)
		resp, err := l.GetGroupInfoList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
