package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gomall/internal/config"
	"gomall/internal/logger"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Client RabbitMQ连接对象
var Client *amqp.Connection
var Channel *amqp.Channel

// QueueName 队列名称常量
const (
	OrderQueue    = "order_queue"    // 订单创建队列
	SeckillQueue  = "seckill_queue"  // 秒杀队列
	PayQueue      = "pay_queue"      // 支付队列
	DelayQueue    = "delay_queue"    // 延迟队列（用于订单超时取消）
)

// Init 初始化RabbitMQ连接
func Init() error {
	rabbitConfig := GetRabbitMQConfig()

	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/",
		rabbitConfig.GetString("username"),
		rabbitConfig.GetString("password"),
		rabbitConfig.GetString("host"),
		rabbitConfig.GetInt("port"),
	)

	var err error
	Client, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("RabbitMQ连接失败: %w", err)
	}

	// 创建通道
	Channel, err = Client.Channel()
	if err != nil {
		return fmt.Errorf("RabbitMQ通道创建失败: %w", err)
	}

	// 声明队列
	queues := []string{OrderQueue, SeckillQueue, PayQueue, DelayQueue}
	for _, queue := range queues {
		_, err = Channel.QueueDeclare(
			queue,                      // 队列名称
			true,                       // 持久化
			false,                      // 不自动删除
			false,                      // 不排他
			false,                      // 不阻塞
			amqp.Table{                 // 额外参数
				"x-message-ttl": 86400000, // 消息24小时过期
			},
		)
		if err != nil {
			return fmt.Errorf("队列[%s]声明失败: %w", queue, err)
		}
	}

	// 声明延迟队列（基于死信交换机）
	delayQueueName := "delay_order_queue"
	delayExchange := "delay_exchange"
	delayRoutingKey := "delay_order"

	// 声明死信交换机
	err = Channel.ExchangeDeclare(
		delayExchange, // 交换机名称
		"direct",      // 交换机类型
		true,          // 持久化
		false,         // 自动删除
		false,         // 排他
		false,         // 阻塞
		nil,           // 参数
	)
	if err != nil {
		return fmt.Errorf("死信交换机声明失败: %w", err)
	}

	// 声明延迟队列，绑定到死信交换机
	_, err = Channel.QueueDeclare(
		delayQueueName, // 队列名称
		true,           // 持久化
		false,          // 不自动删除
		false,          // 不排他
		false,          // 不阻塞
		amqp.Table{
			"x-dead-letter-exchange":    "",                     // 死信交换机
			"x-dead-letter-routing-key": DelayQueue,             // 死信路由键
			"x-message-ttl":             1800000,                // 30分钟后变为死信
		},
	)
	if err != nil {
		return fmt.Errorf("延迟队列声明失败: %w", err)
	}

	// 绑定延迟队列到死信交换机
	err = Channel.QueueBind(
		delayQueueName, // 队列名称
		delayRoutingKey, // 路由键
		delayExchange,   // 交换机名称
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("延迟队列绑定失败: %w", err)
	}

	logger.Info("RabbitMQ初始化成功")
	return nil
}

// Close 关闭RabbitMQ连接
func Close() error {
	if Channel != nil {
		Channel.Close()
	}
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Ping 检查RabbitMQ连接
func Ping() error {
	if Client == nil {
		return fmt.Errorf("RabbitMQ未初始化")
	}
	if Client.IsClosed() {
		return fmt.Errorf("RabbitMQ连接已关闭")
	}
	return nil
}

// GetRabbitMQConfig 获取RabbitMQ配置
func GetRabbitMQConfig() *viper.Viper {
	return config.GetRabbitMQ()
}

// OrderMessage 订单消息结构
type OrderMessage struct {
	OrderNo     string    `json:"order_no"`
	UserID      uint      `json:"user_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

// PublishOrderMessage 发布订单创建消息
func PublishOrderMessage(ctx context.Context, msg *OrderMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	err = Channel.PublishWithContext(
		ctx,
		"",         // 默认交换机
		OrderQueue, // 队列名称
		false,      // 强制
		false,      // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 持久化
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("消息发布失败: %w", err)
	}

	return nil
}

// PublishDelayOrderMessage 发布延迟订单消息（用于超时取消）
func PublishDelayOrderMessage(ctx context.Context, orderNo string, delay time.Duration) error {
	body, err := json.Marshal(map[string]string{"order_no": orderNo})
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	err = Channel.PublishWithContext(
		ctx,
		"delay_exchange",   // 延迟交换机
		"delay_order",      // 延迟路由键
		false,              // 强制
		false,              // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Expiration:   fmt.Sprintf("%d", delay.Milliseconds()), // 延迟时间
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("延迟消息发布失败: %w", err)
	}

	return nil
}

// SeckillMessage 秒杀消息结构
type SeckillMessage struct {
	UserID      uint `json:"user_id"`
	ProductID   uint `json:"product_id"`
	RequestID   int64 `json:"request_id"` // 请求ID，用于去重
}

// PublishSeckillMessage 发布秒杀消息
func PublishSeckillMessage(ctx context.Context, msg *SeckillMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %w", err)
	}

	err = Channel.PublishWithContext(
		ctx,
		"",          // 默认交换机
		SeckillQueue, // 秒杀队列
		false,       // 强制
		false,       // 立即
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("秒杀消息发布失败: %w", err)
	}

	return nil
}

// ConsumeOrderMessage 消费订单消息
func ConsumeOrderMessage(handler func(msg *OrderMessage) error) {
	msgs, err := Channel.Consume(
		OrderQueue, // 队列名称
		"",         // 消费者标签
		false,      // 自动确认
		false,      // 排他
		false,      // 不本地
		false,      // 不阻塞
		nil,        // 参数
	)
	if err != nil {
		logger.Error("消费订单消息失败", zap.Error(err))
		return
	}

	for msg := range msgs {
		var orderMsg OrderMessage
		if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
			logger.Error("消息解析失败", zap.Error(err))
			msg.Nack(false, false) // 拒绝消息，不重新入队
			continue
		}

		if err := handler(&orderMsg); err != nil {
			logger.Error("订单处理失败", zap.String("order_no", orderMsg.OrderNo), zap.Error(err))
			msg.Nack(false, true) // 拒绝消息，重新入队
			continue
		}

		msg.Ack(false) // 确认消息
	}
}

// ConsumeSeckillMessage 消费秒杀消息
func ConsumeSeckillMessage(handler func(msg *SeckillMessage) error) {
	msgs, err := Channel.Consume(
		SeckillQueue, // 秒杀队列
		"",           // 消费者标签
		false,        // 自动确认
		false,        // 排他
		false,        // 不本地
		false,        // 不阻塞
		nil,          // 参数
	)
	if err != nil {
		logger.Error("消费秒杀消息失败", zap.Error(err))
		return
	}

	for msg := range msgs {
		var seckillMsg SeckillMessage
		if err := json.Unmarshal(msg.Body, &seckillMsg); err != nil {
			logger.Error("秒杀消息解析失败", zap.Error(err))
			msg.Nack(false, false)
			continue
		}

		if err := handler(&seckillMsg); err != nil {
			logger.Error("秒杀处理失败", zap.Uint("user_id", seckillMsg.UserID), zap.Uint("product_id", seckillMsg.ProductID), zap.Error(err))
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}
}
