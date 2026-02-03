package repository

import (
	"errors"
	"gomall/internal/database"
	"gomall/internal/model"
	"time"

	"gorm.io/gorm"
)

// 定义错误信息
var (
	ErrUserNotFound     = errors.New("用户不存在")
	ErrUserAlreadyExist = errors.New("用户已存在")
	ErrProductNotFound  = errors.New("商品不存在")
	ErrOrderNotFound    = errors.New("订单不存在")
	ErrInsufficientStock = errors.New("库存不足")
)

// UserRepository 用户数据访问层
type UserRepository struct{}

// NewUserRepository 创建用户仓库实例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return database.DB.Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ProductRepository 商品数据访问层
type ProductRepository struct{}

// NewProductRepository 创建商品仓库实例
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// Create 创建商品
func (r *ProductRepository) Create(product *model.Product) error {
	return database.DB.Create(product).Error
}

// GetByID 根据ID获取商品
func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return &product, nil
}

// GetList 获取商品列表
func (r *ProductRepository) GetList(page, pageSize int, category string) ([]model.Product, int64) {
	var products []model.Product
	var total int64

	query := database.DB.Model(&model.Product{}).Where("status = ?", 1)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products)

	return products, total
}

// Update 更新商品
func (r *ProductRepository) Update(product *model.Product) error {
	return database.DB.Save(product).Error
}

// Delete 删除商品（软删除）
func (r *ProductRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Product{}, id).Error
}

// OrderRepository 订单数据访问层
type OrderRepository struct{}

// NewOrderRepository 创建订单仓库实例
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

// Create 创建订单（带事务）
func (r *OrderRepository) Create(order *model.Order) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 创建订单记录
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 扣减库存（悲观锁）
		var product model.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, order.ProductID).Error; err != nil {
			return err
		}

		if product.Stock < order.Quantity {
			return ErrInsufficientStock
		}

		// 扣减库存
		product.Stock -= order.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetByID 根据ID获取订单
func (r *OrderRepository) GetByID(id uint) (*model.Order, error) {
	var order model.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNo 根据订单号获取订单
func (r *OrderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

// GetByUserID 根据用户ID获取订单列表
func (r *OrderRepository) GetByUserID(userID uint, page, pageSize int) ([]model.Order, int64) {
	var orders []model.Order
	var total int64

	query := database.DB.Model(&model.Order{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	return orders, total
}

// Update 更新订单状态
func (r *OrderRepository) Update(order *model.Order) error {
	return database.DB.Save(order).Error
}

// StockRepository 库存数据访问层
type StockRepository struct{}

// NewStockRepository 创建库存仓库实例
func NewStockRepository() *StockRepository {
	return &StockRepository{}
}

// Create 创建库存记录
func (r *StockRepository) Create(stock *model.Stock) error {
	return database.DB.Create(stock).Error
}

// GetByProductID 根据商品ID获取库存
func (r *StockRepository) GetByProductID(productID uint) (*model.Stock, error) {
	var stock model.Stock
	if err := database.DB.Where("product_id = ?", productID).First(&stock).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，创建一条
			stock = model.Stock{
				ProductID:  productID,
				TotalStock: 0,
				LockStock:  0,
				SoldStock:  0,
			}
			database.DB.Create(&stock)
			return &stock, nil
		}
		return nil, err
	}
	return &stock, nil
}

// DeductStock 扣减库存（带事务和乐观锁）
func (r *StockRepository) DeductStock(productID uint, quantity int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.Stock
		if err := tx.Where("product_id = ?", productID).First(&stock).Error; err != nil {
			return err
		}

		// 检查库存是否充足
		availableStock := stock.TotalStock - stock.LockStock - stock.SoldStock
		if availableStock < quantity {
			return ErrInsufficientStock
		}

		// 扣减可用库存，增加已售库存
		stock.SoldStock += quantity
		stock.UpdatedAt = time.Now()

		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		return nil
	})
}
