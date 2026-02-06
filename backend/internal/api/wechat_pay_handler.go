package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"gomall/backend/internal/middleware"
	"gomall/backend/internal/response"
	"gomall/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// WeChatPayHandler 微信支付接口处理层
type WeChatPayHandler struct {
	wechatPayService *service.WeChatPayService
	orderService     *service.OrderService
}

// NewWeChatPayHandler 创建微信支付处理器
func NewWeChatPayHandler() *WeChatPayHandler {
	return &WeChatPayHandler{
		wechatPayService: service.NewWeChatPayService(),
		orderService:     service.NewOrderService(),
	}
}

// WeChatPayRequest 微信支付统一下单请求
type WeChatPayRequest struct {
	OrderNo string `json:"order_no" binding:"required"`
}

// WeChatPayResponse 微信支付统一下单响应
type WeChatPayResponse struct {
	PrepayID  string `json:"prepay_id"`
	CodeURL   string `json:"code_url"`
	OrderNo   string `json:"order_no"`
	TotalFee  int    `json:"total_fee"`
	TradeType string `json:"trade_type"`
}

// UnifiedOrder 统一下单
// @Summary 微信支付统一下单
// @Description 创建微信支付订单
// @Tags 微信支付
// @Accept json
// @Produce json
// @Param req body WeChatPayRequest true "订单信息"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/pay/wechat/unified-order [post]
func (h *WeChatPayHandler) UnifiedOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	var req WeChatPayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetOrderByNo(req.OrderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderNotFound, "订单不存在")
		return
	}

	// 验证订单用户
	if order.UserID != userID {
		response.Forbidden(c, "无权操作此订单")
		return
	}

	// 检查订单状态
	if order.Status != 1 {
		response.FailWithMsg(c, response.CodeOrderStatusError, "订单状态不允许支付")
		return
	}

	// 计算支付金额（分）
	totalFee := int(order.TotalPrice * 100)

	// 调用微信统一下单
	result, err := h.wechatPayService.UnifiedOrder(c.Request.Context(), req.OrderNo, totalFee, order.ProductName)
	if err != nil {
		response.FailWithMsg(c, response.CodePayCreateFailed, err.Error())
		return
	}

	response.OkWithData(c, WeChatPayResponse{
		PrepayID:  result.PrepayID,
		CodeURL:   result.CodeURL,
		OrderNo:   req.OrderNo,
		TotalFee:  totalFee,
		TradeType: result.TradeType,
	})
}

// Notify 支付回调
// @Summary 微信支付回调
// @Description 处理微信支付结果通知
// @Tags 微信支付
// @Accept xml
// @Produce xml
// @Success 200 {object} response.Response
// @Router /api/pay/wechat/notify [post]
func (h *WeChatPayHandler) Notify(c *gin.Context) {
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.FailWithMsg(c, response.CodePayParamError, "读取请求体失败")
		return
	}

	// 处理支付回调
	result, err := h.wechatPayService.PayNotify(c.Request.Context(), body)
	if err != nil {
		response.FailWithMsg(c, response.CodePayQueryFailed, err.Error())
		return
	}

	// 返回XML响应
	c.XML(http.StatusOK, result)
}

// QueryOrder 查询订单
// @Summary 查询微信支付订单
// @Description 查询订单支付状态
// @Tags 微信支付
// @Accept json
// @Produce json
// @Param order_no query string true "订单号"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/pay/wechat/query [get]
func (h *WeChatPayHandler) QueryOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	orderNo := c.Query("order_no")
	if orderNo == "" {
		response.BadRequest(c, "订单号不能为空")
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderNotFound, "订单不存在")
		return
	}

	// 验证订单用户
	if order.UserID != userID {
		response.Forbidden(c, "无权操作此订单")
		return
	}

	// 调用微信查询订单
	result, err := h.wechatPayService.QueryOrder(c.Request.Context(), orderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodePayQueryFailed, err.Error())
		return
	}

	// 解析支付状态
	tradeState := "UNKNOWN"
	if result.ResultCode == "SUCCESS" {
		switch result.TradeType {
		case "SUCCESS":
			tradeState = "PAID"
		case "NOTPAY":
			tradeState = "NOT_PAID"
		case "CLOSED":
			tradeState = "CLOSED"
		default:
			tradeState = result.TradeType
		}
	}

	response.OkWithData(c, gin.H{
		"order_no":   orderNo,
		"trade_state": tradeState,
		"prepay_id":  result.PrepayID,
		"trade_type": result.TradeType,
	})
}

// CloseOrder 关闭订单
// @Summary 关闭微信支付订单
// @Description 关闭未支付的订单
// @Tags 微信支付
// @Accept json
// @Produce json
// @Param order_no query string true "订单号"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/pay/wechat/close [post]
func (h *WeChatPayHandler) CloseOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	orderNo := c.Query("order_no")
	if orderNo == "" {
		response.BadRequest(c, "订单号不能为空")
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderNotFound, "订单不存在")
		return
	}

	// 验证订单用户
	if order.UserID != userID {
		response.Forbidden(c, "无权操作此订单")
		return
	}

	// 调用微信关闭订单
	if err := h.wechatPayService.CloseOrder(c.Request.Context(), orderNo); err != nil {
		response.FailWithMsg(c, response.CodeOrderCancelFailed, err.Error())
		return
	}

	response.Ok(c)
}

// Refund 申请退款
// @Summary 申请微信支付退款
// @Description 申请订单退款
// @Tags 微信支付
// @Accept json
// @Produce json
// @Param order_no query string true "订单号"
// @Param refund_fee query int true "退款金额"
// @Security Bearer
// @Success 200 {object} response.Response
// @Router /api/pay/wechat/refund [post]
func (h *WeChatPayHandler) Refund(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	orderNo := c.Query("order_no")
	if orderNo == "" {
		response.BadRequest(c, "订单号不能为空")
		return
	}

	refundFeeStr := c.Query("refund_fee")
	if refundFeeStr == "" {
		response.BadRequest(c, "退款金额不能为空")
		return
	}

	refundFee, err := strconv.Atoi(refundFeeStr)
	if err != nil {
		response.BadRequest(c, "退款金额格式错误")
		return
	}

	// 获取订单信息
	order, err := h.orderService.GetOrderByNo(orderNo)
	if err != nil {
		response.FailWithMsg(c, response.CodeOrderNotFound, "订单不存在")
		return
	}

	// 验证订单用户
	if order.UserID != userID {
		response.Forbidden(c, "无权操作此订单")
		return
	}

	// 检查订单状态
	if order.Status != 2 {
		response.FailWithMsg(c, response.CodeOrderStatusError, "只有已支付的订单可以退款")
		return
	}

	// 生成退款单号
	refundNo := fmt.Sprintf("REF%s%s", orderNo[3:], strconv.FormatInt(int64(userID), 10))

	// 计算支付金额（分）
	totalFee := int(order.TotalPrice * 100)

	// 调用微信退款
	if err := h.wechatPayService.Refund(c.Request.Context(), orderNo, refundNo, totalFee, refundFee); err != nil {
		response.FailWithMsg(c, response.CodePayRefundFailed, err.Error())
		return
	}

	response.OkWithData(c, gin.H{
		"order_no":   orderNo,
		"refund_no":  refundNo,
		"refund_fee": refundFee,
	})
}
