package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gomall/backend/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// Client 全局Redis客户端
var Client *redis.Client

// Init 初始化Redis连接
func Init() error {
	redisConfig := GetRedisConfig()

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.GetString("host"), redisConfig.GetInt("port")),
		Password: redisConfig.GetString("password"),
		DB:       redisConfig.GetInt("db"),
		PoolSize: redisConfig.GetInt("pool_size"),
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}

// Close 关闭Redis连接
func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}

// Ping 检查Redis连接
func Ping() error {
	if Client == nil {
		return fmt.Errorf("Redis未初始化")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return Client.Ping(ctx).Err()
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig() *viper.Viper {
	return config.GetRedis()
}

// CacheKey 生成缓存key
// prefix: 缓存前缀
// id: 资源ID
func CacheKey(prefix string, id interface{}) string {
	return fmt.Sprintf("gomall:%s:%v", prefix, id)
}

// 预定义的缓存key前缀
const (
	UserCachePrefix    = "user"
	ProductCachePrefix = "product"
	StockCachePrefix   = "stock"
	OrderCachePrefix   = "order"
	TokenCachePrefix   = "token"
)

// SetUserCache 设置用户缓存
func SetUserCache(ctx context.Context, userID uint, data interface{}, expiration time.Duration) error {
	key := CacheKey(UserCachePrefix, userID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, jsonData, expiration).Err()
}

// GetUserCache 获取用户缓存
func GetUserCache(ctx context.Context, userID uint) (string, error) {
	key := CacheKey(UserCachePrefix, userID)
	return Client.Get(ctx, key).Result()
}

// DeleteUserCache 删除用户缓存
func DeleteUserCache(ctx context.Context, userID uint) error {
	key := CacheKey(UserCachePrefix, userID)
	return Client.Del(ctx, key).Err()
}

// SetProductCache 设置商品缓存
func SetProductCache(ctx context.Context, productID uint, data interface{}, expiration time.Duration) error {
	key := CacheKey(ProductCachePrefix, productID)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, jsonData, expiration).Err()
}

// GetProductCache 获取商品缓存
func GetProductCache(ctx context.Context, productID uint) (string, error) {
	key := CacheKey(ProductCachePrefix, productID)
	return Client.Get(ctx, key).Result()
}

// DeleteProductCache 删除商品缓存
func DeleteProductCache(ctx context.Context, productID uint) error {
	key := CacheKey(ProductCachePrefix, productID)
	return Client.Del(ctx, key).Err()
}

// SetStockCache 设置库存缓存（用于秒杀）
func SetStockCache(ctx context.Context, productID uint, stock int) error {
	key := CacheKey(StockCachePrefix, productID)
	return Client.Set(ctx, key, stock, 0).Err()
}

// GetStockCache 获取库存缓存
func GetStockCache(ctx context.Context, productID uint) (int, error) {
	key := CacheKey(StockCachePrefix, productID)
	return Client.Get(ctx, key).Int()
}

// DecrStock 扣减库存（原子操作）
func DecrStock(ctx context.Context, productID uint, quantity int) (int64, error) {
	key := CacheKey(StockCachePrefix, productID)
	return Client.DecrBy(ctx, key, int64(quantity)).Result()
}

// IncrStock 增加库存
func IncrStock(ctx context.Context, productID uint, quantity int) error {
	key := CacheKey(StockCachePrefix, productID)
	return Client.IncrBy(ctx, key, int64(quantity)).Err()
}

// LuaScriptStockDecrease Lua脚本：原子扣减库存并返回结果
// 返回值：扣减后的库存数量，-1表示库存不足
const LuaScriptStockDecrease = `
local stock = redis.call('GET', KEYS[1])
if stock == false then
    return -1
end
stock = tonumber(stock)
local quantity = tonumber(ARGV[1])
if stock < quantity then
    return -1
end
stock = stock - quantity
redis.call('SET', KEYS[1], stock)
return stock
`

// DecrStockWithLua 使用Lua脚本原子扣减库存
func DecrStockWithLua(ctx context.Context, productID uint, quantity int) (int, error) {
	key := CacheKey(StockCachePrefix, productID)

	script := redis.NewScript(LuaScriptStockDecrease)
	result, err := script.Run(ctx, Client, []string{key}, quantity).Int()
	if err != nil {
		return -1, err
	}

	return result, nil
}

// SetTokenCache 设置Token黑名单缓存（用于登出）
func SetTokenCache(ctx context.Context, token string, expiration time.Duration) error {
	key := CacheKey(TokenCachePrefix, token)
	return Client.Set(ctx, key, "1", expiration).Err()
}

// IsTokenInvalid 检查Token是否无效（黑名单）
func IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	key := CacheKey(TokenCachePrefix, token)
	exists, err := Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

// NewScript 创建Lua脚本
func NewScript(script string) *redis.Script {
	return redis.NewScript(script)
}
