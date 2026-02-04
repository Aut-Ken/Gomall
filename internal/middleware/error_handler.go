package middleware

import (
	"errors"
	"net/http"

	"gomall/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AppError 应用错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

// 预定义错误
var (
	ErrUnauthorized     = &AppError{Code: 401, Message: "未登录或登录已过期"}
	ErrForbidden        = &AppError{Code: 403, Message: "没有权限"}
	ErrNotFound         = &AppError{Code: 404, Message: "资源不存在"}
	ErrValidationFailed = &AppError{Code: 400, Message: "参数验证失败"}
	ErrInternalServer   = &AppError{Code: 500, Message: "服务器内部错误"}
	ErrTooManyRequests  = &AppError{Code: 429, Message: "请求过于频繁"}
)

// ErrorHandlerMiddleware 全局错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 如果有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 检查是否为 AppError
			var appErr *AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.Code, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
					"details": appErr.Details,
				})
				return
			}

			// 其他错误
			logger.Error("请求处理失败",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.String("error", err.Error()),
			)

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "服务器内部错误",
			})
		}
	}
}

// RenderError 渲染错误响应
func RenderError(c *gin.Context, err error) {
	c.Error(err)
}

// RenderValidationError 渲染参数验证错误
func RenderValidationError(c *gin.Context, details string) {
	RenderError(c, &AppError{
		Code:    400,
		Message: "参数验证失败",
		Details: details,
	})
}

// RenderUnauthorized 渲染未授权错误
func RenderUnauthorized(c *gin.Context) {
	RenderError(c, ErrUnauthorized)
}

// RenderForbidden 渲染禁止访问错误
func RenderForbidden(c *gin.Context) {
	RenderError(c, ErrForbidden)
}

// RenderNotFound 渲染资源不存在错误
func RenderNotFound(c *gin.Context, resource string) {
	RenderError(c, &AppError{
		Code:    404,
		Message: resource + "不存在",
	})
}

// RenderServerError 渲染服务器内部错误
func RenderServerError(c *gin.Context, details string) {
	RenderError(c, &AppError{
		Code:    500,
		Message: "服务器内部错误",
		Details: details,
	})
}

// RenderSuccess 渲染成功响应
func RenderSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

// RenderList 渲染列表响应
func RenderList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
