package model

import "time"

type FileRecord struct {
	Id       int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserId   int64     `gorm:"column:user_id;index;not null"`
	FileName string    `gorm:"column:file_name;size:255;not null"`
	FilePath string    `gorm:"column:file_path;size:512;not null"`
	FileUrl  string    `gorm:"column:file_url;size:512;not null"`
	FileType int32     `gorm:"column:file_type;default:0"`
	FileSize int64     `gorm:"column:file_size;default:0"`
	MimeType string    `gorm:"column:mime_type;size:128"`
	Status   int32     `gorm:"column:status;default:0"` // 0=正常 1=已删除
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (FileRecord) TableName() string {
	return "file_record"
}
