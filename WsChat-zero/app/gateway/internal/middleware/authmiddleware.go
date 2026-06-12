package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	userpb "github.com/your-org/ws-chat-zero/app/user/userservice"
)

type AuthMiddleware struct {
	userClient userpb.UserService
}

func NewAuthMiddleware(userClient userpb.UserService) *AuthMiddleware {
	return &AuthMiddleware{userClient: userClient}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 Header 或 Query 获取 token
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		// 去掉 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			httpx.OkJson(w, map[string]interface{}{
				"code":    401,
				"message": "未登录或 token 已过期",
			})
			return
		}

		// 调用 user 服务验证 token
		resp, err := m.userClient.ValidateToken(context.Background(), &userpb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil || resp.Code != 0 {
			logx.Errorf("auth failed: %v", err)
			httpx.OkJson(w, map[string]interface{}{
				"code":    401,
				"message": "token 无效",
			})
			return
		}

		// 将 userId 注入 context
		ctx := context.WithValue(r.Context(), "userId", resp.UserId)
		next(w, r.WithContext(ctx))
	}
}
