package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql      MysqlConf
	FileConfig FileConfig
}

type MysqlConf struct {
	DataSource string
}

type FileConfig struct {
	StoragePath string // 文件存储根路径
	BaseUrl     string // 文件访问基础URL
}
