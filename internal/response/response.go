package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	TraceID string      `json:"trace_id,omitempty"` // 链路追踪ID
}

// PageData 分页数据
type PageData struct {
	List      interface{} `json:"list"`
	Total     int64       `json:"total"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	TotalPage int         `json:"total_page"`
}

// Ok 成功响应（无数据）
func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
	})
}

// OkWithData 成功响应（带数据）
func OkWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// OkWithList 成功响应（带列表数据）
func OkWithList(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data: PageData{
			List:      list,
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			TotalPage: totalPage,
		},
	})
}

// OkWithPage 成功响应（带分页数据）
func OkWithPage(c *gin.Context, data PageData) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, errMsg string) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeServerError,
		Message: errMsg,
	})
}

// FailWithCode 失败响应（带状态码）
func FailWithCode(c *gin.Context, code int, errMsg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: errMsg,
	})
}

// FailWithMsg 失败响应（带详细消息）
func FailWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
	})
}

// FailWithData 失败响应（带数据）
func FailWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = "未登录或登录已过期"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeUnauthorized,
		Message: msg,
	})
}

// Forbidden 无权限
func Forbidden(c *gin.Context, msg string) {
	if msg == "" {
		msg = "没有权限"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeForbidden,
		Message: msg,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	if msg == "" {
		msg = "资源不存在"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeNotFound,
		Message: msg,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, msg string) {
	if msg == "" {
		msg = "请求参数错误"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeBadRequest,
		Message: msg,
	})
}

// ServerError 服务器内部错误
func ServerError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "服务器内部错误"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeServerError,
		Message: msg,
	})
}

// TooManyRequests 请求过于频繁
func TooManyRequests(c *gin.Context, msg string) {
	if msg == "" {
		msg = "请求过于频繁"
	}
	c.JSON(http.StatusOK, Response{
		Code:    CodeTooManyRequests,
		Message: msg,
	})
}
