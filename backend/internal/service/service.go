package service

/**
 * Service 业务逻辑层 (Business Logic Layer)
 *
 * 本模块实现核心业务逻辑，是连接API层和Repository层的桥梁。
 *
 * 架构分层：
 * - Repository（数据访问层）：数据CRUD操作
 * - Service（业务逻辑层）：业务规则处理
 * - Handler/API（接口层）：HTTP请求处理
 *
 * 设计原则：
 * 1. 每个业务实体对应一个 Service 结构体
 * 2. Service 依赖 Repository 进行数据访问
 * 3. Service 封装业务规则和验证逻辑
 * 4. Service 处理事务边界（复杂操作）
 *
 * 请求-响应模式：
 * - 使用 Request 结构体接收请求参数
 * - 使用 Response 结构体返回处理结果
 * - 错误通过 error 类型返回
 */

import (
	"context"                  // 上下文，用于超时控制和取消
	"errors"                   // 错误处理
	"fmt"                      // 格式化
	"gomall/backend/internal/model"    // 数据模型
	"gomall/backend/internal/rabbitmq" // RabbitMQ消息队列
	"gomall/backend/internal/redis"    // Redis缓存
	"gomall/backend/internal/repository" // 数据访问层
	"gomall/backend/pkg/jwt"           // JWT工具包
	"gomall/backend/pkg/password"      // 密码工具包
	"log"                      // 日志
	"time"                     // 时间处理

	"gorm.io/gorm" // GORM ORM框架
)

/**
 * ==================== 业务错误定义 ====================
 *
 * 使用包级变量定义业务相关的错误。
 */

/**
 * ErrInvalidPassword 密码错误
 * 当用户登录时密码不匹配时返回
 */
var ErrInvalidPassword = errors.New("密码错误")

/**
 * ErrUserDisabled 用户已被禁用
 * 当用户账号被禁用时返回
 */
var ErrUserDisabled = errors.New("用户已被禁用")

/**
 * ==================== UserService 用户服务 ====================
 *
 * 负责处理用户相关的业务逻辑，包括：
 * - 用户注册
 * - 用户登录
 * - Token刷新
 * - 密码修改
 * - 用户信息查询
 *
 * 依赖：
 * - UserRepository：用户数据访问
 * - password包：密码加密验证
 * - jwt包：Token生成解析
 */
type UserService struct {
	// userRepo 用户仓储实例
	// 通过组合方式依赖Repository
	userRepo *repository.UserRepository
}

/**
 * NewUserService 创建用户服务实例
 *
 * 工厂函数，创建UserService并初始化其依赖。
 *
 * 依赖注入说明：
 * - 在构造函数中创建Repository实例
 * - 保持Service与Repository的松耦合
 *
 * 返回值：
 *   *UserService - 用户服务实例
 */
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

/**
 * RegisterRequest 用户注册请求结构
 *
 * 用于接收注册接口的请求参数。
 *
 * 字段标签说明：
 * - json: JSON序列化字段名
 * - binding: 参数验证规则
 *   - required: 必填
 *   - min/max: 长度/数值范围
 *   - email: 邮箱格式
 */
type RegisterRequest struct {
	// Username 用户名，3-50字符
	Username string `json:"username" binding:"required,min=3,max=50"`
	// Password 密码，6-20字符
	Password string `json:"password" binding:"required,min=6,max=20"`
	// Email 邮箱，必填且格式正确
	Email string `json:"email" binding:"required,email"`
	// Phone 手机号，可选
	Phone string `json:"phone"`
}

/**
 * LoginRequest 用户登录请求结构
 */
type LoginRequest struct {
	// Username 用户名
	Username string `json:"username" binding:"required"`
	// Password 密码
	Password string `json:"password" binding:"required"`
}

/**
 * UserResponse 用户信息响应结构
 *
 * 用于返回给客户端的用户信息。
 * 注意：敏感信息（如密码）不应包含在此结构体中。
 */
type UserResponse struct {
	// ID 用户ID
	ID uint `json:"id"`
	// Username 用户名
	Username string `json:"username"`
	// Email 邮箱
	Email string `json:"email"`
	// Phone 手机号
	Phone string `json:"phone"`
}

/**
 * Register 用户注册
 *
 * 注册流程：
 * 1. 检查用户名是否已存在
 * 2. 检查邮箱是否已注册
 * 3. 对密码进行加密
 * 4. 创建用户记录
 *
 * 参数：
 *   req *RegisterRequest - 注册请求
 *
 * 返回值：
 *   *UserResponse - 创建的用户信息
 *   error - 错误信息（用户已存在、邮箱已注册等）
 */
func (s *UserService) Register(req *RegisterRequest) (*UserResponse, error) {
	// 1. 检查用户名是否已存在
	// 使用下划线忽略返回的错误，因为只需要判断是否存在
	existUser, _ := s.userRepo.GetByUsername(req.Username)
	if existUser != nil {
		// 用户名已存在，返回已存在错误
		return nil, repository.ErrUserAlreadyExist
	}

	// 2. 检查邮箱是否已存在
	existEmail, _ := s.userRepo.GetByEmail(req.Email)
	if existEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 3. 密码加密
	// 使用bcrypt算法加密密码
	// bcrypt会自动加盐，防止彩虹表攻击
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 4. 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 调用Repository创建用户记录
	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 5. 返回创建成功的用户信息
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

/**
 * Login 用户登录
 *
 * 登录流程：
 * 1. 根据用户名获取用户
 * 2. 验证密码是否正确
 * 3. 生成JWT Token
 *
 * 参数：
 *   req *LoginRequest - 登录请求
 *
 * 返回值：
 *   string - 访问令牌（AccessToken）
 *   *UserResponse - 用户信息
 *   error - 错误信息（用户不存在、密码错误等）
 */
func (s *UserService) Login(req *LoginRequest) (string, *UserResponse, error) {
	// 1. 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", nil, errors.New("用户不存在")
		}
		return "", nil, err
	}

	// 2. 验证密码
	// CheckPassword会比较明文密码和加密后的密码
	if !password.CheckPassword(req.Password, user.Password) {
		return "", nil, ErrInvalidPassword
	}

	// 3. 生成JWT Token
	// 使用jwt包生成访问令牌和刷新令牌
	jwtUtil := jwt.NewJWT()
	tokenPair, err := jwtUtil.GenerateTokenPair(user.ID, user.Username, user.Email)
	if err != nil {
		return "", nil, errors.New("Token生成失败")
	}

	// 返回访问令牌和用户信息
	return tokenPair.AccessToken, &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

/**
 * RefreshToken 刷新访问令牌
 *
 * 使用刷新令牌获取新的访问令牌。
 *
 * 参数：
 *   refreshToken string - 刷新令牌
 *
 * 返回值：
 *   string - 新的访问令牌
 *   error - 错误信息
 */
func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	jwtUtil := jwt.NewJWT()
	claims, err := jwtUtil.ParseToken(refreshToken)
	if err != nil {
		return "", errors.New("无效的刷新Token")
	}

	// 使用原claims生成新的访问令牌
	return jwtUtil.GenerateToken(claims.UserID, claims.Username, claims.Email)
}

/**
 * ChangePassword 修改密码
 *
 * 修改用户密码流程：
 * 1. 获取用户信息
 * 2. 加密新密码
 * 3. 更新用户密码
 *
 * 参数：
 *   userID uint - 用户ID
 *   newPassword string - 新密码（明文）
 *
 * 返回值：
 *   error - 错误信息
 */
func (s *UserService) ChangePassword(userID uint, newPassword string) error {
	// 1. 获取用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// 2. 加密新密码
	hashedPassword, err := password.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 3. 更新密码
	user.Password = hashedPassword
	return s.userRepo.UpdatePassword(user)
}

/**
 * GetUserByID 根据ID获取用户信息
 *
 * 参数：
 *   id uint - 用户ID
 *
 * 返回值：
 *   *UserResponse - 用户信息
 *   error - 错误信息
 */
func (s *UserService) GetUserByID(id uint) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

/**
 * ==================== ProductService 商品服务 ====================
 *
 * 负责商品相关的业务逻辑，包括：
 * - 商品创建
 * - 商品列表查询
 * - 商品详情查询
 * - 商品更新
 * - 商品删除
 */
type ProductService struct {
	productRepo *repository.ProductRepository
}

/**
 * NewProductService 创建商品服务实例
 */
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
	}
}

/**
 * CreateProductRequest 创建商品请求结构
 */
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

/**
 * UpdateProductRequest 更新商品请求结构
 */
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Status      int     `json:"status"`
}

/**
 * ProductResponse 商品响应结构
 */
type ProductResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

/**
 * Create 创建商品
 */
func (s *ProductService) Create(req *CreateProductRequest) (*ProductResponse, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Status:      1, // 默认上架
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, errors.New("商品创建失败")
	}

	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * GetList 获取商品列表
 */
func (s *ProductService) GetList(page, pageSize int, category string) ([]ProductResponse, int64) {
	products, total := s.productRepo.GetList(page, pageSize, category)

	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Stock:       p.Stock,
			Category:    p.Category,
			ImageURL:    p.ImageURL,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return responses, total
}

/**
 * GetByID 根据ID获取商品
 */
func (s *ProductService) GetByID(id uint) (*ProductResponse, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		Status:      product.Status,
		CreatedAt:   product.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * Update 更新商品
 */
func (s *ProductService) Update(id uint, req *UpdateProductRequest) error {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Status > 0 {
		product.Status = req.Status
	}

	return s.productRepo.Update(product)
}

/**
 * Delete 删除商品
 */
func (s *ProductService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}

/**
 * ==================== OrderService 订单服务 ====================
 *
 * 负责订单相关的业务逻辑，包括：
 * - 创建订单（同步/异步）
 * - 订单列表查询
 * - 订单详情查询
 * - 订单支付
 * - 订单取消
 *
 * 订单创建模式：
 * - 同步模式：直接创建订单
 * - 异步模式：通过RabbitMQ队列异步创建（流量削峰）
 */
type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	stockRepo   *repository.StockRepository
}

/**
 * NewOrderService 创建订单服务实例
 */
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:   repository.NewOrderRepository(),
		productRepo: repository.NewProductRepository(),
		stockRepo:   repository.NewStockRepository(),
	}
}

/**
 * CreateOrderRequest 创建订单请求结构
 */
type CreateOrderRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

/**
 * OrderResponse 订单响应结构
 */
type OrderResponse struct {
	ID          uint    `json:"id"`
	OrderNo     string  `json:"order_no"`
	UserID      uint    `json:"user_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Status      int     `json:"status"`
	PayType     int     `json:"pay_type"`
	CreatedAt   string  `json:"created_at"`
}

/**
 * CreateOrder 创建订单（异步模式）
 *
 * 通过RabbitMQ发送订单消息，由消费者异步创建订单，实现流量削峰。
 *
 * 流程：
 * 1. 获取商品信息
 * 2. 检查商品状态
 * 3. 库存预检
 * 4. 生成订单号
 * 5. 发送订单消息到RabbitMQ
 * 6. 返回处理中状态
 *
 * 参数：
 *   userID uint - 用户ID
 *   req *CreateOrderRequest - 创建订单请求
 *
 * 返回值：
 *   *OrderResponse - 订单响应
 *   error - 错误信息
 */
func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	// 1. 获取商品信息
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	// 2. 检查商品状态
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 3. 使用Redis库存预检
	ctx := context.Background()
	stockKey := fmt.Sprintf("gomall:stock:%d", req.ProductID)
	stock, err := redis.Client.Get(ctx, stockKey).Int()
	if err == nil && stock >= 0 {
		// Redis库存存在，使用Redis库存
		if stock < req.Quantity {
			return nil, repository.ErrInsufficientStock
		}
	} else {
		// Redis库存不存在，使用数据库库存
		if product.Stock < req.Quantity {
			return nil, repository.ErrInsufficientStock
		}
	}

	// 4. 生成订单号
	orderNo := generateOrderNo()

	// 5. 构建订单消息
	orderMsg := &rabbitmq.OrderMessage{
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  product.Price * float64(req.Quantity),
	}

	// 6. 发送订单消息到RabbitMQ（异步处理）
	if err := rabbitmq.PublishOrderMessage(ctx, orderMsg); err != nil {
		return nil, errors.New("订单提交失败，请稍后重试")
	}

	// 7. 返回订单信息（订单状态为"处理中"）
	return &OrderResponse{
		ID:          0,
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  orderMsg.TotalPrice,
		Status:      0, // 0: 处理中
		PayType:     1,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * CreateOrderSync 同步创建订单
 */
func (s *OrderService) CreateOrderSync(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	if product.Stock < req.Quantity {
		return nil, repository.ErrInsufficientStock
	}

	totalPrice := product.Price * float64(req.Quantity)
	orderNo := generateOrderNo()

	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  totalPrice,
		Status:      1, // 待支付
		PayType:     1,
	}

	if err := s.orderRepo.Create(order); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrInsufficientStock
		}
		return nil, errors.New("订单创建失败")
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		UserID:      order.UserID,
		ProductID:   order.ProductID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		PayType:     order.PayType,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * GetOrderList 获取用户订单列表
 */
func (s *OrderService) GetOrderList(userID uint, page, pageSize int) ([]OrderResponse, int64) {
	orders, total := s.orderRepo.GetByUserID(userID, page, pageSize)

	responses := make([]OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = OrderResponse{
			ID:          o.ID,
			OrderNo:     o.OrderNo,
			UserID:      o.UserID,
			ProductID:   o.ProductID,
			ProductName: o.ProductName,
			Quantity:    o.Quantity,
			TotalPrice:  o.TotalPrice,
			Status:      o.Status,
			PayType:     o.PayType,
			CreatedAt:   o.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return responses, total
}

/**
 * GetOrderByNo 根据订单号获取订单
 */
func (s *OrderService) GetOrderByNo(orderNo string) (*OrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		UserID:      order.UserID,
		ProductID:   order.ProductID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		PayType:     order.PayType,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * PayOrder 支付订单
 */
func (s *OrderService) PayOrder(orderNo string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	if order.Status != 1 {
		return errors.New("订单状态不允许支付")
	}

	order.Status = 2 // 已支付
	return s.orderRepo.Update(order)
}

/**
 * CancelOrder 取消订单
 */
func (s *OrderService) CancelOrder(orderNo string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	if order.Status != 1 {
		return errors.New("当前订单状态不允许取消")
	}

	order.Status = 5 // 已取消
	return s.orderRepo.Update(order)
}

/**
 * generateOrderNo 生成订单号
 *
 * 格式：ORD + 时间戳 + 随机数
 * 例如：ORD202401011200001234
 */
func generateOrderNo() string {
	return "ORD" + time.Now().Format("20060102150405") + randomString(4)
}

/**
 * randomString 生成随机字符串
 */
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	time.Sleep(time.Nanosecond)
	return string(result)
}

/**
 * StartOrderConsumer 启动订单消费者
 */
func (s *OrderService) StartOrderConsumer() {
	log.Println("订单消费者已启动")

	rabbitmq.ConsumeOrderMessage(func(msg *rabbitmq.OrderMessage) error {
		log.Printf("收到订单消息: %s", msg.OrderNo)

		req := &CreateOrderRequest{
			ProductID: msg.ProductID,
			Quantity:  msg.Quantity,
		}

		_, err := s.CreateOrderSync(msg.UserID, req)
		if err != nil {
			log.Printf("订单创建失败: %s, 错误: %v", msg.OrderNo, err)
			return err
		}

		log.Printf("订单创建成功: %s", msg.OrderNo)
		return nil
	})
}

/**
 * ==================== CartService 购物车服务 ====================
 *
 * 负责购物车相关的业务逻辑，包括：
 * - 添加商品到购物车
 * - 获取购物车列表
 * - 更新商品数量
 * - 删除商品
 * - 清空购物车
 */
type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

/**
 * NewCartService 创建购物车服务实例
 */
func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repository.NewCartRepository(),
		productRepo: repository.NewProductRepository(),
	}
}

/**
 * AddToCartRequest 添加到购物车请求结构
 */
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

/**
 * UpdateCartRequest 更新购物车请求结构
 */
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}

/**
 * CartItemResponse 购物车项响应结构
 */
type CartItemResponse struct {
	ID           uint    `json:"id"`
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
	SubTotal     float64 `json:"sub_total"`
}

/**
 * CartResponse 购物车响应结构
 */
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCount int                `json:"total_count"`
	TotalPrice float64            `json:"total_price"`
}

/**
 * AddToCart 添加商品到购物车
 */
func (s *CartService) AddToCart(userID uint, req *AddToCartRequest) (*CartItemResponse, error) {
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	if product.Stock < req.Quantity {
		return nil, repository.ErrInsufficientStock
	}

	existingCart, err := s.cartRepo.GetByUserAndProduct(userID, req.ProductID)
	if err == nil && existingCart != nil {
		existingCart.Quantity += req.Quantity
		if err := s.cartRepo.Update(existingCart); err != nil {
			return nil, errors.New("更新购物车失败")
		}

		return &CartItemResponse{
			ID:           existingCart.ID,
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductImage: product.ImageURL,
			Price:        product.Price,
			Quantity:     existingCart.Quantity,
			SubTotal:     product.Price * float64(existingCart.Quantity),
		}, nil
	}

	cart := &model.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartRepo.Create(cart); err != nil {
		return nil, errors.New("添加购物车失败")
	}

	return &CartItemResponse{
		ID:           cart.ID,
		ProductID:    product.ID,
		ProductName:  product.Name,
		ProductImage: product.ImageURL,
		Price:        product.Price,
		Quantity:     cart.Quantity,
		SubTotal:     product.Price * float64(cart.Quantity),
	}, nil
}

/**
 * GetCartList 获取购物车列表（优化版：批量查询）
 */
func (s *CartService) GetCartList(userID uint) (*CartResponse, error) {
	carts, err := s.cartRepo.GetListByUserID(userID)
	if err != nil {
		return nil, errors.New("获取购物车失败")
	}

	response := &CartResponse{
		Items:      make([]CartItemResponse, 0),
		TotalCount: 0,
		TotalPrice: 0,
	}

	if len(carts) == 0 {
		return response, nil
	}

	// 批量获取商品信息，避免N+1查询
	productIDs := make([]uint, len(carts))
	for i, cart := range carts {
		productIDs[i] = cart.ProductID
	}

	products, err := s.productRepo.GetByIDs(productIDs)
	if err != nil {
		return nil, errors.New("获取商品信息失败")
	}

	productMap := make(map[uint]*model.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	for _, cart := range carts {
		product, exists := productMap[cart.ProductID]
		if !exists {
			continue
		}

		subTotal := product.Price * float64(cart.Quantity)
		item := CartItemResponse{
			ID:           cart.ID,
			ProductID:    product.ID,
			ProductName:  product.Name,
			ProductImage: product.ImageURL,
			Price:        product.Price,
			Quantity:     cart.Quantity,
			SubTotal:     subTotal,
		}

		response.Items = append(response.Items, item)
		response.TotalCount += cart.Quantity
		response.TotalPrice += subTotal
	}

	return response, nil
}

/**
 * UpdateCartItem 更新购物车商品数量
 */
func (s *CartService) UpdateCartItem(userID, cartID uint, req *UpdateCartRequest) error {
	cart, err := s.cartRepo.GetByUserAndProduct(userID, cartID)
	if err != nil {
		return err
	}

	product, err := s.productRepo.GetByID(cart.ProductID)
	if err != nil {
		return err
	}

	if product.Stock < req.Quantity {
		return repository.ErrInsufficientStock
	}

	cart.Quantity = req.Quantity
	return s.cartRepo.Update(cart)
}

/**
 * RemoveFromCart 从购物车删除商品
 */
func (s *CartService) RemoveFromCart(userID, productID uint) error {
	cart, err := s.cartRepo.GetByUserAndProduct(userID, productID)
	if err != nil {
		return err
	}

	return s.cartRepo.Delete(cart.ID)
}

/**
 * ClearCart 清空购物车
 */
func (s *CartService) ClearCart(userID uint) error {
	return s.cartRepo.DeleteAllByUserID(userID)
}
