package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gomall/backend/internal/config"
	"gomall/backend/pkg/jwt"
)

// AuthMiddleware JWT认证中间件
// 用于保护需要登录才能访问的接口
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		// 验证Token格式 (Bearer token)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析Token
		jwtUtil := jwt.NewJWT()
		claims, err := jwtUtil.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// GetUserID 从上下文中获取当前用户ID
func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文中获取当前用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// AdminAuthMiddleware 管理员权限认证中间件
// 用于保护需要管理员权限才能访问的接口
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行登录认证
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		jwtUtil := jwt.NewJWT()
		claims, err := jwtUtil.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Token无效或已过期",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		// 从配置获取管理员ID列表，支持多个管理员
		appConfig := config.GetApp()
		adminIDsStr := appConfig.GetString("admin_ids")
		var adminIDs []uint
		if adminIDsStr != "" {
			// 支持逗号分隔的多个管理员ID
			for _, idStr := range strings.Split(adminIDsStr, ",") {
				if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
					adminIDs = append(adminIDs, uint(id))
				}
			}
		}

		// 默认管理员ID为1（向后兼容）
		if len(adminIDs) == 0 {
			adminIDs = []uint{1}
		}

		// 检查当前用户是否为管理员
		isAdmin := false
		for _, adminID := range adminIDs {
			if claims.UserID == adminID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
