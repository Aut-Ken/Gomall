package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"gomall/backend/internal/config"
	"gomall/backend/internal/repository"
)

// WeChatPayService 微信支付服务
type WeChatPayService struct {
	orderRepo *repository.OrderRepository
}

// NewWeChatPayService 创建微信支付服务实例
func NewWeChatPayService() *WeChatPayService {
	return &WeChatPayService{
		orderRepo: repository.NewOrderRepository(),
	}
}

// WeChatPayConfig 微信支付配置
type WeChatPayConfig struct {
	AppID     string
	MchID     string
	Key       string
	NotifyURL string
	TradeType string
	Sandbox   bool
}

// GetConfig 获取微信支付配置
func (s *WeChatPayService) GetConfig() *WeChatPayConfig {
	wxConfig := config.Config.Sub("wechat")
	if wxConfig == nil {
		return nil
	}

	return &WeChatPayConfig{
		AppID:     wxConfig.GetString("appid"),
		MchID:     wxConfig.GetString("mch_id"),
		Key:       wxConfig.GetString("key"),
		NotifyURL: wxConfig.GetString("notify_url"),
		TradeType: wxConfig.GetString("trade_type"),
		Sandbox:   wxConfig.GetBool("sandbox"),
	}
}

// UnifiedOrderRequest 统一下单请求
type UnifiedOrderRequest struct {
	AppID          string `xml:"appid"`
	MchID          string `xml:"mch_id"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Body           string `xml:"body"`
	OutTradeNo     string `xml:"out_trade_no"`
	TotalFee       int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyURL      string `xml:"notify_url"`
	TradeType      string `xml:"trade_type"`
}

// UnifiedOrderResponse 统一下单响应
type UnifiedOrderResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	CodeURL    string `xml:"code_url"`
	TradeType  string `xml:"trade_type"`
}

// PayNotifyRequest 支付回调请求
type PayNotifyRequest struct {
	AppID        string `xml:"appid"`
	MchID        string `xml:"mch_id"`
	NonceStr     string `xml:"nonce_str"`
	Sign         string `xml:"sign"`
	ResultCode  string `xml:"result_code"`
	OutTradeNo   string `xml:"out_trade_no"`
	TransactionID string `xml:"transaction_id"`
	TradeType    string `xml:"trade_type"`
	TradeState   string `xml:"trade_state"`
	TotalFee     int    `xml:"total_fee"`
	BankType     string `xml:"bank_type"`
	TimeEnd      string `xml:"time_end"`
}

// PayNotifyResponse 支付回调响应
type PayNotifyResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// UnifiedOrder 统一下单
func (s *WeChatPayService) UnifiedOrder(ctx context.Context, orderNo string, totalFee int, body string) (*UnifiedOrderResponse, error) {
	cfg := s.GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("微信支付配置不存在")
	}

	// 生成随机字符串
	nonceStr := generateNonceStr(32)

	// 构建请求
	req := UnifiedOrderRequest{
		AppID:          cfg.AppID,
		MchID:          cfg.MchID,
		NonceStr:       nonceStr,
		Body:           body,
		OutTradeNo:     orderNo,
		TotalFee:       totalFee,
		SpbillCreateIP: "127.0.0.1",
		NotifyURL:      cfg.NotifyURL,
		TradeType:      cfg.TradeType,
	}

	// 生成签名
	sign := generateSign(req, cfg.Key)
	req.Sign = sign

	// 序列化为XML
	xmlData, err := xml.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("XML序列化失败: %w", err)
	}

	// 添加XML头
	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	// 确定请求URL
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	if cfg.Sandbox {
		url = "https://api.mch.weixin.qq.com/sandbox/pay/unifiedorder"
	}

	// 发送请求
	resp, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("请求微信支付失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result UnifiedOrderResponse
	if err := xml.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("XML解析失败: %w", err)
	}

	// 检查返回状态
	if result.ReturnCode != "SUCCESS" {
		return nil, fmt.Errorf("微信支付返回错误: %s", result.ReturnMsg)
	}

	// 检查业务结果
	if result.ResultCode != "SUCCESS" {
		return nil, fmt.Errorf("微信支付业务处理失败")
	}

	return &result, nil
}

// PayNotify 支付回调处理
func (s *WeChatPayService) PayNotify(ctx context.Context, data []byte) (*PayNotifyResponse, error) {
	cfg := s.GetConfig()
	if cfg == nil {
		return &PayNotifyResponse{
			ReturnCode: "FAIL",
			ReturnMsg:  "配置不存在",
		}, nil
	}

	// 解析请求
	var req PayNotifyRequest
	if err := xml.Unmarshal(data, &req); err != nil {
		return &PayNotifyResponse{
			ReturnCode: "FAIL",
			ReturnMsg:  "XML解析失败",
		}, nil
	}

	// 验证签名
	if !verifySign(req, cfg.Key) {
		return &PayNotifyResponse{
			ReturnCode: "FAIL",
			ReturnMsg:  "签名验证失败",
		}, nil
	}

	// 更新订单状态
	if req.ResultCode == "SUCCESS" {
		order, err := s.orderRepo.GetByOrderNo(req.OutTradeNo)
		if err != nil {
			return &PayNotifyResponse{
				ReturnCode: "FAIL",
				ReturnMsg:  "订单不存在",
			}, nil
		}

		// 更新订单状态为已支付
		order.Status = 2 // 已支付
		order.PayType = 2 // 微信支付
		if err := s.orderRepo.Update(order); err != nil {
			return &PayNotifyResponse{
				ReturnCode: "FAIL",
				ReturnMsg:  "订单更新失败",
			}, nil
		}
	}

	return &PayNotifyResponse{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	}, nil
}

// QueryOrder 查询订单
func (s *WeChatPayService) QueryOrder(ctx context.Context, orderNo string) (*UnifiedOrderResponse, error) {
	cfg := s.GetConfig()
	if cfg == nil {
		return nil, fmt.Errorf("微信支付配置不存在")
	}

	// 生成随机字符串
	nonceStr := generateNonceStr(32)

	// 构建请求
	req := map[string]string{
		"appid":        cfg.AppID,
		"mch_id":       cfg.MchID,
		"out_trade_no": orderNo,
		"nonce_str":    nonceStr,
	}

	// 生成签名
	sign := generateSignFromMap(req, cfg.Key)
	req["sign"] = sign

	// 序列化为XML
	xmlData, _ := xml.Marshal(req)
	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	// 发送请求
	url := "https://api.mch.weixin.qq.com/pay/orderquery"
	resp, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("请求微信支付查询失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var result UnifiedOrderResponse
	if err := xml.Unmarshal(respData, &result); err != nil {
		return nil, fmt.Errorf("XML解析失败: %w", err)
	}

	return &result, nil
}

// CloseOrder 关闭订单
func (s *WeChatPayService) CloseOrder(ctx context.Context, orderNo string) error {
	cfg := s.GetConfig()
	if cfg == nil {
		return fmt.Errorf("微信支付配置不存在")
	}

	// 生成随机字符串
	nonceStr := generateNonceStr(32)

	// 构建请求
	req := map[string]string{
		"appid":        cfg.AppID,
		"mch_id":       cfg.MchID,
		"out_trade_no": orderNo,
		"nonce_str":    nonceStr,
	}

	// 生成签名
	sign := generateSignFromMap(req, cfg.Key)
	req["sign"] = sign

	// 序列化为XML
	xmlData, _ := xml.Marshal(req)
	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	// 发送请求
	url := "https://api.mch.weixin.qq.com/pay/closeorder"
	resp, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewReader(xmlData))
	if err != nil {
		return fmt.Errorf("请求关闭订单失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查返回状态
	var result struct {
		ReturnCode string `xml:"return_code"`
		ReturnMsg  string `xml:"return_msg"`
	}
	if err := xml.Unmarshal(respData, &result); err != nil {
		return fmt.Errorf("XML解析失败: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return fmt.Errorf("关闭订单失败: %s", result.ReturnMsg)
	}

	return nil
}

// Refund 退款
func (s *WeChatPayService) Refund(ctx context.Context, orderNo, refundNo string, totalFee, refundFee int) error {
	cfg := s.GetConfig()
	if cfg == nil {
		return fmt.Errorf("微信支付配置不存在")
	}

	// 生成随机字符串
	nonceStr := generateNonceStr(32)

	// 构建请求
	req := map[string]string{
		"appid":         cfg.AppID,
		"mch_id":        cfg.MchID,
		"nonce_str":     nonceStr,
		"out_trade_no":  orderNo,
		"out_refund_no": refundNo,
		"total_fee":     fmt.Sprintf("%d", totalFee),
		"refund_fee":    fmt.Sprintf("%d", refundFee),
	}

	// 生成签名
	sign := generateSignFromMap(req, cfg.Key)
	req["sign"] = sign

	// 序列化为XML
	xmlData, _ := xml.Marshal(req)
	xmlData = append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`), xmlData...)

	// 发送请求（需要证书）
	url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	resp, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewReader(xmlData))
	if err != nil {
		return fmt.Errorf("请求退款失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查返回状态
	var result struct {
		ReturnCode string `xml:"return_code"`
		ReturnMsg  string `xml:"return_msg"`
	}
	if err := xml.Unmarshal(respData, &result); err != nil {
		return fmt.Errorf("XML解析失败: %w", err)
	}

	if result.ReturnCode != "SUCCESS" {
		return fmt.Errorf("退款失败: %s", result.ReturnMsg)
	}

	return nil
}

// generateNonceStr 生成随机字符串
func generateNonceStr(length int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = str[rand.Intn(len(str))]
	}
	return string(result)
}

// generateSign 生成签名
func generateSign(req UnifiedOrderRequest, key string) string {
	// 获取所有字段
	m := make(map[string]string)
	m["appid"] = req.AppID
	m["mch_id"] = req.MchID
	m["nonce_str"] = req.NonceStr
	m["body"] = req.Body
	m["out_trade_no"] = req.OutTradeNo
	m["total_fee"] = fmt.Sprintf("%d", req.TotalFee)
	m["spbill_create_ip"] = req.SpbillCreateIP
	m["notify_url"] = req.NotifyURL
	m["trade_type"] = req.TradeType

	return generateSignFromMap(m, key)
}

// generateSignFromMap 从map生成签名
func generateSignFromMap(m map[string]string, key string) string {
	// 获取所有key
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signStr string
	for _, k := range keys {
		value := m[k]
		if value == "" || k == "sign" {
			continue
		}
		signStr += fmt.Sprintf("%s=%s&", k, value)
	}

	// 添加密钥
	signStr += "key=" + key

	// MD5签名
	return strings.ToUpper(md5Hash(signStr))
}

// verifySign 验证签名
func verifySign(req PayNotifyRequest, key string) bool {
	m := make(map[string]string)
	m["appid"] = req.AppID
	m["mch_id"] = req.MchID
	m["nonce_str"] = req.NonceStr
	m["result_code"] = req.ResultCode
	m["out_trade_no"] = req.OutTradeNo
	m["transaction_id"] = req.TransactionID
	m["trade_type"] = req.TradeType
	m["trade_state"] = req.TradeState
	m["total_fee"] = fmt.Sprintf("%d", req.TotalFee)
	m["bank_type"] = req.BankType
	m["time_end"] = req.TimeEnd

	calcSign := generateSignFromMap(m, key)
	return calcSign == req.Sign
}

// md5Hash MD5哈希
func md5Hash(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
