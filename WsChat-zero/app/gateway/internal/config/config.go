package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	UserRpc   RpcConf
	MsgRpc    RpcConf
	FriendRpc RpcConf
	FileRpc   RpcConf
	Redis     RedisConf
	JwtAuth   JwtAuthConf
	Mysql     MysqlConf
}

type MysqlConf struct {
	DataSource string
}

type RpcConf struct {
	Target string          `json:",optional"`
	Etcd   discov.EtcdConf `json:",optional,inherit"`
}

type RedisConf struct {
	Host string
	Pass string
}

type JwtAuthConf struct {
	AccessSecret string
}
