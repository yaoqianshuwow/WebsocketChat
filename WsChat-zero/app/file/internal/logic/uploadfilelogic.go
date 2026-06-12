package logic

import (
	"context"

	"github.com/your-org/ws-chat-zero/app/file/internal/model"
	"github.com/your-org/ws-chat-zero/app/file/internal/svc"
	"github.com/your-org/ws-chat-zero/app/file/pb/pb"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadFileLogic) UploadFile(in *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	now := time.Now()
	dateDir := now.Format("2006/01/02")

	// 生成唯一文件名
	ext := filepath.Ext(in.FileName)
	if ext == "" {
		ext = ".bin"
	}
	saveName := uuid.New().String() + ext

	// 存储目录
	storeDir := filepath.Join(l.svcCtx.Config.FileConfig.StoragePath, dateDir)
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return &pb.UploadFileResponse{Code: 1, Message: "创建目录失败"}, nil
	}

	// 写文件
	savePath := filepath.Join(storeDir, saveName)
	if err := os.WriteFile(savePath, in.Data, 0644); err != nil {
		return &pb.UploadFileResponse{Code: 1, Message: "保存文件失败"}, nil
	}

	fileUrl := l.svcCtx.Config.FileConfig.BaseUrl + "/" + dateDir + "/" + saveName

	// 确定MIME类型
	mimeType := detectMimeType(in.FileName)
	fileSize := int64(len(in.Data))

	// 保存记录到DB
	record := model.FileRecord{
		UserId:   in.UserId,
		FileName: in.FileName,
		FilePath: savePath,
		FileUrl:  fileUrl,
		FileType: in.FileType,
		FileSize: fileSize,
		MimeType: mimeType,
	}
	if err := l.svcCtx.DB.Create(&record).Error; err != nil {
		return &pb.UploadFileResponse{Code: 1, Message: "保存记录失败"}, nil
	}

	return &pb.UploadFileResponse{
		Code:     0,
		Message:  "上传成功",
		FileUrl:  fileUrl,
		FileSize: fileSize,
	}, nil
}

func detectMimeType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".pdf":
		return "application/pdf"
	case ".doc", ".docx":
		return "application/msword"
	case ".zip":
		return "application/zip"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
