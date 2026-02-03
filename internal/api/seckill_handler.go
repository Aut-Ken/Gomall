package api

import (
	"net/http"

	"gomall/internal/middleware"
	"gomall/internal/service"

	"github.com/gin-gonic/gin"
)

// SeckillHandler 秒杀接口处理层
type SeckillHandler struct {
	seckillService *service.SeckillService
}

// NewSeckillHandler 创建秒杀处理器
func NewSeckillHandler() *SeckillHandler {
	return &SeckillHandler{
		seckillService: service.NewSeckillService(),
	}
}

// SeckillRequest 秒杀请求结构
type SeckillRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// Seckill 秒杀接口
// @Summary 秒杀接口
// @Description 参与秒杀活动
// @Tags 秒杀
// @Accept json
// @Produce json
// @Param req body SeckillRequest true "秒杀请求"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/seckill [post]
func (h *SeckillHandler) Seckill(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未登录",
		})
		return
	}

	var req SeckillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	svcReq := &service.SeckillRequest{ProductID: req.ProductID}
	response, err := h.seckillService.SeckillWithRedis(c.Request.Context(), userID, svcReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "秒杀成功",
		"data": response,
	})
}

// InitStock 初始化秒杀库存（管理员接口）
// @Summary 初始化秒杀库存
// @Description 将库存预加载到Redis
// @Tags 秒杀
// @Accept json
// @Produce json
// @Param product_id query int true "商品ID"
// @Param stock query int true "库存数量"
// @Security Bearer
// @Success 200 {object} map[string]interface{}
// @Router /api/seckill/init [post]
func (h *SeckillHandler) InitStock(c *gin.Context) {
	productID := c.Query("product_id")
	stock := c.Query("stock")

	if productID == "" || stock == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 这里应该添加管理员权限验证
	// ...

	if err := h.seckillService.InitSeckillStock(c.Request.Context(), 0, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "初始化失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "初始化成功",
	})
}
