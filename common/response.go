package common

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code 		int 					`json:"code"`
	Message 	string 					`json:"message"`
	Data  		interface{} 			`json:"data,omitempty"`
	Errors		[]FieldError			`json:"errors,omitempty"`			
	Details		map[string]interface{} 	`json:"details,omitempty"`
	Timestamp 	string 					`json:"timestamp"`
}

type FieldError struct {
	Field	string `json:"field"`
	Message string `json:"message"`
}

func Success(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, &Response{
		Code: code,
		Message: msg,
		Data: data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func ServerError(c *gin.Context,httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Message: message,
		Errors: []FieldError{},
		Details: map[string]interface{}{},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
// 字段验证错误
func ValidationError(c *gin.Context, errors []FieldError) {
	c.JSON(http.StatusBadRequest, Response{	
		Code:      400,
		Message:   "参数验证失败",
		Errors:    errors,
		Details:   map[string]interface{}{},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
// 业务错误
func ErrorWithDetails(c *gin.Context, httpStatus int, code int, message string, details map[string]interface{}) {
	if details == nil {
		details = map[string]interface{}{}
	}
	c.JSON(httpStatus, Response{
		Code:      code,
		Message:   message,
		Errors:    []FieldError{},
		Details:   details,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func BusinessError(c *gin.Context, err *BizError) {
	ErrorWithDetails(c, err.HTTPStatus, err.Code, err. Message, err.Details)
}