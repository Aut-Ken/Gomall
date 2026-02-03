package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Init 初始化配置
func Init(configPath string) error {
	Config = viper.New()

	// 设置配置文件
	Config.SetConfigFile(configPath)
	Config.SetConfigType("yaml")

	// 读取配置（支持环境变量覆盖）
	Config.SetEnvPrefix("GOMALL")
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Config.AutomaticEnv()

	// 读取配置文件
	if err := Config.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	return nil
}

// Config 全局配置对象
var Config *viper.Viper

// GetDatabase 获取数据库配置
func GetDatabase() *viper.Viper {
	return Config.Sub("database")
}

// GetRedis 获取Redis配置
func GetRedis() *viper.Viper {
	return Config.Sub("redis")
}

// GetApp 获取应用配置
func GetApp() *viper.Viper {
	return Config.Sub("app")
}

// GetJWT 获取JWT配置
func GetJWT() *viper.Viper {
	return Config.Sub("jwt")
}

// GetRabbitMQ 获取RabbitMQ配置
func GetRabbitMQ() *viper.Viper {
	return Config.Sub("rabbitmq")
}

// GetGRPCConfig 获取gRPC配置
func GetGRPCConfig() *viper.Viper {
	return Config.Sub("grpc")
}
