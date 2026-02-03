package service

import (
	"context"
	"errors"
	"fmt"
	"gomall/internal/model"
	"gomall/internal/rabbitmq"
	"gomall/internal/redis"
	"gomall/internal/repository"
	"gomall/pkg/jwt"
	"gomall/pkg/password"
	"log"
	"time"

	"gorm.io/gorm"
)

// 定义错误信息
var (
	ErrInvalidPassword = errors.New("密码错误")
	ErrUserDisabled    = errors.New("用户已被禁用")
)

// UserService 用户业务逻辑层
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应结构
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*UserResponse, error) {
	// 检查用户名是否已存在
	existUser, _ := s.userRepo.GetByUsername(req.Username)
	if existUser != nil {
		return nil, repository.ErrUserAlreadyExist
	}

	// 检查邮箱是否已存在
	existEmail, _ := s.userRepo.GetByEmail(req.Email)
	if existEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 密码加密
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("用户创建失败")
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (string, *UserResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", nil, errors.New("用户不存在")
		}
		return "", nil, err
	}

	// 验证密码
	if !password.CheckPassword(req.Password, user.Password) {
		return "", nil, ErrInvalidPassword
	}

	// 生成JWT Token
	jwtUtil := jwt.NewJWT()
	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return "", nil, errors.New("Token生成失败")
	}

	return token, &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

// GetUserByID 根据ID获取用户信息
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

// ProductService 商品业务逻辑层
type ProductService struct {
	productRepo *repository.ProductRepository
}

// NewProductService 创建商品服务实例
func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepository(),
	}
}

// CreateProductRequest 创建商品请求结构
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
}

// UpdateProductRequest 更新商品请求结构
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Status      int     `json:"status"`
}

// ProductResponse 商品响应结构
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

// Create 创建商品
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

// GetList 获取商品列表
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

// GetByID 根据ID获取商品
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

// Update 更新商品
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

// Delete 删除商品
func (s *ProductService) Delete(id uint) error {
	return s.productRepo.Delete(id)
}

// OrderService 订单业务逻辑层
type OrderService struct {
	orderRepo    *repository.OrderRepository
	productRepo  *repository.ProductRepository
	stockRepo    *repository.StockRepository
}

// NewOrderService 创建订单服务实例
func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:   repository.NewOrderRepository(),
		productRepo: repository.NewProductRepository(),
		stockRepo:   repository.NewStockRepository(),
	}
}

// CreateOrderRequest 创建订单请求结构
type CreateOrderRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// OrderResponse 订单响应结构
type OrderResponse struct {
	ID          uint    `json:"id"`
	OrderNo     string  `json:"order_no"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Status      int     `json:"status"`
	PayType     int     `json:"pay_type"`
	CreatedAt   string  `json:"created_at"`
}

// CreateOrder 创建订单（异步模式）
// 通过 RabbitMQ 发送订单消息，由消费者异步创建订单，实现流量削峰
func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	// 获取商品信息
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	// 检查商品状态
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 使用 Redis 库存预检（与秒杀保持一致）
	ctx := context.Background()
	stockKey := fmt.Sprintf("gomall:stock:%d", req.ProductID)
	stock, err := redis.Client.Get(ctx, stockKey).Int()
	if err == nil && stock >= 0 {
		// Redis 库存存在，使用 Redis 库存
		if stock < req.Quantity {
			return nil, repository.ErrInsufficientStock
		}
	} else {
		// Redis 库存不存在，使用数据库库存
		if product.Stock < req.Quantity {
			return nil, repository.ErrInsufficientStock
		}
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 构建订单消息
	orderMsg := &rabbitmq.OrderMessage{
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  product.Price * float64(req.Quantity),
	}

	// 发送订单消息到 RabbitMQ（异步处理）
	if err := rabbitmq.PublishOrderMessage(ctx, orderMsg); err != nil {
		return nil, errors.New("订单提交失败，请稍后重试")
	}

	// 返回订单信息（订单状态为"处理中"）
	return &OrderResponse{
		ID:          0, // 异步创建，暂无数据库ID
		OrderNo:     orderNo,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  orderMsg.TotalPrice,
		Status:      0, // 0: 处理中
		PayType:     1,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateOrderSync 同步创建订单（保留原有逻辑，用于消费者调用）
func (s *OrderService) CreateOrderSync(userID uint, req *CreateOrderRequest) (*OrderResponse, error) {
	// 获取商品信息
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	// 检查商品状态
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 检查库存
	if product.Stock < req.Quantity {
		return nil, repository.ErrInsufficientStock
	}

	// 计算总价
	totalPrice := product.Price * float64(req.Quantity)

	// 生成订单号
	orderNo := generateOrderNo()

	// 创建订单
	order := &model.Order{
		OrderNo:     orderNo,
		UserID:      userID,
		ProductID:   product.ID,
		ProductName: product.Name,
		Quantity:    req.Quantity,
		TotalPrice:  totalPrice,
		Status:      1, // 待支付
		PayType:     1, // 默认支付宝
	}

	if err := s.orderRepo.Create(order); err != nil {
		// 判断是否为库存不足错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrInsufficientStock
		}
		return nil, errors.New("订单创建失败")
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		ProductID:   order.ProductID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		PayType:     order.PayType,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetOrderList 获取用户订单列表
func (s *OrderService) GetOrderList(userID uint, page, pageSize int) ([]OrderResponse, int64) {
	orders, total := s.orderRepo.GetByUserID(userID, page, pageSize)

	responses := make([]OrderResponse, len(orders))
	for i, o := range orders {
		responses[i] = OrderResponse{
			ID:          o.ID,
			OrderNo:     o.OrderNo,
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

// GetOrderByNo 根据订单号获取订单
func (s *OrderService) GetOrderByNo(orderNo string) (*OrderResponse, error) {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{
		ID:          order.ID,
		OrderNo:     order.OrderNo,
		ProductID:   order.ProductID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		PayType:     order.PayType,
		CreatedAt:   order.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// PayOrder 支付订单
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

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(orderNo string) error {
	order, err := s.orderRepo.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	// 只有待支付的订单可以取消
	if order.Status != 1 {
		return errors.New("当前订单状态不允许取消")
	}

	order.Status = 5 // 已取消
	return s.orderRepo.Update(order)
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	// 格式: 时间戳 + 随机数
	return "ORD" + time.Now().Format("20060102150405") + randomString(4)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	time.Sleep(time.Nanosecond) // 确保不同调用产生不同随机数
	return string(result)
}

// StartOrderConsumer 启动订单消费者
// 从 RabbitMQ 消费订单消息，异步创建订单
func (s *OrderService) StartOrderConsumer() {
	log.Println("订单消费者已启动")

	// 使用 rabbitmq 包的消费者
	rabbitmq.ConsumeOrderMessage(func(msg *rabbitmq.OrderMessage) error {
		log.Printf("收到订单消息: %s", msg.OrderNo)

		// 构建创建订单请求
		req := &CreateOrderRequest{
			ProductID: msg.ProductID,
			Quantity:  msg.Quantity,
		}

		// 调用同步创建订单方法
		_, err := s.CreateOrderSync(msg.UserID, req)
		if err != nil {
			log.Printf("订单创建失败: %s, 错误: %v", msg.OrderNo, err)
			return err
		}

		log.Printf("订单创建成功: %s", msg.OrderNo)
		return nil
	})
}

// CartService 购物车业务逻辑层
type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

// NewCartService 创建购物车服务实例
func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repository.NewCartRepository(),
		productRepo: repository.NewProductRepository(),
	}
}

// AddToCartRequest 添加到购物车请求结构
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// UpdateCartRequest 更新购物车请求结构
type UpdateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=1"`
}

// CartItemResponse 购物车项响应结构
type CartItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductImage string  `json:"product_image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"`
}

// CartResponse 购物车响应结构
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCount int                `json:"total_count"`
	TotalPrice float64            `json:"total_price"`
}

// AddToCart 添加商品到购物车
func (s *CartService) AddToCart(userID uint, req *AddToCartRequest) (*CartItemResponse, error) {
	// 检查商品是否存在
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, err
	}

	// 检查商品是否上架
	if product.Status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 检查库存
	if product.Stock < req.Quantity {
		return nil, repository.ErrInsufficientStock
	}

	// 检查购物车中是否已存在该商品
	existingCart, err := s.cartRepo.GetByUserAndProduct(userID, req.ProductID)
	if err == nil && existingCart != nil {
		// 已存在，增加数量
		existingCart.Quantity += req.Quantity
		if err := s.cartRepo.Update(existingCart); err != nil {
			return nil, errors.New("更新购物车失败")
		}

		return &CartItemResponse{
			ID:          existingCart.ID,
			ProductID:   product.ID,
			ProductName: product.Name,
			ProductImage: product.ImageURL,
			Price:       product.Price,
			Quantity:    existingCart.Quantity,
			SubTotal:    product.Price * float64(existingCart.Quantity),
		}, nil
	}

	// 不存在，创建新记录
	cart := &model.Cart{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.cartRepo.Create(cart); err != nil {
		return nil, errors.New("添加购物车失败")
	}

	return &CartItemResponse{
		ID:          cart.ID,
		ProductID:   product.ID,
		ProductName: product.Name,
		ProductImage: product.ImageURL,
		Price:       product.Price,
		Quantity:   cart.Quantity,
		SubTotal:   product.Price * float64(cart.Quantity),
	}, nil
}

// GetCartList 获取购物车列表
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

	for _, cart := range carts {
		product, err := s.productRepo.GetByID(cart.ProductID)
		if err != nil {
			continue // 跳过不存在的商品
		}

		subTotal := product.Price * float64(cart.Quantity)
		item := CartItemResponse{
			ID:          cart.ID,
			ProductID:   product.ID,
			ProductName: product.Name,
			ProductImage: product.ImageURL,
			Price:       product.Price,
			Quantity:   cart.Quantity,
			SubTotal:   subTotal,
		}

		response.Items = append(response.Items, item)
		response.TotalCount += cart.Quantity
		response.TotalPrice += subTotal
	}

	return response, nil
}

// UpdateCartItem 更新购物车商品数量
func (s *CartService) UpdateCartItem(userID, cartID uint, req *UpdateCartRequest) error {
	cart, err := s.cartRepo.GetByUserAndProduct(userID, cartID)
	if err != nil {
		return err
	}

	// 检查商品库存
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

// RemoveFromCart 从购物车删除商品
func (s *CartService) RemoveFromCart(userID, productID uint) error {
	cart, err := s.cartRepo.GetByUserAndProduct(userID, productID)
	if err != nil {
		return err
	}

	return s.cartRepo.Delete(cart.ID)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(userID uint) error {
	return s.cartRepo.DeleteAllByUserID(userID)
}
