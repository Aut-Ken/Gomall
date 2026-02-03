package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gomall/internal/config"
	"gomall/internal/database"
	"gomall/internal/rabbitmq"
	redispkg "gomall/internal/redis"
	"gomall/internal/router"
	"gomall/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "conf/config.yaml", "配置文件路径")
	flag.Parse()

	// 1. 初始化配置
	log.Println("正在初始化配置...")
	if err := config.Init(*configPath); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}
	log.Println("配置初始化成功")

	// 2. 初始化数据库
	log.Println("正在连接数据库...")
	if err := database.Init(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer database.Close()
	log.Println("数据库连接成功")

	// 3. 初始化Redis（可选，用于秒杀场景）
	log.Println("正在连接Redis...")
	if err := redispkg.Init(); err != nil {
		log.Printf("警告: Redis连接失败，将使用数据库方案: %v", err)
	} else {
		defer redispkg.Close()
		log.Println("Redis连接成功")
	}

	// 4. 初始化RabbitMQ（可选，用于异步订单）
	log.Println("正在连接RabbitMQ...")
	if err := rabbitmq.Init(); err != nil {
		log.Printf("警告: RabbitMQ连接失败，将使用同步方案: %v", err)
	} else {
		defer rabbitmq.Close()
		log.Println("RabbitMQ连接成功")
	}

	// 5. 设置Gin运行模式
	ginMode := config.GetApp().GetString("mode")
	gin.SetMode(ginMode)

	// 6. 创建Gin引擎
	r := gin.Default()

	// 7. 设置路由
	router.Setup(r)

	// 8. 启动秒杀订单处理协程（如果RabbitMQ可用）
	go func() {
		seckillSvc := service.NewSeckillService()
		seckillSvc.ProcessSeckillOrders()
	}()

	// 9. 获取服务配置
	appConfig := config.GetApp()
	host := appConfig.GetString("host")
	port := appConfig.GetInt("port")
	addr := fmt.Sprintf("%s:%d", host, port)

	// 10. 优雅关闭处理
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("正在关闭服务...")
		// 这里可以添加资源清理逻辑
		os.Exit(0)
	}()

	// 11. 启动服务
	log.Printf("服务启动成功，访问地址: http://%s", addr)
	log.Printf("健康检查: http://%s/health", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
