package api

import (
	"strconv"

	"gomall/backend/internal/middleware"
	"gomall/backend/internal/response"
	"gomall/backend/internal/service"
	"gomall/backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户接口处理层
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 新用户注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /api/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeUserAlreadyExist, err.Error())
		return
	}

	// 注册成功后生成 token
	jwtUtil := jwt.NewJWT()
	token, err := jwtUtil.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		response.FailWithMsg(c, response.CodeServerError, "Token生成失败")
		return
	}

	response.OkWithData(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /api/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	token, user, err := h.userService.Login(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeUserPasswordError, err.Error())
		return
	}

	response.OkWithData(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.FailWithMsg(c, response.CodeUserNotFound, err.Error())
		return
	}

	response.OkWithData(c, user)
}

// ProductHandler 商品接口处理层
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler 创建商品处理器
func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

// Create 创建商品
// @Summary 创建商品
// @Description 创建新商品（需要管理员权限）
// @Tags 商品
// @Accept json
// @Produce json
// @Param req body service.CreateProductRequest true "商品信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/product [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeProductCreateFailed, err.Error())
		return
	}

	response.OkWithData(c, product)
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 获取商品列表，支持分页和分类筛选
// @Tags 商品
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "商品分类"
// @Success 200 {object} response.Response
// @Router /api/product [get]
func (h *ProductHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	products, total := h.productService.GetList(page, pageSize, category)

	response.OkWithList(c, products, int64(total), page, pageSize)
}

// Get 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商品
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response
// @Router /api/product/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		response.FailWithMsg(c, response.CodeProductNotFound, err.Error())
		return
	}

	response.OkWithData(c, product)
}

// Update 更新商品
// @Summary 更新商品
// @Description 更新商品信息
// @Tags 商品
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Param req body service.UpdateProductRequest true "商品信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/product/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.productService.Update(uint(id), &req); err != nil {
		response.FailWithMsg(c, response.CodeProductUpdateFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Delete 删除商品
// @Summary 删除商品
// @Description 删除商品（软删除）
// @Tags 商品
// @Produce json
// @Param id path int true "商品ID"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/product/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.productService.Delete(uint(id)); err != nil {
		response.FailWithMsg(c, response.CodeProductDeleteFailed, err.Error())
		return
	}

	response.Ok(c)
}

// OrderHandler 订单接口处理层
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler 创建订单处理器
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
	}
}

// Create 创建订单
// @Summary 创建订单
// @Description 创建新订单
// @Tags 订单
// @Accept json
// @Produce json
// @Param req body service.CreateOrderRequest true "订单信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderCreateFailed, err.Error())
		return
	}

	response.OkWithData(c, order)
}

// List 获取订单列表
// @Summary 获取订单列表
// @Description 获取当前用户的订单列表
// @Tags 订单
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order [get]
func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	orders, total := h.orderService.GetOrderList(userID, page, pageSize)

	response.OkWithList(c, orders, int64(total), page, pageSize)
}

// Get 获取订单详情
// @Summary 获取订单详情
// @Description 根据订单号获取订单详情
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order/{order_no} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderNotFound, err.Error())
		return
	}

	response.OkWithData(c, order)
}

// Pay 支付订单
// @Summary 支付订单
// @Description 模拟支付订单
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order/{order_no}/pay [post]
func (h *OrderHandler) Pay(c *gin.Context) {
	orderNo := c.Param("order_no")

	if err := h.orderService.PayOrder(orderNo); err != nil {
		response.FailWithMsg(c, response.CodeOrderPayFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Cancel 取消订单
// @Summary 取消订单
// @Description 取消未支付的订单
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order/{order_no}/cancel [post]
func (h *OrderHandler) Cancel(c *gin.Context) {
	orderNo := c.Param("order_no")

	if err := h.orderService.CancelOrder(orderNo); err != nil {
		response.FailWithMsg(c, response.CodeOrderCancelFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Checkout 购物车结算
// @Summary 购物车结算
// @Description 将购物车中的商品结算为订单
// @Tags 订单
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/order/checkout [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	orders, err := h.orderService.Checkout(userID)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderCreateFailed, err.Error())
		return
	}

	response.OkWithData(c, orders)
}

// CartHandler 购物车接口处理层
type CartHandler struct {
	cartService *service.CartService
}

// NewCartHandler 创建购物车处理器
func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(),
	}
}

// AddToCart 添加商品到购物车
// @Summary 添加到购物车
// @Description 将商品添加到购物车
// @Tags 购物车
// @Accept json
// @Produce json
// @Param req body service.AddToCartRequest true "商品信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/cart [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	item, err := h.cartService.AddToCart(userID, &req)
	if err != nil {
		response.FailWithMsg(c, response.CodeCartAddFailed, err.Error())
		return
	}

	response.OkWithData(c, item)
}

// List 获取购物车列表
// @Summary 获取购物车列表
// @Description 获取当前用户的购物车列表
// @Tags 购物车
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/cart [get]
func (h *CartHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	cart, err := h.cartService.GetCartList(userID)
	if err != nil {
		response.FailWithMsg(c, response.CodeCartNotFound, err.Error())
		return
	}

	response.OkWithData(c, cart)
}

// Update 更新购物车商品数量
// @Summary 更新购物车
// @Description 更新购物车中商品的数量
// @Tags 购物车
// @Accept json
// @Produce json
// @Param product_id query int true "商品ID"
// @Param req body service.UpdateCartRequest true "数量"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/cart [put]
func (h *CartHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if productID == 0 {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	var req service.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	if err := h.cartService.UpdateCartItem(userID, uint(productID), &req); err != nil {
		response.FailWithMsg(c, response.CodeCartUpdateFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Remove 从购物车删除商品
// @Summary 删除购物车商品
// @Description 从购物车中删除指定商品
// @Tags 购物车
// @Produce json
// @Param product_id query int true "商品ID"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/cart [delete]
func (h *CartHandler) Remove(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 64)
	if productID == 0 {
		response.BadRequest(c, "商品ID不能为空")
		return
	}

	if err := h.cartService.RemoveFromCart(userID, uint(productID)); err != nil {
		response.FailWithMsg(c, response.CodeCartDeleteFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Clear 清空购物车
// @Summary 清空购物车
// @Description 清空当前用户的购物车
// @Tags 购物车
// @Produce json
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/cart/clear [delete]
func (h *CartHandler) Clear(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	if err := h.cartService.ClearCart(userID); err != nil {
		response.FailWithMsg(c, response.CodeCartClearFailed, err.Error())
		return
	}

	response.Ok(c)
}
