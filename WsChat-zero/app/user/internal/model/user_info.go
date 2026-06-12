package model

import "time"

type UserInfo struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Username  string    `gorm:"column:username;size:64;uniqueIndex"`
	Password  string    `gorm:"column:password;size:128"`
	Phone     string    `gorm:"column:phone;size:20"`
	Avatar    string    `gorm:"column:avatar;size:255"`
	Nickname  string    `gorm:"column:nickname;size:64"`
	Sex       string    `gorm:"column:sex;size:10"`
	Age       int       `gorm:"column:age;default:0"`
	Bio       string    `gorm:"column:bio;size:255"`
	Status    int       `gorm:"column:status;default:0"` // 0=正常 1=禁用
	Role      int       `gorm:"column:role;default:0"`   // 0=普通 1=管理员
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
