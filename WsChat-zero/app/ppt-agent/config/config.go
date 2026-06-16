package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/spf13/viper"
	"ppt-agent/pkg/myfile"
)

var env string

func init() {
	flag.StringVar(&env, "env", "", "配置文件后缀（例如：local, prod）")
}

type Config struct {
	Server    ServerConfig    `yaml:"server" mapstructure:"server"`
	Database  DatabaseConfig  `yaml:"database" mapstructure:"database"`
	AI        AIConfig        `yaml:"ai" mapstructure:"ai"`
}

type ServerConfig struct {
	ConfigActive string `yaml:"config-active" mapstructure:"config-active"`
	Port         int    `yaml:"port" mapstructure:"port"`
	ContextPath  string `yaml:"context-path" mapstructure:"context-path"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
}

type AIConfig struct {
	TextModel   ModelConfig `yaml:"text-model" mapstructure:"text-model"`
	AgentModel  ModelConfig `yaml:"agent-model" mapstructure:"agent-model"`
	ImageModel  ModelConfig `yaml:"image-model" mapstructure:"image-model"`
}

type ModelConfig struct {
	BaseURL   string `yaml:"base-url" mapstructure:"base-url"`
	APIKey    string `yaml:"api-key" mapstructure:"api-key"`
	ModelName string `yaml:"model-name" mapstructure:"model-name"`
}

func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.DBName)
}

func mergeConfig(base, override interface{}) {
	baseValue := reflect.ValueOf(base).Elem()
	overrideValue := reflect.ValueOf(override).Elem()
	overrideType := overrideValue.Type()

	for i := 0; i < overrideValue.NumField(); i++ {
		fieldName := overrideType.Field(i).Name
		overrideField := overrideValue.Field(i)
		baseField := baseValue.FieldByName(fieldName)

		switch overrideField.Kind() {
		case reflect.String:
			if overrideField.String() != "" {
				baseField.SetString(overrideField.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if overrideField.Int() != 0 {
				baseField.SetInt(overrideField.Int())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if overrideField.Uint() != 0 {
				baseField.SetUint(overrideField.Uint())
			}
		case reflect.Float32, reflect.Float64:
			if overrideField.Float() != 0 {
				baseField.SetFloat(overrideField.Float())
			}
		case reflect.Bool:
			if overrideField.Bool() {
				baseField.SetBool(overrideField.Bool())
			}
		}
	}
}

func InitConfig() *Config {
	flag.Parse()

	projectRoot, err := myfile.GetProjectRoot()
	if err != nil {
		panic(fmt.Errorf("获取项目根目录失败: %w", err))
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(projectRoot, "config"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败: %w", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("解析配置文件失败: %w", err))
	}

	configActive := config.Server.ConfigActive
	if env != "" {
		configActive = env
	}

	if configActive != "" {
		viper.SetConfigName(fmt.Sprintf("config-%s", configActive))
		if err := viper.MergeInConfig(); err == nil {
			var overrideConfig Config
			if err := viper.Unmarshal(&overrideConfig); err == nil {
				mergeConfig(&config.Server, &overrideConfig.Server)
				mergeAIConfig(&config.AI, &overrideConfig.AI)
			}
		}
	}

	applyEnvOverrides(&config)

	return &config
}

func mergeAIConfig(base, override *AIConfig) {
	mergeConfig(&base.TextModel, &override.TextModel)
	mergeConfig(&base.AgentModel, &override.AgentModel)
	mergeConfig(&base.ImageModel, &override.ImageModel)
}

func applyEnvOverrides(cfg *Config) {
	overrideModelConfig(&cfg.AI.TextModel, "PPT_AGENT_TEXT_MODEL", "DEEPSEEK")
	overrideModelConfig(&cfg.AI.AgentModel, "PPT_AGENT_AGENT_MODEL", "AGNES")
	overrideModelConfig(&cfg.AI.ImageModel, "PPT_AGENT_IMAGE_MODEL", "AGNES_IMAGE")

	if cfg.AI.AgentModel.APIKey == "" {
		cfg.AI.AgentModel.APIKey = os.Getenv("AGNES_API_KEY")
	}
	if cfg.AI.ImageModel.APIKey == "" {
		cfg.AI.ImageModel.APIKey = os.Getenv("AGNES_API_KEY")
	}
}

func overrideModelConfig(target *ModelConfig, prefixes ...string) {
	for _, prefix := range prefixes {
		if value := os.Getenv(prefix + "_BASE_URL"); value != "" {
			target.BaseURL = value
		}
		if value := os.Getenv(prefix + "_API_KEY"); value != "" {
			target.APIKey = value
		}
		if value := os.Getenv(prefix + "_MODEL_NAME"); value != "" {
			target.ModelName = value
		}
	}
}
