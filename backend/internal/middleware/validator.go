package middleware

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"gomall/backend/internal/response"
)

// Validator 全局验证器实例
var (
	validate *validator.Validate
	trans    ut.Translator
)

// InitValidator 初始化验证器
func InitValidator() {
	// 创建中英文翻译器
	zhT := zh.New()
	enT := en.New()

	// 通用翻译器
	uni := ut.New(zhT, enT)

	// 默认使用中文翻译
	trans, _ = uni.GetTranslator("zh")

	// 创建验证器
	validate = validator.New()

	// 注册翻译器
	_ = zhTranslations.RegisterDefaultTranslations(validate, trans)
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	// 自定义错误翻译
	registerCustomTranslations()
}

// registerCustomTranslations 注册自定义错误翻译
func registerCustomTranslations() {
	//Required
	trans.Add("required", "{0}不能为空", true)

	// 最小长度
	trans.Add("min", "{0}长度至少为{1}", true)

	// 最大长度
	trans.Add("max", "{0}长度最大为{1}", true)

	// 邮箱
	trans.Add("email", "{0}格式不正确", true)

	// 数值最小
	trans.Add("gt", "{0}必须大于{1}", true)

	// 数值最小等于
	trans.Add("gte", "{0}必须大于等于{1}", true)

	// 数值最大
	trans.Add("lt", "{0}必须小于{1}", true)

	// 数值最大等于
	trans.Add("lte", "{0}必须小于等于{1}", true)

	// 数字范围
	trans.Add("number", "{0}必须是数字", true)

	// 整数
	trans.Add("integer", "{0}必须是整数", true)
}

// ValidatorMiddleware 参数验证中间件
func ValidatorMiddleware(bean interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 初始化验证器
		if validate == nil {
			InitValidator()
		}

		// 获取绑定的结构体
		if bean == nil {
			c.Next()
			return
		}

		// 根据 Content-Type 获取验证数据
		var err error
		contentType := c.ContentType()

		if strings.Contains(contentType, "multipart/form-data") {
			// multipart 表单使用 ShouldBind
			err = c.ShouldBind(bean)
		} else if strings.Contains(contentType, "application/json") {
			// JSON 使用 ShouldBindJSON
			err = c.ShouldBindJSON(bean)
		} else {
			// 其他使用 ShouldBind
			err = c.ShouldBind(bean)
		}

		if err != nil {
			// 获取错误消息
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				response.BadRequest(c, err.Error())
				c.Abort()
				return
			}

			// 翻译错误消息
			errMsg := formatValidationErrors(errs)
			response.BadRequest(c, errMsg)
			c.Abort()
			return
		}

		// 将验证后的结构体存入上下文
		c.Set("validated_data", bean)

		c.Next()
	}
}

// ValidatorForm 表单验证中间件（用于POST表单）
func ValidatorForm(bean interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validate == nil {
			InitValidator()
		}

		if bean == nil {
			c.Next()
			return
		}

		err := c.ShouldBind(bean)
		if err != nil {
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				response.BadRequest(c, err.Error())
				c.Abort()
				return
			}

			errMsg := formatValidationErrors(errs)
			response.BadRequest(c, errMsg)
			c.Abort()
			return
		}

		c.Set("validated_data", bean)
		c.Next()
	}
}

// ValidatorQuery URL参数验证中间件
func ValidatorQuery(bean interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if validate == nil {
			InitValidator()
		}

		if bean == nil {
			c.Next()
			return
		}

		err := c.ShouldBindQuery(bean)
		if err != nil {
			errs, ok := err.(validator.ValidationErrors)
			if !ok {
				response.BadRequest(c, err.Error())
				c.Abort()
				return
			}

			errMsg := formatValidationErrors(errs)
			response.BadRequest(c, errMsg)
			c.Abort()
			return
		}

		c.Set("validated_data", bean)
		c.Next()
	}
}

// formatValidationErrors 格式化验证错误消息
func formatValidationErrors(errs validator.ValidationErrors) string {
	var errMsgs []string

	for _, err := range errs {
		// 获取字段的中文名称
		fieldName := getFieldName(err.Field(), err.StructField())

		// 根据 tag 生成错误消息
		switch err.Tag() {
		case "required":
			errMsgs = append(errMsgs, fieldName+"不能为空")
		case "min":
			errMsgs = append(errMsgs, fieldName+"长度至少为"+err.Param())
		case "max":
			errMsgs = append(errMsgs, fieldName+"长度最大为"+err.Param())
		case "email":
			errMsgs = append(errMsgs, fieldName+"格式不正确")
		case "gt":
			errMsgs = append(errMsgs, fieldName+"必须大于"+err.Param())
		case "gte":
			errMsgs = append(errMsgs, fieldName+"必须大于等于"+err.Param())
		case "lt":
			errMsgs = append(errMsgs, fieldName+"必须小于"+err.Param())
		case "lte":
			errMsgs = append(errMsgs, fieldName+"必须小于等于"+err.Param())
		case "eq":
			errMsgs = append(errMsgs, fieldName+"必须等于"+err.Param())
		case "ne":
			errMsgs = append(errMsgs, fieldName+"不能等于"+err.Param())
		case "len":
			errMsgs = append(errMsgs, fieldName+"长度必须等于"+err.Param())
		case "alpha":
			errMsgs = append(errMsgs, fieldName+"只能包含字母")
		case "alphanum":
			errMsgs = append(errMsgs, fieldName+"只能包含字母和数字")
		case "numeric":
			errMsgs = append(errMsgs, fieldName+"必须是数字")
		case "oneof":
			errMsgs = append(errMsgs, fieldName+"必须是其中一个值: "+err.Param())
		default:
			errMsgs = append(errMsgs, fieldName+"验证失败")
		}
	}

	return strings.Join(errMsgs, "; ")
}

// getFieldName 获取字段的中文名称
func getFieldName(field, structField string) string {
	// 如果 structField 存在，返回结构体字段名
	if structField != "" {
		return structField
	}

	// 否则返回字段名
	return field
}

// GetValidatedData 从上下文中获取验证后的数据
func GetValidatedData(c *gin.Context, bean interface{}) bool {
	data, exists := c.Get("validated_data")
	if !exists {
		return false
	}

	// 类型断言
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	targetVal := reflect.ValueOf(bean).Elem()
	if val.Type() != targetVal.Type() {
		return false
	}

	targetVal.Set(val)
	return true
}

// CustomValidator 自定义验证函数
func CustomValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}
