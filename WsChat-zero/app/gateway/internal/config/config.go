package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	UserRpc   RpcConf
	MsgRpc    RpcConf
	FriendRpc RpcConf
	FileRpc   RpcConf
	Redis     RedisConf
	JwtAuth   JwtAuthConf
}

type RpcConf struct {
	Target string // 直连地址: localhost:9091
}

type RedisConf struct {
	Host string
	Pass string
}

type JwtAuthConf struct {
	AccessSecret string
}
