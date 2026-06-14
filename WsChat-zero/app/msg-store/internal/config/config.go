package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Kafka KafkaConfig
	Mysql MysqlConfig
	ES    ElasticsearchConfig
}

type KafkaConfig struct {
	Brokers     []string `json:"Brokers"`
	ChatTopic   string   `json:"ChatTopic"`
	LoginTopic  string   `json:"LoginTopic"`
	LogoutTopic string   `json:"LogoutTopic"`
	GroupId     string   `json:"GroupId"`
}

type MysqlConfig struct {
	DataSource string `json:"DataSource"`
}

type ElasticsearchConfig struct {
	Addresses []string `json:"Addresses"`
	Index     string   `json:"Index"`
}
