package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State 熔断器状态
type State int

const (
	StateClosed State = iota
	StateHalfOpen
	StateOpen
)

// String 返回状态字符串
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// ErrCircuitOpen 熔断器开启错误
var ErrCircuitOpen = errors.New("circuit breaker is open")

// Config 熔断器配置
type Config struct {
	FailureThreshold  int           // 失败阈值（连续失败次数）
	SuccessThreshold  int           // 成功阈值（半开状态下需要成功的次数）
	Timeout           time.Duration // 超时时间（open -> half-open）
	Interval          time.Duration // 统计周期
	RequestTimeout    time.Duration // 单次请求超时
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		FailureThreshold:  5,           // 连续5次失败后开启熔断
		SuccessThreshold:  3,           // 半开状态下需要3次成功才能关闭
		Timeout:            60 * time.Second, // 60秒后尝试半开
		Interval:           60 * time.Second, // 统计周期
		RequestTimeout:     10 * time.Second, // 单次请求超时
	}
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	name             string
	config           *Config
	state            State
	failures         int
	successes        int
	lastStateChange  time.Time
	lastIntervalReset time.Time
	mu               sync.RWMutex
}

// New 创建熔断器
func New(name string, config *Config) *CircuitBreaker {
	if config == nil {
		config = DefaultConfig()
	}
	return &CircuitBreaker{
		name:             name,
		config:           config,
		state:            StateClosed,
		lastStateChange:  time.Now(),
		lastIntervalReset: time.Now(),
	}
}

// Name 返回熔断器名称
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// State 返回当前状态
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	state := cb.state

	// 检查是否应该从open变为half-open
	if state == StateOpen && time.Since(cb.lastStateChange) > cb.config.Timeout {
		state = StateHalfOpen
	}

	// 检查周期重置
	if time.Since(cb.lastIntervalReset) > cb.config.Interval {
		cb.failures = 0
		cb.successes = 0
		cb.lastIntervalReset = time.Now()
	}

	return state
}

// Execute 执行受保护的函数
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	// 创建带超时的上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, cb.config.RequestTimeout)
	defer cancel()

	// 执行函数
	err := fn(timeoutCtx)

	// 记录结果
	cb.recordResult(err)

	return err
}

// allowRequest 检查是否允许请求
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	state := cb.state

	switch state {
	case StateOpen:
		// 检查是否超时
		if time.Since(cb.lastStateChange) > cb.config.Timeout {
			cb.state = StateHalfOpen
			cb.successes = 0
			cb.failures = 0
			cb.lastStateChange = time.Now()
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return true
	}
}

// recordResult 记录执行结果
func (cb *CircuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failures++
		cb.successes = 0

		// 检查是否应该开启熔断
		if cb.failures >= cb.config.FailureThreshold {
			cb.state = StateOpen
			cb.lastStateChange = time.Now()
		}
	} else {
		cb.successes++

		// 检查是否应该关闭熔断
		if cb.state == StateHalfOpen && cb.successes >= cb.config.SuccessThreshold {
			cb.state = StateClosed
			cb.failures = 0
			cb.successes = 0
			cb.lastStateChange = time.Now()
		}
	}
}

// Reset 重置熔断器
func (cb *CircuitBreaker) Reset() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateClosed
	cb.failures = 0
	cb.successes = 0
	cb.lastStateChange = time.Now()
	cb.lastIntervalReset = time.Now()
}

// ForceOpen 强制开启熔断
func (cb *CircuitBreaker) ForceOpen() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.state = StateOpen
	cb.lastStateChange = time.Now()
}

// Metrics 获取熔断器指标
func (cb *CircuitBreaker) Metrics() map[string]interface{} {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return map[string]interface{}{
		"name":              cb.name,
		"state":             cb.state.String(),
		"failures":         cb.failures,
		"successes":        cb.successes,
		"failure_threshold": cb.config.FailureThreshold,
		"success_threshold": cb.config.SuccessThreshold,
		"last_state_change": cb.lastStateChange.Format(time.RFC3339),
	}
}

// ==================== 熔断器组 ====================

// BreakerGroup 熔断器组
type BreakerGroup struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
	config   *Config
}

// NewBreakerGroup 创建熔断器组
func NewBreakerGroup(config *Config) *BreakerGroup {
	if config == nil {
		config = DefaultConfig()
	}
	return &BreakerGroup{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
	}
}

// Get 获取或创建熔断器
func (g *BreakerGroup) Get(name string) *CircuitBreaker {
	g.mu.RLock()
	if cb, ok := g.breakers[name]; ok {
		g.mu.RUnlock()
		return cb
	}
	g.mu.RUnlock()

	g.mu.Lock()
	defer g.mu.Unlock()

	// 双重检查
	if cb, ok := g.breakers[name]; ok {
		return cb
	}

	cb := New(name, g.config)
	g.breakers[name] = cb
	return cb
}

// Remove 移除熔断器
func (g *BreakerGroup) Remove(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.breakers, name)
}

// ResetAll 重置所有熔断器
func (g *BreakerGroup) ResetAll() {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, cb := range g.breakers {
		cb.Reset()
	}
}
