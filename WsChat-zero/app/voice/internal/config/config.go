package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	BaiduAsr BaiduAsrConf
}

type BaiduAsrConf struct {
	AppId     string // 百度 AI 应用 AppID
	ApiKey    string // 百度 AI API Key
	SecretKey string // 百度 AI Secret Key
}
