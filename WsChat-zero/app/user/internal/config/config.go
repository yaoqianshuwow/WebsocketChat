package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql       MysqlConf
	BizRedis    RedisConf // 业务 Redis（需加 Tag 避免与 RpcServerConf.Redis 冲突）
	JwtAuth     JwtAuthConf
	ES          ESConfig
}

type ESConfig struct {
	Addresses []string
	Index     string `json:",default=users"`
}

type MysqlConf struct {
	DataSource string
}

type RedisConf struct {
	Host string
	Pass string
}

type JwtAuthConf struct {
	AccessSecret string
	AccessExpire int64
}
