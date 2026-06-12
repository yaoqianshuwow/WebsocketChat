package svc

import (
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/config"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移消息表
	_ = db.AutoMigrate(&model.Message{})

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
