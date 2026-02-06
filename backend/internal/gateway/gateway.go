package gateway

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceAddrGetter 服务地址获取接口
type ServiceAddrGetter interface {
	GetServiceAddr(serviceName string) string
}

// APIGateway API 网关
type APIGateway struct {
	serviceAddr ServiceAddrGetter
	port        int
}

// NewAPIGateway 创建 API 网关
func NewAPIGateway(getter ServiceAddrGetter, port int) *APIGateway {
	return &APIGateway{
		serviceAddr: getter,
		port:        port,
	}
}

// Start 启动网关
func (g *APIGateway) Start() error {
	r := gin.Default()

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API 路由 - 动态路由到各个微服务
	apiGroup := r.Group("/api")
	{
		// 用户服务路由
		apiGroup.POST("/user/register", g.proxyTo("user-service", "/register"))
		apiGroup.POST("/user/login", g.proxyTo("user-service", "/login"))
		apiGroup.GET("/user/profile", g.proxyTo("user-service", "/profile"))

		// 商品服务路由
		apiGroup.GET("/product", g.proxyTo("product-service", "/list"))
		apiGroup.GET("/product/:id", g.proxyTo("product-service", "/detail"))
		apiGroup.POST("/product", g.proxyTo("product-service", "/create"))

		// 订单服务路由
		apiGroup.POST("/order", g.proxyTo("order-service", "/create"))
		apiGroup.GET("/order", g.proxyTo("order-service", "/list"))
		apiGroup.POST("/order/:order_no/pay", g.proxyTo("order-service", "/pay"))

		// 库存服务路由
		apiGroup.POST("/seckill/init", g.proxyTo("stock-service", "/init"))
		apiGroup.POST("/seckill", g.proxyTo("stock-service", "/deduct"))
	}

	addr := fmt.Sprintf(":%d", g.port)
	log.Printf("API Gateway starting on %s", addr)
	return r.Run(addr)
}

// proxyTo 创建代理中间件
func (g *APIGateway) proxyTo(serviceName, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		addr := g.serviceAddr.GetServiceAddr(serviceName)
		if addr == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": fmt.Sprintf("服务 %s 不可用", serviceName),
			})
			return
		}

		// 代理请求到目标服务
		targetURL := fmt.Sprintf("http://%s%s", addr, path)

		// TODO: 实现实际的 HTTP 代理（可使用 httputil.ReverseProxy）
		c.JSON(http.StatusOK, gin.H{
			"message": "代理请求",
			"target":  targetURL,
			"service": serviceName,
		})
	}
}
