package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

// sensitiveFields 敏感字段正则匹配
var sensitiveFields = regexp.MustCompile(`(?i)(password|secret|token|key|authorization|credit_card|cvv|ssn)`)

// MaskedValue 脱敏后的值
const MaskedValue = "***MASKED***"

// LogSanitizerMiddleware 日志脱敏中间件
// 自动脱敏请求体中的敏感信息
func LogSanitizerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取请求体
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 脱敏请求体
		if len(bodyBytes) > 0 {
			sanitizedBody := sanitizeJSON(bodyBytes)
			c.Set("sanitized_body", sanitizedBody)
		}

		c.Next()
	}
}

// sanitizeJSON 脱敏JSON中的敏感信息
func sanitizeJSON(data []byte) []byte {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return data // 如果解析失败，返回原始数据
	}

	sanitizeValue(raw)
	sanitized, _ := json.Marshal(raw)
	return sanitized
}

// sanitizeValue 递归脱敏值
func sanitizeValue(v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v := range val {
			if sensitiveFields.MatchString(k) {
				val[k] = MaskedValue
			} else {
				sanitizeValue(v)
			}
		}
	case []interface{}:
		for i, item := range val {
			sanitizeValue(item)
			val[i] = item
		}
	}
}

// SanitizeLogMessage 脱敏日志消息
func SanitizeLogMessage(msg string) string {
	// 脱敏URL中的敏感参数
	msg = sanitizeURLParams(msg)
	// 脱敏JSON字符串中的敏感信息
	msg = sanitizeString(msg)
	return msg
}

// sanitizeURLParams 脱敏URL中的敏感参数
func sanitizeURLParams(url string) string {
	// 检查URL中是否包含敏感参数
	patterns := []string{
		`([?&])(password|token|key|secret)=([^&]*)`,
		`([?&])(authorization)=([^&]*)`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		url = re.ReplaceAllString(url, fmt.Sprintf("$1%s=%s", "$2", MaskedValue))
	}
	return url
}

// sanitizeString 脱敏字符串中的敏感信息
func sanitizeString(s string) string {
	// 脱敏邮箱
	emailRe := regexp.MustCompile(`([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
	s = emailRe.ReplaceAllStringFunc(s, func(match string) string {
		parts := emailRe.FindStringSubmatch(match)
		if len(parts) >= 3 {
			return parts[1][:2] + "***@" + parts[2]
		}
		return match
	})

	// 脱敏手机号
	phoneRe := regexp.MustCompile(`1[3-9]\d{9}`)
	s = phoneRe.ReplaceAllStringFunc(s, func(match string) string {
		return match[:3] + "****" + match[7:]
	})

	// 脱敏身份证号
	idCardRe := regexp.MustCompile(`\d{17}[\dXx]`)
	s = idCardRe.ReplaceAllStringFunc(s, func(match string) string {
		return match[:6] + "********" + match[14:]
	})

	return s
}

// SecurityHeadersMiddleware 添加安全响应头
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止XSS攻击
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")

		// 严格传输安全
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
// 为每个请求生成唯一ID，便于追踪
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID 生成唯一请求ID
func generateRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		return fmt.Sprintf("%d-%d", time.Now().UnixNano(), time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x", time.Now().UnixNano(), b)
}
