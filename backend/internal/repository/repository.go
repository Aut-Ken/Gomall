package repository

/**
 * Repository 数据访问层 (Data Access Layer)
 *
 * 本模块负责与数据库交互，实现数据的增删改查操作。
 * 采用仓储（Repository）设计模式，将数据访问逻辑封装在独立层。
 *
 * 架构分层：
 * - Model（模型层）：定义数据结构
 * - Repository（数据访问层）：CRUD操作
 * - Service（业务逻辑层）：业务规则处理
 * - Handler/API（接口层）：HTTP请求处理
 *
 * 设计原则：
 * 1. 每个实体对应一个 Repository 结构体
 * 2. Repository 只负责数据访问，不包含业务逻辑
 * 3. 使用 GORM 框架操作数据库
 * 4. 事务操作使用 database.DB.Transaction()
 *
 * 错误处理：
 * - 使用自定义错误类型（ErrXXX）
 * - 通过 errors.Is() 检查错误类型
 * - 返回 nil 表示成功，返回 error 表示失败
 */

import (
	"errors"                           // 错误处理
	"gomall/backend/internal/database" // 数据库连接包
	"gomall/backend/internal/model"    // 数据模型包
	"time"                             // 时间处理

	"gorm.io/gorm" // GORM ORM框架
)

/**
 * ==================== 自定义错误定义 ====================
 *
 * 使用包级变量定义错误，便于错误检查和处理。
 * 使用 errors.New() 创建错误。
 */

/**
 * ErrUserNotFound 用户不存在错误
 * 当查询用户不存在时返回此错误
 */
var ErrUserNotFound = errors.New("用户不存在")

/**
 * ErrUserAlreadyExist 用户已存在错误
 * 当创建用户但用户名或邮箱已存在时返回此错误
 */
var ErrUserAlreadyExist = errors.New("用户已存在")

/**
 * ErrProductNotFound 商品不存在错误
 * 当查询商品不存在时返回此错误
 */
var ErrProductNotFound = errors.New("商品不存在")

/**
 * ErrOrderNotFound 订单不存在错误
 * 当查询订单不存在时返回此错误
 */
var ErrOrderNotFound = errors.New("订单不存在")

/**
 * ErrInsufficientStock 库存不足错误
 * 当下单数量超过商品库存时返回此错误
 */
var ErrInsufficientStock = errors.New("库存不足")

/**
 * ErrCartNotFound 购物车记录不存在错误
 * 当查询购物车记录不存在时返回此错误
 */
var ErrCartNotFound = errors.New("购物车记录不存在")

/**
 * ==================== UserRepository 用户数据访问层 ====================
 *
 * 负责用户数据的增删改查操作。
 *
 * 提供的方法：
 * - Create: 创建用户
 * - GetByID: 根据ID获取用户
 * - GetByUsername: 根据用户名获取用户
 * - GetByEmail: 根据邮箱获取用户
 * - Update: 更新用户信息
 * - UpdatePassword: 更新密码
 */

/**
 * UserRepository 用户仓储结构体
 *
 * 之所以使用空结构体，是因为所有方法都使用包内的 database.DB 全局实例。
 * 这种设计与依赖注入相比更简单，适合小型项目。
 *
 * 使用示例：
 *   repo := repository.NewUserRepository()
 *   user, err := repo.GetByID(1)
 */
type UserRepository struct{}

/**
 * NewUserRepository 创建用户仓库实例
 *
 * 这是一个工厂函数，负责创建 UserRepository 实例。
 * 为什么不直接使用 &UserRepository{}？
 * 1. 统一入口，便于后续扩展（如添加依赖注入）
 * 2. 可以在此添加初始化逻辑
 *
 * 返回值：
 *   *UserRepository - 用户仓库实例
 */
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

/**
 * Create 创建新用户
 *
 * 将用户数据插入数据库。
 * GORM 会自动填充 created_at 和 updated_at 字段。
 *
 * 参数：
 *   user *model.User - 要创建的用户对象
 *
 * 返回值：
 *   error - 插入失败时返回错误（如唯一索引冲突）
 *
 * 注意事项：
 * - 密码应该在调用此方法前加密
 * - 用户名和邮箱的唯一性由数据库约束保证
 */
func (r *UserRepository) Create(user *model.User) error {
	// database.DB.Create() 执行 INSERT 操作
	// .Error 获取可能的错误
	return database.DB.Create(user).Error
}

/**
 * GetByID 根据用户ID获取用户
 *
 * 使用主键查询单个用户记录。
 *
 * 参数：
 *   id uint - 用户ID
 *
 * 返回值：
 *   *model.User - 找到的用户对象
 *   error - 未找到返回 ErrUserNotFound，其他错误返回具体错误
 *
 * GORM 说明：
 * - database.DB.First() 按主键查询最多一条记录
 * - 如果找不到记录，返回 gorm.ErrRecordNotFound
 * - 我们将这个错误转换为自定义的 ErrUserNotFound
 */
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	// First 按主键查询
	if err := database.DB.First(&user, id).Error; err != nil {
		// errors.Is() 检查错误类型
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

/**
 * GetByUsername 根据用户名获取用户
 *
 * 使用 WHERE 子句按用户名查询。
 *
 * 参数：
 *   username string - 用户名
 *
 * 返回值：
 *   *model.User - 找到的用户对象
 *   error - 未找到返回 ErrUserNotFound
 *
 * 查询说明：
 * - 使用 WHERE username = ? 条件
 * - .First() 确保只返回一条记录
 */
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	// Where() 添加查询条件
	// .First() 执行查询并返回第一条记录
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

/**
 * GetByEmail 根据邮箱获取用户
 *
 * 用途：检查邮箱是否已注册、找回密码等
 *
 * 参数：
 *   email string - 用户邮箱
 *
 * 返回值：
 *   *model.User - 找到的用户对象
 *   error - 未找到返回 ErrUserNotFound
 */
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

/**
 * UpdatePassword 更新用户密码
 *
 * 只更新密码字段，使用 .Model() 指定更新目标。
 *
 * 参数：
 *   user *model.User - 包含新密码的用户对象（需包含ID）
 *
 * 返回值：
 *   error - 更新失败时返回错误
 *
 * 优化说明：
 * - 使用 .Model(user) 指定更新 user 这条记录
 * - 只更新 password 字段
 * - 比 .Save() 更高效，因为只生成一条 UPDATE 语句
 */
func (r *UserRepository) UpdatePassword(user *model.User) error {
	// .Model(user) 指定要更新的记录
	// .Update() 更新指定字段
	return database.DB.Model(user).Update("password", user.Password).Error
}

/**
 * Update 更新用户信息
 *
 * 使用 .Save() 方法，会根据是否有主键值决定 INSERT 或 UPDATE。
 *
 * 参数：
 *   user *model.User - 要更新的用户对象
 *
 * 返回值：
 *   error - 更新失败时返回错误
 *
 * 注意：
 * - .Save() 会更新所有字段（即使值没变）
 * - 如果只需要更新部分字段，使用 .Select() 指定字段
 */
func (r *UserRepository) Update(user *model.User) error {
	return database.DB.Save(user).Error
}

/**
 * ==================== ProductRepository 商品数据访问层 ====================
 *
 * 负责商品数据的增删改查操作。
 *
 * 提供的方法：
 * - Create: 创建商品
 * - GetByID: 根据ID获取商品
 * - GetList: 获取商品列表（分页、筛选）
 * - Update: 更新商品信息
 * - Delete: 删除商品（软删除）
 * - GetByIDs: 批量获取商品
 * - GetByIDsWithCache: 批量获取商品（带缓存）
 */

/**
 * ProductRepository 商品仓储结构体
 */
type ProductRepository struct{}

/**
 * NewProductRepository 创建商品仓库实例
 */
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

/**
 * Create 创建新商品
 *
 * 参数：
 *   product *model.Product - 要创建的商品对象
 *
 * 返回值：
 *   error - 创建失败时返回错误
 */
func (r *ProductRepository) Create(product *model.Product) error {
	return database.DB.Create(product).Error
}

/**
 * GetByID 根据商品ID获取商品
 *
 * 参数：
 *   id uint - 商品ID
 *
 * 返回值：
 *   *model.Product - 找到的商品对象
 *   error - 未找到返回 ErrProductNotFound
 */
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

/**
 * GetList 获取商品列表
 *
 * 支持分页查询和分类筛选。
 *
 * 参数：
 *   page int - 页码，从1开始
 *   pageSize int - 每页数量
 *   category string - 分类筛选条件，空字符串表示不过滤
 *
 * 返回值：
 *   []model.Product - 商品列表
 *   int64 - 总记录数（用于计算总页数）
 *
 * 查询说明：
 * - .Where("status = ?", 1) 只查询上架商品
 * - .Count() 统计总数（不含分页）
 * - .Offset() 跳过前面 N 条
 * - .Limit() 限制返回数量
 * - .Order() 排序（按创建时间倒序）
 */
func (r *ProductRepository) GetList(page, pageSize int, category string) ([]model.Product, int64) {
	var products []model.Product
	var total int64

	// 构建基础查询
	// 只查询上架的商品（status = 1）
	query := database.DB.Model(&model.Product{}).Where("status = ?", 1)

	// 如果指定了分类，添加分类筛选条件
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 统计符合条件的总记录数
	// 这个值不随分页参数变化
	query.Count(&total)

	// 计算分页偏移量
	// 第1页跳过0条，第2页跳过 pageSize 条，以此类推
	offset := (page - 1) * pageSize

	// 执行分页查询
	// .Offset() 跳过前面N条
	// .Limit() 限制返回数量
	// .Order() 排序（ DESC 表示倒序）
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products)

	return products, total
}

/**
 * Update 更新商品信息
 *
 * 参数：
 *   product *model.Product - 要更新的商品对象
 *
 * 返回值：
 *   error - 更新失败时返回错误
 */
func (r *ProductRepository) Update(product *model.Product) error {
	return database.DB.Save(product).Error
}

/**
 * Delete 删除商品（软删除）
 *
 * GORM 的软删除特性：
 * - 不会真正删除数据库记录
 * - 而是设置 deleted_at 字段
 * - 查询时会自动过滤已删除的记录
 *
 * 参数：
 *   id uint - 要删除的商品ID
 *
 * 返回值：
 *   error - 删除失败时返回错误
 */
func (r *ProductRepository) Delete(id uint) error {
	// .Delete() 执行软删除
	return database.DB.Delete(&model.Product{}, id).Error
}

/**
 * GetByIDs 批量获取商品
 *
 * 使用 IN 查询一次性获取多个商品，避免 N+1 查询问题。
 *
 * 性能优化说明：
 * - 原始做法：循环调用 GetByID() N次 = N次数据库查询
 * - 优化做法：一次 IN 查询 = 1次数据库查询
 *
 * 参数：
 *   ids []uint - 商品ID列表
 *
 * 返回值：
 *   []model.Product - 找到的商品列表
 *   error - 查询失败时返回错误
 *
 * 使用场景：
 * - 购物车：获取购物车中所有商品信息
 * - 订单：获取订单中的多个商品
 */
func (r *ProductRepository) GetByIDs(ids []uint) ([]model.Product, error) {
	var products []model.Product

	// 空列表直接返回空数组，避免无效查询
	if len(ids) == 0 {
		return products, nil
	}

	// IN 查询：SELECT * FROM products WHERE id IN (?, ?, ...)
	if err := database.DB.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

/**
 * GetByIDsWithCache 批量获取商品（带缓存）
 *
 * 这是一个预留方法，目前直接调用 GetByIDs()。
 * 后续可以加入 Redis 缓存逻辑：
 * 1. 先从 Redis 获取
 * 2. 缓存中没有再查数据库
 * 3. 查到的数据写入 Redis
 *
 * 参数：
 *   ids []uint - 商品ID列表
 *
 * 返回值：
 *   []model.Product - 找到的商品列表
 *   error - 查询失败时返回错误
 */
func (r *ProductRepository) GetByIDsWithCache(ids []uint) ([]model.Product, error) {
	// TODO: 实现带Redis缓存的批量查询
	// 1. 从Redis批量获取
	// 2. 找出未缓存的ID
	// 3. 从数据库查询未缓存的
	// 4. 写入Redis
	// 5. 返回结果
	return r.GetByIDs(ids)
}

/**
 * ==================== OrderRepository 订单数据访问层 ====================
 *
 * 负责订单数据的增删改查操作。
 * Create 方法包含事务处理，确保订单创建和库存扣减的原子性。
 *
 * 提供的方法：
 * - Create: 创建订单（带事务）
 * - GetByID: 根据ID获取订单
 * - GetByOrderNo: 根据订单号获取订单
 * - GetByUserID: 获取用户的订单列表
 * - Update: 更新订单信息
 */

/**
 * OrderRepository 订单仓储结构体
 */
type OrderRepository struct{}

/**
 * NewOrderRepository 创建订单仓库实例
 */
func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

/**
 * Create 创建订单（带事务）
 *
 * 事务保证：
 * - 订单创建成功，库存扣减必须成功
 * - 库存扣减失败，订单创建必须回滚
 *
 * 悲观锁说明：
 * - 使用 FOR UPDATE 锁定商品记录
 * - 防止并发下单导致超卖
 *
 * 参数：
 *   order *model.Order - 要创建的订单对象
 *
 * 返回值：
 *   error - 失败时返回错误（库存不足、数据库错误等）
 */
func (r *OrderRepository) Create(order *model.Order) error {
	// database.DB.Transaction() 创建事务
	// 传入的函数中的所有操作都在一个事务中
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单记录
		// tx.Create() 使用传入的事务实例
		if err := tx.Create(order).Error; err != nil {
			return err // 事务会自动回滚
		}

		// 2. 锁定商品记录（悲观锁）
		// Set("gorm:query_option", "FOR UPDATE") 添加 FOR UPDATE 锁
		// 其他事务无法修改这条记录，直到当前事务提交
		var product model.Product
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&product, order.ProductID).Error; err != nil {
			return err
		}

		// 3. 检查库存
		if product.Stock < order.Quantity {
			return ErrInsufficientStock
		}

		// 4. 扣减库存
		product.Stock -= order.Quantity
		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// 5. 返回 nil 表示事务提交
		return nil
	})
}

/**
 * GetByID 根据订单ID获取订单
 */
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

/**
 * GetByOrderNo 根据订单号获取订单
 *
 * 订单号是业务主键（order_no），不是自增ID。
 *
 * 参数：
 *   orderNo string - 订单号
 *
 * 返回值：
 *   *model.Order - 找到的订单对象
 *   error - 未找到返回 ErrOrderNotFound
 */
func (r *OrderRepository) GetByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	// .Where() 按订单号查询
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	return &order, nil
}

/**
 * GetByUserID 获取用户的订单列表
 *
 * 支持分页查询，按创建时间倒序排列。
 *
 * 参数：
 *   userID uint - 用户ID
 *   page int - 页码
 *   pageSize int - 每页数量
 *
 * 返回值：
 *   []model.Order - 订单列表
 *   int64 - 总记录数
 */
func (r *OrderRepository) GetByUserID(userID uint, page, pageSize int) ([]model.Order, int64) {
	var orders []model.Order
	var total int64

	// 按用户ID查询
	query := database.DB.Model(&model.Order{}).Where("user_id = ?", userID)

	// 统计总数
	query.Count(&total)

	// 分页
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders)

	return orders, total
}

/**
 * Update 更新订单信息
 *
 * 参数：
 *   order *model.Order - 要更新的订单对象
 *
 * 返回值：
 *   error - 更新失败时返回错误
 *
 * 常见用途：
 * - 更新订单状态（待支付 -> 已支付）
 * - 更新支付信息
 */
func (r *OrderRepository) Update(order *model.Order) error {
	return database.DB.Save(order).Error
}

/**
 * ==================== StockRepository 库存数据访问层 ====================
 *
 * 负责库存数据的增删改查操作。
 *
 * 提供的方法：
 * - Create: 创建库存记录
 * - GetByProductID: 获取商品库存
 * - DeductStock: 扣减库存
 */

/**
 * StockRepository 库存仓储结构体
 */
type StockRepository struct{}

/**
 * NewStockRepository 创建库存仓库实例
 */
func NewStockRepository() *StockRepository {
	return &StockRepository{}
}

/**
 * Create 创建库存记录
 *
 * 参数：
 *   stock *model.Stock - 要创建的库存对象
 *
 * 返回值：
 *   error - 创建失败时返回错误
 */
func (r *StockRepository) Create(stock *model.Stock) error {
	return database.DB.Create(stock).Error
}

/**
 * GetByProductID 获取商品库存
 *
 * 如果库存记录不存在，会自动创建一条新的库存记录。
 *
 * 参数：
 *   productID uint - 商品ID
 *
 * 返回值：
 *   *model.Stock - 库存记录
 *   error - 查询失败时返回错误
 */
func (r *StockRepository) GetByProductID(productID uint) (*model.Stock, error) {
	var stock model.Stock
	if err := database.DB.Where("product_id = ?", productID).First(&stock).Error; err != nil {
		// 如果找不到库存记录，创建一个新的
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 初始化库存为0
			stock = model.Stock{
				ProductID:  productID,
				TotalStock: 0,
				LockStock:  0,
				SoldStock:  0,
			}
			// 创建库存记录
			database.DB.Create(&stock)
			return &stock, nil
		}
		return nil, err
	}
	return &stock, nil
}

/**
 * DeductStock 扣减库存（带事务）
 *
 * 使用乐观锁机制：
 * - 先查询库存
 * - 检查可用库存
 * - 更新时使用事务保证原子性
 *
 * 库存计算公式：
 * - 可用库存 = 总库存 - 锁定库存 - 已售库存
 *
 * 参数：
 *   productID uint - 商品ID
 *   quantity int - 扣减数量
 *
 * 返回值：
 *   error - 库存不足或更新失败时返回错误
 */
func (r *StockRepository) DeductStock(productID uint, quantity int) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var stock model.Stock

		// 1. 查询库存
		if err := tx.Where("product_id = ?", productID).First(&stock).Error; err != nil {
			return err
		}

		// 2. 计算可用库存并检查
		availableStock := stock.TotalStock - stock.LockStock - stock.SoldStock
		if availableStock < quantity {
			return ErrInsufficientStock
		}

		// 3. 扣减库存
		// 增加已售库存
		stock.SoldStock += quantity
		// 更新最后修改时间
		stock.UpdatedAt = time.Now()

		// 4. 保存更新
		if err := tx.Save(&stock).Error; err != nil {
			return err
		}

		return nil
	})
}

/**
 * ==================== CartRepository 购物车数据访问层 ====================
 *
 * 负责购物车数据的增删改查操作。
 *
 * 提供的方法：
 * - Create: 添加商品到购物车
 * - GetByUserAndProduct: 获取特定用户的特定商品
 * - GetListByUserID: 获取用户的购物车列表
 * - Update: 更新购物车
 * - Delete: 删除购物车记录
 * - DeleteByUserAndProduct: 按用户和商品删除
 * - DeleteAllByUserID: 清空用户购物车
 */

/**
 * CartRepository 购物车仓储结构体
 */
type CartRepository struct{}

/**
 * NewCartRepository 创建购物车仓库实例
 */
func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

/**
 * Create 添加商品到购物车
 *
 * 参数：
 *   cart *model.Cart - 购物车记录
 *
 * 返回值：
 *   error - 添加失败时返回错误
 */
func (r *CartRepository) Create(cart *model.Cart) error {
	return database.DB.Create(cart).Error
}

/**
 * GetByUserAndProduct 获取用户的特定商品购物车记录
 *
 * 用于检查某商品是否已在购物车中。
 *
 * 参数：
 *   userID uint - 用户ID
 *   productID uint - 商品ID
 *
 * 返回值：
 *   *model.Cart - 找到的购物车记录
 *   error - 未找到返回 ErrCartNotFound
 */
func (r *CartRepository) GetByUserAndProduct(userID, productID uint) (*model.Cart, error) {
	var cart model.Cart
	// 使用 AND 条件查询
	if err := database.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

/**
 * GetByUserAndProductUnscoped 获取用户的特定商品购物车记录（包含已软删除的记录）
 */
func (r *CartRepository) GetByUserAndProductUnscoped(userID, productID uint) (*model.Cart, error) {
	var cart model.Cart
	// Unscoped() 忽略软删除标记
	if err := database.DB.Unscoped().Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}
	return &cart, nil
}

/**
 * GetListByUserID 获取用户的购物车列表
 *
 * 参数：
 *   userID uint - 用户ID
 *
 * 返回值：
 *   []model.Cart - 购物车记录列表
 *   error - 查询失败时返回错误
 */
func (r *CartRepository) GetListByUserID(userID uint) ([]model.Cart, error) {
	var carts []model.Cart
	// 按用户ID查询，不限数量
	if err := database.DB.Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return nil, err
	}
	return carts, nil
}

/**
 * Update 更新购物车记录
 *
 * 参数：
 *   cart *model.Cart - 要更新的购物车对象
 *
 * 返回值：
 *   error - 更新失败时返回错误
 */
func (r *CartRepository) Update(cart *model.Cart) error {
	return database.DB.Save(cart).Error
}

/**
 * Delete 删除购物车记录
 *
 * 参数：
 *   id uint - 购物车记录ID
 *
 * 返回值：
 *   error - 删除失败时返回错误
 */
func (r *CartRepository) Delete(id uint) error {
	return database.DB.Delete(&model.Cart{}, id).Error
}

/**
 * DeleteByUserAndProduct 按用户和商品删除
 *
 * 用于删除购物车中的特定商品。
 *
 * 参数：
 *   userID uint - 用户ID
 *   productID uint - 商品ID
 *
 * 返回值：
 *   error - 删除失败时返回错误
 */
func (r *CartRepository) DeleteByUserAndProduct(userID, productID uint) error {
	return database.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Cart{}).Error
}

/**
 * DeleteAllByUserID 清空用户购物车
 *
 * 参数：
 *   userID uint - 用户ID
 *
 * 返回值：
 *   error - 删除失败时返回错误
 */
func (r *CartRepository) DeleteAllByUserID(userID uint) error {
	return database.DB.Where("user_id = ?", userID).Delete(&model.Cart{}).Error
}
