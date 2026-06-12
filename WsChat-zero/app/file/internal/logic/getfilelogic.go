package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/file/internal/model"
	"github.com/your-org/ws-chat-zero/app/file/internal/svc"
	"github.com/your-org/ws-chat-zero/app/file/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type GetFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileLogic {
	return &GetFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileLogic) GetFile(in *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	var record model.FileRecord
	if err := l.svcCtx.DB.Where("file_url = ? AND status = 0", in.FileUrl).First(&record).Error; err != nil {
		return &pb.GetFileResponse{Code: 1, Message: "文件不存在"}, nil
	}

	data, err := os.ReadFile(record.FilePath)
	if err != nil {
		return &pb.GetFileResponse{Code: 1, Message: "读取文件失败"}, nil
	}

	return &pb.GetFileResponse{
		Code:     0,
		Message:  "ok",
		Data:     data,
		FileName: record.FileName,
		FileSize: record.FileSize,
	}, nil
}
