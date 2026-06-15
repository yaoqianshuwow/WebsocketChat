package svc

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/your-org/ws-chat-zero/app/user/internal/config"
	"github.com/your-org/ws-chat-zero/app/user/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
	ES     *elasticsearch.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.BizRedis.Host,
		Password: c.BizRedis.Pass,
		DB:       0,
	})

	db.AutoMigrate(&model.UserInfo{})

	// 初始化 ES 客户端（地址为空时不启用）
	var es *elasticsearch.Client
	if len(c.ES.Addresses) > 0 {
		es, err = elasticsearch.NewClient(elasticsearch.Config{
			Addresses: c.ES.Addresses,
		})
		if err != nil {
			panic(err)
		}
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  rdb,
		ES:     es,
	}
}
