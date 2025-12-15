package middleware

import(
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"library-system/common"
)

func ErrorHandler() gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 日志
				log.Printf("Panic recovered: %v", err)
				
				// 返回统一的服务器错误
				common.ServerError(c, http.StatusInternalServerError, 500, "服务器内部错误")
				c. Abort()
			}
		}()
		
		c.Next()

		// 处理请求过程中产生的错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			
			// 记录错误日志
			log.Printf("[ERROR] Path: %s, Error: %v", c.Request.URL.Path, err)
			
			// 判断错误类型
			switch e := err.(type) {
			case *common.ValidError:
				// 参数验证错误
				common.ValidationError(c, e. FieldErrors)
				
			case *common.BizError:
				// 业务错误
				common.BusinessError(c, e)
				
			default:
				// 未知错误
				log.Printf("[UNKNOWN_ERROR] %T: %v", err, err)
				common.ServerError(c, http.StatusInternalServerError, 500, "服务器内部错误")
			}
		}
	}
}