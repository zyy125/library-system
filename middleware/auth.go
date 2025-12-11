package middleware

import (
	"strings"
	"library-system/utils"
	"library-system/common"
	"library-system/repository"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(common.ErrInvalidToken)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(common.ErrUnauthorized)
			c.Abort()
			return
		}
		tokenString := parts[1]

		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			c.Error(common.ErrInvalidToken)
			c.Abort()
			return
		}

		inBlacklist, err := repository.Rdb.IsInBlacklist(c. Request.Context(), claims.TokenID)
		if err != nil || inBlacklist {
			c.Error(common.ErrInvalidToken)
			c.Abort()
			return
		}

        c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("token_id", claims.TokenID)

		c.Next()
	}
}