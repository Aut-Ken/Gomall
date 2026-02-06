package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gomall/backend/internal/logger"
	"gomall/backend/internal/model"
	"gomall/backend/internal/rabbitmq"
	"gomall/backend/internal/redis"
	"gomall/backend/internal/repository"

	"go.uber.org/zap"
)

var (
	ErrSeckillStart     = errors.New("秒杀活动未开始")
	ErrSeckillEnd       = errors.New("秒杀活动已结束")
	ErrSeckillRepeat    = errors.New("请勿重复秒杀")
	ErrSeckillStockZero = errors.New("商品已售罄")
	ErrSystemBusy       = errors.New("系统繁忙，请稍后重试")
)

// SeckillService 秒杀服务
// 提供高并发场景下的秒杀功能
type SeckillService struct {
	productRepo *repository.ProductRepository
	orderRepo   *repository.OrderRepository
	stockRepo   *repository.StockRepository
}

// NewSeckillService 创建秒杀服务实例
func NewSeckillService() *SeckillService {
	return &SeckillService{
		productRepo: repository.NewProductRepository(),
		orderRepo:   repository.NewOrderRepository(),
		stockRepo:   repository.NewStockRepository(),
	}
}

// SeckillRequest 秒杀请求
type SeckillRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// SeckillResponse 秒杀响应
type SeckillResponse struct {
	OrderNo     string  `json:"order_no"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
}

// SeckillWithRedis 使用Redis + RabbitMQ实现异步秒杀
// 流程：
// 1. Redis预加载库存（减少数据库压力）
// 2. 用户请求先检查库存（内存级别，快速判断）
// 3. 使用Lua脚本原子扣减库存（保证原子性，防止超卖）
// 4. 扣减成功则发送消息到MQ，立即返回“排队中”
func (s *SeckillService) SeckillWithRedis(ctx context.Context, userID uint, req *SeckillRequest) (*SeckillResponse, error) {
	productID := req.ProductID

	// 1. 获取商品信息 (为了检查状态和价格)
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// 2. 检查商品状态
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 3. 检查用户是否重复秒杀（使用Redis set）
	userKey := fmt.Sprintf("seckill:user:%d:%d", userID, productID)
	exists, err := redis.Client.SIsMember(ctx, userKey, userID).Result()
	if err != nil {
		return nil, fmt.Errorf("检查用户秒杀状态失败: %w", err)
	}
	if exists {
		return nil, ErrSeckillRepeat
	}

	// 4. 使用Lua脚本原子扣减Redis库存
	// 注意：这里只是扣减Redis里的缓存库存，数据库库存稍后由消费者扣减
	result, err := decrStockWithLua(ctx, productID, 1)
	if err != nil {
		return nil, fmt.Errorf("库存扣减失败: %w", err)
	}

	// 5. 库存不足
	// Lua脚本返回的是扣减后的剩余库存，如果小于0说明库存不够
	if result < 0 {
		return nil, ErrSeckillStockZero
	}

	// 6. 记录用户秒杀状态 (防止重复秒杀)
	if err := redis.Client.SAdd(ctx, userKey, userID).Err(); err != nil {
		logger.Error("记录用户秒杀状态失败", zap.Uint("user_id", userID), zap.Uint("product_id", productID), zap.Error(err))
		return nil, ErrSystemBusy
	}
	if err := redis.Client.Expire(ctx, userKey, 24*time.Hour).Err(); err != nil {
		logger.Error("设置秒杀状态过期时间失败", zap.Uint("user_id", userID), zap.Uint("product_id", productID), zap.Error(err))
		return nil, ErrSystemBusy
	}

	// 7. 构造秒杀消息
	msg := &rabbitmq.SeckillMessage{
		UserID:    userID,
		ProductID: productID,
		RequestID: time.Now().UnixNano(),
	}

	// 8. 发送消息到 RabbitMQ (异步下单)
	if err := rabbitmq.PublishSeckillMessage(ctx, msg); err != nil {
		// ⚠️ 关键点：如果发消息失败，必须回滚 Redis 库存和用户状态
		logger.Error("发送秒杀消息失败", zap.Uint("user_id", userID), zap.Uint("product_id", productID), zap.Error(err))

		// 回滚库存
		incrStock(ctx, productID, 1)
		// 删除用户秒杀记录
		redis.Client.Del(ctx, userKey)

		return nil, ErrSystemBusy
	}

	// 9. 立即返回结果
	// 注意：此时订单还没真正创建，OrderNo 为空，前端应提示“排队中”或轮询查询
	return &SeckillResponse{
		OrderNo:     "", // 异步处理，暂无订单号
		ProductID:   product.ID,
		ProductName: product.Name,
		Price:       product.Price,
		Quantity:    1,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// ProcessSeckillOrders 处理秒杀订单（MQ消费者）
// 这是一个后台任务，会持续运行
func (s *SeckillService) ProcessSeckillOrders() {
	logger.Info("秒杀订单消费者已启动，等待消息...")

	// 调用 rabbitmq 包里的消费函数
	rabbitmq.ConsumeSeckillMessage(func(msg *rabbitmq.SeckillMessage) error {
		logger.Info("收到秒杀请求",
			zap.Uint("user_id", msg.UserID),
			zap.Uint("product_id", msg.ProductID),
			zap.Int64("request_id", msg.RequestID),
		)

		ctx := context.Background()

		// 1. 去重检查（使用 Redis 记录用户是否已秒杀过该商品）
		dedupKey := fmt.Sprintf("seckill:processed:%d:%d", msg.UserID, msg.ProductID)
		exists, err := redis.Client.Exists(ctx, dedupKey).Result()
		if err == nil && exists > 0 {
			logger.Warn("订单已处理过，跳过",
				zap.Uint("user_id", msg.UserID),
				zap.Uint("product_id", msg.ProductID),
			)
			return nil // 已处理过，跳过
		}

		// 2. 获取商品信息 (为了拿到最新价格和名称)
		product, err := s.productRepo.GetByID(msg.ProductID)
		if err != nil {
			logger.Error("获取商品失败", zap.Uint("product_id", msg.ProductID), zap.Error(err))
			return err // 返回错误，MQ会重试
		}

		// 3. 生成订单号
		orderNo := generateOrderNo()

		// 4. 构造订单对象
		order := &model.Order{
			OrderNo:     orderNo,
			UserID:      msg.UserID,
			ProductID:   msg.ProductID,
			ProductName: product.Name,
			Quantity:    1,
			TotalPrice:  product.Price,
			Status:      1, // 待支付
			PayType:     1,
		}

		// 5. 写入数据库 (真正的落库操作)
		// OrderRepo.Create 里面包含了事务：创建订单 + 扣减数据库库存
		if err := s.orderRepo.Create(order); err != nil {
			logger.Error("创建订单失败", zap.String("order_no", orderNo), zap.Error(err))
			return err // 返回错误，MQ会重试
		}

		// 6. 标记为已处理（防止重复消费）
		redis.Client.Set(ctx, dedupKey, orderNo, 24*time.Hour)

		logger.Info("秒杀订单创建成功", zap.String("order_no", orderNo))

		return nil
	})
}

// decrStockWithLua 使用Lua脚本原子扣减库存
func decrStockWithLua(ctx context.Context, productID uint, quantity int) (int, error) {
	key := fmt.Sprintf("gomall:stock:%d", productID)

	// Lua脚本：原子扣减并检查库存
	script := redis.NewScript(`
		local stock = redis.call('GET', KEYS[1])
		if stock == false then
			return -1
		end
		stock = tonumber(stock)
		local quantity = tonumber(ARGV[1])
		if stock < quantity then
			return -1
		end
		redis.call('DECRBY', KEYS[1], quantity)
		return stock - quantity
	`)

	result, err := script.Run(ctx, redis.Client, []string{key}, quantity).Int()
	if err != nil {
		return -1, err
	}

	return result, nil
}

// incrStock 增加库存 (用于回滚)
func incrStock(ctx context.Context, productID uint, quantity int) error {
	key := fmt.Sprintf("gomall:stock:%d", productID)
	return redis.Client.IncrBy(ctx, key, int64(quantity)).Err()
}

// InitSeckillStock 初始化秒杀库存到Redis
func (s *SeckillService) InitSeckillStock(ctx context.Context, productID uint, stock int) error {
	key := fmt.Sprintf("gomall:stock:%d", productID)
	return redis.Client.Set(ctx, key, stock, 0).Err()
}
