package svc

import (
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/config"
	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	KafkaWriter *kafka.Writer
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := openDBWithRetry(c.Mysql.DataSource)

	_ = db.AutoMigrate(&model.Message{})

	var writer *kafka.Writer
	if c.StoreMode != "direct" {
		writer = &kafka.Writer{
			Addr:     kafka.TCP(c.Kafka.Brokers...),
			Topic:    c.Kafka.ChatTopic,
			Balancer: &kafka.LeastBytes{},
		}
	}

	return &ServiceContext{
		Config:      c,
		DB:          db,
		KafkaWriter: writer,
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
