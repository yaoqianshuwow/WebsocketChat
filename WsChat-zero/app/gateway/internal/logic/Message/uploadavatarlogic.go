package Message

import (
	"context"
	"io"
	"net/http"

	"github.com/your-org/ws-chat-zero/app/file/fileservice"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/svc"
	"github.com/your-org/ws-chat-zero/app/gateway/internal/types"
	"github.com/your-org/ws-chat-zero/app/user/userservice"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAvatarLogic struct{ logx.Logger; ctx context.Context; svcCtx *svc.ServiceContext }

func NewUploadAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAvatarLogic {
	return &UploadAvatarLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UploadAvatarLogic) UploadAvatar(r *http.Request) (resp *types.CommonResp, err error) {
	// 限制最大 5MB
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		return &types.CommonResp{Code: 1, Message: "解析表单失败"}, nil
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		return &types.CommonResp{Code: 1, Message: "获取头像文件失败"}, nil
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return &types.CommonResp{Code: 1, Message: "读取文件失败"}, nil
	}

	userId := l.ctx.Value("userId").(int64)

	// 调用 file RPC 上传
	fileResp, err := l.svcCtx.FileClient.UploadFile(l.ctx, &fileservice.UploadFileRequest{
		Data:     data,
		FileName: header.Filename,
		FileType: 2, // 头像类型
		UserId:   userId,
	})
	if err != nil || fileResp.Code != 0 {
		msg := "上传头像失败"
		if err != nil {
			msg = err.Error()
		} else if fileResp.Message != "" {
			msg = fileResp.Message
		}
		return &types.CommonResp{Code: 1, Message: msg}, nil
	}

	// 更新用户头像
	_, err = l.svcCtx.UserClient.UpdateUserInfo(l.ctx, &userservice.UpdateUserInfoRequest{
		UserId: userId,
		Avatar: fileResp.FileUrl,
	})
	if err != nil {
		return &types.CommonResp{Code: 1, Message: "更新头像失败"}, nil
	}

	return &types.CommonResp{Code: 0, Message: "头像上传成功"}, nil
}
