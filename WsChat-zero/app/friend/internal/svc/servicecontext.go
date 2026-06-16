package svc

import (
	"time"

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
	db := openDBWithRetry(c.Mysql.DataSource)

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

func openDBWithRetry(dsn string) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, pingErr := db.DB()
			if pingErr == nil {
				pingErr = sqlDB.Ping()
			}
			if pingErr == nil {
				return db
			}
			err = pingErr
		}
		time.Sleep(time.Second)
	}

	panic(err)
}
