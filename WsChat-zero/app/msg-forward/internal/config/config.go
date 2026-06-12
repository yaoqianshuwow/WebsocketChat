package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Kafka   KafkaConf
	Mysql   MysqlConf
}

type KafkaConf struct {
	Brokers   []string
	ChatTopic string
}

type MysqlConf struct {
	DataSource string
}
