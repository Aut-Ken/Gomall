package main

/**
 * Gomall 电商系统主入口
 *
 * 程序启动流程：
 * 1. 解析命令行参数（配置文件路径、运行环境）
 * 2. 初始化配置（必须最先执行，因为其他组件都依赖配置）
 * 3. 初始化日志系统
 * 4. 初始化数据库（MySQL + GORM）
 * 5. 初始化Redis（可选，用于秒杀场景的库存预热和缓存）
 * 6. 初始化RabbitMQ（可选，用于异步订单处理）
 * 7. 初始化链路追踪（可选，用于分布式追踪）
 * 8. 创建Gin引擎并设置中间件
 * 9. 设置路由
 * 10. 启动后台协程（秒杀订单处理、订单消费者）
 * 11. 启动HTTP服务
 *
 * 支持的特性：
 * - 配置热更新（通过SIGHUP信号触发）
 * - 优雅关闭（通过SIGINT/SIGTERM信号）
 * - 多环境配置（dev/prod）
 * - Swagger API文档
 * - Prometheus监控指标
 * - 链路追踪（OpenTelemetry + Jaeger）
 */

import (
	"context"       // 用于创建上下文，支持超时和取消
	"flag"          // 命令行参数解析
	"fmt"           // 格式化输出
	"os"            // 操作系统功能（退出、信号等）
	"os/signal"    // 信号处理
	"syscall"       // 系统调用信号常量

	"gomall/backend/internal/config"     // 配置管理模块
	"gomall/backend/internal/database"    // 数据库连接模块
	"gomall/backend/internal/logger"      // 日志模块
	"gomall/backend/internal/middleware"  // HTTP中间件
	"gomall/backend/internal/rabbitmq"    // RabbitMQ消息队列
	redispkg "gomall/backend/internal/redis" // Redis缓存（重命名避免冲突）
	"gomall/backend/internal/router"      // 路由配置
	"gomall/backend/internal/service"     // 业务逻辑层
	"gomall/backend/internal/tracing"    // 链路追踪

	"github.com/gin-gonic/gin"   // Gin Web框架
	"go.uber.org/zap"           // Uber Zap日志库
)

/**
 * main 函数 - 程序入口点
 *
 * 命令行参数说明：
 * -config: 配置文件路径，默认 conf/config.yaml
 * -env: 运行环境，可选 dev/prod，默认 dev
 *
 * 配置文件选择逻辑：
 * - dev环境: 使用 conf/config-dev.yaml
 * - prod环境: 使用 conf/config-prod.yaml
 */
func main() {
	// ==================== 第一步：解析命令行参数 ====================
	// flag.String 创建字符串指针变量
	// 第一个参数是参数名，第二个是默认值，第三个是参数说明
	configPath := flag.String("config", "conf/config.yaml", "配置文件路径")
	env := flag.String("env", "dev", "运行环境: dev, prod")
	// flag.Parse 解析命令行参数
	flag.Parse()

	// ==================== 第二步：选择配置文件 ====================
	// 根据运行环境选择对应的配置文件
	if *env == "prod" {
		// 生产环境使用生产配置
		*configPath = "conf/config-prod.yaml"
	} else {
		// 开发环境使用开发配置（默认）
		*configPath = "conf/config-dev.yaml"
	}

	// ==================== 第三步：初始化配置 ====================
	// 注意：此时日志系统还没初始化，所以必须用 fmt.Printf 打印，不能用 logger.Info
	// 这是因为Logger需要先读取配置才能初始化，而配置需要先被读取
	fmt.Printf("正在初始化配置... 路径: %s\n", *configPath)
	if err := config.Init(*configPath); err != nil {
		// 配置加载失败是致命错误，直接 Panic 退出
		// fmt.Errorf 使用 %w 包装原始错误，保留错误链
		panic(fmt.Errorf("配置初始化失败: %w", err))
	}
	fmt.Println("配置初始化成功")

	// ==================== 第四步：初始化日志系统 ====================
	// 此时配置已经加载，可以读取到 log.level 等配置项了
	fmt.Printf("正在初始化日志系统... 环境: %s\n", *env)
	if err := logger.Init(); err != nil {
		// 日志初始化失败直接退出，因为后续所有日志都无法输出
		fmt.Printf("日志初始化失败: %v\n", err)
		os.Exit(1)
	}
	// defer 确保程序退出时刷新日志缓冲区
	// 防止某些日志还在缓冲区时程序崩溃导致丢失
	defer logger.Sync()

	// 从这里开始，日志系统就绪，可以愉快地使用 logger 了
	// zap.String 记录字符串类型的字段
	logger.Info("日志系统初始化成功", zap.String("env", *env))

	// ==================== 第五步：初始化数据库 ====================
	// MySQL + GORM 作为主数据库
	logger.Info("正在连接数据库...")
	// database.Init() 内部会读取 config.GetDatabase() 配置
	if err := database.Init(); err != nil {
		// 数据库连接失败是致命错误，调用 logger.Fatal
		// Fatal 会记录日志后调用 os.Exit(1) 退出程序
		logger.Fatal("数据库连接失败", zap.Error(err))
	}
	// defer 确保数据库连接在程序退出时正确关闭
	defer database.Close()
	logger.Info("数据库连接成功")

	// ==================== 第六步：初始化Redis（可选）====================
	// Redis 用于：
	// 1. 秒杀场景的库存预热和原子扣减
	// 2. 分布式锁
	// 3. 缓存（用户信息、商品信息等）
	// 4. Token黑名单（登出时将Token加入黑名单）
	logger.Info("正在连接Redis...")
	if err := redispkg.Init(); err != nil {
		// Redis连接失败不是致命错误，只是某些功能不可用
		// 使用 Warn 级别日志记录，并降级到数据库方案
		logger.Warn("Redis连接失败，将使用数据库方案", zap.Error(err))
	} else {
		defer redispkg.Close()
		logger.Info("Redis连接成功")
	}

	// ==================== 第七步：初始化RabbitMQ（可选）====================
	// RabbitMQ 用于：
	// 1. 异步订单处理（流量削峰）
	// 2. 秒杀结果通知
	// 3. 延迟队列（订单超时取消）
	logger.Info("正在连接RabbitMQ...")
	if err := rabbitmq.Init(); err != nil {
		logger.Warn("RabbitMQ连接失败，将使用同步方案", zap.Error(err))
	} else {
		defer rabbitmq.Close()
		logger.Info("RabbitMQ连接成功")
	}

	// ==================== 第八步：初始化链路追踪（可选）====================
	// OpenTelemetry + Jaeger 用于分布式追踪
	// 可以追踪一个请求在多个服务之间的调用链路
	tracingConfig := config.GetTracing()
	if tracingConfig.GetBool("enabled") {
		logger.Info("正在初始化链路追踪...")
		// 获取链路追踪配置
		serviceName := tracingConfig.GetString("service_name")
		jaegerEndpoint := tracingConfig.GetString("jaeger_endpoint")
		// InitTracing 返回一个 shutdown 函数，用于程序退出时关闭追踪器
		shutdown, err := tracing.InitTracing(serviceName, jaegerEndpoint)
		if err != nil {
			logger.Warn("链路追踪初始化失败", zap.Error(err))
		} else {
			// 使用 defer 确保追踪器正确关闭
			defer func() {
				if err := shutdown(context.Background()); err != nil {
					logger.Error("链路追踪关闭失败", zap.Error(err))
				}
			}()
			logger.Info("链路追踪初始化成功")
		}
	}

	// ==================== 第九步：创建Gin引擎 ====================
	// 设置Gin运行模式
	// debug: 开发模式，输出详细错误信息
	// release: 生产模式，性能更好
	// test: 测试模式
	ginMode := config.GetApp().GetString("mode")
	gin.SetMode(ginMode)

	// 创建Gin引擎
	// gin.New() 创建一个空的Gin引擎
	// gin.Default() 会默认添加 Logger 和 Recovery 中间件
	r := gin.New()

	// 添加全局中间件
	// gin.Recovery() - 捕获Panic，恢复服务（必须）
	// LoggerMiddleware() - 请求日志
	// MetricsMiddleware() - Prometheus监控指标
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.MetricsMiddleware())

	// ==================== 第十步：设置路由 ====================
	// 所有API路由在 router.Setup() 中定义
	router.Setup(r)

	// ==================== 第十一步：启动后台协程 ====================
	// golang的协程（goroutine）是轻量级线程
	// 使用 go 关键字启动新的协程

	// 启动秒杀订单处理协程
	// 这个协程负责从Redis队列中消费秒杀订单请求
	go func() {
		seckillSvc := service.NewSeckillService()
		// ProcessSeckillOrders 不断从Redis队列中读取秒杀请求
		// 并调用 CreateOrderSync 创建订单
		seckillSvc.ProcessSeckillOrders()
	}()

	// 启动订单消费者协程
	// 这个协程负责从RabbitMQ队列中消费订单消息
	go func() {
		orderSvc := service.NewOrderService()
		// StartOrderConsumer 从RabbitMQ接收订单消息
		// 异步创建订单，实现流量削峰
		orderSvc.StartOrderConsumer()
	}()

	// ==================== 第十二步：优雅关闭与配置热更新 ====================
	// 获取应用配置
	appConfig := config.GetApp()
	host := appConfig.GetString("host")
	port := appConfig.GetInt("port")
	addr := fmt.Sprintf("%s:%d", host, port)

	// 启动一个协程处理系统信号
	go func() {
		// 创建一个信号通道
		sigChan := make(chan os.Signal, 1)
		// 订阅信号
		// syscall.SIGINT - Ctrl+C 中断信号
		// syscall.SIGTERM - 终止信号（kill默认发送）
		// syscall.SIGHUP - 挂起信号（通常用于重新加载配置）
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		for {
			// 阻塞等待信号
			sig := <-sigChan
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				// 收到中断或终止信号，直接退出
				logger.Info("收到关闭信号，正在关闭服务...")
				os.Exit(0)
			case syscall.SIGHUP:
				// 收到挂起信号，重新加载配置
				logger.Info("收到 SIGHUP 信号，正在重新加载配置...")
				if err := config.Reload(); err != nil {
					logger.Error("配置重新加载失败", zap.Error(err))
				} else {
					logger.Info("配置重新加载成功")
					// 重新初始化日志级别（配置可能改了日志级别）
					if err := logger.Reload(); err != nil {
						logger.Error("日志配置重新加载失败", zap.Error(err))
					}
				}
			}
		}
	}()

	// ==================== 第十三步：启动HTTP服务 ====================
	// r.Run(addr) 会在指定地址启动HTTP服务器
	// 它内部使用了 http.ListenAndServe(addr, r)
	// 这是一个阻塞调用，会一直运行直到服务关闭
	logger.Info("服务启动成功", zap.String("addr", addr))
	logger.Info("健康检查", zap.String("url", fmt.Sprintf("http://%s/health", addr)))
	logger.Info("API 文档", zap.String("url", fmt.Sprintf("http://%s/swagger/index.html", addr)))
	if err := r.Run(addr); err != nil {
		// Run 返回的错误通常是地址已被占用等致命错误
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}
