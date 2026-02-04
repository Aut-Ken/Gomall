package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// 定义所有指标
var (
	// HTTP 请求指标
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gomall_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HTTPRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gomall_http_requests_in_flight",
			Help: "Current number of HTTP requests being processed",
		},
	)

	// 业务指标
	OrdersCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_orders_created_total",
			Help: "Total number of orders created",
		},
	)

	SeckillRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_seckill_requests_total",
			Help: "Total number of seckill requests",
		},
	)

	SeckillSuccessTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_seckill_success_total",
			Help: "Total number of successful seckill orders",
		},
	)

	SeckillFailTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_seckill_fail_total",
			Help: "Total number of failed seckill requests",
		},
		[]string{"reason"},
	)

	CartItemsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_cart_items_total",
			Help: "Total number of items added to cart",
		},
	)

	// 数据库连接池指标
	DBConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gomall_db_connections_active",
			Help: "Active database connections",
		},
	)

	DBConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gomall_db_connections_idle",
			Help: "Idle database connections",
		},
	)

	DBConnectionsWaitTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_db_connections_wait_total",
			Help: "Total number of connections waited for",
		},
	)

	// Redis 连接指标
	RedisPingDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "gomall_redis_ping_duration_seconds",
			Help:    "Redis ping duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	RedisOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_redis_operations_total",
			Help: "Total number of Redis operations",
		},
		[]string{"operation", "status"},
	)

	// RabbitMQ 指标
	RabbitMQMessagesPublished = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_rabbitmq_messages_published_total",
			Help: "Total number of messages published to RabbitMQ",
		},
		[]string{"queue"},
	)

	RabbitMQMessagesConsumed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_rabbitmq_messages_consumed_total",
			Help: "Total number of messages consumed from RabbitMQ",
		},
		[]string{"queue"},
	)

	RabbitMQMessagesFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gomall_rabbitmq_messages_failed_total",
			Help: "Total number of failed message operations",
		},
		[]string{"queue", "operation"},
	)

	// 用户指标
	UserLoginsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_user_logins_total",
			Help: "Total number of user logins",
		},
	)

	UserRegistrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "gomall_user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "gomall_active_users",
			Help: "Number of active users (logged in)",
		},
	)
)

// RecordOrderCreated 记录订单创建
func RecordOrderCreated() {
	OrdersCreatedTotal.Inc()
}

// RecordSeckillRequest 记录秒杀请求
func RecordSeckillRequest() {
	SeckillRequestsTotal.Inc()
}

// RecordSeckillSuccess 记录秒杀成功
func RecordSeckillSuccess() {
	SeckillSuccessTotal.Inc()
}

// RecordSeckillFail 记录秒杀失败
func RecordSeckillFail(reason string) {
	SeckillFailTotal.WithLabelValues(reason).Inc()
}

// RecordCartItemAdded 记录添加购物车
func RecordCartItemAdded() {
	CartItemsTotal.Inc()
}

// RecordUserLogin 记录用户登录
func RecordUserLogin() {
	UserLoginsTotal.Inc()
}

// RecordUserRegister 记录用户注册
func RecordUserRegister() {
	UserRegistrationsTotal.Inc()
}

// RecordRedisOperation 记录 Redis 操作
func RecordRedisOperation(operation string, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	RedisOperationsTotal.WithLabelValues(operation, status).Inc()
}

// RecordRabbitMQMessagePublished 记录消息发布
func RecordRabbitMQMessagePublished(queue string) {
	RabbitMQMessagesPublished.WithLabelValues(queue).Inc()
}

// RecordRabbitMQMessageConsumed 记录消息消费
func RecordRabbitMQMessageConsumed(queue string) {
	RabbitMQMessagesConsumed.WithLabelValues(queue).Inc()
}

// RecordRabbitMQMessageFailed 记录消息失败
func RecordRabbitMQMessageFailed(queue, operation string) {
	RabbitMQMessagesFailed.WithLabelValues(queue, operation).Inc()
}
