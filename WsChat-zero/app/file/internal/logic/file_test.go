package logic

import (
	"testing"

	"github.com/your-org/ws-chat-zero/app/file/internal/model"
)

func TestFileModel(t *testing.T) {
	r := &model.FileRecord{
		Id:       1,
		UserId:   100,
		FileName: "test.jpg",
		FileUrl:  "/files/test.jpg",
		FileType: 1,
		FileSize: 1024,
		MimeType: "image/jpeg",
		Status:   0,
	}

	if r.TableName() != "file_record" {
		t.Errorf("expected table=file_record, got %s", r.TableName())
	}
	if r.FileName != "test.jpg" {
		t.Errorf("expected FileName=test.jpg, got %s", r.FileName)
	}
}

func TestDetectMimeType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"test.jpg", "image/jpeg"},
		{"test.jpeg", "image/jpeg"},
		{"test.png", "image/png"},
		{"test.gif", "image/gif"},
		{"test.mp4", "video/mp4"},
		{"test.mp3", "audio/mpeg"},
		{"test.pdf", "application/pdf"},
		{"test.zip", "application/zip"},
		{"test.txt", "text/plain"},
		{"unknown.xyz", "application/octet-stream"},
		{"noext", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectMimeType(tt.name)
			if got != tt.expected {
				t.Errorf("detectMimeType(%s) = %s, want %s", tt.name, got, tt.expected)
			}
		})
	}
}
