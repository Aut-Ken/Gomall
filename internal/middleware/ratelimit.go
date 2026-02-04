package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP 级别的限流器
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter 创建 IP 限流器
// rate: 每秒允许的请求数
// burst: 允许的最大突发请求数
func NewIPRateLimiter(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    burst,
	}
}

// getLimiter 获取或创建 IP 限流器
func (i *IPRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// 双重检查
	if limiter, exists = i.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(i.rate, i.burst)
	i.limiters[ip] = limiter
	return limiter
}

// RateLimit 返回 Gin 中间件
func RateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if limiter.getLimiter(ip).Allow() {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
		}
	}
}

// RateLimitByKey 基于自定义 key 的限流
func RateLimitByKey(r rate.Limit, burst int, keyFunc func(c *gin.Context) string) gin.HandlerFunc {
	limiter := NewIPRateLimiter(r, burst)

	return func(c *gin.Context) {
		key := keyFunc(c)
		if key == "" {
			key = c.ClientIP()
		}

		if limiter.getLimiter(key).Allow() {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
		}
	}
}

// GlobalRateLimit 全局限流 - 每秒 1000 请求，突发 2000
func GlobalRateLimit() gin.HandlerFunc {
	return RateLimit(rate.Limit(1000), 2000)
}

// APIRateLimit API 接口限流 - 每秒 100 请求，突发 200
func APIRateLimit() gin.HandlerFunc {
	return RateLimit(rate.Limit(100), 200)
}

// SeckillRateLimit 秒杀接口限流 - 每秒 5 请求，突发 10
func SeckillRateLimit() gin.HandlerFunc {
	return RateLimit(rate.Limit(5), 10)
}

// LoginRateLimit 登录接口限流 - 每秒 10 请求，突发 20
func LoginRateLimit() gin.HandlerFunc {
	return RateLimit(rate.Limit(10), 20)
}

// ==================== Redis 分布式限流器 ====================

// RedisLimiter Redis 分布式限流器
type RedisLimiter struct {
	client    *redis.Client
	rate      int  // 每秒允许的请求数
	burst     int  // 突发上限
	windowSec int  // 滑动窗口秒数
}

// NewRedisLimiter 创建 Redis 分布式限流器
func NewRedisLimiter(client *redis.Client, rate, burst, windowSec int) *RedisLimiter {
	return &RedisLimiter{
		client:    client,
		rate:      rate,
		burst:     burst,
		windowSec: windowSec,
	}
}

// Allow 检查是否允许请求（滑动窗口算法）
func (r *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()
	windowKey := r.getWindowKey(key, now)

	// 使用 Lua 脚本保证原子性
	script := redis.NewScript(`
		local key = KEYS[1]
		local now = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local max = tonumber(ARGV[3])
		local rate = tonumber(ARGV[4])

		-- 删除过期的窗口数据
		redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

		-- 统计当前窗口的请求数
		local count = redis.call('ZCARD', key)

		-- 检查是否超过限制
		if count < max then
			-- 添加当前请求
			redis.call('ZADD', key, now, now .. '-' .. math.random())
			-- 设置过期时间
			redis.call('EXPIRE', key, window + 1)
			return 1
		end

		return 0
	`)

	result, err := script.Run(ctx, r.client, []string{windowKey}, now, r.windowSec, r.burst, r.rate).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

func (r *RedisLimiter) getWindowKey(key string, now int64) string {
	windowStart := now / int64(r.windowSec) * int64(r.windowSec)
	return "ratelimit:" + key + ":" + string(rune(windowStart))
}

// RedisRateLimitMiddleware 创建 Redis 分布式限流中间件
func RedisRateLimitMiddleware(client *redis.Client, rate, burst, windowSec int) gin.HandlerFunc {
	limiter := NewRedisLimiter(client, rate, burst, windowSec)

	return func(c *gin.Context) {
		key := c.ClientIP()
		allowed, err := limiter.Allow(c.Request.Context(), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "限流检查失败",
			})
			c.Abort()
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			return
		}

		c.Next()
	}
}
