package model

import "time"

// Message 聊天消息（与 msg-store 共用同一张表）
type Message struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement"`
	SenderId   int64     `gorm:"column:sender_id;index;not null"`
	ReceiverId int64     `gorm:"column:receiver_id;index;not null"`
	ChatType   int32     `gorm:"column:chat_type;default:1"`  // 1=单聊 2=群聊
	MsgType    int32     `gorm:"column:msg_type;default:1"`   // 1=文本 2=图片 3=文件 4=语音
	Content    string    `gorm:"column:content;type:text"`
	FileUrl    string    `gorm:"column:file_url;size:512"`
	FileSize   int64     `gorm:"column:file_size;default:0"`
	FileName   string    `gorm:"column:file_name;size:255"`
	Status     int32     `gorm:"column:status;default:0"`     // 0=正常 1=撤回 2=已删除
	SessionId  int64     `gorm:"column:session_id;index;not null"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (Message) TableName() string {
	return "message"
}
