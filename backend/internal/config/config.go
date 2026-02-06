package config

/**
 * Config 配置管理模块
 *
 * 本模块负责加载和管理应用程序的所有配置。
 * 使用 Viper 框架，支持：
 * - YAML 配置文件
 * - 环境变量覆盖
 * - 命令行参数
 * - 配置热重载
 *
 * 设计特点：
 * 1. 全局单例模式：通过全局变量 Config 访问配置
 * 2. 层次化配置：支持配置分组（如 database、redis、app 等）
 * 3. 环境变量映射：配置项可以通过环境变量覆盖
 * 4. 热重载支持：运行时重新加载配置
 *
 * 使用示例：
 *   // 获取数据库配置
 *   dbConfig := config.GetDatabase()
 *   host := dbConfig.GetString("host")
 *   port := dbConfig.GetInt("port")
 *
 *   // 获取应用配置
 *   appConfig := config.GetApp()
 *   mode := appConfig.GetString("mode")
 */

import (
	"fmt"       // 格式化错误信息
	"strings"   // 字符串处理，用于环境变量替换

	"github.com/spf13/viper" // Viper 配置管理框架
)

/**
 * Init 初始化配置系统
 *
 * 这是配置模块的入口函数，必须在程序启动时最先调用。
 * 因为其他所有模块（数据库、日志、Redis等）都依赖配置。
 *
 * 执行步骤：
 * 1. 创建新的 Viper 实例
 * 2. 设置配置文件路径和类型
 * 3. 配置环境变量映射规则
 * 4. 读取并解析配置文件
 *
 * 环境变量映射规则：
 * - 前缀：GOMALL_
 * - 分隔符：配置中的 "." 替换为 "_"
 * - 例如：database.host -> GOMALL_DATABASE_HOST
 *
 * 参数：
 *   configPath string - 配置文件路径，如 "conf/config.yaml"
 *
 * 返回值：
 *   error - 加载失败时返回错误，成功时返回 nil
 *
 * 使用示例：
 *   if err := config.Init("conf/config.yaml"); err != nil {
 *       panic(err)
 *   }
 */
func Init(configPath string) error {
	// 1. 创建 Viper 实例
	// Viper 是应用程序的"单例"，一个应用只需要一个实例
	Config = viper.New()

	// 2. 设置配置文件
	// SetConfigFile 指定配置文件路径（包含文件名）
	Config.SetConfigFile(configPath)
	// SetConfigType 指定配置文件类型（yaml、json、toml 等）
	Config.SetConfigType("yaml")

	// 3. 配置环境变量映射
	// SetEnvPrefix 设置环境变量前缀
	// 例如：database.host 会映射到 GOMALL_DATABASE_HOST
	Config.SetEnvPrefix("GOMALL")
	// SetEnvKeyReplacer 配置字符替换规则
	// 将配置键中的 "." 替换为 "_"
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// AutomaticEnv 启用自动环境变量读取
	// 没有在配置文件中设置的项，会自动尝试从环境变量读取
	Config.AutomaticEnv()

	// 4. 读取配置文件
	// ReadInConfig 读取并解析配置文件
	// 如果配置文件不存在或格式错误，会返回错误
	if err := Config.ReadInConfig(); err != nil {
		// 使用 %w 包装原始错误，保留错误链
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	return nil
}

/**
 * Config 全局配置对象
 *
 * 这是一个包级全局变量，存储当前应用的配置实例。
 * 所有配置访问函数（如 GetDatabase、GetRedis）都使用这个变量。
 *
 * 注意事项：
 * - 在调用 Init() 之前不能访问此变量
 * - Init() 会初始化此变量
 * - 其他包可以通过 config.Config 访问，但建议使用封装好的 GetXxx() 函数
 */
var Config *viper.Viper

/**
 * GetDatabase 获取数据库配置子项
 *
 * 返回数据库配置组（database 节）的配置对象。
 * 使用 Sub() 方法可以获取配置中的子节。
 *
 * 返回值：
 *   *viper.Viper - 数据库配置对象
 *
 * 配置项示例（config.yaml）：
 *   database:
 *     host: "localhost"
 *     port: 3306
 *     username: "root"
 *     password: "password"
 *     name: "gomall"
 *     max_idle_conns: 10
 *     max_open_conns: 100
 *
 * 使用示例：
 *   dbConfig := config.GetDatabase()
 *   host := dbConfig.GetString("host")
 *   port := dbConfig.GetInt("port")
 *   maxConns := dbConfig.GetInt("max_open_conns")
 */
func GetDatabase() *viper.Viper {
	// Sub() 返回配置中指定子节的新 viper 实例
	return Config.Sub("database")
}

/**
 * GetRedis 获取 Redis 配置子项
 *
 * 返回 Redis 配置组（redis 节）的配置对象。
 *
 * 返回值：
 *   *viper.Viper - Redis 配置对象
 *
 * 配置项示例（config.yaml）：
 *   redis:
 *     host: "localhost"
 *     port: 6379
 *     password: ""
 *     db: 0
 *     pool_size: 100
 *
 * 使用示例：
 *   redisConfig := config.GetRedis()
 *   host := redisConfig.GetString("host")
 *   port := redisConfig.GetInt("port")
 *   poolSize := redisConfig.GetInt("pool_size")
 */
func GetRedis() *viper.Viper {
	return Config.Sub("redis")
}

/**
 * GetApp 获取应用配置子项
 *
 * 返回应用配置组（app 节）的配置对象。
 * 包含应用运行的基本参数。
 *
 * 返回值：
 *   *viper.Viper - 应用配置对象
 *
 * 配置项示例（config.yaml）：
 *   app:
 *     host: "0.0.0.0"
 *     port: 8080
 *     mode: "debug"
 *     admin_ids: "1"
 *
 * 使用示例：
 *   appConfig := config.GetApp()
 *   host := appConfig.GetString("host")
 *   port := appConfig.GetInt("port")
 *   mode := appConfig.GetString("mode")
 */
func GetApp() *viper.Viper {
	return Config.Sub("app")
}

/**
 * GetJWT 获取 JWT 配置子项
 *
 * 返回 JWT 配置组（jwt 节）的配置对象。
 * 包含 JWT 令牌的密钥和过期时间等配置。
 *
 * 返回值：
 *   *viper.Viper - JWT 配置对象
 *
 * 配置项示例（config.yaml）：
 *   jwt:
 *     secret: "your-secret-key"
 *     expire_hours: 24
 *     refresh_hours: 168
 *
 * 使用示例：
 *   jwtConfig := config.GetJWT()
 *   secret := jwtConfig.GetString("secret")
 *   expireHours := jwtConfig.GetInt("expire_hours")
 */
func GetJWT() *viper.Viper {
	return Config.Sub("jwt")
}

/**
 * GetRabbitMQ 获取 RabbitMQ 配置子项
 *
 * 返回 RabbitMQ 配置组（rabbitmq 节）的配置对象。
 *
 * 返回值：
 *   *viper.Viper - RabbitMQ 配置对象
 *
 * 配置项示例（config.yaml）：
 *   rabbitmq:
 *     host: "localhost"
 *     port: 5672
 *     username: "guest"
 *     password: "guest"
 *     queue_prefix: "gomall_"
 *
 * 使用示例：
 *   mqConfig := config.GetRabbitMQ()
 *   host := mqConfig.GetString("host")
 *   port := mqConfig.GetInt("port")
 */
func GetRabbitMQ() *viper.Viper {
	return Config.Sub("rabbitmq")
}

/**
 * GetGRPCConfig 获取 gRPC 配置子项
 *
 * 返回 gRPC 配置组（grpc 节）的配置对象。
 *
 * 返回值：
 *   *viper.Viper - gRPC 配置对象
 *
 * 配置项示例（config.yaml）：
 *   grpc:
 *     port: 50051
 *     enabled: false
 *
 * 使用示例：
 *   grpcConfig := config.GetGRPCConfig()
 *   port := grpcConfig.GetInt("port")
 *   enabled := grpcConfig.GetBool("enabled")
 */
func GetGRPCConfig() *viper.Viper {
	return Config.Sub("grpc")
}

/**
 * GetTracing 获取链路追踪配置子项
 *
 * 返回链路追踪配置组（tracing 节）的配置对象。
 * 用于配置 OpenTelemetry + Jaeger 链路追踪。
 *
 * 返回值：
 *   *viper.Viper - 链路追踪配置对象
 *
 * 配置项示例（config.yaml）：
 *   tracing:
 *     enabled: false
 *     service_name: "gomall"
 *     jaeger_endpoint: "localhost:4317"
 *
 * 使用示例：
 *   tracingConfig := config.GetTracing()
 *   enabled := tracingConfig.GetBool("enabled")
 *   serviceName := tracingConfig.GetString("service_name")
 */
func GetTracing() *viper.Viper {
	return Config.Sub("tracing")
}

/**
 * GetRateLimit 获取限流配置子项
 *
 * 返回限流配置组（ratelimit 节）的配置对象。
 *
 * 返回值：
 *   *viper.Viper - 限流配置对象
 *
 * 配置项示例（config.yaml）：
 *   ratelimit:
 *     enabled: true
 *     global_rate: 1000
 *     global_burst: 2000
 *     api_rate: 100
 *     seckill_rate: 5
 *     login_rate: 10
 *
 * 使用示例：
 *   limitConfig := config.GetRateLimit()
 *   enabled := limitConfig.GetBool("enabled")
 *   globalRate := limitConfig.GetInt("global_rate")
 */
func GetRateLimit() *viper.Viper {
	return Config.Sub("ratelimit")
}

/**
 * GetLogger 获取日志配置子项
 *
 * 返回日志配置组（logger 节）的配置对象。
 *
 * 返回值：
 *   *viper.Viper - 日志配置对象
 *
 * 配置项示例（config.yaml）：
 *   logger:
 *     level: "info"
 *     format: "json"
 *     output: "stdout"
 *     filename: "app.log"
 *
 * 使用示例：
 *   loggerConfig := config.GetLogger()
 *   level := loggerConfig.GetString("level")
 *   format := loggerConfig.GetString("format")
 */
func GetLogger() *viper.Viper {
	return Config.Sub("logger")
}

/**
 * Reload 重新加载配置文件
 *
 * 此函数用于实现配置热更新。
 * 调用后会重新读取配置文件，更新内存中的配置。
 *
 * 使用场景：
 * - 收到 SIGHUP 信号时重新加载配置
 * - 动态调整日志级别
 *
 * 返回值：
 *   error - 重新加载失败时返回错误
 *
 * 注意事项：
 * - 只会重新加载配置文件，不会重新加载环境变量
 * - 某些配置（如数据库连接）需要重启才能生效
 *
 * 使用示例：
 *   if err := config.Reload(); err != nil {
 *       log.Error("配置重新加载失败", err)
 *   }
 */
func Reload() error {
	// 检查 Config 是否已初始化
	if Config == nil {
		return fmt.Errorf("配置未初始化")
	}
	// ReadInConfig 会重新读取并解析配置文件
	return Config.ReadInConfig()
}
