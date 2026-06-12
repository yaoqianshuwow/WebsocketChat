package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/file/internal/model"
	"github.com/your-org/ws-chat-zero/app/file/internal/svc"
	"github.com/your-org/ws-chat-zero/app/file/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *pb.DeleteFileRequest) (*pb.CommonResponse, error) {
	var record model.FileRecord
	if err := l.svcCtx.DB.Where("file_url = ? AND user_id = ?", in.FileUrl, in.UserId).First(&record).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "文件不存在"}, nil
	}

	// 标记为已删除
	record.Status = 1
	if err := l.svcCtx.DB.Save(&record).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Message: "删除失败"}, nil
	}

	// 删除物理文件（忽略错误）
	os.Remove(record.FilePath)

	return &pb.CommonResponse{Code: 0, Message: "删除成功"}, nil
}
