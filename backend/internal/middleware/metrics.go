package middleware

import (
	"strconv"
	"time"

	"gomall/backend/internal/metrics"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware 创建 Prometheus 指标中间件
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 增加正在处理的请求数
		metrics.HTTPRequestsInFlight.Inc()
		defer metrics.HTTPRequestsInFlight.Dec()

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Seconds()

		// 获取状态码
		status := strconv.Itoa(c.Writer.Status())

		// 记录指标
		metrics.HTTPRequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
