package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFConfig CSRF中间件配置
type CSRFConfig struct {
	ExpireSeconds    int    // CSRF Token过期时间
	TokenLength      int    // Token长度
	HeaderName       string // 自定义Header名称
	FormFieldName    string // 自定义表单字段名称
	SameSiteMode     string // SameSite模式: "strict", "lax", "none"
}

// DefaultCSRFConfig 默认配置
var DefaultCSRFConfig = CSRFConfig{
	ExpireSeconds:  86400, // 24小时
	TokenLength:    32,
	HeaderName:     "X-CSRF-Token",
	FormFieldName:  "csrf_token",
	SameSiteMode:   "lax",
}

// CSRFToken CSRF Token管理
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFStore CSRF Token存储（内存实现，可扩展为Redis）
type CSRFStore struct {
	tokens map[string]*CSRFToken
}

// NewCSRFStore 创建CSRF存储
func NewCSRFStore() *CSRFStore {
	return &CSRFStore{
		tokens: make(map[string]*CSRFToken),
	}
}

// GenerateToken 生成新的CSRF Token
func (s *CSRFStore) GenerateToken(userID string) string {
	token := generateSecureToken(DefaultCSRFConfig.TokenLength)
	s.tokens[token] = &CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(time.Duration(DefaultCSRFConfig.ExpireSeconds) * time.Second),
	}

	// 启动清理过期Token的goroutine
	go s.cleanupExpiredTokens()

	return token
}

// ValidateToken 验证CSRF Token
func (s *CSRFStore) ValidateToken(token string) bool {
	csrfToken, exists := s.tokens[token]
	if !exists {
		return false
	}

	// 检查是否过期
	if time.Now().After(csrfToken.ExpiresAt) {
		delete(s.tokens, token)
		return false
	}

	return true
}

// RemoveToken 移除Token
func (s *CSRFStore) RemoveToken(token string) {
	delete(s.tokens, token)
}

// cleanupExpiredTokens 清理过期的Token
func (s *CSRFStore) cleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for token, csrfToken := range s.tokens {
			if now.After(csrfToken.ExpiresAt) {
				delete(s.tokens, token)
			}
		}
	}
}

var globalCSRFStore = NewCSRFStore()

// generateSecureToken 生成安全的随机Token
func generateSecureToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// CSRFMiddleware 创建CSRF保护中间件
// 注意：此中间件需要与AuthMiddleware配合使用
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于GET请求，生成并返回CSRF Token
		if c.Request.Method == "GET" {
			userID := GetUserID(c)
			if userID > 0 {
				token := globalCSRFStore.GenerateToken(fmt.Sprintf("%d", userID))
				c.Header("X-CSRF-Token", token)
				c.Set("csrf_token", token)
			}
			c.Next()
			return
		}

		// 对于非GET请求，验证CSRF Token
		var token string

		// 优先从Header获取
		token = c.GetHeader(DefaultCSRFConfig.HeaderName)
		if token == "" {
			// 其次从表单获取
			token = c.PostForm(DefaultCSRFConfig.FormFieldName)
		}
		if token == "" {
			// 从Authorization Header中提取（兼容某些前端框架）
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 {
					token = parts[1]
				}
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "缺少CSRF Token",
			})
			return
		}

		if !globalCSRFStore.ValidateToken(token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "CSRF Token无效或已过期",
			})
			return
		}

		c.Next()
	}
}

// GetCSRFToken 获取当前请求的CSRF Token
func GetCSRFToken(c *gin.Context) string {
	if token, exists := c.Get("csrf_token"); exists {
		return token.(string)
	}
	return ""
}
