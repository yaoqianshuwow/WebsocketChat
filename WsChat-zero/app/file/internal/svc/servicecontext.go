package svc

import (
	"github.com/your-org/ws-chat-zero/app/file/internal/config"
	"github.com/your-org/ws-chat-zero/app/file/internal/model"
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

	_ = db.AutoMigrate(&model.FileRecord{})

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
