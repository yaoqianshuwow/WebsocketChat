package Message

import (
	"context"
	"io"
	"net/http"

	"github.com/your-org/ws-chat-zero/app/file/fileservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UploadFileLogic) UploadFile(r *http.Request) (resp *types.CommonResp, err error) {
	// 限制最大 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return &types.CommonResp{Code: 1, Message: "解析表单失败"}, nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return &types.CommonResp{Code: 1, Message: "获取文件失败"}, nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return &types.CommonResp{Code: 1, Message: "读取文件失败"}, nil
	}

	// 从 context 中获取 userId (由 Auth 中间件设置)
	userId := l.ctx.Value("userId").(int64)

	fileType := int32(0)
	if r.FormValue("fileType") != "" {
		// 可以解析 fileType
	}

	// 调用 file RPC
	fileResp, err := l.svcCtx.FileClient.UploadFile(l.ctx, &fileservice.UploadFileRequest{
		Data:     data,
		FileName: header.Filename,
		FileType: fileType,
		UserId:   userId,
	})
	if err != nil || fileResp.Code != 0 {
		msg := "上传失败"
		if err != nil {
			msg = err.Error()
		} else if fileResp.Message != "" {
			msg = fileResp.Message
		}
		return &types.CommonResp{Code: 1, Message: msg}, nil
	}

	return &types.CommonResp{Code: 0, Message: "上传成功"}, nil
}
