package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type MainConfig struct {
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
}

type MysqlConfig struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	User         string `toml:"user"`
	Password     string `toml:"password"`
	DatabaseName string `toml:"databaseName"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Password string `toml:"password"`
	Db       int    `toml:"db"`
}

type AuthCodeConfig struct {
	AccessKeyID     string `toml:"accessKeyID"`
	AccessKeySecret string `toml:"accessKeySecret"`
	SignName        string `toml:"signName"`
	TemplateCode    string `toml:"templateCode"`
}

type LogConfig struct {
	LogPath string `toml:"logPath"`
}

type KafkaConfig struct {
	MessageMode string        `toml:"messageMode"`
	HostPort    string        `toml:"hostPort"`
	LoginTopic  string        `toml:"loginTopic"`
	LogoutTopic string        `toml:"logoutTopic"`
	ChatTopic   string        `toml:"chatTopic"`
	Partition   int           `toml:"partition"`
	Timeout     time.Duration `toml:"timeout"`
}

type StaticSrcConfig struct {
	StaticAvatarPath string `toml:"staticAvatarPath"`
	StaticFilePath   string `toml:"staticFilePath"`
}

type Config struct {
	MainConfig      `toml:"mainConfig"`
	MysqlConfig     `toml:"mysqlConfig"`
	RedisConfig     `toml:"redisConfig"`
	AuthCodeConfig  `toml:"authCodeConfig"`
	LogConfig       `toml:"logConfig"`
	KafkaConfig     `toml:"kafkaConfig"`
	StaticSrcConfig `toml:"staticSrcConfig"`
}

var config *Config

func LoadConfig() error {
	// 尝试多个可能的配置文件路径
	possiblePaths := []string{
		// 1. 当前工作目录
		filepath.Join(".", "configs", "config.toml"),
		// 2. 项目根目录（相对于当前工作目录）
		filepath.Join("..", "configs", "config.toml"),
		filepath.Join("..", "..", "configs", "config.toml"),
		filepath.Join("..", "..", "..", "configs", "config.toml"),
	}

	// 3. 可执行文件所在目录及其父目录
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possiblePaths = append(possiblePaths,
			filepath.Join(exeDir, "configs", "config.toml"),
			filepath.Join(exeDir, "..", "configs", "config.toml"),
			filepath.Join(exeDir, "..", "..", "configs", "config.toml"),
		)
	}

	// 尝试所有可能的路径
	for _, configPath := range possiblePaths {
		if _, err := os.Stat(configPath); err == nil {
			// 找到配置文件，加载它
			if _, err := toml.DecodeFile(configPath, config); err == nil {
				return nil
			} else {
				log.Printf("Error decoding config file %s: %v", configPath, err)
			}
		}
	}

	// 如果所有路径都尝试过仍然失败，使用项目根目录的绝对路径作为最后尝试
	projectRoot := "C:\\Users\\28407\\Desktop\\BaiduSyncdisk\\goproject\\KamaChat-main"
	configPath := filepath.Join(projectRoot, "configs", "config.toml")
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatalf("Failed to load config file from any location: %v", err)
		return err
	}

	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = LoadConfig()
	}
	return config
}
