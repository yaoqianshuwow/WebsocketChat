package config

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Log logx.LogConf

	Kafka KafkaConfig
	Mysql MysqlConfig
	ES    ElasticsearchConfig
}

type KafkaConfig struct {
	Brokers    []string
	ChatTopic  string
	LoginTopic string
	LogoutTopic string
	GroupId    string
}

type MysqlConfig struct {
	DataSource string
}

type ElasticsearchConfig struct {
	Addresses []string
	Index     string
}
