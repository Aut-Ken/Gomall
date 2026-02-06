package middleware

import (
	"context"
	"log"
	"net/http"

	"gomall/backend/internal/circuitbreaker"

	"github.com/gin-gonic/gin"
)

// CircuitBreakerMiddleware 熔断器中间件
type CircuitBreakerMiddleware struct {
	breaker *circuitbreaker.CircuitBreaker
}

// NewCircuitBreakerMiddleware 创建熔断器中间件
func NewCircuitBreakerMiddleware(name string, config *circuitbreaker.Config) *CircuitBreakerMiddleware {
	if config == nil {
		config = circuitbreaker.DefaultConfig()
	}

	return &CircuitBreakerMiddleware{
		breaker: circuitbreaker.New(name, config),
	}
}

// Protect 保护函数执行
func (m *CircuitBreakerMiddleware) Protect(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.breaker.Execute(ctx, fn)
}

// CircuitBreakerGin 熔断器Gin中间件
func CircuitBreakerGin(name string, config *circuitbreaker.Config) gin.HandlerFunc {
	m := NewCircuitBreakerMiddleware(name, config)

	return func(c *gin.Context) {
		fn := func(ctx context.Context) error {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			err := c.Errors.Last()
			if err != nil {
				return err
			}
			return nil
		}

		err := m.breaker.Execute(c.Request.Context(), fn)
		if err != nil {
			if err == circuitbreaker.ErrCircuitOpen {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"code": 503,
					"msg":  "服务暂时不可用，请稍后重试",
				})
				return
			}
		}
	}
}

// BreakerGroupMiddleware 熔断器组中间件
type BreakerGroupMiddleware struct {
	group *circuitbreaker.BreakerGroup
}

// NewBreakerGroupMiddleware 创建熔断器组中间件
func NewBreakerGroupMiddleware(config *circuitbreaker.Config) *BreakerGroupMiddleware {
	return &BreakerGroupMiddleware{
		group: circuitbreaker.NewBreakerGroup(config),
	}
}

// Get 获取熔断器
func (m *BreakerGroupMiddleware) Get(name string) *circuitbreaker.CircuitBreaker {
	return m.group.Get(name)
}

// Middleware 创建熔断器中间件
func (m *BreakerGroupMiddleware) Middleware(name string, config *circuitbreaker.Config) gin.HandlerFunc {
	breaker := m.group.Get(name)

	return func(c *gin.Context) {
		fn := func(ctx context.Context) error {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			err := c.Errors.Last()
			if err != nil {
				return err
			}
			return nil
		}

		err := breaker.Execute(c.Request.Context(), fn)
		if err != nil {
			if err == circuitbreaker.ErrCircuitOpen {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"code": 503,
					"msg":  "服务暂时不可用，请稍后重试",
				})
				log.Printf("[CircuitBreaker] %s is open", name)
				return
			}
		}
	}
}

// WithBreaker 为处理器添加熔断保护
func WithBreaker(name string, config *circuitbreaker.Config) gin.HandlerFunc {
	m := NewCircuitBreakerMiddleware(name, config)
	return m.MiddlewareFunc()
}

// MiddlewareFunc 实现gin.HandlerFunc
func (m *CircuitBreakerMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		fn := func(ctx context.Context) error {
			c.Request = c.Request.WithContext(ctx)
			c.Next()
			return c.Errors.Last()
		}

		err := m.breaker.Execute(c.Request.Context(), fn)
		if err != nil && err != circuitbreaker.ErrCircuitOpen {
			// 记录非熔断错误
			log.Printf("[CircuitBreaker] error: %v", err)
		}
	}
}
