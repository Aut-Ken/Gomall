package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gomall/internal/config"
	"gomall/internal/database"
	"gomall/internal/logger"
	"gomall/internal/middleware"
	"gomall/internal/rabbitmq"
	redispkg "gomall/internal/redis"
	"gomall/internal/router"
	"gomall/internal/service"
	"gomall/internal/tracing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "conf/config.yaml", "配置文件路径")
	env := flag.String("env", "dev", "运行环境: dev, prod")
	flag.Parse()

	// 选择配置文件
	if *env == "prod" {
		*configPath = "conf/config-prod.yaml"
	} else {
		*configPath = "conf/config-dev.yaml"
	}

	logger.Info("初始化日志系统...", zap.String("env", *env))
	if err := logger.Init(); err != nil {
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 初始化配置
	logger.Info("正在初始化配置...", zap.String("config", *configPath))
	if err := config.Init(*configPath); err != nil {
		logger.Fatal("配置初始化失败", zap.Error(err))
	}
	logger.Info("配置初始化成功")

	// 初始化数据库
	logger.Info("正在连接数据库...")
	if err := database.Init(); err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	defer database.Close()
	logger.Info("数据库连接成功")

	// 初始化Redis（可选，用于秒杀场景）
	logger.Info("正在连接Redis...")
	if err := redispkg.Init(); err != nil {
		logger.Warn("Redis连接失败，将使用数据库方案", zap.Error(err))
	} else {
		defer redispkg.Close()
		logger.Info("Redis连接成功")
	}

	// 初始化RabbitMQ（可选，用于异步订单）
	logger.Info("正在连接RabbitMQ...")
	if err := rabbitmq.Init(); err != nil {
		logger.Warn("RabbitMQ连接失败，将使用同步方案", zap.Error(err))
	} else {
		defer rabbitmq.Close()
		logger.Info("RabbitMQ连接成功")
	}

	// 初始化链路追踪（可选）
	tracingConfig := config.GetTracing()
	if tracingConfig.GetBool("enabled") {
		logger.Info("正在初始化链路追踪...")
		serviceName := tracingConfig.GetString("service_name")
		jaegerEndpoint := tracingConfig.GetString("jaeger_endpoint")
		shutdown, err := tracing.InitTracing(serviceName, jaegerEndpoint)
		if err != nil {
			logger.Warn("链路追踪初始化失败", zap.Error(err))
		} else {
			defer func() {
				if err := shutdown(context.Background()); err != nil {
					logger.Error("链路追踪关闭失败", zap.Error(err))
				}
			}()
			logger.Info("链路追踪初始化成功")
		}
	}

	// 设置Gin运行模式
	ginMode := config.GetApp().GetString("mode")
	gin.SetMode(ginMode)

	// 创建Gin引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.MetricsMiddleware())

	// 设置路由
	router.Setup(r)

	// 启动秒杀订单处理协程（如果RabbitMQ可用）
	go func() {
		seckillSvc := service.NewSeckillService()
		seckillSvc.ProcessSeckillOrders()
	}()

	// 启动订单消费者（如果RabbitMQ可用）
	go func() {
		orderSvc := service.NewOrderService()
		orderSvc.StartOrderConsumer()
	}()

	// 获取服务配置
	appConfig := config.GetApp()
	host := appConfig.GetString("host")
	port := appConfig.GetInt("port")
	addr := fmt.Sprintf("%s:%d", host, port)

	// 优雅关闭 + 配置热更新
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				logger.Info("收到关闭信号，正在关闭服务...")
				os.Exit(0)
			case syscall.SIGHUP:
				logger.Info("收到 SIGHUP 信号，正在重新加载配置...")
				if err := config.Reload(); err != nil {
					logger.Error("配置重新加载失败", zap.Error(err))
				} else {
					logger.Info("配置重新加载成功")
					// 重新初始化日志级别
					if err := logger.Reload(); err != nil {
						logger.Error("日志配置重新加载失败", zap.Error(err))
					}
				}
			}
		}
	}()

	// 启动服务
	logger.Info("服务启动成功", zap.String("addr", addr))
	logger.Info("健康检查", zap.String("url", fmt.Sprintf("http://%s/health", addr)))
	logger.Info("API 文档", zap.String("url", fmt.Sprintf("http://%s/swagger/index.html", addr)))
	if err := r.Run(addr); err != nil {
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}
