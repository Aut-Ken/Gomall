package api

import (
	"net/http"
	"strconv"

	"gomall/internal/middleware"
	"gomall/internal/service"

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
// @Success 200 {object} map[string]interface{}
// @Router /api/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": user,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body service.LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	token, user, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// GetProfile 获取当前用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的信息
// @Tags 用户
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": user,
	})
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
// @Success 200 {object} map[string]interface{}
// @Router /api/product [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	product, err := h.productService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": product,
	})
}

// List 获取商品列表
// @Summary 获取商品列表
// @Description 获取商品列表，支持分页和分类筛选
// @Tags 商品
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "商品分类"
// @Success 200 {object} map[string]interface{}
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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":     products,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// Get 获取商品详情
// @Summary 获取商品详情
// @Description 根据ID获取商品详情
// @Tags 商品
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/product/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": product,
	})
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
// @Success 200 {object} map[string]interface{}
// @Router /api/product/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if err := h.productService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// Delete 删除商品
// @Summary 删除商品
// @Description 删除商品（软删除）
// @Tags 商品
// @Produce json
// @Param id path int true "商品ID"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/product/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if err := h.productService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
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
// @Success 200 {object} map[string]interface{}
// @Router /api/order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": order,
	})
}

// List 获取订单列表
// @Summary 获取订单列表
// @Description 获取当前用户的订单列表
// @Tags 订单
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/order [get]
func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":     orders,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// Get 获取订单详情
// @Summary 获取订单详情
// @Description 根据订单号获取订单详情
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/order/{order_no} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	orderNo := c.Param("order_no")

	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": order,
	})
}

// Pay 支付订单
// @Summary 支付订单
// @Description 模拟支付订单
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/order/{order_no}/pay [post]
func (h *OrderHandler) Pay(c *gin.Context) {
	orderNo := c.Param("order_no")

	if err := h.orderService.PayOrder(orderNo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "支付成功",
	})
}

// Cancel 取消订单
// @Summary 取消订单
// @Description 取消未支付的订单
// @Tags 订单
// @Produce json
// @Param order_no path string true "订单号"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/order/{order_no}/cancel [post]
func (h *OrderHandler) Cancel(c *gin.Context) {
	orderNo := c.Param("order_no")

	if err := h.orderService.CancelOrder(orderNo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "取消成功",
	})
}
