package svc

import (
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
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

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
