package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Kafka     KafkaConf
	Mysql     MysqlConf
	StoreMode string `json:",optional,default=kafka,options=kafka|direct"`
}

type KafkaConf struct {
	Brokers   []string
	ChatTopic string
}

type MysqlConf struct {
	DataSource string
}
