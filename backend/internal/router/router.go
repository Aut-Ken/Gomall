package router

import (
	"gomall/backend/internal/api"
	"gomall/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 路由设置（集成所有优化）
func Setup(r *gin.Engine) {
	// 初始化处理器
	userHandler := api.NewUserHandler()
	productHandler := api.NewProductHandler()
	orderHandler := api.NewOrderHandler()
	seckillHandler := api.NewSeckillHandler()
	cartHandler := api.NewCartHandler()
	authHandler := api.NewAuthHandler()
	fileHandler := api.NewFileHandler()
	wechatPayHandler := api.NewWeChatPayHandler()
	healthCheck := api.NewHealthCheck()

	// 全局中间件顺序：
	// 1. RequestID（请求追踪）
	// 2. SecurityHeaders（安全头）
	// 3. LogSanitizer（日志脱敏）
	// 4. Metrics（监控指标）
	// 5. GlobalRateLimit（全局限流）

	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.LogSanitizerMiddleware())
	r.Use(middleware.MetricsMiddleware())
	// r.Use(middleware.GlobalRateLimit()) // 临时禁用全局限流

	// 健康检查（无中间件，快速响应）
	r.GET("/health", healthCheck.Health)
	r.GET("/ready", healthCheck.Ready)

	// Prometheus 指标端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.URL("/swagger/doc.json")))

	// API 版本控制
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version":            "v1",
			"api_prefix":         "/api",
			"recommended_prefix": "/api",
		})
	})

	// 全局限流和熔断器保护（用于 API 组）
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.MetricsMiddleware())
	{
		// 用户模块（无需登录）
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", userHandler.Register)
			// 登录接口暂不限流
			userGroup.POST("/login", userHandler.Login)
		}

		// 商品模块（部分需要登录）
		productGroup := apiGroup.Group("/product")
		{
			productGroup.GET("", productHandler.List)    // 获取商品列表（无需登录）
			productGroup.GET("/:id", productHandler.Get) // 获取商品详情（无需登录）

			// 以下接口需要管理员权限
			productGroup.Use(middleware.AdminAuthMiddleware())
			productGroup.POST("", productHandler.Create)       // 创建商品
			productGroup.PUT("/:id", productHandler.Update)    // 更新商品
			productGroup.DELETE("/:id", productHandler.Delete) // 删除商品
		}

		// 订单模块（需要登录）
		orderGroup := apiGroup.Group("/order")
		orderGroup.Use(middleware.AuthMiddleware())
		{
			orderGroup.POST("", orderHandler.Create)                  // 创建订单
			orderGroup.GET("", orderHandler.List)                     // 获取订单列表
			orderGroup.GET("/:order_no", orderHandler.Get)            // 获取订单详情
			orderGroup.POST("/:order_no/pay", orderHandler.Pay)       // 支付订单
			orderGroup.POST("/:order_no/cancel", orderHandler.Cancel) // 取消订单
		}

		// --- 新增：秒杀模块 ---
		seckillGroup := apiGroup.Group("/seckill")
		seckillGroup.Use(middleware.AuthMiddleware(), middleware.SeckillRateLimit())
		{
			seckillGroup.POST("", seckillHandler.Seckill) // 秒杀接口: POST /api/seckill
		}

		// 秒杀管理（需要管理员权限）
		seckillAdminGroup := apiGroup.Group("/seckill")
		seckillAdminGroup.Use(middleware.AdminAuthMiddleware())
		{
			seckillAdminGroup.POST("/init", seckillHandler.InitStock) // 初始化库存: POST /api/seckill/init
		}

		// --- 新增：购物车模块 ---
		cartGroup := apiGroup.Group("/cart")
		cartGroup.Use(middleware.AuthMiddleware())
		{
			cartGroup.POST("", cartHandler.AddToCart)          // 添加到购物车: POST /api/cart
			cartGroup.GET("", cartHandler.List)               // 获取购物车列表: GET /api/cart
			cartGroup.PUT("", cartHandler.Update)             // 更新购物车: PUT /api/cart
			cartGroup.DELETE("", cartHandler.Remove)          // 删除购物车商品: DELETE /api/cart?product_id=xxx
			cartGroup.DELETE("/clear", cartHandler.Clear)     // 清空购物车: DELETE /api/cart/clear
		}

		// 用户中心（需要登录）
		profileGroup := apiGroup.Group("/user")
		profileGroup.Use(middleware.AuthMiddleware())
		{
			profileGroup.GET("/profile", userHandler.GetProfile) // 获取个人信息
		}

		// --- 新增：认证模块 ---
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/refresh-token", authHandler.RefreshToken) // 刷新Token: POST /api/auth/refresh-token
			authGroup.POST("/change-password", middleware.AuthMiddleware(), authHandler.ChangePassword) // 修改密码: POST /api/auth/change-password
			authGroup.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout) // 退出登录: POST /api/auth/logout
		}

		// --- 新增：文件上传模块 ---
		uploadGroup := apiGroup.Group("/upload")
		uploadGroup.Use(middleware.AuthMiddleware())
		{
			uploadGroup.POST("", fileHandler.Upload) // 单文件上传: POST /api/upload
			uploadGroup.POST("/multi", fileHandler.UploadMulti) // 多文件上传: POST /api/upload/multi
		}

		// 配置静态文件服务
		api.SetupStatic(r)

		// --- 新增：微信支付模块 ---
		wechatPayGroup := apiGroup.Group("/pay/wechat")
		wechatPayGroup.Use(middleware.AuthMiddleware())
		{
			wechatPayGroup.POST("/unified-order", wechatPayHandler.UnifiedOrder)   // 统一下单: POST /api/pay/wechat/unified-order
			wechatPayGroup.GET("/query", wechatPayHandler.QueryOrder)              // 订单查询: GET /api/pay/wechat/query
			wechatPayGroup.POST("/close", wechatPayHandler.CloseOrder)             // 关闭订单: POST /api/pay/wechat/close
			wechatPayGroup.POST("/refund", wechatPayHandler.Refund)                // 申请退款: POST /api/pay/wechat/refund
		}

		// 微信支付回调（无需认证）
		apiGroup.POST("/pay/wechat/notify", wechatPayHandler.Notify) // 支付回调: POST /api/pay/wechat/notify
	}
}

// RegisterShopRoutes 商家管理路由（预留）
func RegisterShopRoutes(r *gin.Engine) {
	shopGroup := r.Group("/api/shop")
	shopGroup.Use(middleware.AuthMiddleware())
	{
		// 商家相关接口
	}
}
