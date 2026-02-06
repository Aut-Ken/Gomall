package middleware

import (
	"time"

	"gomall/backend/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 创建结构化日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(start)

		// 获取状态码
		statusCode := c.Writer.Status()

		// 获取错误
		var errMsg string
		if len(c.Errors) > 0 {
			errMsg = c.Errors.String()
		}

		// 记录请求日志
		logger.Info("HTTP Request",
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Int("body-size", c.Writer.Size()),
			zap.String("error", errMsg),
		)
	}
}
