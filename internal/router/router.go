package router

import (
	"gomall/internal/api"
	"gomall/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 路由设置
func Setup(r *gin.Engine) {
	// 初始化处理器
	userHandler := api.NewUserHandler()
	productHandler := api.NewProductHandler()
	orderHandler := api.NewOrderHandler()
	seckillHandler := api.NewSeckillHandler()
	cartHandler := api.NewCartHandler()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "OK",
		})
	})

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 用户模块（无需登录）
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", userHandler.Register)
			// 登录接口添加限流
			userGroup.POST("/login", middleware.LoginRateLimit(), userHandler.Login)
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
