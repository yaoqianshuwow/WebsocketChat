package Message

import (
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	filepb "github.com/your-org/ws-chat-zero/app/file/fileservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/golang-jwt/jwt/v4"
)

func parseToken(r *http.Request, secret string) int64 {
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		return 0
	}
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || t == nil || !t.Valid {
		return 0
	}
	uid, ok := claims["user_id"].(float64)
	if !ok {
		return 0
	}
	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
		return 0
	}
	return int64(uid)
}

func mimeByExt(name string) string {
	ext := strings.ToLower(path.Ext(name))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".mp4":
		return "video/mp4"
	case ".pdf":
		return "application/pdf"
	case ".mp3":
		return "audio/mpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

func ViewFileHandlerWithAuth(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext) {
	uid := parseToken(r, svcCtx.Config.JwtAuth.AccessSecret)
	if uid <= 0 {
		http.Error(w, `{"code":401,"message":"未授权"}`, http.StatusUnauthorized)
		return
	}

	fileUrl := r.URL.Query().Get("fileUrl")
	if fileUrl == "" {
		http.Error(w, `{"code":1,"message":"fileUrl 不能为空"}`, http.StatusBadRequest)
		return
	}

	// Try multiple URL formats to match DB records
	candidates := []string{
		strings.ReplaceAll(fileUrl, "127.0.0.1", "__SERVER_IP__"),
		strings.ReplaceAll(fileUrl, "localhost", "__SERVER_IP__"),
		fileUrl,
	}
	var resp *filepb.GetFileResponse
	var err error
	for _, url := range candidates {
		resp, err = svcCtx.FileClient.GetFile(r.Context(), &filepb.GetFileRequest{FileUrl: url})
		if err == nil && resp.Code == 0 {
			break
		}
	}
	if err != nil || resp.Code != 0 {
		msg := "文件不存在"
		if err != nil {
			msg = err.Error()
		} else if resp.Message != "" {
			msg = resp.Message
		}
		http.Error(w, `{"code":1,"message":"`+msg+`"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", mimeByExt(resp.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(resp.FileSize, 10))
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write(resp.Data)
}
