package model

import "time"

// Message 兼容 message 表的完整结构
type Message struct {
	Id    int64  `gorm:"column:id;primaryKey;autoIncrement"`
	MsgId *int64 `gorm:"column:msg_id;uniqueIndex"`

	Uuid       string     `gorm:"column:uuid;size:20;not null"`
	SessionId  int64      `gorm:"column:session_id;index;not null"`
	Type       int32      `gorm:"column:type;default:1;not null"`
	Content    string     `gorm:"column:content;type:text"`
	Url        string     `gorm:"column:url;size:255"`
	SendID     string     `gorm:"column:send_id;size:20;index;not null"`
	SendName   string     `gorm:"column:send_name;size:64;not null"`
	SendAvatar string     `gorm:"column:send_avatar;size:255;not null"`
	ReceiveID  string     `gorm:"column:receive_id;size:20;index;not null"`
	FileType   string     `gorm:"column:file_type;size:10"`
	FileName   string     `gorm:"column:file_name;size:255"`
	FileSize   int64      `gorm:"column:file_size;default:0"`
	Status     int32      `gorm:"column:status;default:0;not null"`
	CreatedAt  time.Time  `gorm:"column:created_at"`
	SendAt     *time.Time `gorm:"column:send_at"`
	AvData     string     `gorm:"column:av_data;type:longtext"`

	SenderId   int64  `gorm:"column:sender_id;index"`
	ReceiverId int64  `gorm:"column:receiver_id;index"`
	ChatType   int32  `gorm:"column:chat_type;default:1"`
	MsgType    int32  `gorm:"column:msg_type;default:1"`
	FileUrl    string `gorm:"column:file_url;size:512"`
}

func (Message) TableName() string {
	return "message"
}
