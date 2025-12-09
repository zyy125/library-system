package common

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidError struct {
	FieldErrors []FieldError
}

func (e *ValidError) Error() string {
	return "参数验证失败"
}

// ValidateStruct 验证结构体并返回错误
func ValidateStruct(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		// 判断是否为验证错误
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			fieldErrors := make([]FieldError, 0)
			for _, e := range validationErrors {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   toSnakeCase(e.Field()),
					Message: getValidationMessage(e),
				})
			}
			return &ValidError{FieldErrors: fieldErrors}
		} else {
			// JSON 解析错误
			return &ValidError{
				FieldErrors: []FieldError{
					{Field: "body", Message: "JSON格式错误"},
				},
			}
		}
	}
	return nil
}

// 获取验证错误的友好消息
func getValidationMessage(e validator. FieldError) string {
	field := toSnakeCase(e. Field())
	switch e.Tag() {
	case "required":
		return fmt. Sprintf("%s字段不能为空", field)
	case "email":
		return "邮箱格式不正确"
	case "min":
		return fmt.Sprintf("%s长度不能少于%s", field, e.Param())
	case "max":
		return fmt.Sprintf("%s长度不能超过%s", field, e.Param())
	case "len":
		return fmt.Sprintf("%s长度必须为%s", field, e. Param())
	case "gte":
		return fmt.Sprintf("%s必须大于等于%s", field, e. Param())
	case "lte":
		return fmt. Sprintf("%s必须小于等于%s", field, e. Param())
	case "oneof":
		return fmt. Sprintf("%s必须是以下值之一:  %s", field, e. Param())
	default:
		return fmt.Sprintf("%s验证失败", field)
	}
}

// 驼峰转蛇形
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}