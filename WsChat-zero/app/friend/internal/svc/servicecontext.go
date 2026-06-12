package svc

import (
	"github.com/your-org/ws-chat-zero/app/friend/internal/config"
	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
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

	// 自动迁移表结构
	_ = db.AutoMigrate(
		&model.Contact{},
		&model.ContactApply{},
		&model.Session{},
		&model.GroupInfo{},
		&model.GroupMember{},
	)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
