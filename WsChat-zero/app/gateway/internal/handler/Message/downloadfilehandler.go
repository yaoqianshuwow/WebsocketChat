package Message

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	filepb "github.com/your-org/ws-chat-zero/app/file/fileservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

func parseUserIDForDownload(r *http.Request, accessSecret string) int64 {
	// 优先取 Authorization header
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		// 降级到查询参数
		token = r.URL.Query().Get("token")
	}
	if token == "" {
		return 0
	}

	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
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

func DownloadFileHandlerWithAuth(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext) {
	uid := parseUserIDForDownload(r, svcCtx.Config.JwtAuth.AccessSecret)
	if uid <= 0 {
		http.Error(w, `{"code":401,"message":"未授权"}`, http.StatusUnauthorized)
		return
	}
	_ = uid // 暂不记录谁下载了

	fileUrl := r.URL.Query().Get("fileUrl")
	if fileUrl == "" {
		fileUrl = r.FormValue("fileUrl")
	}
	if fileUrl == "" {
		http.Error(w, `{"code":1,"message":"fileUrl 不能为空"}`, http.StatusBadRequest)
		return
	}

	resp, err := svcCtx.FileClient.GetFile(r.Context(), &filepb.GetFileRequest{FileUrl: fileUrl})
	if err != nil || resp.Code != 0 {
		msg := "文件不存在"
		if err != nil {
			msg = err.Error()
		} else if resp.Message != "" {
			msg = resp.Message
		}
		logx.Errorf("download file failed: fileUrl=%s err=%s", fileUrl, msg)
		http.Error(w, `{"code":1,"message":"`+msg+`"}`, http.StatusNotFound)
		return
	}

	fileName := resp.FileName
	if fileName == "" {
		parts := strings.Split(fileUrl, "/")
		fileName = parts[len(parts)-1]
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(resp.FileSize, 10))
	w.Write(resp.Data)
}
